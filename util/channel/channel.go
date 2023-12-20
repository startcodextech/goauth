package channel

import (
	"errors"
	"github.com/startcodextech/goauth/proto"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	Channels = Channel{
		result: make(map[string]chan ResultChannel),
		mutex:  &sync.Mutex{},
	}
)

type (
	Channel struct {
		result map[string]chan ResultChannel
		mutex  *sync.Mutex
	}

	ResultChannel struct {
		CorrelationID string
		Failed        *proto.EventError
		Success       any
	}

	ChannelHandler struct {
		OnSuccess SuccessFunc
		OnFailed  FailedFunc
		Logger    *zap.Logger
	}

	SuccessFunc func(interface{})
	FailedFunc  func(eventError *proto.EventError)
)

func (r ResultChannel) IsSuccess() bool {
	return r.Failed == nil && r.Success != nil
}

func (c Channel) AddChannel(correlationID string) {
	c.mutex.Lock()
	c.result[correlationID] = make(chan ResultChannel, 1)
	c.mutex.Unlock()
}

func (c Channel) GetChannel(correlationID string) chan ResultChannel {
	channel, exist := c.result[correlationID]
	if !exist {
		return nil
	}

	return channel
}

func (c Channel) GetResult(correlationID string, callback ChannelHandler) error {
	var err error

	channel := c.GetChannel(correlationID)
	if channel == nil {
		return errors.New("")
	}

	select {
	case result := <-channel:
		if result.IsSuccess() {
			callback.OnSuccess(result.Success)
		} else {
			callback.OnFailed(result.Failed)
		}
	case <-time.After(15 * time.Second):
		err = errors.New("waiting time for result 15 sec")
		callback.Logger.Error(
			"An error occurred while waiting for a response.",
			zap.String("correlation_id", correlationID),
			zap.Error(err),
		)
	}

	return err
}

func (c Channel) Success(correlationID string, event interface{}) {
	c.mutex.Lock()
	if resultChan, exists := c.result[correlationID]; exists {
		resultChan <- ResultChannel{
			CorrelationID: correlationID,
			Success:       event,
		}
		close(resultChan)
		delete(c.result, correlationID)
	}
	c.mutex.Unlock()
}

func (c Channel) Failed(correlationID string, event *proto.EventError) {
	c.mutex.Lock()
	if resultChan, exists := c.result[correlationID]; exists {
		resultChan <- ResultChannel{
			CorrelationID: correlationID,
			Failed:        event,
		}
		close(resultChan)
		delete(c.result, correlationID)
	}
	c.mutex.Unlock()
}
