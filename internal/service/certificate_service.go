package service

import "context"

func (s *CertificateService) GetSignedCertificateByCsrID(ctx context.Context, CsrID int) ([]byte, error) {
	cert, err := s.Storage.GetAgentInfo(ctx, CsrID)
	if err != nil {
		return nil, err
	}
	return cert.CertificatePem, nil
}
