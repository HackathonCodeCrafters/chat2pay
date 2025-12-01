package dto

import "chat2pay/internal/entities"

// User is the presenter object which will be passed in the response by Handler
type UserProfile struct {
	ID           string `json:"id,omitempty"`
	Fullname     string `json:"fullname"`
	Email        string `json:"email"`
	Age          int    `json:"age"`
	MobileNumber string `json:"mobile_number"`
}

func GetProfileSuccess(data *entities.User) UserProfile {
	user := UserProfile{
		ID:           data.ID,
		Fullname:     data.FullName,
		Email:        data.Email,
		Age:          data.Age,
		MobileNumber: data.MobileNumber,
	}
	return user
}
