# reflection

Reflection utility functions for type introspection in Go.

This package provides helper functions for analyzing types at runtime using Go's reflection package.

## Features

- Check if a type is a literal (primitive) type
- Check if a type is a custom (non-predeclared) type
- Check if a struct has unexported fields
- Check if a field is unexported
- Check if a type can be used as a map key (hashable)
- Check if a type is an anonymous struct

## Usage

```go
package main

import (
    "fmt"

    "github.com/example/golib/reflection"
    "reflect"
)

type MyStruct struct {
    Public  string
    private int
}

func main() {
    t := reflect.TypeOf(MyStruct{})

    // Check if type is literal
    fmt.Printf("Is literal: %v\n", reflection.IsLiteralType(t)) // false

    // Check if type is custom
    fmt.Printf("Is custom: %v\n", reflection.IsCustomType(t)) // true

    // Check if struct has unexported fields
    fmt.Printf("Has unexported: %v\n", reflection.HasUnexportedField(t)) // true

    // Check if type is hashable
    fmt.Printf("Is hashable: %v\n", reflection.Hashable(t)) // false (contains string)
}
```

## API Reference

### Functions

| Function | Description |
|----------|-------------|
| `IsLiteralType(t reflect.Type) bool` | Returns true if the type is a primitive/literal type (bool, string, numeric types) |
| `IsCustomType(t reflect.Type) bool` | Returns true if the type is not predeclared |
| `HasUnexportedField(t reflect.Type) bool` | Returns true if the struct has unexported fields (panics if not a struct) |
| `IsUnexportedField(field reflect.StructField) bool` | Returns true if the field is unexported |
| `Hashable(in reflect.Type) bool` | Returns true if the type can be used as a map key |
| `IsAnonymousStruct(t reflect.Type) bool` | Returns true if the type is an anonymous struct |