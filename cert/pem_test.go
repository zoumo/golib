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
	"bytes"
	"crypto/rsa"
	"encoding/pem"
	"testing"
)

func createPEMBytes() []byte {
	caKey, _ := NewRSAPrivateKey()
	caCert, _ := NewSelfSignedCACert(Options{}, caKey)
	myKey, _ := NewECDSAPrivateKey("P224")
	myCert, _ := NewSignedCert(Options{}, myKey, caKey, caCert)

	caKeyPEM := NewPEMForRSAKey(caKey)
	myKeyPEM := NewPEMForECDSAKey(myKey)
	caCertPEM := NewPEMForCert(caCert)
	myCertPEM := NewPEMForCert(myCert)

	return bytes.Join([][]byte{
		caKeyPEM.EncodeToMemory(),
		myKeyPEM.EncodeToMemory(),
		caCertPEM.EncodeToMemory(),
		myCertPEM.EncodeToMemory(),
	}, []byte{'\n'})
}

func Test_parsePEM(t *testing.T) {
	pemBytes := createPEMBytes()

	tests := []struct {
		name   string
		filter func(block *pem.Block) bool
		first  bool
		want   int
	}{
		{
			"filter key",
			privateKeyFilter,
			false,
			2,
		},
		{
			"filter first key",
			privateKeyFilter,
			true,
			1,
		},
		{
			"filter cert",
			certsFilter,
			false,
			2,
		},
		{
			"filter cert",
			certsFilter,
			true,
			1,
		},
		{
			"filter rsa key",
			func(block *pem.Block) bool {
				return block.Type == RASPrivateKeyPEMBlockType
			},
			false,
			1,
		},
		{
			"filter ec key",
			func(block *pem.Block) bool {
				return block.Type == ECDSAPrivateKeyPEMBlockType
			},
			false,
			1,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePEM(pemBytes, tt.first, tt.filter); len(got) != tt.want {
				t.Errorf("parsePEM() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestParsePrivateKeyPEM(t *testing.T) {
	pemBytes := createPEMBytes()
	got, err := ParsePrivateKeyPEM(pemBytes)
	if err != nil {
		t.Errorf("ParsePrivateKeyPEM() error = %v", err)
		return
	}
	_, ok := got.(*rsa.PrivateKey)
	if !ok {
		t.Errorf("ParsePrivateKeyPEM() = %v, want rsa.PrivateKey", got)
	}
}

func TestParseCertsPEM(t *testing.T) {
	pemBytes := createPEMBytes()
	got, err := ParseCertsPEM(pemBytes)
	if err != nil {
		t.Errorf("ParseCertsPEM() error = %v", err)
		return
	}
	want := 2
	if len(got) != want {
		t.Errorf("ParseCertsPEM got = %v, want = %v", len(got), want)
	}
}
