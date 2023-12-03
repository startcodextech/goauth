package gochannel

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

func New(logger watermill.LoggerAdapter) *gochannel.GoChannel {
	config := gochannel.Config{
		Persistent:                     false,
		BlockPublishUntilSubscriberAck: false,
	}
	return gochannel.NewGoChannel(config, logger)
}
