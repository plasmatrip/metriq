package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func LoadPrivateKey(filePath string) (*rsa.PrivateKey, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}

	// Reading the private key
	keyData, err := os.ReadFile(filepath.Dir(ex) + "/" + filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read private key: %w", err)
	}

	// Decoding the private key
	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("unexpected PEM block type")
	}

	// Parsing the private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("cannot parse private key: %w", err)
	}

	return privateKey, nil
}

func DecryptData(encryptedData []byte, privKey *rsa.PrivateKey) ([]byte, error) {
	// Decrypting the data
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, encryptedData)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}
