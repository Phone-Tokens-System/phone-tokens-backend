package service

import (
	"context"
	"phone-tokens/internal/model"
)

func (s *CertificateService) AcceptCertificate(ctx context.Context, block []byte) (int, error) {
	csrRequest := model.CsrRequest{
		CSR:    block,
		Status: "PENDING",
	}

	request, err := s.storage.SaveCsrRequest(ctx, csrRequest)
	if err != nil {
		return 0, err
	}
	return request.ID, nil
}

func (s *CertificateService) ApproveCertificate(ctx context.Context, ID string) (*model.CsrRequest, error) {
	err := s.storage.UpdateCsrStatus(ctx, ID, "APPROVED")
	if err != nil {
		return nil, err
	}
	return s.storage.GetCsrRequest(ctx, ID)
}

func (s *CertificateService) GetCertificateRequests(ctx context.Context) ([]model.CsrRequest, error) {
	return s.storage.GetCsrRequests(ctx)
}
