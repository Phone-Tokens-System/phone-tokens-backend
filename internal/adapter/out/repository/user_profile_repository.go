package repository

import (
	"context"
	"fmt"
	"phone-tokens/internal/model"

	"gorm.io/gorm"
)

type UserProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) *UserProfileRepository {
	return &UserProfileRepository{db: db}
}

func (r *UserProfileRepository) SaveProfile(ctx context.Context, userProfile model.UserProfile) error {
	return r.db.WithContext(ctx).Save(userProfile).Error
}

func (r *UserProfileRepository) UpdateProfileByUserId(ctx context.Context, userId string, userProfile *model.UserProfile) error {
	return r.db.WithContext(ctx).
		Model(&model.UserProfile{}).
		Where("user_id = ?", userId).
		Updates(userProfile).
		Error
}

func (r *UserProfileRepository) GetProfileByUserId(ctx context.Context, userId string) (*model.UserProfile, error) {
	var userProfile model.UserProfile
	err := r.db.WithContext(ctx).Find(&userProfile, "user_id = ?", userId).Error
	return &userProfile, err
}

func (r *UserProfileRepository) DeleteProfile(ctx context.Context, userId string) error {
	err := r.db.WithContext(ctx).Delete(&model.UserProfile{}, "user_id = ?", userId).Error
	return err
}

// TODO: filter age. normalize gender values at least.
// maybe countries region city and education too make them a list
func (r *UserProfileRepository) FilterUserProfiles(ctx context.Context, filterName, filterValue string) ([]model.UserProfile, error) {
	allowed := map[string]bool{
		"gender":    true,
		"country":   true,
		"region":    true,
		"city":      true,
		"education": true,
	}

	if !allowed[filterName] {
		return nil, fmt.Errorf("invalid filter field")
	}

	var userProfiles []model.UserProfile
	filter := fmt.Sprintf("%s = ?", filterName)
	err := r.db.WithContext(ctx).Find(&userProfiles, filter, filterValue).Error
	if err != nil {
		return nil, err
	}
	return userProfiles, nil
}
