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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/cenk/backoff"
)

// PostToWebhooks marshals podEnv into JSON and posts it to the webhook
// URLs provided.
func PostToWebhooks(urls []string, podEnv *K8PodEnv, secret string) {
	podSlug := podEnv.Pod.Namespace + "/" + podEnv.Pod.Name
	log := logger.WithField("pod", podSlug)

	objJSON, jsonErr := json.Marshal(podEnv)
	if jsonErr != nil {
		log.Errorf("failed to marshal event to JSON: %v: %+v", jsonErr, podEnv)
		return
	}

	for _, url := range urls {
		go func(u string) {
			log.Infof("posting pod '%s' to '%s'", podSlug, u)
			if err := postToWebhook(podSlug, u, objJSON, secret); err != nil {
				log.Errorf("failed to post pod '%s' to '%s': %s", podSlug, u, err.Error())
			}
		}(url)
	}
}

// postToWebhook post the provided payload to the URL.
func postToWebhook(pod string, url string, payload []byte, secret string) (e error) {
	log := logger.WithField("pod", pod)

	post := func() error {
		client := &http.Client{}
		req, reqErr := http.NewRequest("POST", url, bytes.NewBuffer(payload))
		if reqErr != nil {
			return fmt.Errorf("failed to create POST request to %s: %v", url, reqErr)
		}
		req.Header.Add("content-type", "application/json")
		if secret != "" {
			signature := generateSignature(payload, secret)
			log.Debugf("signing payload with secret: %s", signature)
			req.Header.Add("x-atomist-signature", signature)
		}
		resp, postErr := client.Do(req)
		if postErr != nil {
			return fmt.Errorf("failed to POST event to %s: %v", url, postErr)
		}
		defer resp.Body.Close()
		corrID, corrErr := extractCorrelationID(resp)
		if corrErr != nil {
			log.Warnf("failed to extract correlation ID from %s response: %v", url, corrErr)
		}
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return fmt.Errorf("non-200 response from webhook %s: code:%d,correlation-id:%s", url, resp.StatusCode, corrID)
		}
		log.WithFields(logrus.Fields{
			"code":           resp.StatusCode,
			"correlation-id": corrID,
		}).Infof("posted pod '%s' to '%s'", pod, url)
		return nil
	}

	return backoff.Retry(post, backoff.NewExponentialBackOff())
}

// extractCorrelationID reads the provided response body, parses it as
// JSON, and returns the "correlation-id" element.
func extractCorrelationID(resp *http.Response) (cid string, e error) {
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return "", fmt.Errorf("failed to read response: %v", readErr)
	}
	var respObj map[string]string
	if err := json.Unmarshal(body, &respObj); err != nil {
		return "", fmt.Errorf("failed to parse '%s' as JSON: %v", string(body), err)
	}
	corrID, corrExists := respObj["correlation-id"]
	if !corrExists {
		return "", fmt.Errorf("response '%s' has no 'correlation-id'", string(body))
	}
	return corrID, nil
}
