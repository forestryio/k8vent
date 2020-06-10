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
