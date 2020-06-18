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
var Version = "0.15.1"

// packageSlug returns string containing package name and version.
func packageSlug() string {
	return Pkg + "-" + Version
}

// newerK8sventVersion checks the current k8svent version against the
// provided tags and returns `true` if the tags contain a newer
// version
func newerK8sventVersion(tags []string) bool {
	return newerVersion(Version, tags)
}

// newVersion returns true if a version newer the `version` is
// available in the tags.  If `version` is a release semantic version,
// only semantic versions are considered.  Otherwise, both pre-release
// and release versions are considered.
func newerVersion(version string, tags []string) bool {
	v, vErr := semver.Make(version)
	if vErr != nil {
		logger.Errorf("Version '%s' could not be made into a semantic version: %v", version, vErr)
		return false
	}
	release := len(v.Pre) == 0
	for _, tag := range tags {
		tagVersion, tvErr := semver.ParseTolerant(tag)
		if tvErr != nil {
			logger.Debugf("Tag '%s' is not a semantic version: %v", tag, tvErr)
			continue
		}
		if release && len(tagVersion.Pre) != 0 {
			continue
		}
		if tagVersion.GT(v) {
			return true
		}
	}
	return false
}
