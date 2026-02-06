# exec

Utilities for executing external programs and processes with enhanced functionality.

This package provides a more powerful and flexible alternative to `os/exec` with support for command piping, startup probes, and output handling.

## Features

- Command execution with context support
- Pipeline support for chaining commands
- Startup probe with configurable configuration
- Direct output reading
- Functional closures for repeated execution
- Command mutation capabilities

## Usage

### Basic Execution

```go
package main

import (
    "fmt"

    "github.com/example/golib/exec"
)

func main() {
    // Run a simple command
    cmd := exec.Command("echo", "hello world")
    output, err := cmd.Output()
    if err != nil {
        panic(err)
    }
    fmt.Println(string(output))
}
```

### Pipeline

```go
// Chain commands together like shell piping
cmd := exec.Command("echo", "1\n2\n3").Pipe("sort")
output, err := cmd.Output()
if err != nil {
    panic(err)
}
fmt.Println(string(output))
```

### With Context

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

cmd := exec.CommandContext(ctx, "sleep", "10")
err := cmd.Run()
// Will timeout after 5 seconds
```

### Startup Probe

```go
cmd := exec.Command("long-running-server")

// Run with startup probe configuration
probe := &exec.Probe{
    PeriodSeconds:      1,
    SuccessThreshold:   2,
    FailureThreshold:   3,
    Handler:            exec.IsCmdRunningHandler,
}
err := cmd.RunForever(probe)
```

### Output Closure

```go
// Create a closure for repeated execution with different args
echo := exec.Command("echo").OutputClosure()
output1, _ := echo("first")
output2, _ := echo("second")
```

## API Reference

### Functions

| Function | Description |
|----------|-------------|
| `Command(name string, args ...string) *Cmd` | Creates a new Cmd for executing the named binary |
| `CommandContext(ctx context.Context, name string, args ...string) *Cmd` | Creates a new Cmd with context support |

### Cmd Methods

| Method | Description |
|--------|-------------|
| `Pipe(name string, args ...string) *Cmd` | Chains this command's output to a new command's input |
| `SetIO(in io.Reader, out, err io.Writer)` | Sets stdin, stdout, and stderr |
| `Run() error` | Starts the command and waits for it to complete |
| `RunForever(startup *Probe) error` | Runs command with startup probe |
| `Start() error` | Starts the command without waiting |
| `Wait() error` | Waits for the command to exit |
| `Output() ([]byte, error)` | Runs the command and returns its stdout |
| `CombinedOutput() ([]byte, error)` | Returns combined stdout and stderr |
| `ReadStdout() ([]byte, error)` | Reads stdout after command finishes |
| `ReadStderr() ([]byte, error)` | Reads stderr after command finishes |
| `OutputClosure() func(...string) ([]byte, error)` | Returns a closure for repeated execution |
| `SetCmdMutator(f func(string, []string) (string, []string))` | Sets a mutator for command args |