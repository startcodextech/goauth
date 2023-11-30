package account

import "context"

type (
	UserService interface {
		Create(ctx context.Context, data UserRegisterDto) error
	}
)
