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
	routerMiddleware "github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/go-chi/chi/v5"
	"github.com/message-board/identity-go/internal/interfaces/handlers/createuser"
	"github.com/message-board/identity-go/internal/interfaces/rest"
)

func main() {
	ctx := context.Background()

	logger := watermill.NewStdLogger(false, false)
	cqrsMarshaler := cqrs.JSONMarshaler{}

	commandsConfig := gochannel.Config{}
	commandsPublisher := gochannel.NewGoChannel(commandsConfig, logger)
	commandsSubscriber := NewSubscriber(commandsConfig, logger)

	// You can use any Pub/Sub implementation from here: https://watermill.io/docs/pub-sub-implementations/
	// Detailed RabbitMQ implementation: https://watermill.io/docs/pub-sub-implementations/#rabbitmq-amqp
	// Commands will be send to queue, because they need to be consumed once.
	// commandsAMQPConfig := amqp.NewDurableQueueConfig(amqpAddress)
	// commandsPublisher, err := amqp.NewPublisher(commandsAMQPConfig, logger)
	// if err != nil {
	// 	panic(err)
	// }
	// commandsSubscriber, err := amqp.NewSubscriber(commandsAMQPConfig, logger)
	// if err != nil {
	// 	panic(err)
	// }

	// Events will be published to PubSub configured Rabbit, because they may be consumed by multiple consumers.
	// (in that case BookingsFinancialReport and OrderBeerOnRoomBooked).
	eventsPublisher := gochannel.NewGoChannel(gochannel.Config{}, logger)
	// eventsPublisher, err := amqp.NewPublisher(amqp.NewDurablePubSubConfig(amqpAddress, nil), logger)
	// if err != nil {
	// 	panic(err)
	// }

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
	router.AddMiddleware(routerMiddleware.Recoverer)

	// cqrs.Facade is facade for Command and Event buses and processors.
	// You can use facade, or create buses and processors manually (you can inspire with cqrs.NewFacade)
	cqrsFacade, err := cqrs.NewFacade(cqrs.FacadeConfig{
		GenerateCommandsTopic: func(commandName string) string {
			// we are using queue RabbitMQ config, so we need to have topic per command type
			return commandName
		},
		CommandHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.CommandHandler {
			return []cqrs.CommandHandler{
				createuser.NewCreateUserHandler(eb),
			}
		},
		CommandsPublisher: commandsPublisher,
		CommandsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			// we can reuse subscriber, because all commands have separated topics
			return commandsSubscriber, nil
		},
		GenerateEventsTopic: func(eventName string) string {
			// because we are using PubSub RabbitMQ config, we can use one topic for all events
			return "events"

			// we can also use topic per event type
			// return eventName
		},
		EventHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.EventHandler {
			return []cqrs.EventHandler{}
		},
		EventsPublisher: eventsPublisher,
		EventsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			config := gochannel.Config{}
			return NewSubscriber(config, logger), nil
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
