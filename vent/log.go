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
	"strings"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Entry

// setupLogger creates and configures the global logger.
func setupLogger(logLevel string) {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	level := strings.ToLower(logLevel)
	if level == "debug" {
		l.SetLevel(logrus.DebugLevel)
	} else if level == "error" {
		l.SetLevel(logrus.ErrorLevel)
	} else if level == "fatal" {
		l.SetLevel(logrus.FatalLevel)
	} else if level == "panic" {
		l.SetLevel(logrus.PanicLevel)
	} else if level == "trace" {
		l.SetLevel(logrus.TraceLevel)
	} else if level == "warn" {
		l.SetLevel(logrus.WarnLevel)
	} else {
		l.SetLevel(logrus.InfoLevel)
	}
	fields := logrus.Fields{"service": Pkg}
	if host, hostErr := os.Hostname(); hostErr == nil {
		fields["host"] = host
	}
	logger = l.WithFields(fields)
}
