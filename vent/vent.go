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
	"os/signal"
	"syscall"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Venter contains the information used to send pods to webhook
// endpoints.
type Venter struct {
	env    map[string]string
	secret string
	urls   []string
}

// Vent sets up and starts the listener for pod events, which posts
// them to the provided webhooks when it receives them.  It should
// never return.
func Vent(urls []string, namespace string, secret string, logLevel string) error {

	setupLogger(logLevel)

	logger.Infof("%s version %s starting", Pkg, Version)

	logger.Info("Creating Kubernetes API client set")
	config, configErr := rest.InClusterConfig()
	if configErr != nil {
		logger.Errorf("Failed to load in-cluster config: %v", configErr)
		return configErr
	}
	clientset, clientErr := kubernetes.NewForConfig(config)
	if clientErr != nil {
		logger.Errorf("Failed to create client from config: %v", clientErr)
		return clientErr
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	go func() {
		<-sigterm
		logger.Info("Received signal, exiting")
		os.Exit(0)
	}()

	initiateReleaseCheck()

	env := envMap()
	venter := &Venter{
		env,
		secret,
		urls,
	}

	sleepDuration := 0 * time.Second
	lastPods := map[string]v1.Pod{}
	logger.Info("Starting to vent")
	for {
		time.Sleep(sleepDuration)

		pods, listErr := listPods(clientset, namespace)
		if listErr != nil {
			logger.Errorf("Failed to list pods: %v", listErr)
			sleepDuration = 30 * time.Second
			continue
		} else {
			sleepDuration = 120 * time.Second
		}

		logger.Debugf("Processing %d pods", len(pods))
		lastPods = venter.processPods(pods, lastPods)
	}
}
