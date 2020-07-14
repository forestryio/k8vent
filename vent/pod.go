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
	"github.com/google/go-cmp/cmp"
	v1 "k8s.io/api/core/v1"
)

// processPodsArgs provides the argument fo processPods.
type processPodsArgs struct {
	// pods are the pods to process.
	pods []v1.Pod
	// lastPods are the last set of pods processed.
	lastPods map[string]v1.Pod
	// processor is the function that processes each individual pod
	// that is determined to be either new or unhealthy.
	processor func(v1.Pod) error
}

// processPods iterates through the provided pods and processes those
// that do not have an identical pod in lastPods or are not healthy.
// It returns a map of successfully processed pods.
func processPods(args *processPodsArgs) map[string]v1.Pod {
	newPods := map[string]v1.Pod{}
	for _, pod := range args.pods {
		slug := podSlug(pod)
		log := logger.WithField("pod", slug)
		newPods[slug] = pod
		if lastPod, ok := args.lastPods[slug]; ok {
			delete(args.lastPods, slug)
			if podHealthy(pod) && cmp.Diff(pod, lastPod) == "" {
				log.Debug("Pod is healthy and state is unchanged")
				continue
			}
		}
		if err := args.processor(pod); err != nil {
			log.Errorf("Failed to process pod: %v", err)
			delete(newPods, slug)
			continue
		}
	}
	for slug, deletedPod := range args.lastPods {
		log := logger.WithField("pod", slug)
		deletedPod.Status.Phase = "Deleted"
		if err := args.processor(deletedPod); err != nil {
			log.Errorf("Failed to process pod: %v", err)
			continue
		}
	}
	return newPods
}

// podSlug returns a string uniquely identifying a pod in a Kubernetes
// cluster.
func podSlug(pod v1.Pod) string {
	return pod.ObjectMeta.Namespace + "/" + pod.ObjectMeta.Name
}

// podHealthy determines if a pod is healthy.
func podHealthy(pod v1.Pod) bool {
	if pod.Status.Phase != v1.PodRunning {
		return false
	}
	for _, condition := range pod.Status.Conditions {
		if condition.Status != v1.ConditionTrue {
			return false
		}
	}
	for _, containerStatus := range pod.Status.InitContainerStatuses {
		if !containerHealthy(containerStatus, true) {
			return false
		}
	}
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerHealthy(containerStatus, false) {
			return false
		}
	}
	return true
}

// containerHealthy interrogates the container status to determine if
// the container is fully healthy.
func containerHealthy(containerStatus v1.ContainerStatus, init bool) bool {
	if !containerStatus.Ready {
		return false
	}
	if containerStatus.State.Waiting != nil {
		return false
	}
	if init {
		if containerStatus.State.Terminated != nil {
			if containerStatus.State.Terminated.ExitCode != 0 {
				return false
			}
		} else if containerStatus.State.Running == nil {
			return false
		}
	} else {
		if containerStatus.State.Terminated != nil {
			return false
		}
		if containerStatus.State.Running == nil {
			return false
		}
	}
	return true
}

// ProcessPods iterates through the pods and calls PostToWebhooks for
// each.
func (v *Venter) processPod(pod v1.Pod) error {
	payload := webhookPayload{Pod: pod}
	postToWebhooks(v.urls, &payload, v.secret)
	return nil
}
