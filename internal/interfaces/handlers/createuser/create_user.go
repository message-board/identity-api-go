package createuser

import "github.com/google/uuid"

type CreateUser struct {
	Id           uuid.UUID `json:"id"`
	EmailAddress string    `json:"emailAddress"`
	Password     string    `json:"password"`
}

func NewCreateUser(id uuid.UUID, emailAddress, password string) CreateUser {
	return CreateUser{
		Id:           id,
		EmailAddress: emailAddress,
		Password:     password,
	}
}
