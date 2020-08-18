package message

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type CommandBus interface {
	Send(ctx context.Context, cmd interface{}) error
}

type EventBus interface {
	Publish(ctx context.Context, event interface{}) error
}

type CommandHandlerGenerator func(cb *cqrs.CommandBus, eb *cqrs.EventBus) cqrs.CommandHandler
type EventHandlerGenerator func(cb *cqrs.CommandBus, eb *cqrs.EventBus) cqrs.EventHandler

type CQRSServer struct {
	router             *message.Router
	commandsPublisher  message.Publisher
	commandsSubscriber message.Subscriber
	eventsPublisher    message.Publisher
	eventsSubscriber   message.Subscriber
	cqrsMarshaler      cqrs.CommandEventMarshaler

	commandBus CommandBus
	eventBus   EventBus
	nlog       watermill.LoggerAdapter
}

func NewCQRSServer(logger *logrus.Entry, cqrsMarshaler cqrs.CommandEventMarshaler) (*CQRSServer, error) {
	nlog := &NatsLog{logger: logger}
	// CQRS is built on messages router.
	// Detailed documentation: https://watermill.io/docs/messages-router/
	router, err := message.NewRouter(message.RouterConfig{}, nlog)
	if err != nil {
		return nil, err
	}

	commandsPubSub := gochannel.NewGoChannel(
		gochannel.Config{BlockPublishUntilSubscriberAck: true},
		nlog,
	)

	eventsPubSub := gochannel.NewGoChannel(
		gochannel.Config{BlockPublishUntilSubscriberAck: true},
		nlog,
	)

	return &CQRSServer{
		router:             router,
		commandsPublisher:  commandsPubSub,
		commandsSubscriber: commandsPubSub,
		eventsPublisher:    eventsPubSub,
		eventsSubscriber:   eventsPubSub,
		cqrsMarshaler:      cqrsMarshaler,
		nlog:               nlog,
	}, nil
}

func (s *CQRSServer) Setup(commandHandlerGenerators []CommandHandlerGenerator,
	eventHandlerGenerators []EventHandlerGenerator) error {

	// Simple middleware which will recover panics from event or command handlers.
	// More about router middlewares you can find in the documentation:
	// https://watermill.io/docs/messages-router/#middleware
	//
	// List of available middlewares you can find in message/router/middleware.
	s.router.AddMiddleware(middleware.Recoverer)

	// cqrs.Facade is facade for Command and Event buses and processors.
	// You can use facade, or create buses and processors manually (you can inspire with cqrs.NewFacade)
	cqrsFacade, err := cqrs.NewFacade(cqrs.FacadeConfig{
		GenerateCommandsTopic: func(commandName string) string {
			// we are using queue RabbitMQ config, so we need to have topic per command type
			return commandName
		},
		CommandHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.CommandHandler {
			handlers := []cqrs.CommandHandler{}
			for _, gen := range commandHandlerGenerators {
				handlers = append(handlers, gen(cb, eb))
			}

			return handlers
		},
		CommandsPublisher: s.commandsPublisher,
		CommandsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			// we can reuse subscriber, because all commands have separated topics
			return s.commandsSubscriber, nil
		},
		GenerateEventsTopic: func(eventName string) string {
			// because we are using PubSub RabbitMQ config, we can use one topic for all events
			return "events"

			// we can also use topic per event type
			// return eventName
		},
		EventHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.EventHandler {
			handlers := []cqrs.EventHandler{}
			for _, gen := range eventHandlerGenerators {
				handlers = append(handlers, gen(cb, eb))
			}

			return handlers
		},
		EventsPublisher: s.eventsPublisher,
		EventsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			return s.eventsSubscriber, nil
		},
		Router:                s.router,
		CommandEventMarshaler: s.cqrsMarshaler,
		Logger:                s.nlog,
	})

	if err != nil {
		return err
	}

	s.commandBus = cqrsFacade.CommandBus()
	s.eventBus = cqrsFacade.EventBus()

	return nil
}

func (s *CQRSServer) GetCommandBus() CommandBus {
	return s.commandBus
}

func (s *CQRSServer) GetEventBus() EventBus {
	return s.eventBus
}

func (s *CQRSServer) Start() error {
	if err := s.router.Run(context.Background()); err != nil {
		return errors.Wrap(err, "failed to run cqrs server")
	}
	return nil
}
