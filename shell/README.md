# shell

Shell command utilities for executing Bash commands.

This package provides convenience functions for executing shell commands using Bash as the interpreter.

## Features

- Execute shell commands via `/bin/bash`
- Command argument escaping for safe shell usage
- Context support for timeout/cancellation

## Usage

```go
package main

import (
    "fmt"

    "github.com/example/golib/shell"
)

func main() {
    // Execute a simple shell command
    cmd := shell.Command("echo", "hello world")
    output, err := cmd.Output()
    if err != nil {
        panic(err)
    }
    fmt.Println(string(output))

    // Execute with piped command
    cmd = shell.Command("echo", "1 2 3").Pipe("awk", "{print $2}")
    output, _ = cmd.Output()
    fmt.Println(string(output))

    // Escape arguments for shell safety
    unsafe := "arg with 'quotes' and spaces"
    safe := shell.QueryEscape(unsafe)
    fmt.Println(safe) // 'arg with '\''quotes'\'' and spaces'
}
```

## API Reference

### Functions

| Function | Description |
|----------|-------------|
| `Command(name string, args ...string) *exec.Cmd` | Creates a new command to be executed by shell (uses `/bin/bash -c`) |
| `CommandContext(ctx context.Context, name string, args ...string) *exec.Cmd` | Creates a new command with context support |
| `QueryEscape(arg string) string` | Escapes a string for safe use in shell commands |