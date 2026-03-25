package users

import (
	"context"
	"fmt"
	"phone-tokens/internal/adapter/dto"
	"phone-tokens/internal/adapter/out/repository"
	"phone-tokens/internal/model"
	"phone-tokens/internal/service/tokens"
	"time"
)

type UserProfileService struct {
	repo            Repository
	userProfileRepo *repository.UserProfileRepository
	tokenRepo       tokens.Repository
}

func NewUserProfileService(repo Repository, userProfileRepo *repository.UserProfileRepository,
	tokenRepo tokens.Repository) *UserProfileService {
	return &UserProfileService{repo, userProfileRepo, tokenRepo}
}

func (s *UserProfileService) SaveUserProfile(ctx context.Context, userProfile model.UserProfile) error {
	if !userProfile.BirthDate.IsZero() {
		userProfile.Age = calculateAge(userProfile.BirthDate)
	}

	return s.userProfileRepo.SaveProfile(ctx, userProfile)
}

func (s *UserProfileService) GetUserProfileByToken(ctx context.Context, tokenName string) (*dto.UserProfileToken, error) {
	token, err := s.tokenRepo.GetTokenByToken(ctx, tokenName)
	if err != nil {
		return nil, err
	}

	userProfile, err := s.userProfileRepo.GetProfileByUserId(ctx, token.UserID)
	if err != nil {
		return nil, err
	}
	return &dto.UserProfileToken{Token: tokenName, UserProfile: *userProfile}, nil
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

func (s *UserProfileService) FilterUserProfiles(ctx context.Context, filter dto.FilterRequest) ([]model.UserProfile, error) {
	return s.userProfileRepo.FilterUserProfiles(ctx, filter.Filters)
}

func (s *UserProfileService) FilterUserProfilesByAgentId(ctx context.Context, filter dto.FilterRequest, agentID string) ([]dto.UserProfileToken, error) {
	return s.userProfileRepo.FilterUserProfilesForAgent(ctx, filter.Filters, agentID)
}

func (s *UserProfileService) GetFilters() dto.FilterResponse {
	return dto.FilterResponse{
		Filters: []dto.Filter{
			{Key: "Gender", Type: "select", Options: []string{"male", "female"}},
			{Key: "Country", Type: "select", OptionSource: "api/v1/dictionary/countries"},
			{Key: "Region", Type: "select", OptionSource: "api/v1/dictionary/regions/?country={country}"},
			{Key: "City", Type: "select", OptionSource: "api/v1/dictionary/cities/?country={country}&region={region}"},
			{Key: "education", Type: "select", Options: []string{"school", "bachelor", "master", "phd"}},
			{Key: "age_from", Type: "value"},
			{Key: "age_to", Type: "value"},
		},
	}
}

func (s *UserProfileService) GetFilteredTokensForAgent(ctx context.Context, req dto.FilterRequest, agentID string) ([]model.UserToken, error) {
	profiles, err := s.FilterUserProfilesByAgentId(ctx, req, agentID)
	if err != nil {
		return nil, err
	}
	tokens := make([]model.UserToken, 0)

	for _, profile := range profiles {
		token, err := s.tokenRepo.GetTokensByUserIdAndAgentId(ctx, profile.UserProfile.UserId, agentID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tokens = append(tokens, token...)
	}
	return tokens, nil
}

func (s *UserProfileService) GetUserProfilesForAgent(ctx context.Context, agentID string) ([]dto.UserProfileToken, error) {
	return s.userProfileRepo.GetUserProfilesForAgent(ctx, agentID)
}

func calculateAge(birthDate time.Time) int {
	timeNow := time.Now()
	age := timeNow.Year() - birthDate.Year()
	if timeNow.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}
