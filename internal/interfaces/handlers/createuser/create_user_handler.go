package createuser

import (
	"context"
	"log"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/message-board/identity-go/pkg/events"
)

type CreateUserHandler struct {
	eventBus *cqrs.EventBus
}

func NewCreateUserHandler(eb *cqrs.EventBus) CreateUserHandler {
	return CreateUserHandler{
		eventBus: eb,
	}
}

func (h CreateUserHandler) HandlerName() string {
	return "CreateUserHandler"
}

// NewCommand returns type of command which this handle should handle. It must be a pointer.
func (h CreateUserHandler) NewCommand() interface{} {
	return &CreateUser{}
}

func (h CreateUserHandler) Handle(ctx context.Context, c interface{}) error {
	// c is always the type returned by `NewCommand`, so casting is always safe
	cmd := c.(*CreateUser)

	log.Printf(
		"Create user command with %s, %s and %s",
		cmd.Id,
		cmd.EmailAddress,
		cmd.Password,
	)

	// Publish UserCreated event
	if err := h.eventBus.Publish(ctx, &events.UserCreated{
		Id:           cmd.Id,
		EmailAddress: cmd.EmailAddress,
	}); err != nil {
		return err
	}

	return nil
}
