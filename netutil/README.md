# netutil

Network interface utilities for inspecting network configuration.

## Features

- List and filter network interfaces
- Get interfaces by name or IP address
- Get loopback interfaces
- Convenient methods for working with addresses and interfaces

## Usage

```go
package main

import (
    "fmt"

    "github.com/example/golib/netutil"
)

func main() {
    // Get all interfaces
    ifaces, err := netutil.Interfaces()
    if err != nil {
        panic(err)
    }

    for _, iface := range ifaces {
        fmt.Printf("Interface: %s\n", iface.Name)
        addrs, _ := iface.Addrs()
        for _, addr := range addrs {
            fmt.Printf("  IP: %s (IPv4: %v, IPv6: %v)\n",
                addr.IPAddr(), addr.IsIPv4(), addr.IsIPv6())
        }
    }

    // Get interface by name
    eth0, err := netutil.InterfaceByName("eth0")
    if err != nil {
        fmt.Printf("eth0 not found: %v\n", err)
    } else {
        fmt.Printf("Found eth0: %+v\n", eth0)
    }

    // Get interfaces by IP
    ifaces, err := netutil.InterfacesByIP("192.168.1.100")
    if err != nil {
        fmt.Printf("No interface with this IP: %v\n", err)
    }

    // Get loopback interfaces
    loopback, err := netutil.InterfacesByLoopback()
    if err != nil {
        fmt.Printf("No loopback interface: %v\n", err)
    }
}
```

## API Reference

### Functions

| Function | Description |
|----------|-------------|
| `Interfaces() (InterfaceSlice, error)` | Returns all local network interfaces |
| `InterfacesByLoopback() (InterfaceSlice, error)` | Returns loopback interfaces |
| `InterfacesByIP(ip string) (InterfaceSlice, error)` | Returns interfaces using the specified IP |
| `InterfaceByName(name string) (*Interface, error)` | Returns the interface with the given name |

### Types

#### `Interface`

Wraps `net.Interface` with additional methods.

| Method | Description |
|--------|-------------|
| `Addrs() (AddrSlice, error)` | Returns unicast interface addresses |
| `IsLoopback() bool` | Returns true if it's a loopback interface |

#### `Addr`

Represents a network endpoint address.

| Method | Description |
|--------|-------------|
| `IPAddr() string` | Returns the IP address string |
| `MaskSize() int` | Returns the number of leading ones in mask |
| `IsIPv4() bool` | Returns true if it's an IPv4 address |
| `IsIPv6() bool` | Returns true if it's an IPv6 address |
| `IsLoopback() bool` | Returns true if it's a loopback address |

#### `InterfaceSlice`

Slice of `Interface` with helper methods.

| Method | Description |
|--------|-------------|
| `Get(name string) *Interface` | Returns interface by name |
| `Contains(name string) bool` | Checks if interface exists |
| `Filter(filterFunc func(iface Interface) bool) InterfaceSlice` | Filters interfaces |
| `One() *Interface` | Returns first interface |

#### `AddrSlice`

Slice of `Addr` with helper methods.

| Method | Description |
|--------|-------------|
| `Contains(ip string) bool` | Checks if IP exists in the slice |