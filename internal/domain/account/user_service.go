package account

import "context"

type (
	UserService interface {
		Create(ctx context.Context, data UserRegisterDto) (string, error)
		IsExists(ctx context.Context, email string) (bool, error)
	}
)
