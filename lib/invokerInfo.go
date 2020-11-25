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
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// InvokerInfo represents a invoker configuration.
// It holds information needed for invocation of pull functions.
type InvokerInfo struct {
	NumberOfMessages  string //number of messages pull function will pull and store (only for synchronous pull)
	NumberOfSeconds   int    //time duration of a context for which the messages will be pulled
	NumberOfInstances int    //number of pull function instances that will run in parallel
	InstanceNumber    int    //help parameter used for logging error messages (indicates on which instance the error occurred)
	FunctionURL       string //URL of a Cloud function which will be triggered by invoker
}

// SetInvokerInfo sets the parameters of a invoker configuration by extracting values ​​from the corresponding environment variables.
// An error is returned if any errors occur during the function execution.
func SetInvokerInfo(invokerInfo *InvokerInfo) error {
	var err error

	invokerInfo.NumberOfMessages = os.Getenv("NUM_OF_MESSAGES")

	numOfSeconds, err := getEnvVariable("NUM_OF_SECONDS")
	if err != nil {
		return err
	}

	invokerInfo.NumberOfSeconds, err = strconv.Atoi(numOfSeconds)
	if err != nil {
		return err
	}

	numOfInstances, err := getEnvVariable("NUM_OF_INSTANCES")
	if err != nil {
		return err
	}
	invokerInfo.NumberOfInstances, err = strconv.Atoi(numOfInstances)
	if err != nil {
		return err
	}

	invokerInfo.FunctionURL, err = getEnvVariable("FUNC_URL")
	if err != nil {
		return err
	}

	err = checkURL(invokerInfo.FunctionURL)
	if err != nil {
		return err
	}

	return nil

}

// checkURL represents helper function which checks if given URL is in the form of a Cloud function URL.
// The function returns error message if URL was not a match.
func checkURL(url string) error {
	var err error
	regex := regexp.MustCompile(`^https:\/{2}[a-zA-Z0-9-_]+-[a-zA-Z0-9-_]+\.cloudfunctions\.net\/[a-zA-Z0-9-_]+$`)
	matched := regex.MatchString(url)

	if !matched {
		err = fmt.Errorf("Invalid Cloud function url")
	}
	return err
}
