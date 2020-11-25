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

package invoker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/syntio/aquarium-persistor-gcp/lib"
)

const (
	contentType = "application/json"
)

// InvokerHandler represents the main invoke function which is triggered by HTTP request.
// Function provides information needed for invocation by calling SetInvokerInfo method and
// starts N parallel instances of InvokeFunction (go routines), where N is given by the NumberOfInstances parameter.
// Function logs invocation error messages received through channel.
func InvokerHandler(w http.ResponseWriter, r *http.Request) {

	invokerInfo := lib.InvokerInfo{}
	err := lib.SetInvokerInfo(&invokerInfo)
	if err != nil {
		log.Printf("Error during retrieving environment variable.\n")
		panic(err)
	}

	ch := make(chan string)

	for i := 0; i < invokerInfo.NumberOfInstances; i++ {
		invokerInfo.InstanceNumber = i + 1
		go InvokeFunction(invokerInfo, ch)
	}

	for i := 0; i < invokerInfo.NumberOfInstances; i++ {
		log.Printf("Call to function finished. %v.\n", <-ch)
	}
	fmt.Fprint(w, "Finished execution")

}

// InvokeFunction sends a HTTP post request and checks the call validity.
// The invokerInfo argument contains the target URL and informations which will be sent through the request body.
// Returned result is an error which defines the validity of the function action.
func InvokeFunction(invokerInfo lib.InvokerInfo, ch chan<- string) error {
	var err error

	//NumberOfInstances and NumberOfSeconds will be passed to the invoked function.
	jsonRequest, err := json.Marshal(&invokerInfo)
	if err != nil {
		log.Printf("Error during #%d request parameter marshaling. %v.\n", invokerInfo.InstanceNumber, err)
		return err
	}

	response, err := http.Post(invokerInfo.FunctionURL, contentType, bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Printf("Error during #%d function invocation: %v.\n", invokerInfo.InstanceNumber, err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		ch <- fmt.Sprintf("Function number #%d invoked successfully", invokerInfo.InstanceNumber)
	} else {
		responseBody, _ := ioutil.ReadAll(response.Body)
		ch <- fmt.Sprintf("Error on #%d function invocation. Status Code: %d: %s", invokerInfo.InstanceNumber, response.StatusCode, string(responseBody))
	}

	return err
}
