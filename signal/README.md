# signal

Signal handling utilities for Go programs.

This package provides a helper struct for handling OS signals with configurable exit behavior.

## Features

- Easy signal handling
- Customizable exit behavior on signal receipt
- Helpers for common signals (SIGINT, SIGTERM)

## Usage

```go
package main

import (
    "fmt"
    "os"
    "syscall"

    "github.com/example/golib/signal"
)

func main() {
    // Handle SIGINT (Ctrl+C)
    signal.HandleSigint(func(sig os.Signal) int {
        fmt.Printf("Received signal: %v\n", sig)
        // Cleanup code here
        return 0  // Exit code
    })

    // Handle SIGTERM
    signal.HandleSigterm(func(sig os.Signal) int {
        fmt.Printf("Gracefully shutting down...\n")
        // Cleanup code here
        return 0  // Exit code
    })

    // Custom signal handling - don't exit on signal
    s := signal.New(false, syscall.SIGUSR1)
    s.Handle(func(sig os.Signal) int {
        fmt.Printf("Custom signal received: %v\n", sig)
        // Don't exit, just log
        return 0
    })

    // Your program continues...
    select {} // Wait forever
}
```

## API Reference

### Functions

| Function | Description |
|----------|-------------|
| `New(exit bool, sig ...os.Signal) *Signal` | Creates a new Signal that listens for specified signals |
| `HandleSigterm(handler func(os.Signal) int)` | Helper to handle SIGTERM with exit enabled |
| `HandleSigint(handler func(os.Signal) int)` | Helper to handle SIGINT with exit enabled |

### Signal Methods

| Method | Description |
|--------|-------------|
| `Handle(handler func(os.Signal) int)` | Receives signals and lets the handler process them |

### Signal Struct

| Field | Type | Description |
|-------|------|-------------|
| `exit` | `bool` | If true, the program will exit with the handler's return code |