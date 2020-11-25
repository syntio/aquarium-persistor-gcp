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
)

// getEnvVariable represents helper function which receives an environment variable
// name and attempts to extract its value.
// The function returns an extracted value and an error message if an error occurs.
func getEnvVariable(name string) (string, error) {
	value, ok := os.LookupEnv(name)

	if !ok {
		return value, fmt.Errorf("Environment variable '%s' is not set ", name)
	}
	if value == "" {
		return value, fmt.Errorf("Value for '%s' was not found within the environment variables", name)
	}

	return value, nil
}
