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
	"strings"
	"testing"
)

func TestGetDockerTagDigest(t *testing.T) {
	tags := []string{"latest", "next"}
	for _, tag := range tags {
		digest, digestErr := getDockerTagDigest(tag)
		if digestErr != nil {
			t.Errorf("failed to get digest for %s: %v", tag, digestErr)
		}
		if !strings.HasPrefix(digest, "sha256:") {
			t.Errorf("tag %s digest does not appear to be sha256: %s", tag, digest)
		}
	}
}
