package dto

type UserProfileRequest struct {
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
	Country   string `json:"country"`
	Region    string `json:"region"`
	City      string `json:"city"`
	Education string `json:"education"`
}
