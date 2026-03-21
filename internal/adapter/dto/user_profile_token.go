package dto

import "phone-tokens/internal/model"

type UserProfileToken struct {
	UserProfile model.UserProfile
	Token       string
}
