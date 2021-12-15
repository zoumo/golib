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
	"strings"
)

const (
	// CertificatePEMBlockType is a possible value for pem.Block.Type.
	CertificatePEMBlockType = "CERTIFICATE"
	// CertificateRequestPEMBlockType is a possible value for pem.Block.Type.
	CertificateRequestPEMBlockType = "CERTIFICATE REQUEST"
	// RASPrivateKeyPEMBlockType is a possible value for pem.Block.Type.
	RASPrivateKeyPEMBlockType = "RSA PRIVATE KEY"
	// ECDSAPrivateKeyPEMBlockType is a possible value for pem.Block.Type.
	ECDSAPrivateKeyPEMBlockType = "EC PRIVATE KEY"
)

// PEMBlock contains the raw bytes and a block of pem
type PEMBlock struct {
	Block *pem.Block
}

// Encode writes the PEM encoding of block to out.
func (p *PEMBlock) Encode(out io.Writer) error {
	return pem.Encode(out, p.Block)
}

// EncodeToMemory returns the PEM encoding bytes of p.
func (p *PEMBlock) EncodeToMemory() []byte {
	return pem.EncodeToMemory(p.Block)
}

// WriteFile writes the PEM encoding to a file
func (p *PEMBlock) WriteFile(f string) error {
	return ioutil.WriteFile(f, p.EncodeToMemory(), 0644)
}

// NewPEM creates a new PEM struct from pem.Block
func NewPEM(b *pem.Block) *PEMBlock {
	return &PEMBlock{
		Block: b,
	}
}

// NewPEMForPrivateKey returns a pemBlock for crypto private key
// It returns an error if the key is not *rsa.PrivateKey or *ecdsa.PrivateKey
func NewPEMForPrivateKey(key crypto.Signer) (*PEMBlock, error) {
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
func NewPEMForRSAKey(key *rsa.PrivateKey) *PEMBlock {
	return NewPEM(&pem.Block{Type: RASPrivateKeyPEMBlockType, Bytes: x509.MarshalPKCS1PrivateKey(key)})
}

// NewPEMForECDSAKey returns a pemBlock for ecdsa private key
func NewPEMForECDSAKey(key *ecdsa.PrivateKey) *PEMBlock {
	bytes, _ := x509.MarshalECPrivateKey(key)
	return NewPEM(&pem.Block{Type: ECDSAPrivateKeyPEMBlockType, Bytes: bytes})
}

// NewPEMForCert returns a pemBlock for x509 certificate
func NewPEMForCert(crt *x509.Certificate) *PEMBlock {
	if crt == nil {
		return nil
	}
	return NewPEMForCertDER(crt.Raw)
}

// NewPEMForCertificate returns a pemBlock for x509 certificate
func NewPEMForCertDER(derBytes []byte) *PEMBlock {
	return NewPEM(&pem.Block{Type: CertificatePEMBlockType, Bytes: derBytes})
}

// NewPEMForCSR returns a pemBlock for certificate request
func NewPEMForCSR(csr *x509.CertificateRequest) *PEMBlock {
	if csr == nil {
		return nil
	}
	return NewPEMForCSRDER(csr.Raw)
}

// NewPEMForCSRDER returns a pemBlock for certificate request
func NewPEMForCSRDER(derBytes []byte) *PEMBlock {
	return NewPEM(&pem.Block{Type: CertificateRequestPEMBlockType, Bytes: derBytes})
}

// ParsePEM decode input pem bytes to pem blocks.
func ParsePEM(pemBytes []byte) []*PEMBlock {
	return parsePEM(pemBytes, false, nil)
}

// ParsePEM find valid pem block in bytes and decode the first block.
func ParseFirstPEMBlock(pemBytes []byte) *PEMBlock {
	pems := parsePEM(pemBytes, true, nil)
	if len(pems) == 0 {
		return nil
	}
	return pems[0]
}

// ParsePrivateKeyPEM find and decode the first valid private key pem block, then
// convert it to crypto.PrivateKey(maybe rsa.PrivateKey or ecdsa.PrivateKey)
func ParsePrivateKeyPEM(pemBytes []byte) (crypto.Signer, error) {
	pems := parsePEM(pemBytes, false, privateKeyFilter)
	if len(pems) == 0 {
		return nil, errors.New("data does not contain any valid RSA or ECDSA private key")
	}
	return ParsePrivateKey(pems[0].Block.Bytes)
}

// ParsePrivateKeyPEM decode all valid certificate pem blocks to x509 certificates
func ParseCertsPEM(pemBytes []byte) ([]*x509.Certificate, error) {
	pems := parsePEM(pemBytes, false, certsFilter)
	if len(pems) == 0 {
		return nil, errors.New("data does not contain any valid RSA or ECDSA certificates")
	}
	certs := []*x509.Certificate{}
	for _, pem := range pems {
		cert, err := x509.ParseCertificate(pem.Block.Bytes)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}
	return certs, nil
}

func privateKeyFilter(block *pem.Block) bool {
	return block.Type == "PRIVATE KEY" || strings.HasSuffix(block.Type, " PRIVATE KEY")
}

func certsFilter(block *pem.Block) bool {
	return block.Type == CertificatePEMBlockType
}

// parsePEM decode input pem bytes to pem blocks.
// If filter is set and returns false, the block will be ignored
func parsePEM(pemBytes []byte, first bool, filter func(block *pem.Block) bool) []*PEMBlock {
	pems := []*PEMBlock{}
	for len(pemBytes) > 0 {
		var block *pem.Block
		block, pemBytes = pem.Decode(pemBytes)
		if block == nil {
			break
		}
		if filter != nil && !filter(block) {
			continue
		}
		pems = append(pems, NewPEM(block))

		if first {
			break
		}
	}

	return pems
}
