package certificates

import (
	"context"
	"phone-tokens/internal/model"
)

// TODO: service
type Repository interface {
	SaveCsrRequest(ctx context.Context, request model.CsrRequest) (model.CsrRequest, error)
	GetCsrRequest(ctx context.Context, ID int) (*model.CsrRequest, error)
	UpdateCsrStatus(ctx context.Context, ID int, status string) error
	GetCsrRequests(ctx context.Context) ([]model.CsrRequest, error)
	SaveCertificateInfo(ctx context.Context, info model.CertificateInfo) error
	GetCertificateInfo(ctx context.Context, csrID int) (*model.CertificateInfo, error)
}
