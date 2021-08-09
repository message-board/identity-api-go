package app

import "github.com/message-board/identity/app/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateUserCommandHandler command.CreateUserCommandHandler
}

type Queries struct {
}
