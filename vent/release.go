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
	"os"
	"time"

	"github.com/blang/semver"
)

// initiateReleaseCheck starts a go routine to periodically check for
// a new release.
func initiateReleaseCheck() {
	go func() {
		tag := "next"
		if v, vErr := semver.Make(Version); vErr == nil {
			tag = versionTag(v)
		} else {
			logger.Warnf("Version '%s' could not be made into a semantic version: %v", Version, vErr)
		}
		logger.Infof("Using Docker image tag '%s' for digest check", tag)
		rest := 0 * time.Second
		lastDigest := ""
		for {
			time.Sleep(rest)
			digest, digestErr := getDockerTagDigest(tag)
			if digestErr != nil {
				logger.Errorf("Failed to get Docker image digest for tag %s: %v", tag, digestErr)
				rest = 1 * time.Hour
				continue
			}
			if lastDigest == "" {
				lastDigest = digest
			}
			if digest != lastDigest {
				logger.Info("New version detected, exiting")
				os.Exit(0)
			}
			rest = tagDuration(tag)
		}
	}()
}

// versionTag returns the Docker image tag that maps to the provided
// version.  Specifically, it returns "latest" for release versions
// and "next" for pre-release versions.
func versionTag(v semver.Version) string {
	if isRelease(v) {
		return "latest"
	}
	return "next"
}

// tagDuration returns the duration to sleep for between checks for a
// given Docker image tag.
func tagDuration(tag string) time.Duration {
	if tag == "latest" {
		return 24 * time.Hour
	}
	return 4 * time.Hour
}
