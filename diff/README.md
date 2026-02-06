# diff

Text diff utilities for comparing strings and generating unified diff output.

This package provides an interface for computing diffs and generating unified diff format output with optional color highlighting.

## Features

- Generate unified diff format output
- Optional color highlighting for terminal output
- Compare two text strings
- Support for custom names in diff output

## Usage

```go
package main

import (
    "fmt"

    "github.com/example/golib/diff"
)

func main() {
    text1 := "Hello World\nThis is original text"
    text2 := "Hello Go\nThis is modified text\nNew line added"

    // Create a new diff instance with colored output
    d := diff.New(diff.WithColored(true))

    // Generate unified diff
    diffOutput := d.DiffUnified("original.txt", "modified.txt", text1, text2)
    fmt.Println(diffOutput)
}
```

## API Reference

### `New(options ...DiffOption) Diff`

Creates a new Diff instance with the specified options.

### Options

| Function | Description |
|----------|-------------|
| `WithColored(bool)` | Enable or disable colored output for terminals |

### Diff Methods

| Method | Description |
|--------|-------------|
| `DiffUnified(name1, name2, text1, text2 string) string` | Generates a unified diff between two texts |

## Output Format

The output uses the standard unified diff format:

- Lines starting with `-` are removed (shown in red when colored)
- Lines starting with `+` are added (shown in green when colored)
- Lines without a prefix are unchanged context