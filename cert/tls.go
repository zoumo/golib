/**
 * Copyright 2024 jim.zoumo@gmail.com
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cert

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"time"
)

// TLSCertificate represents the external cert api secret for https
type TLSCertificate struct {
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

// PkixName represents an X.509 distinguished name. This only includes the common
// elements of a DN. When parsing, all elements are stored in Names and
// non-standard elements can be extracted from there. When marshaling, elements
// in ExtraNames are appended and override other values with the same OID.
type PkixName struct {
	Organization []string `json:"organization,omitempty"`
	// CommonName
	CommonName string `json:"commonName,omitempty"`
}

// LoadX509KeyPair reads and parses a public/private key pair from a pair
// of files. The files must contain PEM encoded data.
func LoadX509KeyPair(certFile, keyFile string) (*TLSCertificate, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return convertTLSCertificate(cert)
}

// X509KeyPair parses a public/private key pair from a pair of
// PEM encoded data.
func X509KeyPair(certPEMBlock, keyPEMBlock []byte) (*TLSCertificate, error) {
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, err
	}
	return convertTLSCertificate(cert)
}

// LoadX509KeyPairWithPassword parses a encryption public/private key pair from a pair of
// PEM encoded data.
func LoadX509KeyPairWithPassword(certFile, keyFile, passwd string) (*TLSCertificate, error) {
	certPEMBlock, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	keyPEMBlock, err := os.ReadFile(keyFile)
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
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEM.EncodeToMemory())
	if err != nil {
		return nil, err
	}
	tlsCert, err := convertTLSCertificate(cert)
	return tlsCert, err
}

func convertTLSCertificate(cert tls.Certificate) (*TLSCertificate, error) {
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
