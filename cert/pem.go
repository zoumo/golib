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
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"
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

// WriteFile writes the PEM encoding to a file
func (p *PEM) WriteFile(f string) error {
	return ioutil.WriteFile(f, p.Raw, 0664)
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

// PEMBlockForCert returns a pemBlock for x509 certificate
func PEMBlockForCert(derBytes []byte) *PEM {
	return NewPEM(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
}

// PEMBlockForCertRequest returns a pemBlock for certificate request
func PEMBlockForCertRequest(csrBytes []byte) *PEM {
	return NewPEM(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
}
