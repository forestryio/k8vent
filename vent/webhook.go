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
	"net/http"

	"github.com/cenk/backoff"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

// webhookPayload is the structure serialized and sent to the webhook
// endpoints.
type webhookPayload struct {
	Pod v1.Pod `json:"pod"`
}

// PostToWebhooks marshals payload into JSON and posts it to the webhook
// URLs provided.
func postToWebhooks(urls []string, payload *webhookPayload, secret string) {
	slug := podSlug(payload.Pod)
	log := logger.WithField("pod", slug)

	objJSON, jsonErr := json.Marshal(payload)
	if jsonErr != nil {
		log.Errorf("Failed to marshal event to JSON: %v: %+v", jsonErr, payload)
		return
	}
	log.Tracef("Sending payload: %s", string(objJSON))

	for _, url := range urls {
		go func(u string) {
			log.Infof("Posting to '%s'", u)
			if err := postToWebhook(slug, u, objJSON, secret); err != nil {
				log.Errorf("Failed to post to '%s': %s", u, err.Error())
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
			signature, signErr := generateSignature(payload, secret)
			if signErr != nil {
				return signErr
			}
			log.Debugf("Signing payload with secret: %s", signature)
			req.Header.Add("x-atomist-signature", signature)
		}
		resp, postErr := client.Do(req)
		if postErr != nil {
			return fmt.Errorf("failed to POST event to %s: %v", url, postErr)
		}
		defer resp.Body.Close()
		corrID, corrErr := extractPropertyString(resp, "correlation-id")
		if corrErr != nil {
			log.Warnf("Failed to extract correlation ID from %s response: %v", url, corrErr)
		}
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return fmt.Errorf("non-200 response from webhook %s: code:%d,correlation-id:%s", url, resp.StatusCode, corrID)
		}
		log.WithFields(logrus.Fields{
			"code":           resp.StatusCode,
			"correlation-id": corrID,
		}).Infof("Posted to '%s'", url)
		return nil
	}

	return backoff.Retry(post, backoff.NewExponentialBackOff())
}
