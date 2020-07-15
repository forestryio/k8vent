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
	"time"

	"github.com/blang/semver"
)

func TestVersionTag(t *testing.T) {
	vv := "1.2.3"
	v, vErr := semver.Make(vv)
	if vErr != nil {
		t.Errorf("version '%s' could not be made into a semantic version: %v", vv, vErr)
	}
	vt := versionTag(v)
	if vt != "latest" {
		t.Errorf("expected version '%s' to map to 'latest' but got '%s'", vv, vt)
	}

	pvv := "4.5.6-tmg-ciyc.421"
	pv, pvErr := semver.Make(pvv)
	if pvErr != nil {
		t.Errorf("version '%s' could not be made into a semantic version: %v", pvv, pvErr)
	}
	pvt := versionTag(pv)
	if pvt != "next" {
		t.Errorf("expected version '%s' to map to 'next' but got '%s'", pvv, pvt)
	}
}

func TestTagDuration(t *testing.T) {
	ld := tagDuration("latest")
	if ld != 24*time.Hour {
		t.Errorf("expected latest duration to be 24 hours: %v", ld)
	}
	nd := tagDuration("next")
	if nd != 4*time.Hour {
		t.Errorf("expected next duration to be 4 hours: %v", nd)
	}
}
