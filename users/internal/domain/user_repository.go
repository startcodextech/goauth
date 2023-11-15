package domain

import "context"

type UserRepository interface {
	Load(ctx context.Context, userID string) (user *User, err error)
	Save(ctx context.Context, user *User) (err error)
}
