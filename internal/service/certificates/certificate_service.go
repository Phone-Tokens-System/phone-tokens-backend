package certificates

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
)

func (s *CertificateService) GetSignedCertificateByCsrID(ctx context.Context, CsrID int) ([]byte, error) {
	cert, err := s.Storage.GetAgentInfo(ctx, CsrID)
	if err != nil {
		return nil, err
	}
	return cert.CertificatePem, nil
}

func (s *CertificateService) VerifyCertificate(cert []byte) error {
	certBytes, _ := pem.Decode(cert)
	if certBytes == nil {
		return errors.New("failed to decode certificate")
	}

	certificate, err := x509.ParseCertificate(certBytes.Bytes)
	if err != nil {
		return err
	}
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(s.CACertificatePem)
	if !ok {
		return errors.New("failed to parse root certificate")
	}

	_, err = certificate.Verify(x509.VerifyOptions{Roots: roots})
	if err != nil {
		return err
	}
	return nil
}

func (s *CertificateService) ExtractAgentId(cert []byte) (*big.Int, error) {
	certBytes, _ := pem.Decode(cert)
	if certBytes == nil {
		return nil, errors.New("failed to decode certificate")
	}

	certificate, err := x509.ParseCertificate(certBytes.Bytes)
	if err != nil {
		return nil, err
	}
	return certificate.SerialNumber, nil
}
