package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

func GetPublicKeyFromCert(certFile string) (*rsa.PublicKey, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}

	// Reading the certificate
	certPEM, err := os.ReadFile(filepath.Dir(ex) + "/" + certFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read certificate: %w", err)
	}

	// Decoding the certificate
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("unexpected PEM block type")
	}

	// Parsing the certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	// Extracting the public key
	pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("unexpected public key type")
	}

	return pubKey, nil
}

func EncryptData(data []byte, pubKey *rsa.PublicKey) ([]byte, error) {
	// Encrypting the data
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
	if err != nil {
		return nil, err
	}

	return encryptedData, nil
}
