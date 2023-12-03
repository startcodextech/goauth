package account

import (
	"context"
)

type (
	// UserRepository represents the user's persistence contract.
	// This contract is used to persist and retrieve user's data.
	UserRepository interface {
		// Save function persists a user.
		Save(ctx context.Context, user User)

		// Delete function removes a user.
		Delete(ctx context.Context, ID string)

		// Find function retrieves a user.
		Find(ctx context.Context, criterial map[string]interface{}) ([]User, error)

		// Commit function persists all changes.
		// This function is used to persist all changes in the current context.
		Commit(ctx context.Context) error
	}
)
