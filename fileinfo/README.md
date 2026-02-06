# fileinfo

File information utilities providing a custom `os.FileInfo` implementation.

This package provides an easy way to create `os.FileInfo` objects for testing or when you need to represent file metadata without an actual file.

## Features

- Create `os.FileInfo` objects programmatically
- Specify name, size, mode, modification time, and directory status
- Compatible with standard Go file operations

## Usage

```go
package main

import (
    "time"

    "github.com/example/golib/fileinfo"
    "os"
)

func main() {
    // Create a file info for a regular file
    fileInfo := fileinfo.NewInfo(
        "example.txt",
        1024,
        os.FileMode(0644),
        time.Now(),
        false, // not a directory
    )

    println(fileInfo.Name())
    println(fileInfo.Size())
}

// Create a file info for a directory
dirInfo := fileinfo.NewInfo(
    "mydir",
    0,
    os.FileMode(0755|os.ModeDir),
    time.Now(),
    true, // is a directory
)
```

## API Reference

### Functions

| Function | Description |
|----------|-------------|
| `NewInfo(name string, size int64, mode os.FileMode, mtime time.Time, isDir bool) os.FileInfo` | Creates a new file info |

### Info Methods

| Method | Description |
|--------|-------------|
| `Name() string` | Returns the file name |
| `Size() int64` | Returns the file size |
| `Mode() os.FileMode` | Returns the file mode |
| `ModTime() time.Time` | Returns the modification time |
| `IsDir() bool` | Returns true if it's a directory |
| `Sys() interface{}` | Returns nil (placeholder for system-specific data) |