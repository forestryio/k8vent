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
	"fmt"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

func TestNewerVersion(t *testing.T) {
	var hook *test.Hook
	logger, hook = test.NewNullLogger()

	assert := require.New(t)

	for _, bad := range []string{"x.y.z", "M.N.P", "1", "2.3", "4.5.6.7", "1.02.3", "1.2.3-x.01"} {
		if newerVersion(bad, []string{}) {
			t.Errorf("successfully parsed bad version: '%s'", bad)
		}
		assert.Equal(1, len(hook.Entries))
		assert.Equal(logrus.ErrorLevel, hook.LastEntry().Level)
		em := fmt.Sprintf("Version '%s' could not be made into a semantic version: ", bad)
		assert.True(strings.HasPrefix(hook.LastEntry().Message, em))
		hook.Reset()
	}

	releaseSlices := [][]string{
		[]string{"0.1.0", "0.2.0", "0.3.0", "0.4.0", "0.8.0", "0.12.0", "1.42.12", "2.12.22", "latest"},
		[]string{"0.12.0", "0.1.0", "2.12.22", "0.2.0", "0.4.0", "0.3.0", "1.42.12", "0.8.0", "latest"},
	}
	prereleaseSlices := [][]string{
		[]string{"0.9.111-far.far.away", "2.12.23-dark-side.2017", "1.99.999-alpha.23", "latest"},
		[]string{"0.12.0+sdm.14", "2.12.22-bright+wilco.19", "1.42.12+being.there.7", "2.12.23-before.11+something.22", "latest"},
	}
	tagSlices := [][]string{[]string{"0.9.111-far.far.away", "0.12.0", "1.42.12", "2.12.22", "2.12.23-dark-side.2017", "latest"}}
	tagSlices = append(tagSlices, releaseSlices...)
	tagSlices = append(tagSlices, prereleaseSlices...)
	newVersions := []string{"2.14.9", "2.12.23-zbranch.20200101", "2.12.23", "3.1.4", "3.0.0-taj-mahal.3754"}
	for _, v := range newVersions {
		for _, tags := range tagSlices {
			if newerVersion(v, tags) {
				t.Errorf("erroneously found version newer than '%s' in %v", v, tags)
			}
		}
	}
	oldVersions := []string{"2.4.99", "2.12.22-before.10", "2.12.21", "0.11.4"}
	for _, v := range oldVersions {
		for _, tags := range releaseSlices {
			if !newerVersion(v, tags) {
				t.Errorf("failed to find version newer than '%s' in %v", v, tags)
			}
		}
	}
	oldPrereleaseVersions := []string{"2.4.99-smog+along", "2.12.23-before.10", "2.12.21-weird.tales.99", "0.11.4-until-you"}
	for _, v := range oldPrereleaseVersions {
		for _, tags := range prereleaseSlices {
			if !newerVersion(v, tags) {
				t.Errorf("failed to find version newer than '%s' in %v", v, tags)
			}
		}
	}

	logger.SetLevel(logrus.DebugLevel)
	hook.Reset()
	notSemVer := []string{"should", "ignore", "tags", "like", "latest", "that", "are", "not", "semver"}
	if newerVersion("0.1.0", notSemVer) {
		t.Errorf("treated non-semver tags as semver: %v", notSemVer)
	}
	assert.Equalf(len(notSemVer), len(hook.Entries), "debug log messages")
	for i, log := range hook.Entries {
		assert.Equal(logrus.DebugLevel, log.Level)
		em := fmt.Sprintf("Tag '%s' is not a semantic version: ", notSemVer[i])
		assert.Truef(strings.HasPrefix(log.Message, em), "Expect '%s' to start with: %s", log.Message, em)
	}
	hook.Reset()

}
