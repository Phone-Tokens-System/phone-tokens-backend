package service

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"phone-tokens/internal/adapter/out/repository"
	"time"
)

type CertificateService struct {
	CAKey         *ecdsa.PrivateKey   `json:"ca_key"`
	CACertificate *x509.Certificate   `json:"ca_certificate"`
	storage       *repository.Storage `json:"storage"`
}

func NewCertificateService() (*CertificateService, error) {
	certFile, err := os.ReadFile("cert.pem")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	cert, err := x509.ParseCertificate(certFile)

	if err != nil {
		return nil, err
	}

	keyFile, err := os.ReadFile("key.pem")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	key, err := x509.ParseECPrivateKey(keyFile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &CertificateService{
		CAKey:         key,
		CACertificate: cert,
	}, nil
}

func CreateOurCert() error {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2025),
		Subject: pkix.Name{
			Organization:  []string{"Company, PAKOSTIN INCORPARATED"},
			Country:       []string{"Russia"},
			Province:      []string{"GREAT SIBERIA"},
			Locality:      []string{"Novosibirsk"},
			StreetAddress: []string{"Great street"},
			PostalCode:    []string{"32423"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	caPrivKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}
	cert, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return err
	}

	caPEM, err := os.Create("cert.pem")

	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})
	if err != nil {
		return err
	}

	caPrivKeyPEM, _ := os.Create("key.pem")
	bytesKey, err := x509.MarshalECPrivateKey(caPrivKey)

	err = pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "ECDSA PRIVATE KEY",
		Bytes: bytesKey,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *CertificateService) signCertificateForAgent(block pem.Block) *bytes.Buffer {
	CSR, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		fmt.Println(err)
	}

	if err = CSR.CheckSignature(); err != nil {
		fmt.Println(err)
		return nil
	}

	randNum, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		fmt.Println(err)
	}

	serviceCert := &x509.Certificate{
		SerialNumber: big.NewInt(randNum.Int64()),
		Subject: pkix.Name{
			Organization:  CSR.Subject.Organization,
			Country:       CSR.Subject.Country,
			Province:      CSR.Subject.Province,
			Locality:      CSR.Subject.Locality,
			StreetAddress: CSR.Subject.StreetAddress,
			PostalCode:    CSR.Subject.PostalCode,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  false,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, serviceCert, s.CACertificate, s.CAKey.PublicKey, s.CAKey)

	certPem := new(bytes.Buffer)
	err = pem.Encode(certPem, &pem.Block{
		Type:  "Certificate",
		Bytes: cert,
	})

	//save certificate to our db
	return certPem
}
