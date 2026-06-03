package certificates

import (
	"context"
	"errors"
	"testing"

	"phone-tokens/internal/model"

	"github.com/google/uuid"
)

type certificateRepoStub struct {
	info *model.CertificateInfo
	err  error
}

func (r certificateRepoStub) SaveCsrRequest(ctx context.Context, request model.CsrRequest) (model.CsrRequest, error) {
	return request, nil
}

func (r certificateRepoStub) GetCsrRequest(ctx context.Context, ID int) (*model.CsrRequest, error) {
	return nil, model.ErrNotFound
}

func (r certificateRepoStub) UpdateCsrStatus(ctx context.Context, ID int, status string) error {
	return nil
}

func (r certificateRepoStub) GetCsrRequests(ctx context.Context) ([]model.CsrRequest, error) {
	return nil, nil
}

func (r certificateRepoStub) SaveCertificateInfo(ctx context.Context, info model.CertificateInfo) error {
	return nil
}

func (r certificateRepoStub) GetCertificateInfo(ctx context.Context, csrID int) (*model.CertificateInfo, error) {
	return r.info, r.err
}

func (r certificateRepoStub) GetActiveCertificateInfoByAgentID(ctx context.Context, agentID string) (*model.CertificateInfo, error) {
	return r.info, r.err
}

type approveRepoStub struct {
	request       *model.CsrRequest
	statusUpdates []string
}

func (r *approveRepoStub) SaveCsrRequest(ctx context.Context, request model.CsrRequest) (model.CsrRequest, error) {
	return request, nil
}

func (r *approveRepoStub) GetCsrRequest(ctx context.Context, ID int) (*model.CsrRequest, error) {
	return r.request, nil
}

func (r *approveRepoStub) UpdateCsrStatus(ctx context.Context, ID int, status string) error {
	r.statusUpdates = append(r.statusUpdates, status)
	return nil
}

func (r *approveRepoStub) GetCsrRequests(ctx context.Context) ([]model.CsrRequest, error) {
	return nil, nil
}

func (r *approveRepoStub) SaveCertificateInfo(ctx context.Context, info model.CertificateInfo) error {
	return nil
}

func (r *approveRepoStub) GetCertificateInfo(ctx context.Context, csrID int) (*model.CertificateInfo, error) {
	return nil, model.ErrNotFound
}

func (r *approveRepoStub) GetActiveCertificateInfoByAgentID(ctx context.Context, agentID string) (*model.CertificateInfo, error) {
	return nil, model.ErrNotFound
}

func TestGetSignedCertificateByCsrIDForAgentAllowsOwner(t *testing.T) {
	agentID := uuid.MustParse("11111111-2222-3333-4444-555555555555")
	svc := CertificateService{
		Storage: certificateRepoStub{
			info: &model.CertificateInfo{
				ID:             agentID,
				CertificatePem: []byte("cert-pem"),
			},
		},
	}

	cert, err := svc.GetSignedCertificateByCsrIDForAgent(context.Background(), 7, agentID.String())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if string(cert) != "cert-pem" {
		t.Fatalf("expected certificate pem, got %q", cert)
	}
}

func TestGetSignedCertificateByCsrIDForAgentRejectsAnotherAgent(t *testing.T) {
	svc := CertificateService{
		Storage: certificateRepoStub{
			info: &model.CertificateInfo{
				ID:             uuid.MustParse("11111111-2222-3333-4444-555555555555"),
				CertificatePem: []byte("cert-pem"),
			},
		},
	}

	_, err := svc.GetSignedCertificateByCsrIDForAgent(
		context.Background(),
		7,
		"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
	)
	if !errors.Is(err, model.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestApproveCertificateRequestDoesNotApproveInvalidCSR(t *testing.T) {
	repo := &approveRepoStub{
		request: &model.CsrRequest{
			ID:      9,
			AgentID: "11111111-2222-3333-4444-555555555555",
			CSR:     []byte("not a pem csr"),
			Status:  "PENDING",
		},
	}
	svc := CertificateService{Storage: repo}

	_, err := svc.ApproveCertificateRequest(context.Background(), 9)
	if err == nil {
		t.Fatal("expected invalid CSR error")
	}
	if len(repo.statusUpdates) != 1 || repo.statusUpdates[0] != "SIGN_FAILED" {
		t.Fatalf("expected only SIGN_FAILED status update, got %#v", repo.statusUpdates)
	}
}
