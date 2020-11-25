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

// StorageInfo represents storage configuration.
// It holds information needed for storing messages to GCS.
type StorageInfo struct {
	MessageID string // ID of a message (used for naming a file in which the message will be written)
	BucketID  string // ID of a bucket in which messages will be stored
	Prefix    string // prefix of a file name
	Extension string // file extension (txt, json, yaml, etc.)
}

// SetStorageInfo sets the parameters of a storage config.
// An error is returned if any errors occur during the function execution.
func SetStorageInfo(storageInfo *StorageInfo) error {
	var err error

	storageInfo.BucketID, err = getEnvVariable("BUCKET_ID")
	if err != nil {
		return err
	}

	storageInfo.Prefix, err = getEnvVariable("MSG_PREFIX")
	if err != nil {
		return err
	}

	storageInfo.Extension, err = getEnvVariable("MSG_EXTENSION")
	if err != nil {
		return err
	}

	return err
}
