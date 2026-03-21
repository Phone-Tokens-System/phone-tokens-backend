package repository

import (
	"context"
	"phone-tokens/internal/model"

	"gorm.io/gorm"
)

type AgentPackageRepository struct {
	db *gorm.DB
}

func NewAgentPackageRepository(db *gorm.DB) *AgentPackageRepository {
	return &AgentPackageRepository{db: db}
}

func (r *AgentPackageRepository) AddAgentPackage(ctx context.Context, req *model.AgentPackages) error {
	return r.db.WithContext(ctx).Save(req).Error
}

func (r *AgentPackageRepository) RemoveAgentPackage(ctx context.Context, req *model.AgentPackages) error {
	return r.db.WithContext(ctx).Delete(req).Error
}

func (r *AgentPackageRepository) GetAgentPackageById(ctx context.Context, id int) (*model.AgentPackages, error) {
	var pkg model.AgentPackages
	err := r.db.WithContext(ctx).Find(&pkg, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *AgentPackageRepository) GetAgentPackagesByAgentId(ctx context.Context, agentId string) ([]model.AgentPackages, error) {
	var pkg []model.AgentPackages
	err := r.db.WithContext(ctx).Find(&pkg, "agent_id = ?", agentId).Error
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

func (r *AgentPackageRepository) UpdateAgentPackage(ctx context.Context, id int, pkg *model.AgentPackages) (*model.AgentPackages, error) {
	err := r.db.WithContext(ctx).
		Model(&model.AgentPackages{}).
		Where("id = ?", id).
		Updates(pkg).
		Error
	if err != nil {
		return nil, err
	}

	var updated model.AgentPackages
	if err := r.db.WithContext(ctx).First(&updated, id).Error; err != nil {
		return nil, err
	}

	return &updated, nil
}
