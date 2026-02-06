# ssh

SSH client utilities with shell and SCP support.

This package provides a convenient SSH client with methods for executing shell commands and transferring files via SCP.

## Features

- SSH connection with custom configuration
- Remote shell execution
- SCP file upload and download
- Configurable host key callback, timeout, and user

## sub-packages

- `shell` - Remote shell execution
- `scp` - SCP file transfer
- `krb5` - Kerberos authentication (GSSAPI)

## Usage

```go
package main

import (
    "golang.org/x/crypto/ssh"

    "github.com/example/golib/ssh"
)

func main() {
    // Configure SSH client
    config := &ssh.ClientConfig{
        User: "ubuntu",
        Auth: []ssh.AuthMethod{
            ssh.Password("password"),
        },
    }

    // Create SSH client
    client, err := ssh.DialTCP("myserver", "22", config)
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // Execute remote shell command
    err = client.Shell(strings.NewReader("echo hello"), os.Stdout, os.Stderr)

    // Upload file via SCP
    err = client.Upload(context.Background(), "local.txt", "/remote/path.txt")

    // Download file via SCP
    err = client.Download(context.Background(), "/remote/path.txt", "local.txt")
}
```

## API Reference

### Functions

| Function | Description |
|----------|-------------|
| `DialTCP(host, port string, cfg *ssh.ClientConfig) (*Client, error)` | Creates an SSH connection via TCP |
| `Dial(network, host, port string, cfg *ssh.ClientConfig) (*Client, error)` | Creates an SSH connection with custom network type |
| `SetClientConfigDefaults(cfg *ssh.ClientConfig)` | Sets defaults for SSH client config (user, timeout, host key callback) |

### Client Methods

| Method | Description |
|--------|-------------|
| `Dial(network, host, port string, cfg *ssh.ClientConfig) (*Client, error)` | Creates a new SSH connection from an existing client |
| `Shell(stdin io.Reader, stdout, stderr io.Writer) error` | Executes a remote shell command |
| `Upload(ctx context.Context, local, remote string) error` | Uploads a file via SCP |
| `Download(ctx context.Context, remote, local string) error` | Downloads a file via SCP |

### Client Config Defaults

When using `Dial` or `DialTCP`, the following defaults are applied:

| Setting | Default Value |
|---------|---------------|
| `User` | `"root"` |
| `Timeout` | `5 * time.Second` |
| `HostKeyCallback` | `ssh.InsecureIgnoreHostKey()` |