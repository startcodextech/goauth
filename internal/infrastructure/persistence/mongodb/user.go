package mongodb

import (
	"context"
	"encoding/json"
	"github.com/startcodextech/goauth/internal/domain/account"
	"github.com/startcodextech/goauth/internal/infrastructure/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User represents a document in the database.
// This struct is used to map the document in the database to a struct.
type User struct {
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

// Marshal marshals the User into a slice of bytes.
// This is needed to save the user in the database.
func (u User) Marshal() []byte {
	raw, _ := json.Marshal(u)
	return raw
}

// Ensure User implements the Entity interface.
// If the interface changes and the implementation doesn't implement the
var _ persistence.Model = (*account.User)(nil)

// Ensure User implements the Entity interface.
// If the interface changes and the implementation doesn't implement the
var _ persistence.Entity = (*User)(nil)

type MongoUserRepository struct {
	MongoRepository[User]
}

// Ensure MongoUserRepository implements the UserRepository interface.
// If the interface changes and the implementation doesn't implement the
var _ account.UserRepository = (*MongoUserRepository)(nil)

// NewMongoUserRepository creates a new UserRepository backed by MongoDB.
// It also creates the indexes for the collection if they don't exist.
func NewMongoUserRepository(ctx context.Context, db *mongo.Database, name string) *MongoUserRepository {
	repo := &MongoUserRepository{
		NewMongoRepository[User](db, name),
	}

	model := mongo.IndexModel{
		Keys: bson.M{
			"email": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := repo.collection.Indexes().CreateOne(ctx, model)
	if err != nil {
		panic(err)
	}

	return repo
}

// Save persists the user.
// If the user doesn't exist in the database it's inserted.
func (r *MongoUserRepository) Save(ctx context.Context, user account.User) {
	r.MongoRepository.Save(ctx, user)
}

// Find retrieves a user from the database.
// If the user doesn't exist in the database it returns nil.
func (r *MongoUserRepository) Find(ctx context.Context, filter map[string]interface{}) ([]account.User, error) {

	var result []account.User

	factory := func() persistence.Model {
		return account.NewUser()
	}

	rows, err := r.MongoRepository.Find(ctx, filter, factory)
	if err != nil {
		return nil, err
	}

	result = make([]account.User, len(rows))

	for i, user := range rows {
		result[i] = user.(account.User)
	}

	return result, nil
}
