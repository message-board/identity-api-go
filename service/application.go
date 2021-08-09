package service

import (
	"context"

	"github.com/message-board/identity/app"
	"github.com/message-board/identity/app/command"
)

func NewApplication(ctx context.Context) app.Application {
	return app.Application{
		Commands: app.Commands{
			CreateUserCommandHandler: command.NewCreateUserCommandHandler(),
		},
	}
}
