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
	if err := os.Unsetenv(webhookEnv); err != nil {
		t.Errorf("failed to set environment variable %s: %v", webhookEnv, err)
	}
	initConfig()
	if !reflect.DeepEqual(webhookURLs, []string{}) {
		t.Errorf("unset %s did not result in default webhook: %v", webhookEnv, webhookURLs)
	}

	urlCheck := map[string][]string{
		"http://one":             []string{"http://one"},
		"http://one,https://two": []string{"http://one", "https://two"},
	}
	for k, v := range urlCheck {
		if err := os.Setenv(webhookEnv, k); err != nil {
			t.Errorf("failed to set environment variable %s: %v", webhookEnv, err)
		}
		initConfig()
		if !reflect.DeepEqual(webhookURLs, v) {
			t.Errorf("webhook (%v) not equal to expected (%v)", webhookURLs, v)
		}
	}

	if err := os.Unsetenv(namespaceEnv); err != nil {
		t.Errorf("failed to set environment variable %s: %v", webhookEnv, err)
	}
	initConfig()
	if namespace != "" {
		t.Errorf("unset %s did not result in empty namespace: %v", namespaceEnv, namespace)
	}

	nsCheck := []string{
		"police",
		"reggatta-de-blanc",
	}
	for _, ns := range nsCheck {
		if err := os.Setenv(namespaceEnv, ns); err != nil {
			t.Errorf("failed to set environment variable %s: %v", namespaceEnv, err)
		}
		initConfig()
		if namespace != ns {
			t.Errorf("namespace (%v) not equal to expected (%v)", namespace, ns)
		}
	}
}
