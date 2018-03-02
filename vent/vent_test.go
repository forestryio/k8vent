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
	"testing"

	"github.com/Sirupsen/logrus"
	"k8s.io/api/core/v1"
)

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

	logger := logrus.WithField("pkg", "k8vent-test")
	logrus.SetOutput(ioutil.Discard)

	for i := 0; i < len(objects); i++ {
		podName := fmt.Sprintf("sleep-85576868c9-jvtzb-%d", i)
		containerName := fmt.Sprintf("sleep%d", i)
		containerImage := fmt.Sprintf("atomist/sleep:0.2.%d", i)
		hostIP := fmt.Sprintf("192.168.99.10%d", i)
		podIP := fmt.Sprintf("172.17.0.%d", i)
		pod, annot, err := extractPod(objects[i], logger)
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

	logger := logrus.WithField("pkg", "k8vent-test")
	logrus.SetOutput(ioutil.Discard)

	for i := 0; i < len(objects); i++ {
		podName := fmt.Sprintf("sleep-85576868c9-jvtzb-%d", i)
		containerName := fmt.Sprintf("sleep%d", i)
		containerImage := fmt.Sprintf("atomist/sleep:0.2.%d", i)
		hostIP := fmt.Sprintf("192.168.99.10%d", i)
		podIP := fmt.Sprintf("172.17.0.%d", i)
		pod, annot, err := extractPod(objects[i], logger)
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

func TestPakePod(t *testing.T) {
	fakeAnnotationCache := map[string]string{
		"white/stripes": `{"webhooks":["https://webhook.atomist.com/teams/TELEPHANT"]}`,
		"16/horsepower": `{"webhooks":["https://webhook.atomist.com/teams/TLOWESTATE"],"environment":"ditch-digger"}`,
		"stevie/wonder": `{"webhooks":["https://webhook.atomist.com/teams/TALK1NGB00K"],"environment":"i-believe"}`,
	}
	for k, v := range fakeAnnotationCache {
		annotationCache[k] = v
	}
	p1 := fakePod("white/stripes")
	if p1.ObjectMeta.Name != "stripes" {
		t.Errorf("fake pod name '%s' does not match expected 'stripes'", p1.ObjectMeta.Name)
	}
	if p1.ObjectMeta.Namespace != "white" {
		t.Errorf("fake pod namespace '%s' does not match expected 'white'", p1.ObjectMeta.Namespace)
	}
	if p1.Status.Phase != "Deleted" {
		t.Errorf("fake pod phase '%s' does not match expected 'Deleted'", p1.Status.Phase)
	}
	if p1.Annotations["atomist.com/k8vent"] != fakeAnnotationCache["white/stripes"] {
		t.Errorf("fake pod annotation '%s' does not match expected '%s'", p1.Annotations["atomist.com/k8vent"],
			fakeAnnotationCache["white/stripes"])
	}
	if len(annotationCache) != len(fakeAnnotationCache)-1 {
		t.Errorf("annotation did not get deleted: %v", annotationCache)
	}
	if _, ok := annotationCache["white/stripes"]; ok {
		t.Errorf("annotation did not get deleted: %v", annotationCache)
	}
	p2 := fakePod("stevie/wonder")
	if p2.ObjectMeta.Name != "wonder" {
		t.Errorf("fake pod name '%s' does not match expected 'wonder'", p2.ObjectMeta.Name)
	}
	if p2.ObjectMeta.Namespace != "stevie" {
		t.Errorf("fake pod namespace '%s' does not match expected 'stevie'", p2.ObjectMeta.Namespace)
	}
	if p2.Status.Phase != "Deleted" {
		t.Errorf("fake pod phase '%s' does not match expected 'Deleted'", p2.Status.Phase)
	}
	if p2.Annotations["atomist.com/k8vent"] != fakeAnnotationCache["stevie/wonder"] {
		t.Errorf("fake pod annotation '%s' does not match expected '%s'", p2.Annotations["atomist.com/k8vent"],
			fakeAnnotationCache["stevie/wonder"])
	}
	if len(annotationCache) != len(fakeAnnotationCache)-2 {
		t.Errorf("annotation did not get deleted: %v", annotationCache)
	}
	if _, ok := annotationCache["stevie/wonder"]; ok {
		t.Errorf("annotation did not get deleted: %v", annotationCache)
	}
	p3 := fakePod("dj/kool")
	if p3.ObjectMeta.Name != "kool" {
		t.Errorf("fake pod name '%s' does not match expected 'kool'", p3.ObjectMeta.Name)
	}
	if p3.ObjectMeta.Namespace != "dj" {
		t.Errorf("fake pod namespace '%s' does not match expected 'sj'", p3.ObjectMeta.Namespace)
	}
	if p3.Status.Phase != "Deleted" {
		t.Errorf("fake pod phase '%s' does not match expected 'Deleted'", p3.Status.Phase)
	}
	if p3.Annotations != nil {
		t.Errorf("fake pod has annotation when it should not '%v'", p3)
	}
	if len(annotationCache) != len(fakeAnnotationCache)-2 {
		t.Errorf("no annotations should have been deleted: %v", annotationCache)
	}
}
