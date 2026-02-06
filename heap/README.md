# heap

Priority queue implementation with key-based access and update capabilities.

This package provides a thread-safe heap data structure that implements the standard `container/heap` interface with additional features for key-based access and updates.

## Features

- Key-based object storage and retrieval
- Add, update, or add-if-not-present operations
- Peek at top elements without removing them
- Get objects by key
- Iterate over all objects
- Compatible with standard Go heap interface

## Usage

```go
package main

import (
    "fmt"

    "github.com/example/golib/heap"
)

type Item struct {
    Name  string
    Value int
}

func main() {
    // Create a new heap
    h := heap.New(
        func(obj interface{}) (string, error) {
            return obj.(*Item).Name, nil
        },
        func(x, y interface{}) bool {
            // Priority based on Value (ascending)
            return x.(*Item).Value < y.(*Item).Value
        },
    )

    // Add items
    h.AddOrUpdate(&Item{Name: "a", Value: 5})
    h.AddOrUpdate(&Item{Name: "b", Value: 3})
    h.AddOrUpdate(&Item{Name: "c", Value: 7})

    // Peek at the top element (lowest value)
    top := h.Peek().(*Item)
    fmt.Printf("Top: %s with value %d\n", top.Name, top.Value)

    // Pop the top element
    popped := h.Pop().(*Item)
    fmt.Printf("Popped: %s\n", popped.Name)

    // Update an existing item
    h.UpdateIfPresent(&Item{Name: "c", Value: 1})

    // Get item by key
    if item, ok := h.GetByKey("a"); ok {
        fmt.Printf("Found: %+v\n", item)
    }

    // Remove an item
    h.Remove(&Item{Name: "b"})
}
```

## API Reference

### Functions

| Function | Description |
|----------|-------------|
| `New(keyFunc KeyFunc, lessFunc LessFunc) *Heap` | Creates a new Heap |

### Heap Methods

| Method | Description |
|--------|-------------|
| `AddOrUpdate(obj interface{}) error` | Adds or updates an item |
| `AddIfNotPresent(obj interface{}) error` | Adds item only if key doesn't exist |
| `UpdateIfPresent(obj interface{}) error` | Updates item only if key exists |
| `Remove(obj interface{}) error` | Removes an item |
| `Pop() interface{}` | Removes and returns the top element |
| `Peek() interface{}` | Returns the top element without removing |
| `PeekSecond() interface{}` | Returns the second element without removing |
| `GetByKey(key string) (interface{}, bool)` | Retrieves item by key |
| `Len() int` | Returns the number of items |
| `List() []interface{}` | Returns all items |
| `Range(f func(i int, key string, obj interface{}) bool)` | Iterates over items |

### Types

| Type | Description |
|------|-------------|
| `KeyFunc` | Function to extract a key from an object |
| `LessFunc` | Function to compare two objects for ordering |
| `KeyError` | Error type for key extraction failures |