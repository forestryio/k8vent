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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/google/go-cmp/cmp"
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

// listPods lists all pods in the provided namespace.  Kubernetes
// convention is that if the namespace is an empty string, pods from
// all namespaces are returned.
func listPods(clientset *kubernetes.Clientset, namespace string) ([]v1.Pod, error) {
	pods := []v1.Pod{}
	options := metav1.ListOptions{}
	for ok := true; ok; ok = (options.Continue != "") {
		podList, listErr := clientset.CoreV1().Pods(namespace).List(options)
		if listErr != nil {
			return nil, listErr
		}
		pods = append(pods, podList.Items...)
		options.Continue = podList.Continue
	}
	return pods, nil
}

// processPods iterates through the provided pods and processes those
// that do not have an identical pod in lastPods.  It returns a map of
// successfully processed pods.
func (v *Venter) processPods(pods []v1.Pod, lastPods map[string]v1.Pod) map[string]v1.Pod {
	newPods := map[string]v1.Pod{}
	for _, pod := range pods {
		slug := podSlug(pod)
		log := logger.WithField("pod", slug)
		newPods[slug] = pod
		if lastPod, ok := lastPods[slug]; ok {
			if podHealthy(pod) && cmp.Diff(pod, lastPod) == "" {
				log.Debug("Pod is healthy and state is unchanged")
				continue
			}
		}
		if err := v.processPod(pod); err != nil {
			log.Errorf("Failed to process pod: %v", err)
			delete(newPods, slug)
			continue
		}
	}
	return newPods
}

// K8PodEnv is the structure serialized and sent to the webhook
// endpoints.
type K8PodEnv struct {
	Pod v1.Pod            `json:"pod"`
	Env map[string]string `json:"env"`
}

// ProcessPods iterates through the pods and calls PostToWebhooks for
// each.
func (v *Venter) processPod(pod v1.Pod) error {
	podEnv := K8PodEnv{
		Pod: pod,
		Env: v.env,
	}
	PostToWebhooks(v.urls, &podEnv, v.secret)
	return nil
}
