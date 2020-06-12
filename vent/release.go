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
)

// initiateReleaseCheck starts a go routine to periodically check for
// a new release.
func initiateReleaseCheck() {
	go func() {
		for {
			time.Sleep(24 * time.Hour)
			if newReleaseAvailable() {
				logger.Info("New version detected, exiting")
				os.Exit(0)
			}
		}
	}()
}

// newReleaseAvailable queries the Docker Hub API for tags and sees if
// a newer tag is available.
func newReleaseAvailable() bool {
	tags, tagsErr := getDockerTags()
	if tagsErr != nil {
		logger.Errorf("Failed to get Docker tags: %v", tagsErr)
		return false
	}
	if len(tags) < 1 {
		return false
	}
	return newerK8sventVersion(tags)
}
