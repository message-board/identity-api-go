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
	"github.com/message-board/identity-go/internal/interfaces/handlers/usercreated"
	"github.com/message-board/identity-go/internal/interfaces/rest"
)

// @title Message Board Identity Api
// @version 1.0
// @description TODO.
// @termsOfService http://todo.io/terms

// @contact.name TODO
// @contact.url http://todo.io/support
// @contact.email support@todo.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v1
func main() {
	// ctx := context.Background()

	logger := watermill.NewStdLogger(false, false)
	cqrsMarshaler := cqrs.JSONMarshaler{}

	// Since we are using the go channel implementation we could remove commandsPublisher, commandsSubscriber and  eventsPublisher, to be simple.
	// And then we need to replace the commandsPublisher, commandsSubscriber and  eventsPublisher with the channel 
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
			return commandName
		},
		CommandHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.CommandHandler {
			return []cqrs.CommandHandler{
				createuser.NewCreateUserHandler(eb),
			}
		},
		CommandsPublisher: ch, // <-- Here
		CommandsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			return ch, nil // <-- Here  
		},
		GenerateEventsTopic: func(eventName string) string {
			return "events"
		},
		EventHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.EventHandler {
			return []cqrs.EventHandler{
				usercreated.NewUserCreatedHandler(cb),
			}
		},
		EventsPublisher: ch, // <-- Here
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
	if err := router.Run(context.Background()); err != nil {
		panic(err)
	}
}
