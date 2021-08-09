package command

import "context"

type CreateUserCommandHandler struct {
}

func NewCreateUserCommandHandler() CreateUserCommandHandler {
	return CreateUserCommandHandler{}
}

func (c CreateUserCommandHandler) Handle(ctx context.Context, command CreateUserCommand) error {
	return nil
}
