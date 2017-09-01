// Copyright © 2017 Atomist
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
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"k8s.io/kubernetes/pkg/api/v1"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_3"
)

func addPodSpec(event v1.Event, k8Client *client.Clientset) v1.Event {
	if event.InvolvedObject.Kind != "Pod" || event.Reason == "Killing" {
		return event
	}
	podJSON, podErr := getPodSpecJSON(event.InvolvedObject.Name, event.InvolvedObject.Namespace, k8Client)
	if podErr != nil {
		log.Printf("failed to retrieve Pod: %v: %v\n", podErr, event)
		return event
	}
	annotations := event.Annotations
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations["spec"] = podJSON
	event.SetAnnotations(annotations)
	return event
}

func getPodSpecJSON(name string, namespace string, k8Client *client.Clientset) (p string, e error) {
	pod, podErr := k8Client.Core().Pods(namespace).Get(name)
	if podErr == nil {
		return podSpecToString(pod.Spec)
	}

	return getPodSpecFromAncestor(name, namespace, k8Client)
}

func podSpecToString(podSpec v1.PodSpec) (string, error) {
	podJSON, podJSONErr := json.Marshal(podSpec)
	if podJSONErr != nil {
		return "", fmt.Errorf("failed to serialize Pod to json: %v: %v", podJSONErr, podSpec)
	}
	return string(podJSON), nil
}

func getPodSpecFromAncestor(name string, namespace string, k8Client *client.Clientset) (string, error) {
	podRE, reErr := regexp.Compile(`^([-\w]+)(-\d+)?-[a-z\d]{5}$`)
	if reErr != nil {
		return "", reErr
	}
	sms := podRE.FindStringSubmatch(name)
	if sms == nil {
		return "", fmt.Errorf("failed to get spec for pod %s and it does not appear to be part of a deployment/rs/rc/ds", name)
	}

	var getParents []func() *v1.PodSpec
	var parentName string
	if sms[2] == "" {
		parentName = sms[1]
		getParents = []func() *v1.PodSpec{
			func() *v1.PodSpec {
				t, err := k8Client.Extensions().ReplicaSets(namespace).Get(parentName)
				if err != nil {
					return nil
				}
				return &t.Spec.Template.Spec
			},
			func() *v1.PodSpec {
				t, err := k8Client.Core().ReplicationControllers(namespace).Get(parentName)
				if err != nil {
					return nil
				}
				return &t.Spec.Template.Spec
			},
			func() *v1.PodSpec {
				t, err := k8Client.Extensions().DaemonSets(namespace).Get(parentName)
				if err != nil {
					return nil
				}
				return &t.Spec.Template.Spec
			},
		}
	} else {
		parentName = sms[1] + sms[2]
		grandParentName := sms[1]
		getParents = []func() *v1.PodSpec{
			func() *v1.PodSpec {
				t, err := k8Client.Extensions().ReplicaSets(namespace).Get(parentName)
				if err != nil {
					return nil
				}
				return &t.Spec.Template.Spec
			},
			func() *v1.PodSpec {
				t, err := k8Client.Extensions().Deployments(namespace).Get(grandParentName)
				if err != nil {
					return nil
				}
				return &t.Spec.Template.Spec
			},
		}
	}
	var podSpecRef *v1.PodSpec
	for _, getParent := range getParents {
		podSpecRef = getParent()
		if podSpecRef != nil {
			break
		}
	}
	if podSpecRef == nil {
		return "", fmt.Errorf("failed to get spec for pod %s from any parent", name)
	}
	return podSpecToString(*podSpecRef)
}
