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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestPostToWebhooks(t *testing.T) {
	nullLogger, _ := test.NewNullLogger()
	logger = nullLogger.WithField("test", "webhook")

	objects, loadErr := loadObjects("testdata/vent.json")
	if loadErr != nil {
		t.Error(loadErr.Error())
	}

	// should accept empty list of webhook URLs
	postToWebhooks([]string{}, &objects[0], "")

	store := map[string]interface{}{}
	m := &sync.Mutex{}
	stopCh := make(chan bool, len(objects))
	defer close(stopCh)
	tail := "/k8svent"
	http.HandleFunc(tail, func(w http.ResponseWriter, r *http.Request) {
		if err := storeObject(m, store, w, r, stopCh); err != nil {
			t.Errorf("failed to store event %v: %v", r, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		resp := []byte(`{"correlation_id":"d95f0bc3-76c7-49a9-8eb3-6c427a44478d","message":"successfully posted event"}`)
		if _, err := w.Write(resp); err != nil {
			t.Errorf("failed to write server response: %v", err)
			return
		}
	})
	addr := "127.0.0.1:32866"
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil && err != http.ErrServerClosed {
			t.Errorf("event server process failed: %v", err)
		}
	}()
	urls := []string{fmt.Sprintf("http://%s%s", addr, tail)}

	for _, o := range objects {
		postToWebhooks(urls, &o, "")
	}
	for i := 0; i < len(objects); i++ {
		<-stopCh
	}

	if len(store) != len(objects) {
		t.Errorf("number of objects processed (%d) does not equal the number sent (%d)", len(store), len(objects))
	} else {
		for _, o := range objects {
			k := extractObjectKey(o)
			if _, ok := store[k]; !ok {
				t.Errorf("object %s did not get stored", k)
			}
		}
	}
}

func TestPostToWebhook(t *testing.T) {
	nullLogger, hook := test.NewNullLogger()
	nullLogger.SetLevel(logrus.InfoLevel)
	logger = nullLogger.WithField("test", "webhook")

	payload := []byte(`{
  "pod": {
    "metadata": {
      "name": "sleep-85576868c9-jvtzb",
      "generateName": "sleep-85576868c9-",
      "namespace": "default",
      "selfLink": "/api/v1/namespaces/default/pods/sleep-85576868c9-jvtzb",
      "uid": "36498f39-f0c9-11e7-8b0c-080027815bd2",
      "resourceVersion": "164858",
      "creationTimestamp": "2018-01-03T21:01:04Z",
      "deletionTimestamp": "2018-01-03T21:07:01Z",
      "deletionGracePeriodSeconds": 0,
      "labels": {
        "app": "sleep",
        "pod-template-hash": "4113242475"
      },
      "ownerReferences": [
        {
          "apiVersion": "extensions/v1beta1",
          "kind": "ReplicaSet",
          "name": "sleep-85576868c9",
          "uid": "3647fb49-f0c9-11e7-8b0c-080027815bd2",
          "controller": true,
          "blockOwnerDeletion": true
        }
      ]
    },
    "spec": {
      "volumes": [
        {
          "name": "default-token-s97ln",
          "secret": {
            "secretName": "default-token-s97ln",
            "defaultMode": 420
          }
        }
      ],
      "containers": [
        {
          "name": "sleep",
          "image": "atomist/sleep:0.2.0",
          "resources": {},
          "volumeMounts": [
            {
              "name": "default-token-s97ln",
              "readOnly": true,
              "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
            }
          ],
          "terminationMessagePath": "/dev/termination-log",
          "terminationMessagePolicy": "File",
          "imagePullPolicy": "IfNotPresent"
        }
      ],
      "restartPolicy": "Always",
      "terminationGracePeriodSeconds": 1,
      "dnsPolicy": "ClusterFirst",
      "serviceAccountName": "default",
      "serviceAccount": "default",
      "nodeName": "minikube",
      "securityContext": {},
      "schedulerName": "default-scheduler"
    },
    "status": {
      "phase": "Running",
      "conditions": [
        {
          "type": "Initialized",
          "status": "True",
          "lastProbeTime": null,
          "lastTransitionTime": "2018-01-03T21:01:04Z"
        },
        {
          "type": "Ready",
          "status": "False",
          "lastProbeTime": null,
          "lastTransitionTime": "2018-01-03T21:07:03Z",
          "reason": "ContainersNotReady",
          "message": "containers with unready status: [sleep]"
        },
        {
          "type": "PodScheduled",
          "status": "True",
          "lastProbeTime": null,
          "lastTransitionTime": "2018-01-03T21:01:04Z"
        }
      ],
      "hostIP": "192.168.99.100",
      "startTime": "2018-01-03T21:01:04Z",
      "containerStatuses": [
        {
          "name": "sleep",
          "state": {
            "terminated": {
              "exitCode": 137,
              "reason": "Error",
              "startedAt": "2018-01-03T21:06:06Z",
              "finishedAt": "2018-01-03T21:07:03Z",
              "containerID": "docker://a0b7328e6231905f70a8db9a54a6db02510c1142b71c2c7ae0bcdf58cad1923b"
            }
          },
          "lastState": {},
          "ready": false,
          "restartCount": 1,
          "image": "atomist/sleep:0.2.0",
          "imageID": "docker://sha256:25a42919981849670f942c9d842c29197dcf935a743bb0f3b45ae4f1dfab8074",
          "containerID": "docker://a0b7328e6231905f70a8db9a54a6db02510c1142b71c2c7ae0bcdf58cad1923b"
        }
      ],
      "qosClass": "BestEffort"
    }
  }
}`)

	tail := "/k8svent"
	first := true
	mux := http.NewServeMux()
	mux.HandleFunc(tail, func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, bodyErr := ioutil.ReadAll(r.Body)
		if bodyErr != nil {
			t.Errorf("failed to read request body: %v", bodyErr)
		}
		if string(bodyBytes) != string(payload) {
			t.Error("sent and received body are not identical")
		}
		contentType := r.Header.Get("content-type")
		if contentType != "application/json" {
			t.Errorf("request content-type header is not 'application/json': '%s'", contentType)
		}
		w.Header().Set("content-type", "application/json")
		resp := []byte(`{"status":"ok"}`)
		if first {
			resp = []byte(`{"correlation_id":"472c0bab-be3a-4e96-8cac-569ad9d612a5"}`)
			signature := r.Header.Get("x-atomist-signature")
			eSignature := "sha1=6e9adaa75d8deb8f893ddd9f557c7de8ab9e2dcf"
			if signature != eSignature {
				t.Errorf("request x-atomist-signature header is not '%s': '%s'", eSignature, signature)
			}
			first = false
		}
		if _, err := w.Write(resp); err != nil {
			t.Errorf("failed to write server response: %v", err)
			return
		}
	})
	addr := "127.0.0.1:32867"
	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Errorf("event server process failed: %v", err)
		}
	}()
	defer func() {
		if err := s.Shutdown(context.Background()); err != nil {
			t.Errorf("failed to shut down HTTP server: %v", err)
		}
	}()
	url := fmt.Sprintf("http://%s%s", addr, tail)
	hook.Reset()
	if err := postToWebhook("some/pod", url, payload, "Coast2Coast"); err != nil {
		t.Errorf("failed to handle server response: %v", err)
	}
	if len(hook.Entries) != 1 {
		logEntries := ""
		for i, entry := range hook.Entries {
			logEntries += " " + strconv.Itoa(i) + ":" + entry.Message + ";"
		}
		t.Errorf("expected 1 log entries, got %d: %s", len(hook.Entries), logEntries)
	}
	le := hook.LastEntry()
	if le.Level != logrus.InfoLevel {
		t.Errorf("last log level should be info: %v", le.Level)
	}
	if !strings.HasPrefix(le.Message, "Posted to ") {
		t.Errorf("expected final log to be post info: %s", le.Message)
	}
	corrID, corrIDOk := le.Data["correlation_id"]
	if !corrIDOk {
		t.Error("no correlation ID found")
	}
	eCorrID := "472c0bab-be3a-4e96-8cac-569ad9d612a5"
	if corrID != eCorrID {
		t.Errorf("correlation ID does not match: %s != %s", corrID, eCorrID)
	}
	hook.Reset()
	if err := postToWebhook("some/pod", url, payload, ""); err != nil {
		t.Errorf("failed to handle invalid server response: %v", err)
	}
	if len(hook.Entries) != 2 {
		logEntries := ""
		for i, entry := range hook.Entries {
			logEntries += " " + strconv.Itoa(i) + ":" + entry.Message + ";"
		}
		t.Errorf("expected 2 log entries, got %d:%s", len(hook.Entries), logEntries)
	}
	if hook.Entries[0].Level != logrus.WarnLevel {
		t.Errorf("first log level should be warn: %v", hook.Entries[0].Level)
	}
	if !strings.HasPrefix(hook.Entries[0].Message, "Failed to extract correlation ID from ") {
		t.Errorf(
			"expected first log entry to be correlation ID extraction failure warning: %s",
			hook.Entries[0].Message,
		)
	}
	le = hook.LastEntry()
	if le.Level != logrus.InfoLevel {
		t.Errorf("last log level should be info: %v", le.Level)
	}
	corrID, corrIDOk = le.Data["correlation_id"]
	if !corrIDOk {
		t.Error("no correlation ID found")
	}
	if corrID != "" {
		t.Errorf("correlation ID is not empty: %s", corrID)
	}
}

func loadObjects(objFile string) (o []webhookPayload, e error) {
	objBytes, readErr := ioutil.ReadFile(objFile)
	if readErr != nil {
		return o, fmt.Errorf("failed to open event JSON file %s: %v", objFile, readErr)
	}

	objects := []webhookPayload{}
	if err := json.Unmarshal(objBytes, &objects); err != nil {
		return o, fmt.Errorf("failed to unmarshal objects JSON into []interface{}: %v", err)
	}
	return objects, nil
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

	key := extractObjectKey(obj)
	m.Lock()
	store[key] = obj
	m.Unlock()

	return nil
}

func extractObjectKey(obj interface{}) string {
	objJSON, jsonErr := json.Marshal(obj)
	if jsonErr != nil {
		fmt.Printf("failed to marshal object to JSON:%v\n", jsonErr)
		return fmt.Sprintf("non/pod:%d", rand.Int63())
	}
	wp := webhookPayload{}
	if err := json.Unmarshal(objJSON, &wp); err != nil {
		fmt.Printf("failed to unmarshal object to webhookPayload:%v\n", err)
		return fmt.Sprintf("non/webhookPayload:%d", rand.Int63())
	}
	return wp.Pod.ObjectMeta.Namespace + "/" + wp.Pod.ObjectMeta.Name + ":" + wp.Pod.ObjectMeta.ResourceVersion
}
