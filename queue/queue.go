package queue

import (
	"sync"
	"time"

	"github.com/golang/glog"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
)

const (
	maxRetries = 3
)

// Queue is a wrapper of kubernetes workqueue to do asynchronous work easily.
// It requires a Handler and an optional key function.
// After starting the Queue, you can call the Enqueque function to enqueue items.
// Queue will get key from the items by keyFunc, and add the key to the rate limit workqueue.
// The worker will be invoked to call the Handler.
type Queue struct {
	// queue is the work queue the worker polls
	queue workqueue.RateLimitingInterface
	// Handler is called for each item in the queue
	Handler func(obj interface{}) error

	waitGroup sync.WaitGroup

	maxRetries int
	stopCh     chan struct{}
}

// NewQueue returns a new Queue
func NewQueue(handler func(obj interface{}) error) *Queue {
	sq := &Queue{
		queue:      workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		Handler:    handler,
		waitGroup:  sync.WaitGroup{},
		maxRetries: maxRetries,
		stopCh:     make(chan struct{}),
	}

	return sq
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

// SetMaxRetries sets the max retry times of the queue
func (sq *Queue) SetMaxRetries(max int) {
	if max > 0 {
		sq.maxRetries = max
	}
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
	defer sq.queue.Done(obj)

	err := sq.Handler(obj)
	sq.handleSyncError(err, obj)

	return true
}

// HandleSyncError handles error when sync obj error and retry n times
func (sq *Queue) handleSyncError(err error, obj interface{}) {
	if err == nil {
		// no err
		sq.queue.Forget(obj)
		return
	}

	if sq.queue.NumRequeues(obj) < sq.maxRetries {
		glog.Warningf("Error handling obj %v retry: %v, err: %v", obj, sq.queue.NumRequeues(obj), err)
		sq.queue.AddRateLimited(obj)
		return
	}

	utilruntime.HandleError(err)
	glog.Warningf("Dropping object %v from the queue, err: %v", obj, err)
	sq.queue.Forget(obj)
}
