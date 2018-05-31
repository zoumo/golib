package cert

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
)

// PEM contains the raw bytes and a block of pem
type PEM struct {
	Raw   []byte
	Block *pem.Block
}

// NewPEM creates a new PEM struct from pem.Block
func NewPEM(b *pem.Block) *PEM {
	return &PEM{
		Raw:   pem.EncodeToMemory(b),
		Block: b,
	}
}

// NewPEMFromBytes creates a new PEM struct from ray bytes
func NewPEMFromBytes(raw []byte) *PEM {
	b, _ := pem.Decode(raw)
	return &PEM{
		Raw:   raw,
		Block: b,
	}
}

// Encode writes the PEM encoding of block to out.
func (p *PEM) Encode(out io.Writer) error {
	return pem.Encode(out, p.Block)
}

// PEMBlockForKey returns a pemBlock for ras private key
func PEMBlockForKey(key *rsa.PrivateKey) *PEM {
	return NewPEM(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
}

// PEMBlockForECDSAKey returns a pemBlock for ecdsa private key
func PEMBlockForECDSAKey(key *ecdsa.PrivateKey) *PEM {
	bytes, _ := x509.MarshalECPrivateKey(key)
	return NewPEM(&pem.Block{Type: "EC PRIVATE KEY", Bytes: bytes})
}

// PEMBlockForCert returns  a pemBlock for x509 certificate
func PEMBlockForCert(derBytes []byte) *PEM {
	return NewPEM(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
}
