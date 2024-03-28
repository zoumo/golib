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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

const (
	privateKeySize = 2048
)

// NewRSAPrivateKey creates a new RSA private key
func NewRSAPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, privateKeySize)
}

type EllipticCurve string

const (
	CurveP224 EllipticCurve = "P224"
	CurveP256 EllipticCurve = "P256"
	CurveP384 EllipticCurve = "P384"
	CurveP521 EllipticCurve = "P521"
)

// NewECPrivateKey create a new ECDSA provate key by curve
func NewECPrivateKey(curve EllipticCurve) (*ecdsa.PrivateKey, error) {
	var priv *ecdsa.PrivateKey
	var err error
	switch curve {
	case CurveP224:
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case CurveP256:
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case CurveP384:
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case CurveP521:
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return nil, fmt.Errorf("unrecognized elliptic curve: %q", curve)
	}
	if err != nil {
		return nil, err
	}
	return priv, nil
}

// DecryptPrivateKeyFile takes a password encrypted key file and the password
//
//	used to encrypt it and returns a slice of decrypted DER encoded bytes.
func DecryptPrivateKeyFile(keyFile, passwd string) (*PEMBlock, error) {
	keyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	return DecryptPrivateKeyBytes(keyBytes, passwd)
}

// DecryptPrivateKeyBytes takes a password encrypted PEM block and the password
// used to encrypt it and returns a slice of decrypted DER encoded bytes. It inspects
// the DEK-Info header to determine the algorithm used for decryption. If no
// DEK-Info header is present, an error is returned. If an incorrect password
// is detected an IncorrectPasswordError is returned. Because of deficiencies
// in the encrypted-PEM format, it's not always possible to detect an incorrect
// password. In these cases no error will be returned but the decrypted DER
// bytes will be random noise.
func DecryptPrivateKeyBytes(keyPEMBlock []byte, passwd string) (*PEMBlock, error) {
	pemBlock, err := findPrivateKeyInPEMBlock(keyPEMBlock)
	if err != nil {
		return nil, err
	}
	// nolint
	if !x509.IsEncryptedPEMBlock(pemBlock.Block) {
		return pemBlock, nil
	}
	// nolint
	// It's encrypted, decrypt it
	der, err := x509.DecryptPEMBlock(pemBlock.Block, []byte(passwd))
	if err != nil {
		return nil, err
	}

	newPem := &pem.Block{
		Type:  pemBlock.Block.Type,
		Bytes: der,
	}

	return NewPEMBlock(newPem), nil
}

func findPrivateKeyInPEMBlock(keyPEMBlock []byte) (*PEMBlock, error) {
	blocks := DecodePEMs(keyPEMBlock)
	for _, b := range blocks {
		if filterPrivateKey(b.Block) {
			return b, nil
		}
	}
	return nil, errors.New("no private key found in pem blocks")
}
