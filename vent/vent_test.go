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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"k8s.io/api/core/v1"
)

func TestMain(m *testing.M) {
	logrus.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestExtractPod(t *testing.T) {

	objFile := "testdata/extract.json"
	objBytes, readErr := ioutil.ReadFile(objFile)
	if readErr != nil {
		t.Errorf("failed to open event JSON file %s: %v", objFile, readErr)
	}

	objects := []v1.Pod{}
	if err := json.Unmarshal(objBytes, &objects); err != nil {
		t.Errorf("failed to unmarshal objects JSON into []interface{}: %v", err)
	}

	for i := 0; i < len(objects); i++ {
		podName := fmt.Sprintf("sleep-85576868c9-jvtzb-%d", i)
		containerName := fmt.Sprintf("sleep%d", i)
		containerImage := fmt.Sprintf("atomist/sleep:0.2.%d", i)
		hostIP := fmt.Sprintf("192.168.99.10%d", i)
		podIP := fmt.Sprintf("172.17.0.%d", i)
		pod, annot, err := extractPod(objects[i])
		if err != nil {
			t.Errorf("failed to extract object %d: %v", i, err)
			continue
		}
		if annot != nil {
			t.Errorf("erroneously extracted webhooks from object %d: %v", i, annot)
		}
		if pod.Name != podName {
			t.Errorf("pod name (%s) for object %d does not match expected (%s)", pod.Name, i, podName)
		}
		if pod.Namespace != "testing" {
			t.Errorf("pod namespace (%s) for object %d does not match expected (testing)", pod.Namespace, i)
		}
		if len(pod.Spec.Containers) != 1 {
			t.Errorf("number of containers (%d) for object %d does not match expected (1)",
				len(pod.Spec.Containers), i)
		}
		if pod.Spec.Containers[0].Name != containerName {
			t.Errorf("container name (%s) for object %d does not match expected (%s)",
				pod.Spec.Containers[0].Name, i, containerName)
		}
		if pod.Spec.Containers[0].Image != containerImage {
			t.Errorf("container image (%s) for object %d does not match expected (%s)",
				pod.Spec.Containers[0].Image, i, containerImage)
		}
		if pod.Status.HostIP != hostIP {
			t.Errorf("pod host IP (%s) for object %d does not match expected (%s)", pod.Status.HostIP, i, hostIP)
		}
		if pod.Status.PodIP != podIP {
			t.Errorf("pod IP (%s) for object %d does not match expected (%s)", pod.Status.PodIP, i, podIP)
		}
		if pod.Status.Phase != "Running" {
			t.Errorf("pod phase (%s) for object %d does not match expected (Running)", pod.Status.Phase, i)
		}
	}
}

func TestExtractPodAnnotation(t *testing.T) {

	objFile := "testdata/extract-annot.json"
	objBytes, readErr := ioutil.ReadFile(objFile)
	if readErr != nil {
		t.Errorf("failed to open event JSON file %s: %v", objFile, readErr)
	}

	objects := []v1.Pod{}
	if err := json.Unmarshal(objBytes, &objects); err != nil {
		t.Errorf("failed to unmarshal objects JSON into []interface{}: %v", err)
	}

	for i := 0; i < len(objects); i++ {
		podName := fmt.Sprintf("sleep-85576868c9-jvtzb-%d", i)
		containerName := fmt.Sprintf("sleep%d", i)
		containerImage := fmt.Sprintf("atomist/sleep:0.2.%d", i)
		hostIP := fmt.Sprintf("192.168.99.10%d", i)
		podIP := fmt.Sprintf("172.17.0.%d", i)
		pod, annot, err := extractPod(objects[i])
		if err != nil {
			t.Errorf("failed to extract object %d: %v", i, err)
			continue
		}
		if annot == nil {
			t.Errorf("failed to extract k8vent annotation from object %d: %v", i, annot)
		} else {
			if i > 0 {
				env := fmt.Sprintf("env-%d", i)
				if annot.Environment != env {
					t.Errorf("environment annotation for object %d '%s' does not match '%s'",
						i, annot.Environment, env)
				}
			}
			if len(annot.Webhooks) != i+1 {
				t.Errorf("number of webhooks (%d) from object %d does not match expected (%d)",
					len(annot.Webhooks), i, i)
			}
			for j := 0; j < len(annot.Webhooks); j++ {
				wh := fmt.Sprintf("https://webhook.atomist.com/atomist/kube/teams/TEAM_ID%d", j)
				if annot.Webhooks[j] != wh {
					t.Errorf("webhook %d (%s) for object %d does not match expected (%s)",
						j, annot.Webhooks[j], i, wh)
				}
			}
		}
		if pod.Name != podName {
			t.Errorf("pod name (%s) for object %d does not match expected (%s)", pod.Name, i, podName)
		}
		if pod.Namespace != "testing" {
			t.Errorf("pod namespace (%s) for object %d does not match expected (testing)", pod.Namespace, i)
		}
		if len(pod.Spec.Containers) != 1 {
			t.Errorf("number of containers (%d) for object %d does not match expected (1)",
				len(pod.Spec.Containers), i)
		}
		if pod.Spec.Containers[0].Name != containerName {
			t.Errorf("container name (%s) for object %d does not match expected (%s)",
				pod.Spec.Containers[0].Name, i, containerName)
		}
		if pod.Spec.Containers[0].Image != containerImage {
			t.Errorf("container image (%s) for object %d does not match expected (%s)",
				pod.Spec.Containers[0].Image, i, containerImage)
		}
		if pod.Status.HostIP != hostIP {
			t.Errorf("pod host IP (%s) for object %d does not match expected (%s)", pod.Status.HostIP, i, hostIP)
		}
		if pod.Status.PodIP != podIP {
			t.Errorf("pod IP (%s) for object %d does not match expected (%s)", pod.Status.PodIP, i, podIP)
		}
		if pod.Status.Phase != "Running" {
			t.Errorf("pod phase (%s) for object %d does not match expected (Running)", pod.Status.Phase, i)
		}
	}
}
