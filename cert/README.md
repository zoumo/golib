# cert

Certificate management utilities for TLS/SSL in Go.

This package provides utilities for generating X.509 certificates, certificate signing requests (CSRs), and working with TLS certificates.

## Features

- Generate self-signed certificates
- Generate self-signed CA certificates
- Generate certificates signed by a CA
- Generate certificate signing requests (CSRs)
- Support for Subject Alternative Names (SANs) including DNS names and IP addresses
- TLS certificate and key management

## Usage

```go
package main

import (
    "crypto/rand"
    "crypto/rsa"
    "log"

    "github.com/example/golib/cert"
)

func main() {
    // Generate a new RSA key
    key, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        log.Fatal(err)
    }

    // Create a self-signed certificate
    cfg := cert.Config{
        CommonName:   "example.com",
        Organization: []string{"My Organization"},
        AltNames: cert.AltNames{
            DNSNames: []string{"example.com", "www.example.com"},
            IPs:      []net.IP{net.ParseIP("127.0.0.1")},
        },
    }

    certPEM, err := cert.NewSelfSignedCert(cfg, key)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Generated certificate for %s", cfg.CommonName)
}
```

## API Reference

### `Config`

Configuration for creating certificates.

| Field | Type | Description |
|-------|------|-------------|
| `CommonName` | `string` | The common name for the certificate |
| `Organization` | `[]string` | Organization name(s) |
| `AltNames` | `AltNames` | Subject Alternative Names |
| `Usages` | `[]x509.ExtKeyUsage` | Extended key usage values |

### `AltNames`

Subject Alternative Names for a certificate.

| Field | Type | Description |
|-------|------|-------------|
| `DNSNames` | `[]string` | DNS names to include in the certificate |
| `IPs` | `[]net.IP` | IP addresses to include in the certificate |

### Functions

- `NewSignedCert(cfg Config, key crypto.Signer, caKey crypto.Signer, caCert *x509.Certificate) (*x509.Certificate, error)` - Creates a new certificate signed by a CA
- `NewSelfSignedCert(cfg Config, key crypto.Signer) (*x509.Certificate, error)` - Creates a new self-signed certificate
- `NewSelfSignedCACert(cfg Config, key crypto.Signer) (*x509.Certificate, error)` - Creates a new self-signed CA certificate
- `NewCSR(cfg Config, key crypto.Signer) (*x509.CertificateRequest, error)` - Creates a new certificate signing request