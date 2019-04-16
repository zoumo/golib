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
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"io/ioutil"
	"math/big"
	"time"
)

const (
	privateKeySize = 2048
	oneYear        = time.Hour * 24 * 365
)

// LoadX509KeyPair reads and parses a public/private key pair from a pair
// of files. The files must contain PEM encoded data.
func LoadX509KeyPair(certFile, keyFile string) (*TLSCertificate, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return convert(cert)
}

// X509KeyPair parses a public/private key pair from a pair of
// PEM encoded data.
func X509KeyPair(certPEMBlock, keyPEMBlock []byte) (*TLSCertificate, error) {
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, err
	}
	return convert(cert)
}

// LoadX509KeyPairWithPassword parses a encryption public/private key pair from a pair of
// PEM encoded data.
func LoadX509KeyPairWithPassword(certFile, keyFile, passwd string) (*TLSCertificate, error) {
	certPEMBlock, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	keyPEMBlock, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	return X509KeyPairWithPassword(certPEMBlock, keyPEMBlock, passwd)
}

// X509KeyPairWithPassword parses a public/private key pair from a pair of
// PEM encoded data.
func X509KeyPairWithPassword(certPEMBlock, keyPEMBlock []byte, passwd string) (*TLSCertificate, error) {
	keyPEM, err := DecryptPrivateKeyBytes(keyPEMBlock, passwd)
	if err != nil {
		return nil, err
	}
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEM.Raw)
	if err != nil {
		return nil, err
	}
	tlsCert, err := convert(cert)
	return tlsCert, err
}

func convert(cert tls.Certificate) (*TLSCertificate, error) {
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}
	return &TLSCertificate{
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

// NewCertificateRequest returns a new x509 certificate request
func NewCertificateRequest(cfg Options, key crypto.Signer) (*x509.CertificateRequest, error) {
	csrDERBytes, err := NewCertificateRequestBytes(cfg, key)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificateRequest(csrDERBytes)
}

// NewCertificateRequestBytes returns a new certificate bytes in DER encoding
func NewCertificateRequestBytes(cfg Options, key crypto.Signer) ([]byte, error) {
	if len(cfg.Organization) == 0 {
		cfg.Organization = []string{
			"Acme Co",
		}
	}

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		IPAddresses: cfg.IPs,
		DNSNames:    cfg.DNSNames,
	}

	return x509.CreateCertificateRequest(rand.Reader, &template, key)
}

// NewSelfSignedCertificate returns a new self-signed x509 certificate
//
// All keys types that are implemented via crypto.Signer are supported (This
// includes *rsa.PrivateKey and *ecdsa.PrivateKey.)
func NewSelfSignedCertificate(cfg Options, key crypto.Signer) (*x509.Certificate, error) {
	certDERBytes, err := newSelfSignedCertificateBytes(cfg, key, false)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}

// NewSelfSignedCertificateBytes returns a new self-signed certificate in DER encoding
//
// All keys types that are implemented via crypto.Signer are supported (This
// includes *rsa.PrivateKey and *ecdsa.PrivateKey.)
func NewSelfSignedCertificateBytes(cfg Options, key crypto.Signer) ([]byte, error) {
	return newSelfSignedCertificateBytes(cfg, key, false)
}

// NewSelfSignedCACert returns a new self-signed CA x509 certificate
//
// All keys types that are implemented via crypto.Signer are supported (This
// includes *rsa.PrivateKey and *ecdsa.PrivateKey.)
func NewSelfSignedCACert(cfg Options, key crypto.Signer) (*x509.Certificate, error) {
	certDERBytes, err := newSelfSignedCertificateBytes(cfg, key, true)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}

// NewSelfSignedCACertBytes returns a new self-signed CA certificate in DER encoding
//
// All keys types that are implemented via crypto.Signer are supported (This
// includes *rsa.PrivateKey and *ecdsa.PrivateKey.)
func NewSelfSignedCACertBytes(cfg Options, key crypto.Signer) ([]byte, error) {
	return newSelfSignedCertificateBytes(cfg, key, true)
}

// NewSignedCert returns a new certificate signed by given ca key and certificate
//
// All keys types that are implemented via crypto.Signer are supported (This
// includes *rsa.PrivateKey and *ecdsa.PrivateKey.)
func NewSignedCert(cfg Options, key crypto.Signer, caKey crypto.Signer, caCert *x509.Certificate) (*x509.Certificate, error) {
	certBytes, err := newSignedCertBytes(cfg, key, caKey, caCert)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certBytes)
}

// Based in the code https://golang.org/src/crypto/tls/generate_cert.go
func newSignedCertBytes(cfg Options, key crypto.Signer, caKey crypto.Signer, caCert *x509.Certificate) ([]byte, error) {
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
func newSelfSignedCertificateBytes(cfg Options, key crypto.Signer, isCA bool) ([]byte, error) {
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
