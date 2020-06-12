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
	"testing"
)

func TestGetDockerTags(t *testing.T) {
	tags, err := getDockerTags()
	if err != nil {
		t.Errorf("failed to get Docker tags: %v", err)
	}
	found := map[string]bool{
		"0.11.0":                false,
		"0.12.0":                false,
		"0.13.0":                false,
		"0.13.1":                false,
		"0.14.0":                false,
		"0.14.1-20200612183032": false,
		"latest":                false,
	}
	for _, v := range tags {
		if _, ok := found[v]; ok {
			found[v] = true
		}
	}
	for v, f := range found {
		if !f {
			t.Errorf("tags did not include '%s': %v", v, tags)
		}
	}
}
