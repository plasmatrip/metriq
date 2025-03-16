package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

// Программа для генерации приватного ключа и сертификата TLS.
// Она генерирует приватный ключ RSA размером 4096 бит и сертификат
// с информацией о владельце и IP-адресами 127.0.0.1 и ::1.
// Программа создает файлы cert.pem и key.pem, содержащие сертификат
// и приватный ключ соответственно.
func main() {
	cl := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	var path string
	cl.StringVar(&path, "path", "", "path to save key and cert")

	if err := cl.Parse(os.Args[1:]); err != nil {
		log.Fatalf("failed to parse flags: %v", err)
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"Metriq.Inc"},
			Country:      []string{"RU"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: nil,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	// Генерация приватного ключа
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	// Создание сертификата
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Кодирование сертификата
	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		log.Fatalf("Ошибка при кодировании сертификата: %v", err)
	}

	// Запись сертификата в файл
	certFile, err := os.Create(path + "cert.pem")
	if err != nil {
		log.Fatalf("Ошибка при создании файла сертификата: %v", err)
	}
	defer certFile.Close()

	if _, err := certFile.Write(certPEM.Bytes()); err != nil {
		log.Fatalf("Ошибка при записи сертификата в файл: %v", err)
	}
	log.Println("Сертификат успешно записан в файл cert.pem")

	// Кодирование приватного ключа
	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		log.Fatalf("Ошибка при кодировании приватного ключа: %v", err)
	}

	// Запись приватного ключа в файл
	privateKeyFile, err := os.Create(path + "key.pem")
	if err != nil {
		log.Fatalf("Ошибка при создании файла приватного ключа: %v", err)
	}
	defer privateKeyFile.Close()

	if _, err := privateKeyFile.Write(privateKeyPEM.Bytes()); err != nil {
		log.Fatalf("Ошибка при записи приватного ключа в файл: %v", err)
	}
	log.Println("Приватный ключ успешно записан в файл key.pem")
}
