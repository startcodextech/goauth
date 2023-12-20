package channel

import (
	"errors"
	"github.com/startcodextech/goauth/proto"
	"go.uber.org/zap"
	"sync"
	"time"
)

var ResultMap = make(map[string]chan ResultChannel)
var MapMutex = &sync.Mutex{}

type (
	ResultChannel struct {
		CorrelationID string
		Failed        *proto.EventError
		Success       any
	}

	ResultCallback struct {
		CorrelationID string
		OnSuccess     SuccessFunc
		OnFailed      FailedFunc
		Logger        *zap.Logger
	}

	SuccessFunc func(interface{})
	FailedFunc  func(eventError *proto.EventError)
)

func (r ResultChannel) IsSuccess() bool {
	return r.Failed == nil && r.Success != nil
}

func AddChannel(correlationID string, channel chan ResultChannel) {
	MapMutex.Lock()
	ResultMap[correlationID] = channel
	MapMutex.Unlock()
}

func GetResult(channel chan ResultChannel, callback ResultCallback) error {
	var err error
	select {
	case result := <-channel:
		if result.IsSuccess() {
			callback.OnSuccess(result.Success)
		} else {
			callback.OnFailed(result.Failed)
		}
	case <-time.After(15 * time.Second):
		err = errors.New("waiting time for result 15 sec")
	}

	if err != nil {
		callback.Logger.Error(
			"An error occurred while waiting for a response.",
			zap.String("correlation_id", callback.CorrelationID),
			zap.Error(err),
		)
	}

	return err
}

func FinishSuccess(correlationID string, event interface{}) {
	MapMutex.Lock()
	if resultChan, exists := ResultMap[correlationID]; exists {
		resultChan <- ResultChannel{
			CorrelationID: correlationID,
			Success:       event,
		}
		close(resultChan)
		delete(ResultMap, correlationID)
	}
	MapMutex.Unlock()
}

func FinishFailed(correlationID string, event *proto.EventError) {
	MapMutex.Lock()
	if resultChan, exists := ResultMap[correlationID]; exists {
		resultChan <- ResultChannel{
			CorrelationID: correlationID,
			Failed:        event,
		}
		close(resultChan)
		delete(ResultMap, correlationID)
	}
	MapMutex.Unlock()
}
