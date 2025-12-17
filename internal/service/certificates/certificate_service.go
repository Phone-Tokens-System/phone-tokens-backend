package certificates

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"

	"github.com/google/uuid"
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

func (s *CertificateService) ExtractAgentId(cert []byte) (uuid.UUID, error) {
	certBytes, _ := pem.Decode(cert)
	if certBytes == nil {
		return uuid.UUID{}, errors.New("failed to decode certificate")
	}

	certificate, err := x509.ParseCertificate(certBytes.Bytes)
	if err != nil {
		return uuid.UUID{}, err
	}

	uuidId, err := SerialToUUID(certificate.SerialNumber)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuidId, nil
}

func SerialToUUID(serial *big.Int) (uuid.UUID, error) {
	b := serial.Bytes()

	if len(b) != 16 {
		if len(b) < 16 {
			padded := make([]byte, 16)
			copy(padded[16-len(b):], b)
			b = padded
		} else {
			return uuid.Nil, fmt.Errorf("invalid serial length: %d", len(b))
		}
	}

	return uuid.FromBytes(b)
}
