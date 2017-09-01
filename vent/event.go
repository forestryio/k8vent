// Copyright © 2017 Atomist
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
	"log"
	"net/http"
	"strings"
	"time"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/v1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_3"
    "k8s.io/kubernetes/pkg/client/restclient"

	"github.com/cenk/backoff"
)

const (
	interval = 5 * time.Second
)

type eventCache struct {
	count     int32
	timestamp time.Time
}

// Vent polls the kubernetes service for events and posts those events
// as JSON to the webhook at `url`.  It should never return.
func Vent(urls []string) (e error) {
	k8Config, configErr := restclient.InClusterConfig()
	if configErr != nil {
		return configErr
	}

	k8Client, clientErr := client.NewForConfig(k8Config)
	if clientErr != nil {
		return clientErr
	}

	eventCache := make(map[string]eventCache)
	errs := []string{}

	log.Println("starting polling loop")
	for {
		time.Sleep(interval)

		opts := api.ListOptions{
		//Watch: true,
		}
		list, listErr := k8Client.Core().Events(api.NamespaceAll).List(opts)
		if listErr != nil {
			log.Printf("failed to list events: %v\n", listErr)
			errs = append(errs, listErr.Error())
			if len(errs) > 12 {
				return fmt.Errorf("too many consecutive failures to reach k8 service: %s", strings.Join(errs, "; "))
			}
			continue
		}
		errs = []string{}

		emitNewEvents(urls, list.Items, eventCache, k8Client)
	}
}

func emitNewEvents(urls []string, events []v1.Event, cache map[string]eventCache, k8Client *client.Clientset) {
	for _, event := range events {
		name := event.ObjectMeta.Name
		if ci, ok := cache[name]; ok {
			count := ci.count
			ts := ci.timestamp
			if event.Count > count && event.LastTimestamp.Time.After(ts.Add(1*time.Minute)) {
				emit(urls, event, cache, k8Client)
			}
		} else {
			emit(urls, event, cache, k8Client)
		}
	}
}

func emit(urls []string, event v1.Event, cache map[string]eventCache, k8Client *client.Clientset) {
	eventSpec := addPodSpec(event, k8Client)
	for _, url := range urls {
		go func(u string) {
			if err := postEvent(u, eventSpec); err != nil {
				log.Println(err.Error())
			}
		}(url)
	}
	cache[event.ObjectMeta.Name] = eventCache{
		count:     event.Count,
		timestamp: event.LastTimestamp.Time,
	}
}

func postEvent(url string, event v1.Event) (e error) {

	eventJSON, jsonErr := json.Marshal(event)
	if jsonErr != nil {
		return fmt.Errorf("failed to marshal event to JSON: %v: %v", jsonErr, event)
	}

	post := func() error {
		resp, postErr := http.Post(url, "application/json", bytes.NewBuffer(eventJSON))
		if postErr != nil {
			return fmt.Errorf("failed to POST event to %s: %v", url, postErr)
		}
		_ = resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return fmt.Errorf("non-200 response from webhook: %d", resp.StatusCode)
		}
		return nil
	}

	return backoff.Retry(post, backoff.NewExponentialBackOff())
}
