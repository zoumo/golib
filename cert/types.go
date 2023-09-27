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
	"crypto/tls"
	"crypto/x509"
	"net"
	"time"
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

// Options contains various common Options for creating a certificate
type Options struct {
	CommonName   string
	Organization []string
	DNSNames     []string
	IPs          []net.IP
	Usages       []x509.ExtKeyUsage
}
