// Copyright © 2020 Atomist
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

// extractPropertyString reads the provided response body, parses it
// as JSON, and returns the value of `key` property.  If the property
// does not exist, an error is returned.
func extractPropertyString(resp *http.Response, key string) (v string, e error) {
	value, extractErr := extractProperty(resp, key)
	if extractErr != nil {
		return v, extractErr
	}
	stringValue, ok := value.(string)
	if !ok {
		return v, fmt.Errorf("response property '%s' is not a string: %v (%s)", key, value, reflect.TypeOf(value))
	}
	return stringValue, nil
}

// extractPropertyStringSlice reads the provided response body, parses
// it as JSON, and returns the value of `key` property as a []string.
// If the property does not exist, an error is returned.
func extractPropertyStringSlice(resp *http.Response, key string) (v []string, e error) {
	value, extractErr := extractProperty(resp, key)
	if extractErr != nil {
		return v, extractErr
	}
	interfaceSliceValue, ok := value.([]interface{})
	if !ok {
		return v, fmt.Errorf("response property '%s' is not a slice: %v (%s)", key, value, reflect.TypeOf(value))
	}
	stringSliceValue := make([]string, len(interfaceSliceValue))
	for i, x := range interfaceSliceValue {
		s, sOk := x.(string)
		if !sOk {
			return v, fmt.Errorf("response property '%s' element is not a string: %v (%s)", key, x, reflect.TypeOf(x))
		}
		stringSliceValue[i] = s
	}
	return stringSliceValue, nil
}

func extractProperty(resp *http.Response, key string) (v interface{}, e error) {
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return v, fmt.Errorf("failed to read response: %v", readErr)
	}
	var respObj map[string]interface{}
	if err := json.Unmarshal(body, &respObj); err != nil {
		return v, fmt.Errorf("failed to parse '%s' as JSON: %v", string(body), err)
	}
	value, exists := respObj[key]
	if !exists {
		return v, fmt.Errorf("response '%s' has no property '%s'", string(body), key)
	}
	return value, nil
}
