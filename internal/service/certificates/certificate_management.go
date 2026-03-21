package certificates

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"phone-tokens/internal/model"
	"time"

	"github.com/google/uuid"
)

type CertificateService struct {
	CAKeyPem         []byte     `json:"ca_key"`
	CACertificatePem []byte     `json:"ca_certificate"`
	Storage          Repository `json:"storage"`
}

func NewCertificateService(storage Repository) (*CertificateService, error) {
	certFile, err := os.ReadFile("cert.pem")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	keyFile, err := os.ReadFile("key.pem")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &CertificateService{
		CAKeyPem:         keyFile,
		CACertificatePem: certFile,
		Storage:          storage,
	}, nil
}

func createOurCert() error {
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

func (s *CertificateService) parseCertFromPem(certFile []byte) (*x509.Certificate, error) {
	certDer, _ := pem.Decode(certFile)
	if certDer == nil {
		return nil, fmt.Errorf("failed to decode PEM certificate")
	}
	cert, err := x509.ParseCertificate(certDer.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func (s *CertificateService) parseKeyFromPem(keyFile []byte) (*ecdsa.PrivateKey, error) {
	keyDer, _ := pem.Decode(keyFile)
	if keyDer == nil {
		return nil, fmt.Errorf("failed to decode PEM key")
	}
	key, err := x509.ParseECPrivateKey(keyDer.Bytes)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return key, nil
}
func (s *CertificateService) signCertificateForAgent(ctx context.Context, csr *model.CsrRequest) (*bytes.Buffer, error) {
	fmt.Println("certificate")
	block := csr.CSR

	fmt.Println(string(block))
	csrDer, _ := pem.Decode(block)
	fmt.Println(csrDer.Type)
	if csrDer == nil || csrDer.Type != "CERTIFICATE REQUEST" {
		return nil, fmt.Errorf("failed to decode certificate request block")
	}
	fmt.Printf("DER len: %d\n", len(csrDer.Bytes))
	fmt.Println(string(csrDer.Bytes))
	CSR, err := x509.ParseCertificateRequest(csrDer.Bytes)
	if err != nil {
		return nil, err
	}

	if err = CSR.CheckSignature(); err != nil {
		fmt.Printf("sign %v", err)
		return nil, err
	}
	fmt.Println(string(s.CACertificatePem))
	certCA, err := s.parseCertFromPem(s.CACertificatePem)
	if err != nil {
		fmt.Printf("ourcert %v", err)
		return nil, err
	}
	keyCA, err := s.parseKeyFromPem(s.CAKeyPem)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	parsedUUID, err := uuid.Parse(csr.AgentID)
	if err != nil {
		return nil, err
	}

	serialNumber := new(big.Int).SetBytes(parsedUUID[:])

	serviceCert := &x509.Certificate{
		SerialNumber: serialNumber,
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

	cert, err := x509.CreateCertificate(rand.Reader, serviceCert, certCA, CSR.PublicKey, keyCA)

	certPem := new(bytes.Buffer)
	err = pem.Encode(certPem, &pem.Block{
		Type:  "Certificate",
		Bytes: cert,
	})

	agentInfo := model.CertificateInfo{
		OrganizationID: CSR.Subject.Organization[0],
		CertificatePem: certPem.Bytes(),
		IsActive:       true,
		CsrID:          csr.ID,
		ID:             parsedUUID,
	}
	err = s.Storage.SaveCertificateInfo(ctx, agentInfo)
	if err != nil {
		return nil, err
	}
	return certPem, nil
}
