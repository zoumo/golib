// Copyright 2023 jim.zoumo@gmail.com
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

package queue

import (
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
)

const (
	ErrRetryForever = -1
	ErrRetryNone    = 0
)

type HandleResult struct {
	// RequeueRateLimited re-enqueue the object after the rate limiter says it's ok. Defaults to false.
	RequeueRateLimited bool

	// RequeueImmediately tells the Queue to requeue the object immediately. Defaults to false.
	RequeueImmediately bool

	// RequeueAfter tells the Queue to re-enqueue the object after the Duration if it is greater than 0.
	RequeueAfter time.Duration

	// MaxRequeueTimes tells the Queue the limit count of requeueing object. Defaults to 1.
	// If you want to requeue forever, please use set it to MaxErrRetryForever
	MaxRequeueTimes int
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
	queue            workqueue.RateLimitingInterface
	queueRateLimiter workqueue.RateLimiter

	waitGroup sync.WaitGroup

	maxErrRetries int

	stopCh chan struct{}
}

// NewQueue returns a new Queue
func NewQueue(handler Handler) *Queue {
	rateLimiter := workqueue.DefaultControllerRateLimiter()
	return &Queue{
		queue:            workqueue.NewRateLimitingQueue(rateLimiter),
		queueRateLimiter: rateLimiter,
		handler:          handler,
		waitGroup:        sync.WaitGroup{},
		stopCh:           make(chan struct{}),
	}
}

// Run starts n workers to sync
func (q *Queue) Run(workers int) {
	for i := 0; i < workers; i++ {
		go wait.Until(q.worker, time.Second, q.stopCh)
	}
}

// SetMaxRetries sets the max retry times of the queue
func (q *Queue) SetMaxErrRetries(max int) *Queue {
	if max >= -1 {
		q.maxErrRetries = max
	}

	return q
}

// Len returns the unprocessed item length
func (q *Queue) Len() int {
	return q.queue.Len()
}

// ShutDown shuts down the work queue and waits for the worker to ACK
func (q *Queue) ShutDown() {
	close(q.stopCh)

	// q shutdown the queue, then worker can't get key from queue
	// processNextWorkItem return false, and then waitGroup -1
	q.queue.ShutDown()
	q.waitGroup.Wait()
}

// IsShuttingDown returns if the method Shutdown was invoked
func (q *Queue) IsShuttingDown() bool {
	return q.queue.ShuttingDown()
}

// Queue returns the rate limit work queue
func (q *Queue) Queue() workqueue.RateLimitingInterface {
	return q.queue
}

// Enqueue wraps queue.Add
func (q *Queue) Enqueue(obj interface{}) {
	if q.IsShuttingDown() {
		return
	}
	q.queue.Add(obj)
}

// EnqueueRateLimited wraps queue.AddRateLimited. It adds an item to the workqueue
// after the rate limiter says its ok
func (q *Queue) EnqueueRateLimited(obj interface{}) {
	if q.IsShuttingDown() {
		return
	}
	q.queue.AddRateLimited(obj)
}

// EnqueueAfter wraps queue.AddAfter. It adds an item to the workqueue after the indicated duration has passed
func (q *Queue) EnqueueAfter(obj interface{}, after time.Duration) {
	if q.IsShuttingDown() {
		return
	}
	q.queue.AddAfter(obj, after)
}

// Worker is a common worker for controllers
// worker runs a work thread that just dequeues items, processes them, and marks them done.
// It enforces that the Handler is never invoked concurrently with the same key.
func (q *Queue) worker() {
	q.waitGroup.Add(1)
	defer q.waitGroup.Done()
	// invoked oncely and process any until exhausted
	for q.processNextWorkItem() {
	}
}

// ProcessNextWorkItem processes next item in queue by Handler
func (q *Queue) processNextWorkItem() bool {
	obj, quit := q.queue.Get()
	if quit {
		return false
	}

	// We call Done here so the workqueue knows we have finished
	// processing this item. We also must remember to call Forget if we
	// do not want this work item being re-queued. For example, we do
	// not call Forget if a transient error occurs, instead the item is
	// put back on the workqueue and attempted again after a back-off
	// period.
	defer q.queue.Done(obj)

	return q.handle(obj)
}

func (q *Queue) handle(obj interface{}) bool {
	result, err := q.handler(obj)
	if err != nil {
		q.handleError(obj, err)
		return false
	}

	q.handleRequeue(obj, result)
	return true
}

func (q *Queue) handleError(obj interface{}, err error) {
	if err == nil {
		return
	}
	if q.maxErrRetries == ErrRetryForever ||
		(q.maxErrRetries != ErrRetryNone && q.queue.NumRequeues(obj) < q.maxErrRetries) {
		q.queue.AddRateLimited(obj)
		return
	}
	q.queue.Forget(obj)
}

func (q *Queue) handleRequeue(obj interface{}, result HandleResult) {
	var requeueAfter time.Duration

	if result.MaxRequeueTimes == 0 {
		// 0 means only requeue this time, fix to 1
		result.MaxRequeueTimes = 1
	}

	// let requeueAfter > 0 means we need requeue it
	if result.RequeueAfter > 0 {
		requeueAfter = result.RequeueAfter
	} else if result.RequeueImmediately {
		requeueAfter = time.Millisecond
	} else if result.RequeueRateLimited {
		requeueAfter = time.Microsecond
	}

	if result.MaxRequeueTimes > 0 && q.queue.NumRequeues(obj) >= result.MaxRequeueTimes {
		// more than maximum requeue times
		// skip requeue
		requeueAfter = 0
	}

	if requeueAfter > 0 {
		if result.RequeueRateLimited {
			q.EnqueueRateLimited(obj)
		} else {
			// EnqueueAfter does not record object requeues times, we need to
			// call rateLimiter.When to add 1 time explicitly.
			q.queueRateLimiter.When(obj)
			q.EnqueueAfter(obj, requeueAfter)
		}
		return
	}
	// we should forget this obj if there is no need to requeue this obj
	q.queue.Forget(obj)
}
