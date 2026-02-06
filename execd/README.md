# execd

Daemon process execution and management utilities.

This package provides utilities for running and managing long-running background processes (daemons) with automatic restart and graceful shutdown capabilities.

## Features

- Start daemon processes that run in the background
- Automatic process restart on failure
- Graceful shutdown with configurable grace period
- Custom graceful shutdown handlers
- Process state tracking
- Signal handling

## Usage

### Basic Daemon

```go
package main

import (
    "log"
    "time"

    "github.com/example/golib/execd"
)

func main() {
    daemon := execd.Daemon("my-server", "--config", "config.yaml")

    // Start the daemon and run it forever
    err := daemon.RunForever()
    if err != nil {
        log.Fatal(err)
    }
}
```

### With Grace Period

```go
daemon := execd.Daemon("my-server")
daemon.SetGracePeriod(10 * time.Second)

err := daemon.RunForever()
if err != nil {
    log.Fatal(err)
}

// When stopping, daemon will receive SIGTERM first,
// then SIGKILL after 10 seconds
```

### Custom Graceful Shutdown

```go
daemon := execd.Daemon("my-server")
daemon.SetGracefulShutDown(func(cmd *exec.Cmd) error {
    // Custom shutdown logic
    return cmd.Process.Signal(syscall.SIGWINCH)
})

err := daemon.RunForever()
```

### Stop the Daemon

```go
if err := daemon.Stop(); err != nil {
    log.Printf("Error stopping daemon: %v", err)
}
```

## API Reference

### Functions

| Function | Description |
|----------|-------------|
| `Daemon(name string, arg ...string) *D` | Creates a new D for the named program |
| `DaemonFrom(c *exec.Cmd) *D` | Creates a D from an existing exec.Cmd |

### D Methods

| Method | Description |
|--------|-------------|
| `RunForever() error` | Starts the daemon and keeps it running in background |
| `Stop() error` | Stops the daemon process |
| `IsRunning() bool` | Returns true if the daemon is running |
| `Pid() (int, error)` | Returns the running process PID |
| `Signal(signal os.Signal) error` | Sends a signal to the daemon |
| `Name() string` | Returns the daemon name |
| `SetGracePeriod(d time.Duration)` | Sets graceful shutdown grace period |
| `SetGracefulShutDown(f func(*exec.Cmd) error)` | Sets custom graceful shutdown handler |
| `Command() *exec.Cmd` | Returns the running exec.Cmd |

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `Path` | `string` | Path to the command to run |
| `Args` | `[]string` | Command line arguments |
| `Env` | `[]string` | Environment variables |
| `Dir` | `string` | Working directory |
| `Stdin` | `io.Reader` | Standard input |
| `Stdout` | `io.Writer` | Standard output |
| `Stderr` | `io.Writer` | Standard error |
| `SysProcAttr` | `*syscall.SysProcAttr` | OS-specific attributes |