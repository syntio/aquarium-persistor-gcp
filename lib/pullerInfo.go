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
	"encoding/json"
	"net/http"
	"strconv"
)

// PullInfo represents pull configuration.
// It holds information needed for Pull function.
type PullInfo struct {
	ProjectID        string // the ID of a project in which the topic is located
	SubID            string // the ID of a subscription the messages will be pulled from
	NumberOfMessages int    // the number of messages pull function will persist in one call (this applies only to the synchronous version)
	NumberOfSeconds  int    // the time duration in which the messages will be received
}

// SubConf represents subscriber configuration.
// It holds receive information.
type SubConf struct {
	Synchronous            bool
	MaxExtension           int
	MaxOutstandingMessages int
	MaxOutstandingBytes    int
	NumOfGoroutines        int
}

// recievedInfo is a helper structure which is used for extraction of information received through an HTTP request.
type recievedInfo struct {
	NumberOfMessages string
	NumberOfSeconds  int
}

// SetPullInfo sets the values of the pull configuration by extracting the values ​​from the corresponding environment variables.
// An error is returned if any error occurs during the function execution.
func SetPullInfo(pullInfo *PullInfo) error {
	var err error

	pullInfo.ProjectID, err = getEnvVariable("PROJECT_ID")
	if err != nil {
		return err
	}

	pullInfo.SubID, err = getEnvVariable("SUB_ID")
	if err != nil {
		return err
	}

	return err
}

// ExtractReceivedInfo extracts the information received through the HTTP request body needed for the pull configuration.
// An error is returned if any errors occur during the function execution.
func ExtractReceivedInfo(r *http.Request, pullInfo *PullInfo) error {
	decoder := json.NewDecoder(r.Body)
	var recInfo recievedInfo

	decoder.DisallowUnknownFields()
	err := decoder.Decode(&recInfo)
	if err != nil {
		return err
	}

	pullInfo.NumberOfSeconds = recInfo.NumberOfSeconds

	pullInfo.NumberOfMessages, err = strconv.Atoi(recInfo.NumberOfMessages)
	if err != nil {
		return err
	}

	return nil
}

// SetSubscriberConf sets the parameters of a subscriber configuration by extracting values ​​from the corresponding environment variables.
// The synchronous and maxExtension parameters depend on the type of pull which is determined by passed bool variable.
// An error is returned if any errors occur during the function execution.
func SetSubscriberConf(subscriberConf *SubConf, synchronous bool) error {
	var err error

	subscriberConf.Synchronous = synchronous

	if !synchronous {
		subscriberConf.MaxExtension = -1
	} else {

		maxExtension, err := getEnvVariable("MAX_EXTENSION")
		if err != nil {
			return err
		}
		// Converting strings to integers.
		subscriberConf.MaxExtension, err = strconv.Atoi(maxExtension)
		if err != nil {
			return err
		}

	}

	maxOutstandingMessages, err := getEnvVariable("MAX_OUTSTANDING_MSGS")
	if err != nil {
		return err
	}
	subscriberConf.MaxOutstandingMessages, err = strconv.Atoi(maxOutstandingMessages)
	if err != nil {
		return err
	}

	maxOutstandingBytes, err := getEnvVariable("MAX_OUTSTANDING_BYTES")
	if err != nil {
		return err
	}
	subscriberConf.MaxOutstandingBytes, err = strconv.Atoi(maxOutstandingBytes)
	if err != nil {
		return err
	}

	numOfGoroutines, err := getEnvVariable("NUM_OF_GOROUTINS")
	if err != nil {
		return err
	}
	subscriberConf.NumOfGoroutines, err = strconv.Atoi(numOfGoroutines)
	if err != nil {
		return err
	}

	return err
}
