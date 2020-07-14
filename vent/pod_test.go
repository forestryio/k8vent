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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	v1 "k8s.io/api/core/v1"
)

func TestPodHealthy(t *testing.T) {
	pods, loadErr := loadPods("testdata/pod.json")
	if loadErr != nil {
		t.Error(loadErr.Error())
	}
	healthy := 1
	for i, pod := range pods {
		if (i < healthy) != podHealthy(pod) {
			t.Errorf("pod %d did not give expected health result", i)
		}
	}
}

func TestProcessPods(t *testing.T) {
	p0 := &testPods{
		deletedPods: []v1.Pod{},
		sentPods:    []v1.Pod{},
	}
	a0 := &processPodsArgs{
		pods:      []v1.Pod{},
		lastPods:  map[string]v1.Pod{},
		processor: p0.testProcessor,
	}
	n0 := processPods(a0)
	if len(n0) != 0 {
		t.Errorf("Expected no pods but got some: %v", n0)
	}
	if len(p0.sentPods) != 0 {
		t.Errorf("Expected no sent pods but got some: %v", p0.sentPods)
	}
	if len(p0.deletedPods) != 0 {
		t.Errorf("Expected no deleted pods but got some: %v", p0.deletedPods)
	}

	nullLogger, _ := test.NewNullLogger()
	logger = nullLogger.WithField("test", "webhook")

	pods, loadErr := loadPods("testdata/pod.json")
	if loadErr != nil {
		t.Error(loadErr.Error())
	}

	p1 := &testPods{
		deletedPods: []v1.Pod{},
		sentPods:    []v1.Pod{},
	}
	a1 := &processPodsArgs{
		pods:      []v1.Pod{pods[0]},
		lastPods:  map[string]v1.Pod{},
		processor: p1.testProcessor,
	}
	n1 := processPods(a1)
	if len(n1) != 1 {
		t.Errorf("Expected one pod but got %d: %v", len(n1), n1)
	}
	if _, ok := n1["brian-fallon/local-honey-0"]; !ok {
		t.Errorf("Expected sent pod to be 'brian-fallon/local-honey-0' but got: %v", n1)
	}
	if len(p1.sentPods) != 1 {
		t.Errorf("Expected one sent pod but got %d: %v", len(p1.sentPods), p1.sentPods)
	}
	if podSlug(p1.sentPods[0]) != "brian-fallon/local-honey-0" {
		t.Errorf("Expected sent pod to be 'brian-fallon/local-honey-0': %s", podSlug(p1.sentPods[0]))
	}
	if len(p1.deletedPods) != 0 {
		t.Errorf("Expected no deleted pods but got some: %v", p1.deletedPods)
	}

	p2 := &testPods{
		deletedPods: []v1.Pod{},
		sentPods:    []v1.Pod{},
	}
	a2 := &processPodsArgs{
		pods:      []v1.Pod{pods[0]},
		lastPods:  map[string]v1.Pod{"brian-fallon/local-honey-0": pods[0]},
		processor: p2.testProcessor,
	}
	n2 := processPods(a2)
	if len(n2) != 1 {
		t.Errorf("Expected one pod but got %d: %v", len(n2), n2)
	}
	if _, ok := n2["brian-fallon/local-honey-0"]; !ok {
		t.Errorf("Expected sent pod to be 'brian-fallon/local-honey-0' but got: %v", n2)
	}
	if len(p2.sentPods) != 0 {
		t.Errorf("Expected no sent pods but got some: %v", p2.sentPods)
	}
	if len(p2.deletedPods) != 0 {
		t.Errorf("Expected no deleted pods but got some: %v", p2.deletedPods)
	}

	p3 := &testPods{
		deletedPods: []v1.Pod{},
		sentPods:    []v1.Pod{},
	}
	a3 := &processPodsArgs{
		pods:      []v1.Pod{},
		lastPods:  map[string]v1.Pod{"brian-fallon/local-honey-0": pods[0]},
		processor: p3.testProcessor,
	}
	n3 := processPods(a3)
	if len(n3) != 0 {
		t.Errorf("Expected no pods but got %d: %v", len(n3), n3)
	}
	if len(p3.sentPods) != 0 {
		t.Errorf("Expected no sent pods but got some: %v", p3.sentPods)
	}
	if len(p3.deletedPods) != 1 {
		t.Errorf("Expected one deleted pod but got %d: %v", len(p3.deletedPods), p3.deletedPods)
	}
	if podSlug(p3.deletedPods[0]) != "brian-fallon/local-honey-0" {
		t.Errorf("Expected deleted pod to be 'brian-fallon/local-honey-0': %s", podSlug(p3.deletedPods[0]))
	}

	p4 := &testPods{
		deletedPods: []v1.Pod{},
		sentPods:    []v1.Pod{},
	}
	a4 := &processPodsArgs{
		pods: []v1.Pod{pods[0], pods[1], pods[2]},
		lastPods: map[string]v1.Pod{
			"brian-fallon/local-honey-0": pods[0],
			"brian-fallon/local-honey-1": pods[1],
			"brian-fallon/local-honey-3": pods[3],
		},
		processor: p4.testProcessor,
	}
	n4 := processPods(a4)
	if len(n4) != 3 {
		t.Errorf("Expected no pods but got %d: %v", len(n4), n4)
	}
	if len(p4.sentPods) != 2 {
		t.Errorf("Expected no sent pods but got some: %v", p4.sentPods)
	}
	if podSlug(p4.sentPods[0]) != "brian-fallon/local-honey-1" {
		t.Errorf("Expected sent pod to be 'brian-fallon/local-honey-1': %s", podSlug(p4.sentPods[0]))
	}
	if podSlug(p4.sentPods[1]) != "brian-fallon/local-honey-2" {
		t.Errorf("Expected sent pod to be 'brian-fallon/local-honey-2': %s", podSlug(p4.sentPods[1]))
	}
	if len(p4.deletedPods) != 1 {
		t.Errorf("Expected one deleted pod but got %d: %v", len(p4.deletedPods), p4.deletedPods)
	}
	if podSlug(p4.deletedPods[0]) != "brian-fallon/local-honey-3" {
		t.Errorf("Expected deleted pod to be 'brian-fallon/local-honey': %s", podSlug(p4.deletedPods[0]))
	}
}

func loadPods(podFile string) (o []v1.Pod, e error) {
	podBytes, readErr := ioutil.ReadFile(podFile)
	if readErr != nil {
		return o, fmt.Errorf("failed to read pod JSON file %s: %v", podFile, readErr)
	}

	pods := []v1.Pod{}
	if err := json.Unmarshal(podBytes, &pods); err != nil {
		return o, fmt.Errorf("failed to unmarshal pod JSON: %v", err)
	}
	return pods, nil
}

type testPods struct {
	deletedPods []v1.Pod
	sentPods    []v1.Pod
}

func (p *testPods) testProcessor(pod v1.Pod) error {
	if pod.Status.Phase == "Deleted" {
		p.deletedPods = append(p.deletedPods, pod)
	} else {
		p.sentPods = append(p.sentPods, pod)
	}
	return nil
}
