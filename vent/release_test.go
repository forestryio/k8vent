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

	"github.com/blang/semver"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestNewReleaseAvailable(t *testing.T) {
	nullLogger, _ := test.NewNullLogger()
	logger = nullLogger.WithField("test", "release")
	v, vErr := semver.Make(Version)
	if vErr != nil {
		t.Errorf("k8svent version '%s' could not be made into a semantic version: %v", Version, vErr)
	}
	if newReleaseAvailable(v) {
		t.Errorf("found newer version of k8svent than current unreleased version: %s", Version)
	}
}
