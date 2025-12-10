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
	"phone-tokens/internal/adapter/out/repository"
	"phone-tokens/internal/model"
	"time"

	"github.com/google/uuid"
)

type CertificateService struct {
	CAKeyPem         []byte              `json:"ca_key"`
	CACertificatePem []byte              `json:"ca_certificate"`
	Storage          *repository.Storage `json:"storage"`
}

func NewCertificateService() (*CertificateService, error) {
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
	cert, err := x509.ParseCertificate(certFile)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func (s *CertificateService) parseKeyFromPem(keyFile []byte) (*ecdsa.PrivateKey, error) {
	key, err := x509.ParseECPrivateKey(keyFile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return key, nil
}
func (s *CertificateService) signCertificateForAgent(ctx context.Context, block []byte, csrID int) (*bytes.Buffer, error) {
	CSR, err := x509.ParseCertificateRequest(block)
	if err != nil {
		return nil, err
	}

	if err = CSR.CheckSignature(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	certCA, err := s.parseCertFromPem(s.CACertificatePem)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	keyCA, err := s.parseKeyFromPem(s.CAKeyPem)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	uid := uuid.New()
	serialNumber := new(big.Int).SetBytes(uid[:])

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

	cert, err := x509.CreateCertificate(rand.Reader, serviceCert, certCA, keyCA.PublicKey, keyCA)

	certPem := new(bytes.Buffer)
	err = pem.Encode(certPem, &pem.Block{
		Type:  "Certificate",
		Bytes: cert,
	})

	agentInfo := model.ExternalAgentInfo{
		OrganizationID: CSR.Subject.Organization[0],
		CertificatePem: certPem.Bytes(),
		IsActive:       true,
		CsrID:          csrID,
		ID:             uid,
	}
	err = s.Storage.SaveAgentInfo(ctx, agentInfo)
	if err != nil {
		return nil, err
	}
	return certPem, nil
}
