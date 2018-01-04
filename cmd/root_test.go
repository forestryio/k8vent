// Copyright © 2018 Atomist
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

package cmd

import (
	"os"
	"reflect"
	"testing"
)

func TestEnv(t *testing.T) {
	if err := os.Setenv(webhookEnv, ""); err != nil {
		t.Errorf("failed to set environment variable %s: %v", webhookEnv, err)
	}
	initConfig()
	if !reflect.DeepEqual(webhookURLs, []string{}) {
		t.Errorf("empty value for %s did not result in default webhook: %v", webhookEnv, webhookURLs)
	}

	envCheck := map[string][]string{
		"http://one":             []string{"http://one"},
		"http://one,https://two": []string{"http://one", "https://two"},
	}
	for k, v := range envCheck {
		if err := os.Setenv(webhookEnv, k); err != nil {
			t.Errorf("failed to set environment variable %s: %v", webhookEnv, err)
		}
		initConfig()
		if !reflect.DeepEqual(webhookURLs, v) {
			t.Errorf("webhook (%v) not equal to expected (%v)", webhookURLs, v)
		}
	}
}
