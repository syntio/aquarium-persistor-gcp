// Copyright 2020 Syntio Inc.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lib

import (
	"context"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Pull function pulls messages from provided Pub/Sub subscription and calls storing function on each pulled message.
// As a part of a process, Pull creates a client that will receive blocks of messages.
// Received blocks will be of a limited size if synchronous option is enabled.
// Synchronous pull stores fixed number of messages and cancels the context which prevents further receiving.
// If the streaming pull option is chosen, the client receives blocks of a variable sizes until context duration expires.
// An error is returned if any errors occur during the function execution.
func Pull(ctx context.Context, info *PullInfo, storageInfo StorageInfo, subConf *SubConf) error {

	var err error

	client, err := pubsub.NewClient(ctx, info.ProjectID)
	if err != nil {
		return err
	}
	defer client.Close()

	// Configuring subscriber.
	// Depending on which of two pull option is chosen, receiving settings are set to different values.
	sub := client.Subscription(info.SubID)
	sub.ReceiveSettings.Synchronous = subConf.Synchronous
	sub.ReceiveSettings.MaxExtension = time.Duration(subConf.MaxExtension) * time.Second
	sub.ReceiveSettings.MaxOutstandingMessages = subConf.MaxOutstandingMessages
	sub.ReceiveSettings.MaxOutstandingBytes = subConf.MaxOutstandingBytes
	sub.ReceiveSettings.NumGoroutines = subConf.NumOfGoroutines

	// Receive messages for NumberOfSeconds period.
	ctxx, cancel := context.WithTimeout(ctx, time.Duration(info.NumberOfSeconds)*time.Second)
	defer cancel()

	// Create a channel to handle messages to as they come in.
	cm := make(chan *pubsub.Message)

	go func() {
		var mtx = &sync.Mutex{}
		var messageCounter int = 0

		for {
			select {
			case msg := <-cm:
				mtx.Lock()

				storageInfo.MessageID = msg.ID

				if err := PersistData(ctx, msg.Data, storageInfo); err != nil {
					panic(err)
				}

				messageCounter++
				msg.Ack()

				if subConf.Synchronous {
					// If max message count is exceeded then cancel the context.
					if messageCounter >= info.NumberOfMessages {
						cancel()
					}
				}
				mtx.Unlock()
			case <-ctxx.Done():
				return
			}
		}
	}()

	// Receive blocks until the passed in context exceeds.
	recvErr := sub.Receive(ctxx, func(ctxx context.Context, msg *pubsub.Message) {
		cm <- msg
	})

	if recvErr != nil && status.Code(recvErr) != codes.Canceled {
		log.Printf("Receive: %v.\n", recvErr)
	}

	close(cm)
	return err
}
