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
	"sync"
	"testing"

	"github.com/satori/go.uuid"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPostToWebhooks(t *testing.T) {
	objFile := "testdata/pods.json"
	objBytes, readErr := ioutil.ReadFile(objFile)
	if readErr != nil {
		t.Errorf("failed to open event JSON file %s: %v", objFile, readErr)
	}

	objects := []v1.Pod{}
	if err := json.Unmarshal(objBytes, &objects); err != nil {
		t.Errorf("failed to unmarshal objects JSON into []interface{}: %v", err)
	}

	store := map[string]interface{}{}
	m := &sync.Mutex{}
	stopCh := make(chan bool, len(objects))
	defer close(stopCh)
	http.HandleFunc("/k8event", func(w http.ResponseWriter, r *http.Request) {
		if err := storeObject(m, store, w, r, stopCh); err != nil {
			t.Errorf("failed to store event %v: %v", r, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		resp := []byte(fmt.Sprintf(`{"status":"ok"}`))
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

	for _, o := range objects {
		PostToWebhooks(urls, o)
	}
	for i := 0; i < len(objects); i++ {
		<-stopCh
	}

	if len(store) != len(objects) {
		t.Errorf("number of objects processed (%d) does not equal the number sent (%d)", len(store), len(objects))
	} else {
		for _, o := range objects {
			k := extracObjectKey(o)
			if _, ok := store[k]; !ok {
				t.Errorf("object %s did not get stored", k)
			}
		}
	}
}

func storeObject(m *sync.Mutex, store map[string]interface{}, w http.ResponseWriter, r *http.Request, stopCh chan bool) (e error) {
	defer func() { stopCh <- true }()

	objBytes, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		return readErr
	}
	obj := map[string]interface{}{}
	if err := json.Unmarshal(objBytes, &obj); err != nil {
		return err
	}

	key := extracObjectKey(obj)
	m.Lock()
	store[key] = obj
	m.Unlock()

	return nil
}

func extracObjectKey(obj interface{}) string {
	objJSON, jsonErr := json.Marshal(obj)
	if jsonErr != nil {
		fmt.Printf("failed to marshal object to JSON:%v\n", jsonErr)
		return "non/pod:" + uuid.NewV4().String()
	}
	type K8Object struct {
		metav1.ObjectMeta `json:"metadata"`
	}
	pod := K8Object{}
	if err := json.Unmarshal(objJSON, &pod); err != nil {
		fmt.Printf("failed to unmarshal object to Pod:%v\n", err)
		return "non/pod:" + uuid.NewV4().String()
	}
	return pod.ObjectMeta.Namespace + "/" + pod.ObjectMeta.Name + ":" + pod.ObjectMeta.ResourceVersion
}
