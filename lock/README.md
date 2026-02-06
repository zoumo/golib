# lock

Concurrency control utilities including token buckets for limiting concurrent operations.

## sub-packages

### maxinflight

Token bucket implementation for limiting concurrent executions (max in-flight concurrency).

#### Features

- Multiple implementation strategies (atomic, channel, mutex)
- Resize capacity at runtime
- Infinity token bucket for no limits
- Try-acquire pattern

#### Usage

```go
package main

import (
    "fmt"

    "github.com/example/golib/lock/maxinflight"
)

func main() {
    // Create a token bucket with capacity of 10
    bucket := maxinflight.New(10)

    // Try to acquire a token
    if bucket.TryAcquire() {
        defer bucket.Release()
        // Do work here
        fmt.Println("Got token, doing work")
    } else {
        fmt.Println("Capacity reached")
    }

    // Resize the bucket
    bucket.Resize(20)
}

// Using different implementations
atomicBucket := maxinflight.New(10)    // Default: atomic
channelBucket := maxinflight.New(10)   // Could switch to channel
mutexBucket := maxinflight.New(10)     // Could switch to mutex
infinity := maxinflight.NewInfinity()  // No limits, always returns true
```

#### API Reference (maxinflight)

| Function | Description |
|----------|-------------|
| `New(size uint32) TokenBucket` | Creates a new token bucket (atomic implementation) |
| `NewInfinity() TokenBucket` | Creates a token bucket with unlimited capacity |

#### TokenBucket Interface

| Method | Description |
|--------|-------------|
| `TryAcquire() bool` | Tries to acquire a token, returns true if successful |
| `Release()` | Releases a token back to the bucket |
| `Resize(n uint32)` | Changes the bucket's capacity |

#### TokenBucketType

| Type | Implementation |
|------|----------------|
| `Atomic` | Atomic operations using sync/atomic |
| `Channel` | Buffered channel-based |
| `Mutex` | Mutex-based locking |
| `Infinity` | Unlimited capacity |