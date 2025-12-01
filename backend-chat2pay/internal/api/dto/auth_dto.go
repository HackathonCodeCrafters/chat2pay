package dto

import "chat2pay/internal/entities"

// User is the presenter object which will be passed in the response by Handler
type User struct {
	ID          string `json:"id,omitempty"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

func AuthSucces(data *entities.User, token string) User {
	user := User{
		ID:          data.ID,
		Email:       data.Email,
		AccessToken: token,
	}
	return user
}
