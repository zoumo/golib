# pointer

Pointer utility functions for creating pointers to primitive types.

This package provides convenience functions for creating pointers to primitive Go types, which is especially useful when you need optional pointer parameters or when working with APIs that require pointers.

## Features

- Create pointers from primitive type values
- Supports all Go primitive types

## Usage

```go
package main

import (
    "fmt"

    "github.com/example/golib/pointer"
)

func main() {
    // Create pointers to primitive types
    p := pointer.Bool(true)
    i := pointer.Int(42)
    s := pointer.String("hello")
    f := pointer.Float64(3.14)

    fmt.Println(*p) // true
    fmt.Println(*i) // 42
    fmt.Println(*s) // hello
    fmt.Println(*f) // 3.14
}
```

## API Reference

| Function | Returns | Description |
|----------|---------|-------------|
| `Bool(t bool)` | `*bool` | Creates a pointer to a bool |
| `Int(i int)` | `*int` | Creates a pointer to an int |
| `Int8(i int8)` | `*int8` | Creates a pointer to an int8 |
| `Int16(i int16)` | `*int16` | Creates a pointer to an int16 |
| `Int32(i int32)` | `*int32` | Creates a pointer to an int32 |
| `Int64(i int64)` | `*int64` | Creates a pointer to an int64 |
| `Uint(i uint)` | `*uint` | Creates a pointer to a uint |
| `Uint8(i uint8)` | `*uint8` | Creates a pointer to a uint8 |
| `Uint16(i uint16)` | `*uint16` | Creates a pointer to a uint16 |
| `Uint32(i uint32)` | `*uint32` | Creates a pointer to a uint32 |
| `Uint64(i uint64)` | `*uint64` | Creates a pointer to a uint64 |
| `Uintptr(i uintptr)` | `*uintptr` | Creates a pointer to a uintptr |
| `Float32(f float32)` | `*float32` | Creates a pointer to a float32 |
| `Float64(f float64)` | `*float64` | Creates a pointer to a float64 |
| `String(s string)` | `*string` | Creates a pointer to a string |
| `Complex64(c complex64)` | `*complex64` | Creates a pointer to a complex64 |
| `Complex128(c complex128)` | `*complex128` | Creates a pointer to a complex128 |