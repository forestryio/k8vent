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
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExtractPropertyString(t *testing.T) {
	validResponses := [][]string{
		[]string{"correlation-id", `{"correlation-id":"0"}`, "0"},
		[]string{"correlation-id", `{"correlation-id":"1","message":"successfully posted event"}`, "1"},
		[]string{"correlation-id", `{"status":"ok","correlation-id":"2","message":"successfully posted event"}`, "2"},
		[]string{"correlation-id", `{"error":null,"status":"ok","correlation-id":"3","message":"successfully posted event"}`, "3"},
		[]string{"token", `{"token":"t0k3n","access_token":"ACCESS","expires_in":"3","issued_at":"2"}`, "t0k3n"},
	}
	for _, kbe := range validResponses {
		resp := &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader([]byte(kbe[1]))),
		}
		if v, err := extractPropertyString(resp, kbe[0]); err == nil {
			if v != kbe[2] {
				t.Errorf("extracted %s (%s) is not expected value (%s)", kbe[0], v, kbe[2])
			}
		} else {
			t.Errorf("failed to extract %s from '%s': %v", kbe[0], kbe[1], err)
		}
	}

	invalidResponses := []string{
		"",
		"{}",
		`{"status":"ok"}`,
		`{"message":"successfully posted event"}`,
	}
	for _, r := range invalidResponses {
		resp := &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader([]byte(r))),
		}
		if v, err := extractPropertyString(resp, "correlation-id"); err == nil {
			t.Errorf("unexpectedly extracted correlation ID from invalid response '%s': %s", r, v)
		}
	}
}

func TestExtractPropertyStringSlice(t *testing.T) {
	valid := `{"name":"atomist/k8svent","tags":["0.11.0","0.12.0","0.13.0","0.13.1","0.14.0-poll-10.20200610093614","0.14.0-poll-10.20200610102332","0.14.0-poll-10.20200611073244","0.14.0-poll-10.20200611092344","0.14.0-poll-10.20200611142122","0.14.0-poll-10.20200612073112","0.14.0","0.14.1-20200612183032","latest"]}`
	e := []string{
		"0.11.0",
		"0.12.0",
		"0.13.0",
		"0.13.1",
		"0.14.0-poll-10.20200610093614",
		"0.14.0-poll-10.20200610102332",
		"0.14.0-poll-10.20200611073244",
		"0.14.0-poll-10.20200611092344",
		"0.14.0-poll-10.20200611142122",
		"0.14.0-poll-10.20200612073112",
		"0.14.0",
		"0.14.1-20200612183032",
		"latest",
	}
	resp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(valid))),
	}
	if v, err := extractPropertyStringSlice(resp, "tags"); err == nil {
		d := cmp.Diff(v, e)
		if d != "" {
			t.Errorf("extracted tags is not expected value: %s", d)
		}
	} else {
		t.Errorf("failed to extract tags: %v", err)
	}
}
