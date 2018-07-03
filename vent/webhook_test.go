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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestPostToWebhooks(t *testing.T) {
	log.Info("starting TestPostToWebhooks")
	objects, loadErr := loadObjects("testdata/vent.json")
	if loadErr != nil {
		t.Error(loadErr.Error())
	}

	// should accept empty list of webhook URLs
	PostToWebhooks([]string{}, &objects[0])

	store := map[string]interface{}{}
	m := &sync.Mutex{}
	stopCh := make(chan bool, len(objects))
	defer close(stopCh)
	tail := "/k8vent"
	http.HandleFunc(tail, func(w http.ResponseWriter, r *http.Request) {
		if err := storeObject(m, store, w, r, stopCh); err != nil {
			t.Errorf("failed to store event %v: %v", r, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		resp := []byte(fmt.Sprintf(`{"correlation-id":"d95f0bc3-76c7-49a9-8eb3-6c427a44478d","message":"successfully posted event"}`))
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
		PostToWebhooks(urls, &o)
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

func TestPostToWebhook(t *testing.T) {
	payload := []byte(`{
  "env": {
    "ATOMIST_ENVIRONMENT": "dips",
    "HOME": "/root",
    "HOSTNAME": "k8vent-65bc5b5c56-9kfkk",
    "KUBERNETES_PORT": "tcp://10.96.0.1:443",
    "KUBERNETES_PORT_443_TCP": "tcp://10.96.0.1:443",
    "KUBERNETES_PORT_443_TCP_ADDR": "10.96.0.1",
    "KUBERNETES_PORT_443_TCP_PORT": "443",
    "KUBERNETES_PORT_443_TCP_PROTO": "tcp",
    "KUBERNETES_SERVICE_HOST": "10.96.0.1",
    "KUBERNETES_SERVICE_PORT": "443",
    "KUBERNETES_SERVICE_PORT_HTTPS": "443",
    "PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
  },
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
      "annotations": {
        "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"default\",\"name\":\"sleep-85576868c9\",\"uid\":\"3647fb49-f0c9-11e7-8b0c-080027815bd2\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"164270\"}}\n"
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

	tail := "/k8vent"
	mux := http.NewServeMux()
	mux.HandleFunc(tail, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := []byte(fmt.Sprintf(`{"status":"ok"}`))
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

	if err := postToWebhook("some/pod", url, payload); err != nil {
		t.Errorf("failed to handle invalid server response: %v", err)
	}
}

func TestExtractCorrelationID(t *testing.T) {
	validResponses := []string{
		`{"correlation-id":"0"}`,
		`{"correlation-id":"1","message":"successfully posted event"}`,
		`{"status":"ok","correlation-id":"2","message":"successfully posted event"}`,
		`{"error":null,"status":"ok","correlation-id":"3","message":"successfully posted event"}`,
	}
	for i, r := range validResponses {
		resp := &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader([]byte(r))),
		}
		if cid, err := extractCorrelationID(resp); err == nil {
			if cid != fmt.Sprintf("%d", i) {
				t.Errorf("extracted correlation ID (%s) is not expected value (%d)", cid, i)
			}
		} else {
			t.Errorf("failed to extract correlation ID from '%s': %v", r, err)
		}
	}

	invalidResponses := []string{
		"",
		"{}",
		`{"status":"ok"}`,
		`{"message":"successfully posted event"}`,
	}
	for _, r := range invalidResponses {
		resp := &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader([]byte(r))),
		}
		if cid, err := extractCorrelationID(resp); err == nil {
			t.Errorf("unexpectedly extracted correlation ID from invalid response '%s': %s", r, cid)
		}
	}
}

func loadObjects(objFile string) (o []K8PodEnv, e error) {
	objBytes, readErr := ioutil.ReadFile(objFile)
	if readErr != nil {
		return o, fmt.Errorf("failed to open event JSON file %s: %v", objFile, readErr)
	}

	objects := []K8PodEnv{}
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
		return fmt.Sprintf("non/pod:%d", rand.Int63())
	}
	k8pe := K8PodEnv{}
	if err := json.Unmarshal(objJSON, &k8pe); err != nil {
		fmt.Printf("failed to unmarshal object to K8PodEnv:%v\n", err)
		return fmt.Sprintf("non/K8PodEnv:%d", rand.Int63())
	}
	return k8pe.Pod.ObjectMeta.Namespace + "/" + k8pe.Pod.ObjectMeta.Name + ":" + k8pe.Pod.ObjectMeta.ResourceVersion
}
