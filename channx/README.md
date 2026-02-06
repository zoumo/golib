# channx

A self-adaptive channel with a ring buffer that can act as an unbounded channel.

The channel buffer capacity automatically increases according to excessive input and restores to the original when the buffer is empty.

## Features

- Self-adaptive channel that grows when needed
- Ring buffer for efficient storage
- Can function as an unbounded channel
- Configurable initial and maximum buffer sizes
- Configurable input/output channel buffer sizes
- Graceful shutdown with optional data retention

## Usage

```go
package main

import (
    "fmt"

    "github.com/example/golib/chanx"
)

func main() {
    // Create a new ChannX with default options
    ch := chanx.New()

    // Send values to the channel
    ch.In() <- "hello"
    ch.In() <- "world"

    // Receive values from the channel
    fmt.Println(<-ch.Out())
    fmt.Println(<-ch.Out())

    ch.Close()
}
```

### With Options

```go
ch := chanx.New(
    chanx.InChanSize(100),           // Input channel buffer size
    chanx.OutChanSize(100),          // Output channel buffer size
    chanx.InitBufferSize(128),       // Initial ring buffer size
    chanx.MaxBufferSize(1024),       // Maximum ring buffer size (0 = unlimited)
    chanx.DropClosedBufferData(),    // Drop all data when closed
)
```

## API Reference

### `New(opts ...Options) *ChannX`

Creates a new ChannX with the specified options.

### Options Functions

| Function | Description |
|----------|-------------|
| `InChanSize(size int)` | Sets input channel buffer size |
| `OutChanSize(size int)` | Sets output channel buffer size |
| `InitBufferSize(size int)` | Sets initial ring buffer size |
| `MaxBufferSize(size int)` | Sets maximum ring buffer size (0 = unlimited) |
| `DropClosedBufferData()` | Drops data in ring buffer when Close() is called |

### ChannX Methods

| Method | Description |
|--------|-------------|
| `In() chan<- interface{}` | Returns the input channel |
| `Out() <-chan interface{}` | Returns the output channel |
| `Close()` | Closes the channel gracefully |