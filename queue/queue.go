/*
Copyright 2018 Jim Zhang (jim.zoumo@gmail.com). All rights reserved.

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

package queue

import (
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

type HandleResult struct {
	// Requeue tells the Controller to requeue the reconcile key.  Defaults to false.
	Requeue bool

	// RequeueAfter if greater than 0, tells the Controller to requeue the reconcile key after the Duration.
	// Implies that Requeue is true, there is no need to set Requeue to true at the same time as RequeueAfter.
	RequeueAfter time.Duration
}

type Handler func(obj interface{}) (HandleResult, error)

// Queue is a wrapper of kubernetes workqueue to do asynchronous work easily.
// It requires a Handler and an optional key function.
// After starting the Queue, you can call the Enqueque function to enqueue items.
// Queue will get key from the items by keyFunc, and add the key to the rate limit workqueue.
// The worker will be invoked to call the Handler.
type Queue struct {
	// handler is called for each item in the queue
	handler Handler

	// queue is the work queue the worker polls
	queue workqueue.RateLimitingInterface

	waitGroup sync.WaitGroup

	stopCh chan struct{}
}

// NewQueue returns a new Queue
func NewQueue(handler Handler) *Queue {
	return &Queue{
		queue:     workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		handler:   handler,
		waitGroup: sync.WaitGroup{},
		stopCh:    make(chan struct{}),
	}
}

// Run starts n workers to sync
func (sq *Queue) Run(workers int) {
	for i := 0; i < workers; i++ {
		go wait.Until(sq.worker, time.Second, sq.stopCh)
	}
}

// Len returns the unprocessed item length
func (sq *Queue) Len() int {
	return sq.queue.Len()
}

// ShutDown shuts down the work queue and waits for the worker to ACK
func (sq *Queue) ShutDown() {
	close(sq.stopCh)

	// sq shutdown the queue, then worker can't get key from queue
	// processNextWorkItem return false, and then waitGroup -1
	sq.queue.ShutDown()
	sq.waitGroup.Wait()
}

// IsShuttingDown returns if the method Shutdown was invoked
func (sq *Queue) IsShuttingDown() bool {
	return sq.queue.ShuttingDown()
}

// Queue returns the rate limit work queue
func (sq *Queue) Queue() workqueue.RateLimitingInterface {
	return sq.queue
}

// Enqueue wraps queue.Add
func (sq *Queue) Enqueue(obj interface{}) {
	if sq.IsShuttingDown() {
		return
	}
	sq.queue.Add(obj)
}

// EnqueueRateLimited wraps queue.AddRateLimited. It adds an item to the workqueue
// after the rate limiter says its ok
func (sq *Queue) EnqueueRateLimited(obj interface{}) {
	if sq.IsShuttingDown() {
		return
	}
	sq.queue.AddRateLimited(obj)
}

// EnqueueAfter wraps queue.AddAfter. It adds an item to the workqueue after the indicated duration has passed
func (sq *Queue) EnqueueAfter(obj interface{}, after time.Duration) {
	if sq.IsShuttingDown() {
		return
	}
	sq.queue.AddAfter(obj, after)
}

// Worker is a common worker for controllers
// worker runs a work thread that just dequeues items, processes them, and marks them done.
// It enforces that the Handler is never invoked concurrently with the same key.
func (sq *Queue) worker() {
	sq.waitGroup.Add(1)
	defer sq.waitGroup.Done()
	// invoked oncely and process any until exhausted
	for sq.processNextWorkItem() {
	}
}

// ProcessNextWorkItem processes next item in queue by Handler
func (sq *Queue) processNextWorkItem() bool {
	obj, quit := sq.queue.Get()
	if quit {
		return false
	}

	// We call Done here so the workqueue knows we have finished
	// processing this item. We also must remember to call Forget if we
	// do not want this work item being re-queued. For example, we do
	// not call Forget if a transient error occurs, instead the item is
	// put back on the workqueue and attempted again after a back-off
	// period.
	defer sq.queue.Done(obj)

	return sq.handle(obj)
}

func (sq *Queue) handle(obj interface{}) bool {
	result, err := sq.handler(obj)
	if err != nil {
		klog.Warningf("Error handling obj %v retry: %v, err: %v", obj, sq.queue.NumRequeues(obj), err)
		return false
	} else if result.RequeueAfter > 0 {
		sq.queue.Forget(obj)
		sq.queue.AddAfter(obj, result.RequeueAfter)
		return true
	} else if result.Requeue {
		sq.queue.AddRateLimited(obj)
		return true
	}
	// Finally, if no error occurs we Forget this item so it does not
	// get queued again until another change happens.
	sq.queue.Forget(obj)

	klog.V(1).Infof("Successfully handled obj %v", obj)

	return true
}
