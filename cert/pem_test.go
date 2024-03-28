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
	"crypto/ed25519"
	"crypto/rsa"
	"encoding/pem"
	"testing"
)

func createPEMBytes() []byte {
	caKey, _ := NewRSAPrivateKey()
	caCert, _ := NewSelfSignedCACert(Config{}, caKey)
	myKey, _ := NewECPrivateKey("P224")
	myCert, _ := NewSignedCert(Config{}, myKey, caKey, caCert)

	caKeyPEM := MarshalRSAPrivateKeyToPEM(caKey)
	myKeyPEM, _ := MarshalECPrivateKeyToPEM(myKey)
	caCertPEM := MarshalCertToPEM(caCert)
	myCertPEM := MarshalCertToPEM(myCert)

	return bytes.Join([][]byte{
		caKeyPEM.EncodeToMemory(),
		myKeyPEM.EncodeToMemory(),
		caCertPEM.EncodeToMemory(),
		myCertPEM.EncodeToMemory(),
	}, []byte{'\n'})
}

func Test_decodePEMs(t *testing.T) {
	pemBytes := createPEMBytes()

	tests := []struct {
		name   string
		filter func(block *pem.Block) bool
		first  bool
		want   int
	}{
		{
			"filter key",
			filterPrivateKey,
			false,
			2,
		},
		{
			"filter first key",
			filterPrivateKey,
			true,
			1,
		},
		{
			"filter cert",
			filterCert,
			false,
			2,
		},
		{
			"filter first cert",
			filterCert,
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
				return block.Type == ECPrivateKeyPEMBlockType
			},
			false,
			1,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := decodePEMs(pemBytes, tt.first, tt.filter); len(got) != tt.want {
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

func TestMarshalPrivateKeyToPEM(t *testing.T) {
	tests := []struct {
		name     string
		key      crypto.Signer
		wantType string
		wantErr  bool
	}{
		{
			name: "rsa",
			key: func() crypto.Signer {
				key, _ := NewRSAPrivateKey()
				return key
			}(),
			wantType: RASPrivateKeyPEMBlockType,
			wantErr:  false,
		},
		{
			name: "ec",
			key: func() crypto.Signer {
				key, _ := NewECPrivateKey(CurveP224)
				return key
			}(),
			wantType: ECPrivateKeyPEMBlockType,
			wantErr:  false,
		},
		{
			name:    "error",
			key:     ed25519.PrivateKey{},
			wantErr: true,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalPrivateKeyToPEM(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalPrivateKeyToPEM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Type != tt.wantType {
					t.Errorf("MarshalPrivateKeyToPEM() = %v, want %v", got.Type, tt.wantType)
					return
				}
			}
		})
	}
}
