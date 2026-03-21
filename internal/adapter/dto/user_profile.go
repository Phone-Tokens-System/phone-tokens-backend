package dto

import (
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
	Token     string    `json:"token"`
}

func ToUserProfile(u *UserProfileToken) *UserProfile {
	return &UserProfile{
		BirthDate: u.UserProfile.BirthDate,
		Age:       u.UserProfile.Age,
		Gender:    u.UserProfile.Gender,
		Country:   u.UserProfile.Country,
		Region:    u.UserProfile.Region,
		City:      u.UserProfile.City,
		Education: u.UserProfile.Education,
		Token:     u.Token,
	}
}
