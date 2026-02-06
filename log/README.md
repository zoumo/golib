# log

Logging interface with logr adapter support.

This package provides a deferred logging interface that delegates to a concrete logging implementation via the `github.com/go-logr/logr` interface.

## Features

- Deferred logging through a placeholder interface
- logr adapter support
- Configurable concrete logger via `SetLogger` or `SetLogrLogger`
- Global logger instance for convenience

## Usage

```go
package main

import (
    "github.com/go-logr/logr"
    "github.com/go-logr/zerologr"
    "github.com/rs/zerolog"

    "github.com/example/golib/log"
)

func main() {
    // Set up a concrete logger (e.g., zerolog)
    logger := zerolog.New(os.Stdout)
    logrLogger := zerologr.New(&logger)

    // Set the logger for all deferred Loggers
    log.SetLogrLogger(logrLogger)

    // Now you can use the global log.Log
    log.Log.Info("This is an info message")
    log.Log.Error(nil, "This is an error message")

    // Or set a custom logger
    log.SetLogger(customLogger) // customLogger must implement the Logger interface
}
```

## API Reference

### Functions

| Function | Description |
|----------|-------------|
| `SetLogrLogger(l logr.Logger)` | Sets a logr.Logger as the concrete logging implementation |
| `SetLogger(l Logger)` | Sets a custom Logger as the concrete implementation |

### Variables

| Variable | Type | Description |
|----------|------|-------------|
| `Log` | `Logger` | The global logger, delegates to the concrete implementation |

### Interfaces

#### `Logger`

The logging interface that concrete implementations must implement. Compatible with logr.Logger v0.4.0 and v1.0.0+.