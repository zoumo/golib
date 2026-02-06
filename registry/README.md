# registry

Thread-safe registry for binding names to interfaces.

This package provides a thread-safe registry for registering and retrieving objects by name. It's useful for implementing plugin systems, factory patterns, or any scenario where you need to dynamically register and look up implementations.

## Features

- Thread-safe operations using `sync.Map`
- Optional override protection (panic or allow overwrites)
- Range iteration over all registered items
- Get all keys or values at once

## Usage

### Basic Example

```go
package main

import (
    "fmt"

    "github.com/example/golib/registry"
)

var (
    clouds = registry.New(nil)
)

type Cloud interface {
    Name() string
}

type Config struct {
    // Config fields
}

type awsCloud struct{}

func (a *awsCloud) Name() string { return "AWS" }

type CloudFactory func(Config) (Cloud, error)

func RegisterCloud(name string, factory CloudFactory) error {
    return clouds.Register(name, factory)
}

func GetCloud(name string, config Config) (Cloud, error) {
    v, found := clouds.Get(name)
    if !found {
        return nil, fmt.Errorf("cloud %s not found", name)
    }
    factory := v.(CloudFactory)
    return factory(config)
}

func main() {
    // Register cloud providers
    RegisterCloud("aws", func(cfg Config) (Cloud, error) {
        return &awsCloud{}, nil
    })
    RegisterCloud("gce", func(cfg Config) (Cloud, error) {
        return &gceCloud{}, nil
    })

    // Get and use a cloud provider
    cloud, err := GetCloud("aws", Config{})
    if err != nil {
        panic(err)
    }
    fmt.Println(cloud.Name())
}
```

### With Override Allowed

```go
// Allow overwriting existing registrations
cfg := &registry.Config{
    OverrideAllowed: true,
}
reg := registry.New(cfg)

// This will overwrite the previous "aws" registration
reg.Register("aws", newFactory)
```

## API Reference

### Functions

| Function | Description |
|----------|-------------|
| `New(config *Config) Registry` | Creates a new Registry with the specified config (nil uses defaults) |

### Registry Interface

| Method | Description |
|--------|-------------|
| `Register(name string, v interface{}) error` | Registers an interface by name. Returns error if already registered and override not allowed |
| `Get(name string) (interface{}, bool)` | Retrieves a registered interface by name |
| `Range(f func(key string, value interface{}) bool)` | Iterates over all entries. Iteration stops if f returns false |
| `Keys() []string` | Returns all registered keys |
| `Values() []interface{}` | Returns all registered values |

### Types

| Type | Description |
|------|-------------|
| `Registry` | Interface for registry operations |
| `Config` | Configuration for new registries |

### Config Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `OverrideAllowed` | `bool` | `false` | If true, allows overwriting existing registrations; otherwise returns error |