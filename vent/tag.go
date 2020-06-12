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
	"fmt"
	"net/http"
)

// getDockerTags queries the Docker Hub API for tags.
func getDockerTags() (t []string, e error) {
	client := &http.Client{}
	repository := "atomist/k8svent"
	clientID := packageSlug()
	authURL := "https://auth.docker.io/token?service=registry.docker.io&scope=repository:" + repository +
		":pull&client_id=" + clientID
	authReq, authReqErr := http.NewRequest("GET", authURL, nil)
	if authReqErr != nil {
		return t, fmt.Errorf("failed to create GET request to %s: %v", authURL, authReqErr)
	}
	authReq.Header.Add("content-type", "application/json")
	authResp, authRespErr := client.Do(authReq)
	if authRespErr != nil {
		return t, fmt.Errorf("failed to GET %s: %v", authURL, authRespErr)
	}
	defer authResp.Body.Close()
	token, tokenErr := extractPropertyString(authResp, "token")
	if tokenErr != nil {
		return t, fmt.Errorf("failed to extract token from %s response: %v", authURL, tokenErr)
	}

	tagURL := "https://index.docker.io/v2/" + repository + "/tags/list"
	tagReq, tagReqErr := http.NewRequest("GET", tagURL, nil)
	if tagReqErr != nil {
		return t, fmt.Errorf("failed to create GET request to %s: %v", tagURL, tagReqErr)
	}
	tagReq.Header.Add("content-type", "application/json")
	tagReq.Header.Add("authorization", "Bearer "+token)
	tagResp, tagRespErr := client.Do(tagReq)
	if tagRespErr != nil {
		return t, fmt.Errorf("failed to GET %s: %v", tagURL, tagRespErr)
	}
	defer tagResp.Body.Close()
	tags, tagsErr := extractPropertyStringSlice(tagResp, "tags")
	if tagsErr != nil {
		return t, fmt.Errorf("failed to extract token from %s response: %v", tagURL, tagsErr)
	}

	return tags, nil
}
