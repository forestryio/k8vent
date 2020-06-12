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
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// generateSignature creates a HMAC/SHA-1 signature for payload using key.
func generateSignature(payload []byte, key string) (s string, e error) {
	mac := hmac.New(sha1.New, []byte(key))
	if _, err := mac.Write(payload); err != nil {
		return s, fmt.Errorf("failed to write payload to HMAC: %v", err)
	}
	sum := mac.Sum(nil)
	sig := hex.EncodeToString(sum)
	return "sha1=" + sig, nil
}
