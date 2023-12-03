package services

import (
	"context"
	"github.com/startcodextech/goauth/internal/application/services/account"
	"github.com/startcodextech/goauth/internal/infrastructure/persistence/mongodb"
)

// Services is a struct that holds all services
type Services struct {
	account account.AccountServices
}

// New creates new services
// It requires a mongodb driver to be injected.
func New(ctx context.Context, dbDrive mongodb.MongoDriver) Services {
	db := dbDrive.Database()

	return Services{
		account: account.New(mongodb.NewMongoUserRepository(ctx, db, "goauth")),
	}
}

// Account returns account services
func (s *Services) Account() account.AccountServices {
	return s.account
}
