// Copyright 2023 jim.zoumo@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cert

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	// CertificatePEMBlockType is a possible value for pem.Block.Type.
	CertificatePEMBlockType = "CERTIFICATE"
	// CertificateRequestPEMBlockType is a possible value for pem.Block.Type.
	CertificateRequestPEMBlockType = "CERTIFICATE REQUEST"
	// RASPrivateKeyPEMBlockType is a possible value for pem.Block.Type.
	RASPrivateKeyPEMBlockType = "RSA PRIVATE KEY"
	// ECPrivateKeyPEMBlockType is a possible value for pem.Block.Type.
	ECPrivateKeyPEMBlockType = "EC PRIVATE KEY"
	// PrivateKeyBlockType is a possible value for pem.Block.Type.
	PrivateKeyPEMBlockType = "PRIVATE KEY"
)

// PEMBlock contains the raw bytes and a block of pem
type PEMBlock struct {
	*pem.Block

	onece  sync.Once
	buffer bytes.Buffer
	err    error
}

// EncodeToMemory returns the PEM encoding bytes of p.
//
// If b has invalid headers and cannot be encoded,
// EncodeToMemory returns nil. If it is important to
// report details about this error case, use Encode instead.
func (p *PEMBlock) EncodeToMemory() []byte {
	p.writeToBuffer()
	if p.err != nil {
		return nil
	}
	return p.buffer.Bytes()
}

// WriteTo writes the PEM encoding of block to out.
func (p *PEMBlock) WriteTo(out io.Writer) (int64, error) {
	p.writeToBuffer()
	if p.err != nil {
		return 0, p.err
	}
	n, err := out.Write(p.buffer.Bytes())
	return int64(n), err
}

// WriteFile writes the PEM encoding to a file
func (p *PEMBlock) WriteFile(f string) error {
	p.writeToBuffer()
	if p.err != nil {
		return p.err
	}
	return os.WriteFile(f, p.buffer.Bytes(), 0o644)
}

func (p *PEMBlock) writeToBuffer() {
	p.onece.Do(func() {
		p.err = pem.Encode(&p.buffer, p.Block)
	})
}

// NewPEMBlock creates a new PEM struct from pem.Block
func NewPEMBlock(b *pem.Block) *PEMBlock {
	return &PEMBlock{
		Block: b,
	}
}

// DecodePEMs decode input pem bytes to pem blocks.
func DecodePEMs(pemBytes []byte) []*PEMBlock {
	return decodePEMs(pemBytes, false, nil)
}

// DecodeFirstPEM find valid pem block in bytes and decode the first block.
func DecodeFirstPEM(pemBytes []byte) *PEMBlock {
	pems := decodePEMs(pemBytes, true, nil)
	if len(pems) == 0 {
		return nil
	}
	return pems[0]
}

// MarshalPrivateKeyToPEM converts the private key to PEM block.
func MarshalPrivateKeyToPEM(key crypto.Signer) (*PEMBlock, error) {
	switch pkey := key.(type) {
	case *rsa.PrivateKey:
		return MarshalRSAPrivateKeyToPEM(pkey), nil
	case *ecdsa.PrivateKey:
		return MarshalECPrivateKeyToPEM(pkey)
	default:
		return nil, errors.New("the key must be *rsa.PrivateKey or *ecdsa.PrivateKey")
	}
}

// MarshalRSAPrivateKeyToPEM converts an RSA private key to PKCS #1, ASN.1 DER form.
func MarshalRSAPrivateKeyToPEM(key *rsa.PrivateKey) *PEMBlock {
	return NewPEMBlock(&pem.Block{
		Type:  RASPrivateKeyPEMBlockType,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
}

// MarshalECPrivateKeyToPEM converts an EC private key to SEC 1, ASN.1 DER form.
func MarshalECPrivateKeyToPEM(key *ecdsa.PrivateKey) (*PEMBlock, error) {
	bytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, err
	}
	return NewPEMBlock(&pem.Block{
		Type:  ECPrivateKeyPEMBlockType,
		Bytes: bytes,
	}), nil
}

// MarshalCertToPEM returns a pemBlock for x509 certificate
func MarshalCertToPEM(crt *x509.Certificate) *PEMBlock {
	if crt == nil {
		return nil
	}
	return NewPEMBlock(&pem.Block{
		Type:  CertificatePEMBlockType,
		Bytes: crt.Raw,
	})
}

// MarshalCSRToPEM returns a pemBlock for certificate request
func MarshalCSRToPEM(csr *x509.CertificateRequest) *PEMBlock {
	if csr == nil {
		return nil
	}
	return NewPEMBlock(&pem.Block{
		Type:  CertificateRequestPEMBlockType,
		Bytes: csr.Raw,
	})
}

// ParsePrivateKeyPEM find and decode the first valid private key pem block, then
// convert it to crypto.PrivateKey(maybe rsa.PrivateKey or ecdsa.PrivateKey)
func ParsePrivateKeyPEM(pemBytes []byte) (crypto.Signer, error) {
	pems := decodePEMs(pemBytes, true, filterPrivateKey)
	if len(pems) == 0 {
		return nil, errors.New("data does not contain any valid RSA or ECDSA private key")
	}
	return parsePrivateKey(pems[0].Block)
}

// ParseCertPEM decode first valid certificate pem blocks to x509 certificate
func ParseCertPEM(pemBytes []byte) (*x509.Certificate, error) {
	pems := decodePEMs(pemBytes, true, filterCert)
	if len(pems) == 0 {
		return nil, errors.New("pem data does not contain any valid RSA or ECDSA certificates")
	}
	cert, err := x509.ParseCertificate(pems[0].Block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// ParseCertsPEM decode all valid certificate pem blocks to x509 certificates
func ParseCertsPEM(pemBytes []byte) ([]*x509.Certificate, error) {
	pems := decodePEMs(pemBytes, false, filterCert)
	if len(pems) == 0 {
		return nil, errors.New("pem data does not contain any valid RSA or ECDSA certificates")
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

// parsePrivateKey attempts to parse the given private key pem.Block.
// OpenSSL 0.9.8 generates PKCS#1 private keys by default, while OpenSSL 1.0.0
// generates PKCS#8 keys. OpenSSL ecparam generates SEC1 EC private keys for ECDSA.
// We try all three.
func parsePrivateKey(block *pem.Block) (crypto.Signer, error) {
	switch block.Type {
	case RASPrivateKeyPEMBlockType:
		// RSA Private Key in PKCS#1 format
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case ECPrivateKeyPEMBlockType:
		// ECDSA Private Key in ASN.1 format
		return x509.ParseECPrivateKey(block.Bytes)
	case PrivateKeyPEMBlockType:
		// RSA or ECDSA Private Key in unencrypted PKCS#8 format
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		signer, ok := key.(crypto.Signer)
		if !ok {
			return nil, fmt.Errorf("parsed private key from PKCS#8 is not crypto.Signer")
		}
		return signer, nil
	}
	return nil, fmt.Errorf("unknown private key type")
}

// decodePEMs decode input pem bytes to pem blocks.
// If filter is set and returns false, the block will be ignored
func decodePEMs(pemBytes []byte, getOne bool, filter func(block *pem.Block) bool) []*PEMBlock {
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
		pems = append(pems, NewPEMBlock(block))

		if getOne {
			break
		}
	}
	return pems
}

func filterPrivateKey(block *pem.Block) bool {
	switch block.Type {
	case RASPrivateKeyPEMBlockType, ECPrivateKeyPEMBlockType, PrivateKeyPEMBlockType:
		return true
	}
	return false
}

func filterCert(block *pem.Block) bool {
	return block.Type == CertificatePEMBlockType
}
