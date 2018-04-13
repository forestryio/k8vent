// Copyright © 2018 Atomist
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

	"github.com/Sirupsen/logrus"
	"github.com/cenk/backoff"
)

// PostToWebhooks marshals podEnv into JSON and posts it to the webhook
// URLs provided.
func PostToWebhooks(urls []string, podEnv *K8PodEnv) {
	objJSON, jsonErr := json.Marshal(podEnv)
	if jsonErr != nil {
		logrus.Errorf("failed to marshal event to JSON: %v: %+v", jsonErr, podEnv)
		return
	}

	podSlug := podEnv.Pod.Namespace + "/" + podEnv.Pod.Name
	for _, url := range urls {
		go func(u string) {
			logrus.Infof("posting pod '%s' to '%s'", podSlug, u)
			if err := postToWebhook(u, objJSON); err != nil {
				logrus.Errorf("failed to post pod '%s' to '%s': %s", podSlug, u, err.Error())
			}
		}(url)
	}
}

func postToWebhook(url string, payload []byte) (e error) {

	post := func() error {
		resp, postErr := http.Post(url, "application/json", bytes.NewBuffer(payload))
		if postErr != nil {
			return fmt.Errorf("failed to POST event to %s: %v", url, postErr)
		}
		_ = resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return fmt.Errorf("non-200 response from webhook %s: %d", url, resp.StatusCode)
		}
		return nil
	}

	return backoff.Retry(post, backoff.NewExponentialBackOff())
}
