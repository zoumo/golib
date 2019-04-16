/*
Copyright 2018 Jim Zhang (jim.zoumo@gmail.com). All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cert

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
)

// PEM contains the raw bytes and a block of pem
type PEM struct {
	Raw   []byte
	Block *pem.Block
}

// Encode writes the PEM encoding of block to out.
func (p *PEM) Encode(out io.Writer) error {
	return pem.Encode(out, p.Block)
}

// WriteFile writes the PEM encoding to a file
func (p *PEM) WriteFile(f string) error {
	return ioutil.WriteFile(f, p.Raw, 0644)
}

// NewPEM creates a new PEM struct from pem.Block
func NewPEM(b *pem.Block) *PEM {
	return &PEM{
		Raw:   pem.EncodeToMemory(b),
		Block: b,
	}
}

// NewPEMForBytes creates a new PEM struct from raw bytes
func NewPEMForBytes(raw []byte) *PEM {
	derBlock, _ := pem.Decode(raw)
	return &PEM{
		Raw:   raw,
		Block: derBlock,
	}
}

// NewPEMForPrivateKey returns a pemBlock for crypto private key
// It returns an error if the key is not *rsa.PrivateKey or *ecdsa.PrivateKey
func NewPEMForPrivateKey(key crypto.PrivateKey) (*PEM, error) {
	switch pkey := key.(type) {
	case *rsa.PrivateKey:
		return NewPEMForRSAKey(pkey), nil
	case *ecdsa.PrivateKey:
		return NewPEMForECDSAKey(pkey), nil
	default:
		return nil, errors.New("tls: found unknown private key type in PKCS#8 wrapping")
	}
}

// NewPEMForRSAKey returns a pemBlock for ras private key
func NewPEMForRSAKey(key *rsa.PrivateKey) *PEM {
	return NewPEM(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
}

// NewPEMForECDSAKey returns a pemBlock for ecdsa private key
func NewPEMForECDSAKey(key *ecdsa.PrivateKey) *PEM {
	bytes, _ := x509.MarshalECPrivateKey(key)
	return NewPEM(&pem.Block{Type: "EC PRIVATE KEY", Bytes: bytes})
}

// NewPEMForCertificate returns a pemBlock for x509 certificate
func NewPEMForCertificate(derBytes []byte) *PEM {
	return NewPEM(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
}

// NewPEMForCertificateRequest returns a pemBlock for certificate request
func NewPEMForCertificateRequest(derBytes []byte) *PEM {
	return NewPEM(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: derBytes})
}
