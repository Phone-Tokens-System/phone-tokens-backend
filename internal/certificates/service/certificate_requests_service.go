package service

import (
	"context"
	"phone-tokens/internal/adapter/dto"
	"phone-tokens/internal/certificates/model"
)

func (s *CertificateService) AcceptCertificateRequest(ctx context.Context, requestDTO dto.CSRRequest) (int, error) {
	csrRequest := model.CsrRequest{
		CSR:    requestDTO.CSR,
		Status: "PENDING",
		Email:  requestDTO.Email,
	}

	request, err := s.Storage.SaveCsrRequest(ctx, csrRequest)
	if err != nil {
		return 0, err
	}
	return request.ID, nil
}

func (s *CertificateService) ApproveCertificateRequest(ctx context.Context, ID int) (*model.CsrRequest, error) {
	err := s.Storage.UpdateCsrStatus(ctx, ID, "APPROVED")
	if err != nil {
		return nil, err
	}

	csr, err := s.Storage.GetCsrRequest(ctx, ID)
	if err != nil {
		return nil, err
	}

	_, err = s.signCertificateForAgent(ctx, csr.CSR, ID)
	if err != nil {
		return nil, err
	}
	return csr, nil
}

func (s *CertificateService) GetCertificateRequests(ctx context.Context) ([]model.CsrRequest, error) {
	return s.Storage.GetCsrRequests(ctx)
}
