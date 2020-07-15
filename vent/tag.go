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

// getDockerAuthToken gets a read-only Docker Hub registry token for
// the provided repository.
func getDockerAuthToken(repository string) (t string, e error) {
	client := &http.Client{}
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
	return token, nil
}

// getDockerTagDigest retrieves the digest for atomist/k8svent:tag.
func getDockerTagDigest(tag string) (d string, e error) {
	repository := "atomist/k8svent"
	token, tokenErr := getDockerAuthToken(repository)
	if tokenErr != nil {
		return d, tokenErr
	}
	client := &http.Client{}
	digestURL := "https://index.docker.io/v2/" + repository + "/manifests/" + tag
	digestReq, digestReqErr := http.NewRequest("HEAD", digestURL, nil)
	if digestReqErr != nil {
		return d, fmt.Errorf("failed to create HEAD request to %s: %v", digestURL, digestReqErr)
	}
	digestReq.Header.Add("authorization", "Bearer "+token)
	digestResp, digestRespErr := client.Do(digestReq)
	if digestRespErr != nil {
		return d, fmt.Errorf("failed to HEAD %s: %v", digestURL, digestRespErr)
	}
	defer digestResp.Body.Close()
	digest, digestOK := digestResp.Header[http.CanonicalHeaderKey("docker-content-digest")]
	if !digestOK {
		return d, fmt.Errorf("manifest HEAD response did not contain digest header")
	}
	if len(digest) < 1 {
		return d, fmt.Errorf("manifest HEAD response contained empty digest header: %v", digest)
	}
	return digest[0], nil
}
