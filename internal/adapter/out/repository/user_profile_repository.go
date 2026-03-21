package repository

import (
	"context"
	"fmt"
	"phone-tokens/internal/adapter/dto"
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
func (r *UserProfileRepository) FilterUserProfiles(ctx context.Context, filters map[string]string) ([]model.UserProfile, error) {
	query := r.db.WithContext(ctx).Model(&model.UserProfile{})

	allowed := map[string]bool{
		"gender":    true,
		"country":   true,
		"region":    true,
		"city":      true,
		"education": true,
	}

	for key, value := range filters {
		if !allowed[key] {
			continue
		}
		if key == "age_from" {
			query = query.Where("age >= ?", value)
		}
		if key == "age_to" {
			query = query.Where("age <= ?", value)
		}
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	if val, ok := filters["age_from"]; ok {
		query = query.Where("age >= ?", val)
	}
	if val, ok := filters["age_to"]; ok {
		query = query.Where("age <= ?", val)
	}

	var users []model.UserProfile
	err := query.Find(&users).Error
	return users, err
}

func (r *UserProfileRepository) FilterUserProfilesForAgent(
	ctx context.Context,
	filters map[string]string,
	agentID string,
) ([]dto.UserProfileToken, error) {

	type result struct {
		model.UserProfile
		Token string
	}

	query := r.db.WithContext(ctx).
		Table("user_profile").
		Select("user_profile.*, user_tokens.token").
		Joins("JOIN user_tokens ON user_tokens.user_id = user_profile.user_id").
		Where("user_tokens.agent_id = ?", agentID)

	allowed := map[string]bool{
		"gender":    true,
		"country":   true,
		"region":    true,
		"city":      true,
		"education": true,
	}

	for key, value := range filters {
		if key == "age_from" {
			query = query.Where("user_profile.age >= ?", value)
			continue
		}
		if key == "age_to" {
			query = query.Where("user_profile.age <= ?", value)
			continue
		}

		if allowed[key] {
			query = query.Where(fmt.Sprintf("user_profile.%s = ?", key), value)
		}
	}

	var rows []result
	err := query.Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	// мапим
	out := make([]dto.UserProfileToken, 0, len(rows))
	for _, r := range rows {
		out = append(out, dto.UserProfileToken{
			UserProfile: r.UserProfile,
			Token:       r.Token,
		})
	}

	return out, nil
}

func (r *UserProfileRepository) GetUserProfilesForAgent(ctx context.Context, agentID string) ([]dto.UserProfileToken, error) {
	type result struct {
		model.UserProfile
		Token string
	}

	query := r.db.WithContext(ctx).
		Table("user_profile").
		Select("user_profile.*, user_tokens.token").
		Joins("JOIN user_tokens ON user_tokens.user_id = user_profile.user_id").
		Where("user_tokens.agent_id = ?", agentID)

	var rows []result
	err := query.Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	// мапим
	out := make([]dto.UserProfileToken, 0, len(rows))
	for _, r := range rows {
		out = append(out, dto.UserProfileToken{
			UserProfile: r.UserProfile,
			Token:       r.Token,
		})
	}

	return out, nil

	//query := r.db.WithContext(ctx).
	//	Model(&model.UserProfile{}).
	//	Joins("JOIN user_tokens ON user_tokens.user_id = user_profile.user_id").
	//	Where("user_tokens.agent_id = ?", agentID)
	//
	//var users []model.UserProfile
	//err := query.Find(&users).Error
	//return users, err
}
