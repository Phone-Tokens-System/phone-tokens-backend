package repository

import (
	"context"
	"phone-tokens/internal/model"

	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

func (r *Storage) SaveCsrRequest(ctx context.Context, request model.CsrRequest) (model.CsrRequest, error) {
	err := r.db.WithContext(ctx).Save(&request).Error
	if err != nil {
		return model.CsrRequest{}, err
	}
	return request, nil
}

func (r *Storage) GetCsrRequest(ctx context.Context, ID string) (*model.CsrRequest, error) {
	var request model.CsrRequest
	err := r.db.WithContext(ctx).First(&request, "id = ?", ID).Error
	return &request, err
}

func (r *Storage) UpdateCsrStatus(ctx context.Context, ID string, status string) error {
	var request model.CsrRequest
	err := r.db.WithContext(ctx).First(&request, "id = ?", ID).Error
	if err != nil {
		return err
	}

	request.Status = status
	return r.db.WithContext(ctx).Save(&request).Error
}

func (r *Storage) GetCsrRequests(ctx context.Context) ([]model.CsrRequest, error) {
	requests := []model.CsrRequest{}
	err := r.db.WithContext(ctx).Find(&requests).Error

	if err != nil {
		return nil, err
	}

	return requests, nil
}
