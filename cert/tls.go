/*
Copyright 2016 caicloud authors. All rights reserved.

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
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"time"
)

const (
	privateKeySize = 2048
	oneYear        = time.Hour * 24 * 365
)

// PkixName represents an X.509 distinguished name. This only includes the common
// elements of a DN. When parsing, all elements are stored in Names and
// non-standard elements can be extracted from there. When marshaling, elements
// in ExtraNames are appended and override other values with the same OID.
type PkixName struct {
	Organization []string `json:"organization,omitempty"`
	// CommonName
	CommonName string `json:"commonName,omitempty"`
}

// TLSCert represents the external cert api secret for https
type TLSCert struct {
	// certificate is not valid before this time
	NotBefore time.Time `json:"notBefore,omitempty"`
	// certificate is not valid after this time
	NotAfter time.Time `json:"notAfter,omitempty"`
	// Issuer information extracted from X.509 cert
	Issuer PkixName `json:"issuer,omitempty"`
	// Subject information extracted from X.509 cert
	Subject PkixName `json:"subject,omitempty"`

	// Subject Alternate Name values
	DNSNames    []string `json:"dnsNames,omitempty"`
	IPAddresses []net.IP `json:"ipAddresses,omitempty"`

	Cert     tls.Certificate   `json:"-"`
	X509Cert *x509.Certificate `json:"-"`
}

// TLSCertConfig contains various common config for creating a certificate
type TLSCertConfig struct {
	CommonName   string
	Organization []string
	DNSNames     []string
	IPs          []net.IP
	Usages       []x509.ExtKeyUsage
}

// LoadX509KeyPair reads and parses a public/private key pair from a pair
// of files. The files must contain PEM encoded data.
func LoadX509KeyPair(certFile, keyFile string) (*TLSCert, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return convert(cert)
}

// X509KeyPair parses a public/private key pair from a pair of
// PEM encoded data.
func X509KeyPair(certPEMBlock, keyPEMBlock []byte) (*TLSCert, error) {
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, err
	}
	return convert(cert)
}

func convert(cert tls.Certificate) (*TLSCert, error) {
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}
	return &TLSCert{
		NotBefore: x509Cert.NotAfter,
		NotAfter:  x509Cert.NotAfter,
		Issuer: PkixName{
			CommonName:   x509Cert.Issuer.CommonName,
			Organization: x509Cert.Issuer.Organization,
		},
		Subject: PkixName{
			CommonName:   x509Cert.Subject.CommonName,
			Organization: x509Cert.Subject.Organization,
		},
		DNSNames:    x509Cert.DNSNames,
		IPAddresses: x509Cert.IPAddresses,
		Cert:        cert,
		X509Cert:    x509Cert,
	}, nil
}

func x509ToTLSCert(x509Cert *x509.Certificate) *TLSCert {
	return &TLSCert{
		NotBefore: x509Cert.NotAfter,
		NotAfter:  x509Cert.NotAfter,
		Issuer: PkixName{
			CommonName:   x509Cert.Issuer.CommonName,
			Organization: x509Cert.Issuer.Organization,
		},
		Subject: PkixName{
			CommonName:   x509Cert.Subject.CommonName,
			Organization: x509Cert.Subject.Organization,
		},
		DNSNames:    x509Cert.DNSNames,
		IPAddresses: x509Cert.IPAddresses,
		X509Cert:    x509Cert,
	}
}

// NewPrivateKey creates a new RSA private key
func NewPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, privateKeySize)
}

// NewECDSAPrivateKey create a new ECDSA provate key by curve
func NewECDSAPrivateKey(curve string) (*ecdsa.PrivateKey, error) {
	var priv *ecdsa.PrivateKey
	var err error
	switch curve {
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return nil, fmt.Errorf("Unrecognized elliptic curve: %q", curve)
	}
	if err != nil {
		return nil, err
	}
	return priv, nil
}

// NewSelfSignedCert returns a new self-signed x509 certificate
//
// All keys types that are implemented via crypto.Signer are supported (This
// includes *rsa.PrivateKey and *ecdsa.PrivateKey.)
func NewSelfSignedCert(cfg TLSCertConfig, key crypto.Signer) (*x509.Certificate, error) {
	certDERBytes, err := newSelfSignedCertBytes(cfg, key, false)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}

// NewSelfSignedCertBytes returns a new self-signed certificate in DER encoding
//
// All keys types that are implemented via crypto.Signer are supported (This
// includes *rsa.PrivateKey and *ecdsa.PrivateKey.)
func NewSelfSignedCertBytes(cfg TLSCertConfig, key crypto.Signer) ([]byte, error) {
	return newSelfSignedCertBytes(cfg, key, false)
}

// NewSelfSignedCACert returns a new self-signed CA x509 certificate
//
// All keys types that are implemented via crypto.Signer are supported (This
// includes *rsa.PrivateKey and *ecdsa.PrivateKey.)
func NewSelfSignedCACert(cfg TLSCertConfig, key crypto.Signer) (*x509.Certificate, error) {
	certDERBytes, err := newSelfSignedCertBytes(cfg, key, true)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}

// NewSelfSignedCACertBytes returns a new self-signed CA certificate in DER encoding
//
// All keys types that are implemented via crypto.Signer are supported (This
// includes *rsa.PrivateKey and *ecdsa.PrivateKey.)
func NewSelfSignedCACertBytes(cfg TLSCertConfig, key crypto.Signer) ([]byte, error) {
	return newSelfSignedCertBytes(cfg, key, true)
}

// NewSignedCert returns a new certificate signed by given ca key and certificate
//
// All keys types that are implemented via crypto.Signer are supported (This
// includes *rsa.PrivateKey and *ecdsa.PrivateKey.)
func NewSignedCert(cfg TLSCertConfig, key crypto.Signer, caKey crypto.Signer, caCert *x509.Certificate) (*x509.Certificate, error) {
	certBytes, err := newSignedCertBytes(cfg, key, caKey, caCert)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certBytes)
}

// Based in the code https://golang.org/src/crypto/tls/generate_cert.go
func newSignedCertBytes(cfg TLSCertConfig, key crypto.Signer, caKey crypto.Signer, caCert *x509.Certificate) ([]byte, error) {
	if len(cfg.Organization) == 0 {
		cfg.Organization = []string{
			"Acme Co",
		}
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		DNSNames:              cfg.DNSNames,
		IPAddresses:           cfg.IPs,
		SerialNumber:          serial,
		NotBefore:             caCert.NotBefore,
		NotAfter:              time.Now().Add(oneYear * 10).UTC(),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           cfg.Usages,
		BasicConstraintsValid: true,
	}

	// The parameter pub is the public key of the signee and priv is the private key of the signer.
	return x509.CreateCertificate(rand.Reader, &template, caCert, key.Public(), caKey)
}

// Based in the code https://golang.org/src/crypto/tls/generate_cert.go
func newSelfSignedCertBytes(cfg TLSCertConfig, key crypto.Signer, isCA bool) ([]byte, error) {
	if len(cfg.Organization) == 0 {
		cfg.Organization = []string{
			"Acme Co",
		}
	}

	now := time.Now()
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}
	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		NotBefore:             now.UTC(),
		NotAfter:              now.Add(oneYear * 10).UTC(),
		IPAddresses:           cfg.IPs,
		DNSNames:              cfg.DNSNames,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           cfg.Usages,
		BasicConstraintsValid: true,
	}
	if isCA {
		template.IsCA = isCA
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	return x509.CreateCertificate(rand.Reader, &template, &template, key.Public(), key)
}
