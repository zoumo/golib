# queue

Rate-limited work queue wrapper for asynchronous processing based on Kubernetes workqueue.

This package provides a convenient wrapper around `k8s.io/client-go/util/workqueue` for processing asynchronous work with rate limiting and retry capabilities.

## Features

- Rate-limited work queue
- Configurable retry behavior (including infinite retry)
- Requeue with immediate, rate-limited, or delayed options
- Graceful shutdown
- Multiple worker support
- Handler-based processing

## Usage

```go
package main

import (
    "fmt"
    "time"

    "github.com/example/golib/queue"
)

func main() {
    // Create a queue with a handler
    q := queue.NewQueue(func(obj interface{}) (queue.HandleResult, error) {
        // Process the item
        fmt.Printf("Processing: %v\n", obj)

        // Return result to control requeuing
        return queue.HandleResult{
            RequeueRateLimited: false,
            RequeueImmediately: false,
        }, nil
    })

    // Set max error retries (optional)
    q.SetMaxErrRetries(queue.ErrRetryForever) // or 1, 2, etc.

    // Start workers
    q.Run(5) // 5 workers

    // Enqueue items
    q.Enqueue("item1")
    q.Enqueue("item2")

    // Enqueue with rate limiting
    q.EnqueueRateLimited("item3")

    // Enqueue after delay
    q.EnqueueAfter("item4", 10*time.Second)

    // Shutdown
    q.ShutDown()
}
```

## API Reference

### Functions

| Function | Description |
|----------|-------------|
| `NewQueue(handler Handler) *Queue` | Creates a new Queue with the given handler |

### Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `ErrRetryForever` | -1 | Retry forever on errors |
| `ErrRetryNone` | 0 | No retry on errors |

### Queue Methods

| Method | Description |
|--------|-------------|
| `Run(workers int)` | Starts the specified number of workers |
| `SetMaxErrRetries(max int) *Queue` | Sets max retry times for errors |
| `Len() int` | Returns the number of unprocessed items |
| `ShutDown()` | Shuts down the work queue |
| `IsShuttingDown() bool` | Returns true if ShutDown was invoked |
| `Queue() workqueue.RateLimitingInterface` | Returns the underlying workqueue |
| `Enqueue(obj interface{})` | Adds an item to the queue |
| `EnqueueRateLimited(obj interface{})` | Adds an item with rate limiting |
| `EnqueueAfter(obj interface{}, after time.Duration)` | Adds an item after a delay |

### Types

#### `Handler`

Handler function type for processing queue items.

```go
type Handler func(obj interface{}) (HandleResult, error)
```

#### `HandleResult`

Result type for controlling requeuing behavior.

| Field | Type | Description |
|-------|------|-------------|
| `RequeueRateLimited` | `bool` | Re-enqueue after rate limiter says it's ok |
| `RequeueImmediately` | `bool` | Re-enqueue immediately |
| `RequeueAfter` | `time.Duration` | Re-enqueue after specified duration |
| `MaxRequeueTimes` | `int` | Limit count of requeuing (default: 1) |