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
	"crypto"
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateKeyAndCert() (caKey crypto.Signer, caCert *x509.Certificate, key crypto.Signer, cert *x509.Certificate) {
	caKey, _ = NewRSAPrivateKey()
	caCert, _ = NewSelfSignedCACert(Config{
		CommonName:   "ca.example.com",
		Organization: []string{"ca"},
	}, caKey)

	key, _ = NewRSAPrivateKey()
	cert, _ = NewSignedCert(Config{
		CommonName:   "test.example.com",
		Organization: []string{"server"},
	}, key, caKey, caCert)
	return caKey, caCert, key, cert
}

func TestX509KeyPair(t *testing.T) {
	_, _, key, cert := generateKeyAndCert()

	keyPEM, _ := MarshalPrivateKeyToPEM(key)
	certPEM := MarshalCertToPEM(cert)

	tlsCert, err := X509KeyPair(certPEM.EncodeToMemory(), keyPEM.EncodeToMemory())
	assert.Nil(t, err)
	assert.Equal(t, tlsCert.Issuer.CommonName, "ca.example.com")
	assert.Equal(t, tlsCert.Issuer.Organization, []string{"ca"})
	assert.Equal(t, tlsCert.Subject.CommonName, "test.example.com")
	assert.Equal(t, tlsCert.Subject.Organization, []string{"server"})
	assert.Equal(t, tlsCert.X509Cert.Raw, cert.Raw)

	// assert.Equal(t, tlsCert.)

	// type args struct {
	// 	certPEMBlock []byte
	// 	keyPEMBlock  []byte
	// }
	// tests := []struct {
	// 	name    string
	// 	args    args
	// 	want    *TLSCertificate
	// 	wantErr bool
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		got, err := X509KeyPair(tt.args.certPEMBlock, tt.args.keyPEMBlock)
	// 		if (err != nil) != tt.wantErr {
	// 			t.Errorf("X509KeyPair() error = %v, wantErr %v", err, tt.wantErr)
	// 			return
	// 		}
	// 		if !reflect.DeepEqual(got, tt.want) {
	// 			t.Errorf("X509KeyPair() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
}
