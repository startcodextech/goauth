package mongodb

import (
	"context"
	"encoding/json"
	"github.com/startcodextech/goauth/internal/domain/account"
	"github.com/startcodextech/goauth/internal/infrastructure/persistence"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// User represents a document in the database.
	// This struct is used to map the document in the database to a struct.
	User struct {
		ID           string `bson:"_id" json:"id"`
		Name         string `bson:"name" json:"name"`
		Lastname     string `bson:"lastname" json:"lastname"`
		Email        string `bson:"email" json:"email"`
		PasswordHash string `bson:"password_hash,omitempty" json:"password_hash,omitempty"`
		PhotoUrl     string `bson:"photo_url,omitempty" json:"photo_url,omitempty"`
		Verified     bool   `bson:"verified" json:"verified"`
		FacebookID   string `bson:"facebook_id,omitempty" json:"facebook_id,omitempty"`
		GoogleID     string `bson:"google_id,omitempty" json:"google_id,omitempty"`
		MicrosoftID  string `bson:"microsoft_id,omitempty" json:"microsoft_id,omitempty"`
		AppleID      string `bson:"apple_id,omitempty" json:"apple_id,omitempty"`
	}

	MongoUserRepository[T User] struct {
		MongoRepository[T]
	}
)

// Ensure User implements the Entity interface.
// If the interface changes and the implementation doesn't implement the
var _ persistence.Model = (*account.User)(nil)

// Ensure User implements the Entity interface.
// If the interface changes and the implementation doesn't implement the
var _ persistence.Entity = (*User)(nil)

// Ensure MongoUserRepository implements the UserRepository interface.
// If the interface changes and the implementation doesn't implement the
var _ account.UserRepository = (*MongoUserRepository[User])(nil)

// NewMongoUserRepository creates a new UserRepository backed by MongoDB.
// It also creates the indexes for the collection if they don't exist.
func NewMongoUserRepository(db *mongo.Database, name string) MongoUserRepository[User] {
	return MongoUserRepository[User]{
		NewMongoRepository[User](db, name),
	}
}

// Marshal marshals the User into a slice of bytes.
// This is needed to save the user in the database.
func (u User) Marshal() []byte {
	raw, _ := json.Marshal(u)
	return raw
}

// Save persists the user.
// If the user doesn't exist in the database it's inserted.
func (r *MongoUserRepository[T]) Save(ctx context.Context, user account.User) {
	r.MongoRepository.Save(ctx, user)
}
