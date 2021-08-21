package main

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type Subscriber struct {
	config gochannel.Config
}

func NewSubscriber(config gochannel.Config, logger watermill.LoggerAdapter) *Subscriber {
	return &Subscriber{config}
}

// Subscribe consumes messages from AMQP broker.
//
// Watermill's topic in Subscribe is not mapped to AMQP's topic, but depending on configuration it can be mapped
// to exchange, queue or routing key.
// For detailed description of nomenclature mapping, please check "Nomenclature" paragraph in doc.go file.
func (s *Subscriber) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	return nil, nil
}

func (s *Subscriber) Close() error {
	return nil
}
