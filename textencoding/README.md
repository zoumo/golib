# textencoding

Text encoding transformation utilities supporting multiple character encodings.

This package provides functions for transforming text between different character encodings, supporting a wide range of international encodings.

## Features

- Support for Chinese encodings (GB2312, GB18030, Big5)
- Support for Japanese encodings (Shift_JIS, EUC-JP)
- Support for Korean encodings (EUC-KR)
- Support for Western encodings (ISO-8859 series, Windows-1252)
- Support for Unicode encodings (UTF-8, UTF-16, UTF-32)
- Encode, decode, and transform between encodings

## Usage

```go
package main

import (
    "fmt"

    "github.com/example/golib/textencoding"
)

func main() {
    // Decode bytes from GBK to UTF-8
    gbkBytes := []byte{0xC4, 0xE3, 0xBA, 0xC3} // "你好" in GBK
    utf8Bytes, err := textencoding.Decode(gbkBytes, "GBK")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(utf8Bytes)) // 你好

    // Encode UTF-8 to GBK
    original := "你好"
    gbkBytes, err = textencoding.Encode([]byte(original), "GBK")

    // Transform between encodings
    result, err := textencoding.TransformString("Hello", "UTF-8", "SHIFT_JIS")

    // Check if encoding is supported
    if textencoding.IsEncodingSupported("BIG5") {
        fmt.Println("BIG5 encoding is supported")
    }
}
```

## Supported Encodings

The following encodings are supported (case-insensitive):

- **Unicode**: UTF-8, UTF-16, UTF-32
- **Chinese (Simplified)**: GB18030, GBK, HZGB2312
- **Chinese (Traditional)**: Big5
- **Japanese**: Shift_JIS, EUC-JP, ISO-2022-JP
- **Korean**: EUC-KR
- **Western**: ISO-8859-1 through ISO-8859-16, Windows-1252
- **More**: Various other single-byte and multi-byte encodings

## API Reference

### Functions

| Function | Description |
|----------|-------------|
| `IsEncodingSupported(name string) bool` | Checks if the encoding is supported |
| `Decode(s []byte, from string) ([]byte, error)` | Decodes bytes to UTF-8 |
| `Encode(s []byte, to string) ([]byte, error)` | Encodes UTF-8 bytes to target encoding |
| `TransformString(s string, from, to string) (string, error)` | Transforms a string between encodings |
| `Transform(s []byte, from, to string) ([]byte, error)` | Transforms bytes between encodings |