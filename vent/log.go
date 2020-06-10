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

var logger *logrus.Logger

// setupLogger creates and configures the global logger.
func setupLogger(logLevel string) {
	logger = logrus.New()
	logger.WithFields(logrus.Fields{
		"service":     "k8svent",
		"environment": os.Getenv("ATOMIST_ENVIRONMENT"),
	})
	logger.SetFormatter(&logrus.JSONFormatter{})
	level := strings.ToLower(logLevel)
	if level == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	} else if level == "error" {
		logger.SetLevel(logrus.ErrorLevel)
	} else if level == "fatal" {
		logger.SetLevel(logrus.FatalLevel)
	} else if level == "panic" {
		logger.SetLevel(logrus.PanicLevel)
	} else if level == "trace" {
		logger.SetLevel(logrus.TraceLevel)
	} else if level == "warn" {
		logger.SetLevel(logrus.WarnLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
}
