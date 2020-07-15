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
)

func TestIsRelease(t *testing.T) {
	releases := []string{"0.1.0", "0.3.0+build.tag.19", "0.8.0", "0.12.0", "1.42.12", "2.12.22+some-build.537", "319.518.333"}
	for _, version := range releases {
		v, vErr := semver.Make(version)
		if vErr != nil {
			t.Errorf("Version '%s' could not be made into a semantic version: %v", version, vErr)
		}
		if !isRelease(v) {
			t.Errorf("Release version '%s' not recognized as release version: %s", version, v.Pre)
		}
	}

	preReleases := []string{"0.9.111-far.far.away", "2.12.23-dark-side.2017", "1.99.999-alpha.23", "2.12.22-bright+wilco.19",
		"2.12.23-before.11+something.22"}
	for _, version := range preReleases {
		v, vErr := semver.Make(version)
		if vErr != nil {
			t.Errorf("Version '%s' could not be made into a semantic version: %v", version, vErr)
		}
		if isRelease(v) {
			t.Errorf("Prerelease version '%s' recognized as release version: %s", version, v.Pre)
		}
	}
}
