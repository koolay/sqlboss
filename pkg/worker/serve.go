package worker

import (
	"context"
	"log"

	"github.com/koolay/sqlboss/pkg/conf"
	"github.com/sirupsen/logrus"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type Worker struct {
	cfg    *conf.Config
	pubsub *gochannel.GoChannel
	logger *logrus.Logger
}

func NewWorker(cfg *conf.Config, pubsub *gochannel.GoChannel, logger *logrus.Logger) *Worker {
	return &Worker{
		cfg:    cfg,
		logger: logger,
		pubsub: pubsub,
	}
}

func (w *Worker) Setup() error {
	return nil
}

func (w *Worker) Run(ctx context.Context) error {
	log.Println("start consume", w.cfg.Stream.Topic)
	messages, err := w.pubsub.Subscribe(ctx, w.cfg.Stream.Topic)
	if err != nil {
		return err
	}

	go w.process(messages)
	return nil
}

func (w *Worker) process(messages <-chan *message.Message) {
	for msg := range messages {
		log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
