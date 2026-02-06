# retry

Retry mechanisms with exponential backoff.

This package provides utilities for retrying operations with configurable exponential backoff, based on Kubernetes' wait package.

## Features

- Exponential backoff with configurable parameters
- Continue retrying based on specific error conditions
- Ignore specific errors to stop retrying
- Pre-configured retry strategies

## Usage

```go
package main

import (
    "fmt"
    "time"

    "github.com/example/golib/retry"
)

func main() {
    // Basic retry
    err := retry.Retry(retry.DefaultRetry, func() error {
        // Your operation here
        return nil // or return error to retry
    })
    if err != nil {
        fmt.Printf("Failed after retries: %v\n", err)
    }

    // Retry with error filtering
    err = retry.Retry(retry.DefaultRetry, func() error {
        return fmt.Errorf("temporary error")
    })
}
```

## Pre-configured Strategies

| Constant | Description |
|----------|-------------|
| `DefaultRetry` | Recommended retry for conflicts where multiple clients make changes |
| `DefaultBackoff` | Recommended backoff for conflicts with controllers |

## API Reference

### Types

| Type | Description |
|------|-------------|
| `Backoff` | Parameters for exponential backoff (steps, duration, factor, jitter) |

### Functions

| Function | Description |
|----------|-------------|
| `Retry(backoff Backoff, condition func() error) error` | Executes condition with exponential backoff |
| `RetryContined(backoff Backoff, condition func() error, continued func(error) bool) error` | Keeps retrying if continued returns true |
| `RetryIgnored(backoff Backoff, condition func() error, ignored func(error) bool) error` | Stops retrying if ignored returns true |

### Variables

| Variable | Type | Description |
|----------|------|-------------|
| `ErrWaitTimeout` | `error` | Error returned when condition exits without success |

## Backoff Parameters

| Field | Type | Description |
|-------|------|-------------|
| `Steps` | `int` | Number of times to retry |
| `Duration` | `time.Duration` | Base duration for exponential backoff |
| `Factor` | `float64` | Multiplier for each step |
| `Jitter` | `float64` | Random amount to add (between 0 and duration*(1+jitter)) |