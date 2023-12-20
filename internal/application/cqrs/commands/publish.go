package commands

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"go.uber.org/zap"
)

func Publish(ctx context.Context, bus *cqrs.CommandBus, correlationID string, cmd interface{}, logger *zap.Logger) error {
	err := bus.Send(ctx, cmd)
	if err != nil {
		logger.Error(
			"An error occurred while sending the command",
			zap.String("command_id", correlationID),
			zap.Error(err),
		)
		return err
	}
	return nil
}
