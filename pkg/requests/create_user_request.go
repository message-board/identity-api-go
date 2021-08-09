package requests

import (
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Id           uuid.UUID `json:"id"`
	EmailAddress string    `json:"emailAddress"`
	Password     string    `json:"password"`
}

func NewCreateUserRequest(emailAddress string, password string) CreateUserRequest {
	return CreateUserRequest{
		Id:           uuid.New(),
		EmailAddress: emailAddress,
		Password:     password,
	}
}
