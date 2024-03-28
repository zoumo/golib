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
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"time"
)

const (
	oneYear = time.Hour * 24 * 365
)

// Config contains various common Config for creating a certificate
type Config struct {
	CommonName   string
	Organization []string
	AltNames     AltNames
	Usages       []x509.ExtKeyUsage
}

// AltNames contains the domain names and IP addresses that will be added
// to the API Server's x509 certificate SubAltNames field. The values will
// be passed directly to the x509.Certificate object.
type AltNames struct {
	DNSNames []string
	IPs      []net.IP
}

// NewSignedCert returns a new certificate signed by given ca key and certificate
func NewSignedCert(cfg Config, key crypto.Signer, caKey crypto.Signer, caCert *x509.Certificate) (*x509.Certificate, error) {
	template, err := generateCertTemplate(cfg, false)
	if err != nil {
		return nil, err
	}
	// The parameter pub is the public key of the signee and priv is the private key of the signer.
	certDerBytes, err := x509.CreateCertificate(rand.Reader, template, caCert, key.Public(), caKey)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDerBytes)
}

// NewSelfSignedCert returns a new self-signed x509 certificate
func NewSelfSignedCert(cfg Config, key crypto.Signer) (*x509.Certificate, error) {
	certDERBytes, err := newSelfSignedCert(cfg, key, false)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}

// NewSelfSignedCACert returns a new self-signed CA x509 certificate
func NewSelfSignedCACert(cfg Config, key crypto.Signer) (*x509.Certificate, error) {
	certDERBytes, err := newSelfSignedCert(cfg, key, true)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}

// NewCSR returns a new x509 certificate request
func NewCSR(cfg Config, key crypto.Signer) (*x509.CertificateRequest, error) {
	template := generateCSRTemplate(cfg)
	csrDerBytes, err := x509.CreateCertificateRequest(rand.Reader, template, key)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificateRequest(csrDerBytes)
}

func newSelfSignedCert(cfg Config, key crypto.Signer, isCA bool) ([]byte, error) {
	template, err := generateCertTemplate(cfg, isCA)
	if err != nil {
		return nil, err
	}

	return x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
}

// Based in the code https://golang.org/src/crypto/tls/generate_cert.go
func generateCertTemplate(cfg Config, isCA bool) (*x509.Certificate, error) {
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
	template := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		NotBefore:             now.UTC(),
		NotAfter:              now.Add(oneYear * 100).UTC(),
		IPAddresses:           cfg.AltNames.IPs,
		DNSNames:              cfg.AltNames.DNSNames,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           cfg.Usages,
		BasicConstraintsValid: true,
	}
	if isCA {
		// add ca flag and keyUsage
		template.IsCA = isCA
		template.KeyUsage |= x509.KeyUsageCertSign
	}
	return template, nil
}

func generateCSRTemplate(cfg Config) *x509.CertificateRequest {
	if len(cfg.Organization) == 0 {
		cfg.Organization = []string{
			"Acme Co",
		}
	}

	template := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		IPAddresses: cfg.AltNames.IPs,
		DNSNames:    cfg.AltNames.DNSNames,
	}

	return template
}
