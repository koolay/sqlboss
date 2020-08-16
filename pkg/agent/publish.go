package agent

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/pkg/errors"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/koolay/sqlboss/pkg/conf"
)

type Publisher struct {
	cfg    *conf.Config
	pubSub *gochannel.GoChannel
}

func NewPublisher(cfg *conf.Config, pubSub *gochannel.GoChannel) (*Publisher, error) {
	return &Publisher{
		pubSub: pubSub,
		cfg:    cfg,
	}, nil
}

func (p *Publisher) Publish(topic string, data []byte) error {
	msg := message.NewMessage(watermill.NewUUID(), data)
	if err := p.pubSub.Publish(topic, msg); err != nil {
		return errors.Wrap(err, "failed to publish message")
	}

	return nil
}
