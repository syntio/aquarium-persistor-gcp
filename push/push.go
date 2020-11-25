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

package push

import (
	"context"
	"log"

	cfmetadata "cloud.google.com/go/functions/metadata"

	"github.com/syntio/aquarium-persistor-gcp/lib"
)

// PubsubMessage is a helper structure used for fetching incoming messages data
type PubsubMessage struct {
	Data []byte `json:"data"`
}

// PushHandler represents entry point for processing Pub/Sub push trigger.
// The function creates storage configuration and calls helper function which stores the message.
// Returned result is an error which defines the validity of the function action.
func PushHandler(ctx context.Context, message PubsubMessage) error {
	var err error

	metadata, err := cfmetadata.FromContext(ctx)
	if err != nil {
		log.Printf("Error during metadata unmarshaling. %s.\n", err)
		return err
	}

	//EventID is a unique ID for the event (message).
	storageInfo := lib.StorageInfo{
		MessageID: metadata.EventID,
	}

	err = lib.SetStorageInfo(&storageInfo)
	if err != nil {
		log.Printf("Error during retrieving environment variables. %s.\n", err)
		return err
	}

	err = lib.PersistData(ctx, message.Data, storageInfo)
	if err != nil {
		log.Printf("Error during data storage. %s.\n", err)
		return err
	}

	return nil
}
