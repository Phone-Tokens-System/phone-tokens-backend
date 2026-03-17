package repository

import (
	"context"
	"phone-tokens/internal/model"

	"gorm.io/gorm"
)

type UsageRepository struct {
	db *gorm.DB
}

func NewUsageRepository(db *gorm.DB) *UsageRepository {
	return &UsageRepository{db: db}
}

func (r *UsageRepository) SaveUsage(ctx context.Context, tx *gorm.DB, usage *model.Usage) error {
	return tx.Create(usage).Error
}

func (r *UsageRepository) GetUsageByAgentID(ctx context.Context, agentID string) ([]model.Usage, error) {
	var usage []model.Usage

	if err := r.db.WithContext(ctx).Find(&usage, "agent_id = ?", agentID).Error; err != nil {
		return nil, err
	}
	return usage, nil
}

func (r *UsageRepository) GetUsageByPhoneNumber(ctx context.Context, phoneNumber string) ([]model.Usage, error) {
	var usage []model.Usage
	if err := r.db.WithContext(ctx).Find(&usage, "phone_number = ?", phoneNumber).Error; err != nil {
		return nil, err
	}
	return usage, nil
}

func (r *UsageRepository) DeleteUsage(ctx context.Context, ID string) (*model.Usage, error) {
	var usage model.Usage
	if err := r.db.WithContext(ctx).Delete(&usage, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &usage, nil
}
