package cert

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

// NewRSAPrivateKey creates a new RSA private key
func NewRSAPrivateKey() (*rsa.PrivateKey, error) {
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

// ParsePrivateKey attempts to parse the given private key DER block. OpenSSL 0.9.8 generates
// PKCS#1 private keys by default, while OpenSSL 1.0.0 generates PKCS#8 keys.
// OpenSSL ecparam generates SEC1 EC private keys for ECDSA. We try all three.
func ParsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("tls: found unknown private key type in PKCS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}

	return nil, errors.New("tls: failed to parse private key")
}

// DecryptPrivateKeyFile takes a password encrypted key file and the password
//  used to encrypt it and returns a slice of decrypted DER encoded bytes.
func DecryptPrivateKeyFile(keyFile, passwd string) (*PEM, error) {
	keyPEMBlock, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	return DecryptPrivateKeyBytes(keyPEMBlock, passwd)
}

// DecryptPrivateKeyBytes takes a password encrypted PEM block and the password
// used to encrypt it and returns a slice of decrypted DER encoded bytes. It inspects
// the DEK-Info header to determine the algorithm used for decryption. If no
// DEK-Info header is present, an error is returned. If an incorrect password
// is detected an IncorrectPasswordError is returned. Because of deficiencies
// in the encrypted-PEM format, it's not always possible to detect an incorrect
// password. In these cases no error will be returned but the decrypted DER
// bytes will be random noise.
func DecryptPrivateKeyBytes(keyPEMBlock []byte, passwd string) (*PEM, error) {
	pemBlock, err := findPrivateKeyInPEMBlock(keyPEMBlock)
	if err != nil {
		return nil, err
	}
	if !x509.IsEncryptedPEMBlock(pemBlock.Block) {
		return pemBlock, nil
	}
	// It's encrypted, decrypt it
	der, err := x509.DecryptPEMBlock(pemBlock.Block, []byte(passwd))
	if err != nil {
		return nil, err
	}

	newPem := &pem.Block{
		Type:  pemBlock.Block.Type,
		Bytes: der,
	}

	return NewPEM(newPem), nil
}

func findPrivateKeyInPEMBlock(keyPEMBlock []byte) (*PEM, error) {
	var skippedBlockTypes []string
	var keyDERBlock *pem.Block
	for {
		keyDERBlock, keyPEMBlock = pem.Decode(keyPEMBlock)
		if keyDERBlock == nil {
			if len(skippedBlockTypes) == 0 {
				return nil, errors.New("tls: failed to find any PEM data in key input")
			}
			if len(skippedBlockTypes) == 1 && skippedBlockTypes[0] == "CERTIFICATE" {
				return nil, errors.New("tls: found a certificate rather than a key in the PEM for the private key")
			}
			return nil, fmt.Errorf("tls: failed to find PEM block with type ending in \"PRIVATE KEY\" in key input after skipping PEM blocks of the following types: %v", skippedBlockTypes)
		}
		if keyDERBlock.Type == "PRIVATE KEY" || strings.HasSuffix(keyDERBlock.Type, " PRIVATE KEY") {
			break
		}
		skippedBlockTypes = append(skippedBlockTypes, keyDERBlock.Type)
	}

	return NewPEM(keyDERBlock), nil
}
