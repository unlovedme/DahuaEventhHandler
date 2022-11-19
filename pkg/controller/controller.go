/*
Copyright 2016 Skippbox, Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
	"github.com/bitnami-labs/kubewatch/pkg/handlers"
	"github.com/bitnami-labs/kubewatch/pkg/utils"
	"github.com/sirupsen/logrus"

	apps_v1 "k8s.io/api/apps/v1"
	batch_v1 "k8s.io/api/batch/v1"
	api_v1 "k8s.io/api/core/v1"
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
	rbac_v1beta1 "k8s.io/api/rbac/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const maxRetries = 5

var serverStartTime time.Time

// Event indicate the informerEvent
type Event struct {
	key          string
	eventType    string
	namespace    string
	resourceType string
}

// Controller object
type Controller struct {
	logger       *logrus.Entry
	clientset    kubernetes.Interface
	queue        workqueue.RateLimitingInterface
	informer     cache.SharedIndexInformer
	eventHandler handlers.Handler
}

// Start prepares watchers and run their controllers, then waits for process termination signals
func Start(conf *config.Config, eventHandler handlers.Handler) {
	var kubeClient kubernetes.Interface

	if _, err := rest.InClusterConfig(); err != nil {
		kubeClient = utils.GetClientOutOfCluster()
	} else {
		kubeClient = utils.GetClient()
	}

	// Adding Default Critical Alerts
	// For Capturing Critical Event NodeNotReady in Nodes
	nodeNotReadyInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				options.FieldSelector = "involvedObject.kind=Node,type=Normal,reason=NodeNotReady"
				return kubeClient.CoreV1().Events(conf.Namespace).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				options.FieldSelector = "involvedObject.kind=Node,type=Normal,reason=NodeNotReady"
				return kubeClient.CoreV1().Events(conf.Namespace).Watch(options)
			},
		},
		&api_v1.Event{},
		0, //Skip resync
		cache.Indexers{},
	)

	nodeNotReadyController := newResourceController(kubeClient, eventHandler, nodeNotReadyInformer, "NodeNotReady")
	stopNodeNotReadyCh := make(chan struct{})
	defer close(stopNodeNotReadyCh)

	go nodeNotReadyController.Run(stopNodeNotReadyCh)

	// For Capturing Critical Event NodeReady in Nodes
	nodeReadyInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				options.FieldSelector = "involvedObject.kind=Node,type=Normal,reason=NodeReady"
				return kubeClient.CoreV1().Events(conf.Namespace).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				options.FieldSelector = "involvedObject.kind=Node,type=Normal,reason=NodeReady"
				return kubeClient.CoreV1().Events(conf.Namespace).Watch(options)
			},
		},
		&api_v1.Event{},
		0, //Skip resync
		cache.Indexers{},
	)

	nodeReadyController := newResourceController(kubeClient, eventHandler, nodeReadyInformer, "NodeReady")
	stopNodeReadyCh := make(chan struct{})
	defer close(stopNodeReadyCh)

	go nodeReadyController.Run(stopNodeReadyCh)

	// For Capturing Critical Event NodeRebooted in Nodes
	nodeRebootedInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				options.FieldSelector = "involvedObject.kind=Node,type=Warning,reason=Rebooted"
				return kubeClient.CoreV1().Events(conf.Namespace).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				options.FieldSelector = "involvedObject.kind=Node,type=Warning,reason=Rebooted"
				return kubeClient.CoreV1().Events(conf.Namespace).Watch(options)
			},
		},
		&api_v1.Event{},
		0, //Skip resync
		cache.Indexers{},
	)

	nodeRebootedController := newResourceController(kubeClient, eventHandler, nodeRebootedInformer, "NodeRebooted")
	stopNodeRebootedCh := make(chan struct{})
	defer close(stopNodeRebootedCh)

	go nodeRebootedController.Run(stopNodeRebootedCh)

	// User Configured Events
	if conf.Resource.Pod {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Pods(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Pods(conf.Namespace).Watch(options)
				},
			},
			&api_v1.Pod{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "pod")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)

		// For Capturing CrashLoopBackOff Events in pods
		backoffInformer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					options.FieldSelector = "involvedObject.kind=Pod,type=Warning,reason=BackOff"
					return kubeClient.CoreV1().Events(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					options.FieldSelector = "involvedObject.kind=Pod,type=Warning,reason=BackOff"
					return kubeClient.CoreV1().Events(conf.Namespace).Watch(options)
				},
			},
			&api_v1.Event{},
			0, //Skip resync
			cache.Indexers{},
		)

		backoffcontroller := newResourceController(kubeClient, eventHandler, backoffInformer, "Backoff")
		stopBackoffCh := make(chan struct{})
		defer close(stopBackoffCh)

		go backoffcontroller.Run(stopBackoffCh)

	}

	if conf.Resource.DaemonSet {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.AppsV1().DaemonSets(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.AppsV1().DaemonSets(conf.Namespace).Watch(options)
				},
			},
			&apps_v1.DaemonSet{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "daemon set")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.ReplicaSet {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.AppsV1().ReplicaSets(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.AppsV1().ReplicaSets(conf.Namespace).Watch(options)
				},
			},
			&apps_v1.ReplicaSet{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "replica set")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Services {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Services(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Services(conf.Namespace).Watch(options)
				},
			},
			&api_v1.Service{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "service")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Deployment {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.AppsV1().Deployments(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.AppsV1().Deployments(conf.Namespace).Watch(options)
				},
			},
			&apps_v1.Deployment{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "deployment")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Namespace {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Namespaces().List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Namespaces().Watch(options)
				},
			},
			&api_v1.Namespace{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "namespace")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.ReplicationController {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().ReplicationControllers(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().ReplicationControllers(conf.Namespace).Watch(options)
				},
			},
			&api_v1.ReplicationController{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "replication controller")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Job {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.BatchV1().Jobs(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.BatchV1().Jobs(conf.Namespace).Watch(options)
				},
			},
			&batch_v1.Job{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "job")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.Node {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Nodes().List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Nodes().Watch(options)
				},
			},
			&api_v1.Node{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "node")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.ServiceAccount {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().ServiceAccounts(conf.Namespace).List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().ServiceAccounts(conf.Namespace).Watch(options)
				},
			},
			&api_v1.ServiceAccount{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "service account")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.ClusterRole {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.RbacV1beta1().ClusterRoles().List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.RbacV1beta1().ClusterRoles().Watch(options)
				},
			},
			&rbac_v1beta1.ClusterRole{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandler, informer, "cluster role")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh)
	}

	if conf.Resource.PersistentVolume {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().PersistentVolumes().List(options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Int