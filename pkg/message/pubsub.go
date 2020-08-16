package message

import (
	"sync"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

var (
	pubSub     *gochannel.GoChannel
	pubSubOnce sync.Once
)

func NewPubSub() *gochannel.GoChannel {
	pubSubOnce.Do(func() {
		pubSub = gochannel.NewGoChannel(
			gochannel.Config{},
			watermill.NewStdLogger(false, false),
		)
	})

	return pubSub
}
