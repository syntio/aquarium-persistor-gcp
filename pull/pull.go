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

package pull

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/syntio/aquarium-persistor-gcp/lib"
)

// synchronous identifies whether it is pull or streaming pull, for pull synchronous is set to true
const synchronous = true

// PullHandler represents the main pull function which is triggered by the HTTP request.
// It creates pull, subscriber and storage configurations that are passed to Puller for
// pulling and storing messages from Pub/Sub, using synchronous pull.
func PullHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	ctx := context.Background()

	var storageInfo lib.StorageInfo
	err = lib.SetStorageInfo(&storageInfo)
	if err != nil {
		errorMessage(err)
		panic(err)
	}

	var pullInfo lib.PullInfo
	err = lib.SetPullInfo(&pullInfo)
	if err != nil {
		errorMessage(err)
		panic(err)
	}

	err = lib.ExtractReceivedInfo(r, &pullInfo)
	if err != nil {
		log.Printf("Error while reading received info. %v", err)
		panic(err)
	}

	var subscriberConf lib.SubConf
	err = lib.SetSubscriberConf(&subscriberConf, synchronous)
	if err != nil {
		errorMessage(err)
		panic(err)
	}

	err = lib.Pull(ctx, &pullInfo, storageInfo, &subscriberConf)
	if err != nil {
		log.Printf("Error during pubsub pulling. %s.\n", err)
		panic(err)
	}
	fmt.Fprint(w, "Finished execution")

}

// errorMessage represents a helper function for logging errors that occurred during pulling and storing messages.
func errorMessage(err error) {
	log.Printf("Error during retrieving environment variables. %s.\n", err)
}
