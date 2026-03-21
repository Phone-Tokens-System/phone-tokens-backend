package certificates

import (
	"context"
	"phone-tokens/internal/adapter/dto"
	"phone-tokens/internal/model"
)

func (s *CertificateService) AcceptCertificateRequest(ctx context.Context, requestDTO dto.CSRRequest) (int, error) {
	csrRequest := model.CsrRequest{
		CSR:     []byte(requestDTO.CSR),
		Status:  "PENDING",
		Email:   requestDTO.Email,
		AgentID: requestDTO.AgentID,
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

	_, err = s.signCertificateForAgent(ctx, csr)
	if err != nil {
		return nil, err
	}
	return csr, nil
}

func (s *CertificateService) GetCertificateRequests(ctx context.Context) ([]model.CsrRequest, error) {
	return s.Storage.GetCsrRequests(ctx)
}
