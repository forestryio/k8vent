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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	"k8s.io/kubernetes/pkg/api/v1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_3"
)

func TestVent(t *testing.T) {
	eventFile := "testdata/events.json"
	eventBytes, readErr := ioutil.ReadFile(eventFile)
	if readErr != nil {
		t.Errorf("failed to open event JSON file %s: %v", eventFile, readErr)
	}

	eventList := &v1.EventList{}
	if err := json.Unmarshal(eventBytes, eventList); err != nil {
		t.Errorf("failed to unmarshal event JSON into EventList: %v", err)
	}

	eventStore := make(map[string]eventCache)
	m := &sync.Mutex{}
	http.HandleFunc("/k8event", func(w http.ResponseWriter, r *http.Request) {
		name, storeErr := storeEvent(m, eventStore, w, r)
		if storeErr != nil {
			t.Errorf("failed to store event %v: %v", r, storeErr)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		resp := []byte(fmt.Sprintf(`{"event":"%s","status":"ok"}`, name))
		if _, err := w.Write(resp); err != nil {
			t.Errorf("failed to write server response: %v", err)
			return
		}
	})
	go func() {
		if err := http.ListenAndServe(":30256", nil); err != nil {
			t.Errorf("event server process failed: %v", err)
		}
	}()
	urls := []string{"http://127.0.0.1:30256/k8event"}

	cache := make(map[string]eventCache)
	events := eventList.Items
	var k8Client *client.Clientset
	emitNewEvents(urls, events, cache, k8Client)
	if len(events) != len(cache) {
		t.Errorf("number of events returned in cache (%d) do not match the number of events (%d)", len(cache), len(events))
	}
	time.Sleep(1 * time.Second)
	if !reflect.DeepEqual(cache, eventStore) {
		t.Errorf("cache (%v) and eventStore (%v) differ", cache, eventStore)
	}

	emitNewEvents(urls, events, cache, k8Client)
	if len(events) != len(cache) {
		t.Errorf("number of events returned in cache (%d) do not match the number of events (%d)", len(cache), len(events))
	}
	time.Sleep(1 * time.Second)
	if !reflect.DeepEqual(cache, eventStore) {
		t.Errorf("cache (%v) and eventStore (%v) differ", cache, eventStore)
	}
}

func storeEvent(m *sync.Mutex, store map[string]eventCache, w http.ResponseWriter, r *http.Request) (n string, e error) {
	eventBytes, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		return "", readErr
	}
	event := v1.Event{}
	if err := json.Unmarshal(eventBytes, &event); err != nil {
		return "", err
	}

	name := event.ObjectMeta.Name
	m.Lock()
	store[name] = eventCache{
		count:     event.Count,
		timestamp: event.LastTimestamp.Time,
	}
	m.Unlock()

	return name, nil
}
