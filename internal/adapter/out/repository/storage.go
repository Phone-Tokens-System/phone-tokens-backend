package repository

import (
	"context"
	model2 "phone-tokens/internal/certificates/model"

	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

func (r *Storage) SaveCsrRequest(ctx context.Context, request model2.CsrRequest) (model2.CsrRequest, error) {
	err := r.db.WithContext(ctx).Save(&request).Error
	if err != nil {
		return model2.CsrRequest{}, err
	}
	return request, nil
}

func (r *Storage) GetCsrRequest(ctx context.Context, ID int) (*model2.CsrRequest, error) {
	var request model2.CsrRequest
	err := r.db.WithContext(ctx).First(&request, "id = ?", ID).Error
	return &request, err
}

func (r *Storage) UpdateCsrStatus(ctx context.Context, ID int, status string) error {
	var request model2.CsrRequest
	err := r.db.WithContext(ctx).First(&request, "id = ?", ID).Error
	if err != nil {
		return err
	}

	request.Status = status
	return r.db.WithContext(ctx).Save(&request).Error
}

func (r *Storage) GetCsrRequests(ctx context.Context) ([]model2.CsrRequest, error) {
	requests := []model2.CsrRequest{}
	err := r.db.WithContext(ctx).Find(&requests).Error

	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *Storage) SaveAgentInfo(ctx context.Context, info model2.ExternalAgentInfo) error {
	err := r.db.WithContext(ctx).Save(&info).Error
	return err
}

func (r *Storage) GetAgentInfo(ctx context.Context, csrID int) (*model2.ExternalAgentInfo, error) {
	var info model2.ExternalAgentInfo
	err := r.db.WithContext(ctx).First(&info, "csr_id = ?", csrID).Error
	return &info, err
}
