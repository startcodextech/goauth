package account

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/startcodextech/goauth/internal/domain/account"
)

func CreateUserCommandHandler(msg *message.Message) ([]message.Message, error) {
	var payload account.UserCreateDto

	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return nil, err
	}

	payload.ID = msg.UUID

	user := account.NewUser()
	err := user.Create(payload)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
