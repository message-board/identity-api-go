package command

import "github.com/google/uuid"

type CreateUserCommand struct {
	Id           uuid.UUID
	EmailAddress string
	Password     string
}

func NewCreateUserCommand(id uuid.UUID, emailAddress string, password string) CreateUserCommand {
	return CreateUserCommand{
		Id:           id,
		EmailAddress: emailAddress,
		Password:     password,
	}
}
