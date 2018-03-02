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

package vent

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/Sirupsen/logrus"
)

const maxRetries = 5

// Controller object
// Based on Controller from github.com/skippbox/kubewatch
type Controller struct {
	logger    *logrus.Entry
	clientset kubernetes.Interface
	queue     workqueue.RateLimitingInterface
	informer  cache.SharedIndexInformer
	urls      []string
	env       map[string]string
}

// Vent sets up and starts the listener for pod events, which posts
// them to the provided webhooks when it receives them.  It should
// never return.
func Vent(urls []string) (e error) {

	config, configErr := rest.InClusterConfig()
	if configErr != nil {
		return configErr
	}

	clientset, clientErr := kubernetes.NewForConfig(config)
	if clientErr != nil {
		return clientErr
	}

	c := newController(clientset, urls)
	stopCh := make(chan struct{})
	defer close(stopCh)

	go c.Run(stopCh)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm

	return nil
}

func newController(client kubernetes.Interface, urls []string) *Controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return client.CoreV1().Pods(metav1.NamespaceAll).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return client.CoreV1().Pods(metav1.NamespaceAll).Watch(options)
			},
		},
		&v1.Pod{},
		0, // Skip resync
		cache.Indexers{},
	)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	})

	env := map[string]string{}
	for _, envVar := range os.Environ() {
		keyVal := strings.SplitN(envVar, "=", 2)
		env[keyVal[0]] = keyVal[1]
	}

	return &Controller{
		logger:    logrus.WithField("pkg", "k8vent-pod"),
		clientset: client,
		informer:  informer,
		queue:     queue,
		urls:      urls,
		env:       env,
	}
}

// Run starts the k8vent controller
func (c *Controller) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	c.logger.Info("Starting k8vent controller")

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	c.logger.Info("k8vent controller synced and ready")

	wait.Until(c.runWorker, time.Second, stopCh)
}

// HasSynced is required for the cache.Controller interface.
func (c *Controller) HasSynced() bool {
	return c.informer.HasSynced()
}

// LastSyncResourceVersion is required for the cache.Controller interface.
func (c *Controller) LastSyncResourceVersion() string {
	return c.informer.LastSyncResourceVersion()
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
		// continue looping
	}
}

func (c *Controller) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.processItem(key.(string))
	if err == nil {
		// No error, reset the ratelimit counters
		c.queue.Forget(key)
	} else if c.queue.NumRequeues(key) < maxRetries {
		c.logger.Errorf("Error processing %s (will retry): %v", key, err)
		c.queue.AddRateLimited(key)
	} else {
		// err != nil and too many retries
		c.logger.Errorf("Error processing %s (giving up): %v", key, err)
		c.queue.Forget(key)
		utilruntime.HandleError(err)
	}

	return true
}

// K8PodEnv is the structure serialized and sent to the webhook
// endpoints.
type K8PodEnv struct {
	Pod v1.Pod            `json:"pod"`
	Env map[string]string `json:"env"`
}

const k8ventAnnotationKey = "atomist.com/k8vent"

// K8VentPodAnnotation defines the valid structure of the
// "atomist.com/k8vent" pod annotation.
type K8VentPodAnnotation struct {
	Webhooks    []string `json:"webhooks"`
	Environment string   `json:"environment"`
}

// annotationCache provides a way for deleted pods to still use their
// k8vent annotations.
var annotationCache = map[string]string{}

// processItem looks up key in the indexer, converts it into a v1.Pod,
// and calls PostToWebhooks.
func (c *Controller) processItem(key string) error {
	c.logger.Infof("Processing change to Pod %s", key)

	obj, exists, err := c.informer.GetIndexer().GetByKey(key)
	if err != nil {
		return fmt.Errorf("Error fetching object with key %s from store: %v", key, err)
	}
	pod := v1.Pod{}
	env := c.env
	webhookURLs := c.urls
	if exists {
		var extractErr error
		var annot *K8VentPodAnnotation
		pod, annot, extractErr = extractPod(obj, c.logger)
		if extractErr != nil {
			return extractErr
		}
		if annot != nil {
			if annot.Environment != "" {
				env = map[string]string{}
				for k, v := range c.env {
					env[k] = v
				}
				env["ATOMIST_ENVIRONMENT"] = annot.Environment
			}
			if annot.Webhooks != nil && len(annot.Webhooks) > 0 {
				webhookURLs = annot.Webhooks
			}
			annotJSON, jsonErr := json.Marshal(annot)
			if jsonErr != nil {
				c.logger.Infof("failed to marshal k8vent annotation '%v': %v", annot, jsonErr)
			} else {
				annotationCache[key] = string(annotJSON)
			}
		} else {
			if _, ok := annotationCache[key]; ok {
				// old version had an annotation and this one does not
				delete(annotationCache, key)
			}
		}
	} else {
		splitName := strings.SplitN(key, "/", 2)
		pod = v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      splitName[1],
				Namespace: splitName[0],
			},
			Status: v1.PodStatus{
				Phase: "Deleted",
			},
		}
		if annot, ok := annotationCache[key]; ok {
			pod.Annotations[k8ventAnnotationKey] = annot
			delete(annotationCache, key)
		}
	}

	postIt := K8PodEnv{
		Pod: pod,
		Env: env,
	}
	PostToWebhooks(webhookURLs, postIt)

	return nil
}

// extractPod tries to convert object into a v1.Pod by marshaling it
// to JSON and back again, returning it as p.  If there is an
// "atomist.com/k8vent" annotation on the pod, it parses it as JSON
// and returns the value of the webhooks key, if it exists, as w.  If
// an error occurs, e will be non-nil.
func extractPod(obj interface{}, logger *logrus.Entry) (p v1.Pod, d *K8VentPodAnnotation, e error) {
	objJSON, jsonErr := json.Marshal(obj)
	if jsonErr != nil {
		return p, nil, fmt.Errorf("failed to marshal object to JSON: %v", jsonErr)
	}
	pod := v1.Pod{}
	if err := json.Unmarshal(objJSON, &pod); err != nil {
		return p, nil, fmt.Errorf("failed to unmarshal object as Pod: %v", err)
	}

	ventAnnot, annotOK := pod.Annotations[k8ventAnnotationKey]
	if !annotOK {
		return pod, nil, nil
	}

	annot := &K8VentPodAnnotation{}
	if err := json.Unmarshal([]byte(ventAnnot), annot); err != nil {
		logger.Infof("failed to unmarshal k8vent annotation '%s': %v", ventAnnot, err)
		return pod, nil, nil
	}

	return pod, annot, nil
}
