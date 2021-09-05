package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	middleware "github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/go-chi/chi/v5"
	"github.com/message-board/identity-go/internal/interfaces/handlers/createuser"
	"github.com/message-board/identity-go/internal/interfaces/rest"
)

func main() {
	ctx := context.Background()

	logger := watermill.NewStdLogger(false, false)
	cqrsMarshaler := cqrs.JSONMarshaler{}

	ch := gochannel.NewGoChannel(gochannel.Config{Persistent: true}, logger)

	// CQRS is built on messages router. Detailed documentation: https://watermill.io/docs/messages-router/
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	// Simple middleware which will recover panics from event or command handlers.
	// More about router middlewares you can find in the documentation:
	// https://watermill.io/docs/messages-router/#middleware
	//
	// List of available middlewares you can find in message/router/middleware.
	router.AddMiddleware(middleware.Recoverer)

	// cqrs.Facade is facade for Command and Event buses and processors.
	// You can use facade, or create buses and processors manually (you can inspire with cqrs.NewFacade)
	cqrsFacade, err := cqrs.NewFacade(cqrs.FacadeConfig{
		GenerateCommandsTopic: func(commandName string) string {
			parts := strings.Split(commandName, ".")
			return "identity." + parts[len(parts)-1]
		},
		CommandHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.CommandHandler {
			return []cqrs.CommandHandler{
				createuser.NewCreateUserHandler(eb),
			}
		},
		CommandsPublisher: ch,
		CommandsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			// we can reuse subscriber, because all commands have separated topics
			return ch, nil
		},
		GenerateEventsTopic: func(eventName string) string {
			parts := strings.Split(eventName, ".")
			return "identity." + parts[len(parts)-1]
		},
		EventHandlers: nil,
		// EventHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.EventHandler {
		// 	return []cqrs.EventHandler{}
		// },
		EventsPublisher: ch,
		EventsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			return ch, nil // <-- Here
		},
		Router:                router,
		CommandEventMarshaler: cqrsMarshaler,
		Logger:                logger,
	})
	if err != nil {
		panic(err)
	}

	// application := service.NewApplication(ctx)

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		rest.RunServer(func(router chi.Router) http.Handler {
			return rest.HandlerFromMux(
				rest.NewUserResource(cqrsFacade.CommandBus()),
				router,
			)
		})
	case "grpc":
	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}

	// processors are based on router, so they will work when router will start
	if err := router.Run(ctx); err != nil {
		panic(err)
	}
}
