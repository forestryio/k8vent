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
	"github.com/blang/semver"
)

// Pkg is the canonical package name of this application.
const Pkg = "k8svent"

// Version is the version of this application.  It must be a var and
// initialized with a constant expression so we can set it during the
// linking stage of build.
var Version = "0.17.1"

// packageSlug returns string containing package name and version.
func packageSlug() string {
	return Pkg + "-" + Version
}

// isRelease returns true if provided version is a release, i.e., not
// a prerelease, false otherwise.
func isRelease(version semver.Version) bool {
	return len(version.Pre) == 0
}
