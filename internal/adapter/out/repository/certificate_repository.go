package repository

import (
	"context"
	"log"
	"phone-tokens/internal/model"

	"gorm.io/gorm"
)

type CertificateRepository struct {
	db *gorm.DB
}

func NewCertificateRepository(db *gorm.DB) *CertificateRepository {
	return &CertificateRepository{db: db}
}

func (r *CertificateRepository) SaveCsrRequest(ctx context.Context, request model.CsrRequest) (model.CsrRequest, error) {
	log.Println(r)
	err := r.db.WithContext(ctx).Save(&request).Error
	log.Println(err)
	if err != nil {
		return model.CsrRequest{}, err
	}
	return request, nil
}

func (r *CertificateRepository) GetCsrRequest(ctx context.Context, ID int) (*model.CsrRequest, error) {
	var request model.CsrRequest
	err := r.db.WithContext(ctx).First(&request, "id = ?", ID).Error
	return &request, err
}

func (r *CertificateRepository) UpdateCsrStatus(ctx context.Context, ID int, status string) error {
	var request model.CsrRequest
	err := r.db.WithContext(ctx).First(&request, "id = ?", ID).Error
	if err != nil {
		return err
	}

	request.Status = status
	return r.db.WithContext(ctx).Save(&request).Error
}

func (r *CertificateRepository) GetCsrRequests(ctx context.Context) ([]model.CsrRequest, error) {
	requests := []model.CsrRequest{}
	err := r.db.WithContext(ctx).Find(&requests).Error

	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *CertificateRepository) SaveCertificateInfo(ctx context.Context, info model.CertificateInfo) error {
	err := r.db.WithContext(ctx).Save(&info).Error
	return err
}

func (r *CertificateRepository) GetCertificateInfo(ctx context.Context, csrID int) (*model.CertificateInfo, error) {
	var info model.CertificateInfo
	err := r.db.WithContext(ctx).First(&info, "csr_id = ?", csrID).Error
	return &info, err
}

func (r *CertificateRepository) GetActiveCertificateInfoByAgentID(ctx context.Context, agentID string) (*model.CertificateInfo, error) {
	var info model.CertificateInfo
	err := r.db.WithContext(ctx).
		Where("id = ? AND is_active = ?", agentID, true).
		Order("csr_id DESC").
		First(&info).Error
	return &info, err
}
