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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

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
