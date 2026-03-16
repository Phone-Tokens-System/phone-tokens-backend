package users

import (
	"context"
	"phone-tokens/internal/adapter/out/repository"
	"phone-tokens/internal/model"
)

type UserProfileService struct {
	repo            Repository
	userProfileRepo repository.UserProfileRepository
	tokenRepo       repository.TokenRepository
}

func NewUserProfileService(repo Repository, userProfileRepo repository.UserProfileRepository,
	tokenRepo repository.TokenRepository) *UserProfileService {
	return &UserProfileService{repo, userProfileRepo, tokenRepo}
}

func (s *UserProfileService) SaveUserProfile(ctx context.Context, userProfile model.UserProfile) error {
	return s.userProfileRepo.SaveProfile(ctx, userProfile)
}

func (s *UserProfileService) GetUserProfileByToken(ctx context.Context, tokenName string) (*model.UserProfile, error) {
	token, err := s.tokenRepo.GetTokenByToken(ctx, tokenName)
	if err != nil {
		return nil, err
	}

	userProfile, err := s.userProfileRepo.GetProfileByUserId(ctx, token.UserID)
	if err != nil {
		return nil, err
	}
	return userProfile, nil
}

func (s *UserProfileService) GetUserProfileByUserId(ctx context.Context, userId string) (*model.UserProfile, error) {
	userProfile, err := s.userProfileRepo.GetProfileByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return userProfile, nil
}

func (s *UserProfileService) DeleteUserProfile(ctx context.Context, userId string) error {
	return s.userProfileRepo.DeleteProfile(ctx, userId)
}

func (s *UserProfileService) UpdateUserProfile(ctx context.Context, userId string, userProfile *model.UserProfile) error {
	return s.userProfileRepo.UpdateProfileByUserId(ctx, userId, userProfile)
}

func (s *UserProfileService) FilterUserProfiles(ctx context.Context, filterName, filterValue string) ([]model.UserProfile, error) {
	return nil, nil
}
