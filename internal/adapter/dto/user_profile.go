package dto

import (
	"phone-tokens/internal/model"
	"time"
)

type UserProfile struct {
	BirthDate time.Time `json:"birth_date"`
	Age       int       `json:"age" `
	Gender    string    `json:"gender" `
	Country   string    `json:"country"`
	Region    string    `json:"region"`
	City      string    `json:"city" `
	Education string    `json:"education"`
}

func ToUserProfile(userProfile *model.UserProfile) *UserProfile {
	return &UserProfile{
		BirthDate: userProfile.BirthDate,
		Age:       userProfile.Age,
		Gender:    userProfile.Gender,
		Country:   userProfile.Country,
		Region:    userProfile.Region,
		City:      userProfile.City,
		Education: userProfile.Education,
	}
}
