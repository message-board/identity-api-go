package usercreated

import (
	"context"
	"log"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/message-board/identity-go/pkg/events"
)

type UserCreatedHandler struct {
	commandBus *cqrs.CommandBus
}

func NewUserCreatedHandler(cb *cqrs.CommandBus) UserCreatedHandler {
	return UserCreatedHandler{
		commandBus: cb,
	}
}

func (h UserCreatedHandler) HandlerName() string {
	// this name is passed to EventsSubscriberConstructor and used to generate queue name
	return "UserCreatedHandler"
}

func (h UserCreatedHandler) NewEvent() interface{} {
	return &events.UserCreated{}
}

func (h UserCreatedHandler) Handle(ctx context.Context, e interface{}) error {
	event := e.(*events.UserCreated)

	log.Printf(
		"UserCreated event with id: %s, emailAddress: %s",
		event.Id,
		event.EmailAddress,
	)

	// return o.commandBus.Send(ctx, orderBeerCmd)
	return nil
}
