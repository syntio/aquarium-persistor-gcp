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
	"fmt"
	"time"

	"cloud.google.com/go/storage"
)

// PersistData stores a message in GCS bucket. Before storing, the unique file is created
// using information from storage configuration. Each message is written in a separate file.
// Returned result is an error which defines the validity of the function action.
func PersistData(ctx context.Context, data []byte, info StorageInfo) error {
	var err error

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	objectName := FileName(info)

	ctxx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objectWriter := client.Bucket(info.BucketID).Object(objectName).NewWriter(ctxx)
	if _, err := objectWriter.Write(data); err != nil {
		_ = ResourceCloser(client, objectWriter)
		return err
	}

	err = ResourceCloser(client, objectWriter)
	return err
}

// FileName constructs a name of file in which the message will be written.
// Returned value is a string which consists of current date and hour, followed by chosen prefix, message ID and file extension.
func FileName(info StorageInfo) string {
	currentTime := time.Now()
	formattedTime := fmt.Sprintf("%02d/%02d/%02d/%02d", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour())
	objectName := fmt.Sprintf("%s/%s-%s.%s", formattedTime, info.Prefix, info.MessageID, info.Extension)
	return objectName
}

// ResourceCloser closes the client and writer components of the GCP storage service.
// Returned result is an error which defines the validity of the function action.
func ResourceCloser(client *storage.Client, writer *storage.Writer) error {
	var err error

	err = writer.Close()
	if err != nil {
		return err
	}
	err = client.Close()

	return err
}
