package account

import (
	"context"
	"github.com/startcodextech/goauth/internal/domain/account"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	// User represents a document in the database.
	// This struct is used to map the document in the database to a struct.
	User struct {
		ID           string `bson:"_id"`
		Name         string `bson:"name"`
		Lastname     string `bson:"lastname"`
		Email        string `bson:"email"`
		PasswordHash string `bson:"password_hash,omitempty"`
		PhotoUrl     string `bson:"photo_url,omitempty"`
		Verified     bool   `bson:"verified"`
		FacebookID   string `bson:"facebook_id,omitempty"`
		GoogleID     string `bson:"google_id,omitempty"`
		MicrosoftID  string `bson:"microsoft_id,omitempty"`
		AppleID      string `bson:"apple_id,omitempty"`
	}

	// UserRepositoryMongo represents the repository for the User entity.
	// It implements the UserRepository interface.
	UserRepositoryMongo struct {
		collection *mongo.Collection
		operations []func(ctx mongo.SessionContext) error
	}
)

// Ensure UserRepositoryMongo implements the UserRepository interface.
// If the interface changes and the implementation doesn't implement the
var _ account.UserRepository = (*UserRepositoryMongo)(nil)

// NewUserRepositoryMongo creates a new UserRepository backed by MongoDB.
// It also creates the indexes for the collection if they don't exist.
func NewUserRepositoryMongo(db *mongo.Database, name string) UserRepositoryMongo {
	return UserRepositoryMongo{
		collection: db.Collection(name),
		operations: []func(mongo.SessionContext) error{},
	}
}

// Save persists the user.
// If the user doesn't exist in the database it's inserted.
// If the user already exists in the database it's updated.
func (r *UserRepositoryMongo) Save(ctx context.Context, user account.User) {
	opts := options.Update().SetUpsert(true)
	op := func(ctx mongo.SessionContext) error {

		document, err := r.marshalUser(user)
		if err != nil {
			return err
		}

		filter := bson.M{"_id": user.ID()}

		_, err = r.collection.UpdateOne(ctx, filter, document, opts)
		return err
	}
	r.operations = append(r.operations, op)
}

// Find retrieves users from the repository.
// It can retrieve all the users or filter them by a criteria.
func (r *UserRepositoryMongo) Find(ctx context.Context, filter map[string]interface{}) (users []account.User, err error) {

	users = []account.User{}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = cursor.Close(ctx); err != nil {
			return
		}
	}()

	for cursor.Next(ctx) {
		var user User
		if err = cursor.Decode(&user); err != nil {
			return
		}

		users = append(users, r.decode(user))
	}

	return
}

// Delete removes the user from the repository.
// If the user doesn't exist in the database the function returns an error.
func (r *UserRepositoryMongo) Delete(ctx context.Context, ID string) {
	op := func(ctx mongo.SessionContext) error {
		_, err := r.collection.DeleteOne(ctx, bson.M{"_id": ID})
		return err
	}

	r.operations = append(r.operations, op)
}

// Commit commits the transaction.
// If the transaction fails the function returns an error.
func (r *UserRepositoryMongo) Commit(ctx context.Context) error {
	session, err := r.collection.Database().Client().StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		for _, op := range r.operations {
			if err := op(sessCtx); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		if session.AbortTransaction(ctx) != nil {
			return err
		}
		return err
	}

	err = session.CommitTransaction(ctx)
	if err != nil {
		return err
	}

	r.operations = []func(mongo.SessionContext) error{}

	return nil
}

// marshalUser converts a User to a bson document.
// This function is used to prepare the data to be saved in MongoDB.
func (r *UserRepositoryMongo) marshalUser(user account.User) ([]byte, error) {
	document := User{
		ID:           user.ID(),
		Name:         user.Name(),
		Lastname:     user.Lastname(),
		Email:        user.Email(),
		PasswordHash: user.PasswordHash(),
		PhotoUrl:     user.PhotoUrl(),
		Verified:     user.Verified(),
		FacebookID:   user.FacebookID(),
		GoogleID:     user.GoogleID(),
		MicrosoftID:  user.MicrosoftID(),
		AppleID:      user.AppleID(),
	}

	return bson.Marshal(document)
}

// decode converts a bson document to a User.
// This function is used to retrieve data from MongoDB.
func (r *UserRepositoryMongo) decode(doc User) (user account.User) {
	user.SetID(doc.ID)
	user.SetName(doc.Name)
	user.SetLastname(doc.Lastname)
	user.SetEmail(doc.Email)
	user.SetPasswordHash(doc.PasswordHash)
	user.SetPhotoUrl(doc.PhotoUrl)
	user.SetVerified(doc.Verified)
	user.SetFacebookID(doc.FacebookID)
	user.SetGoogleID(doc.GoogleID)
	user.SetMicrosoftID(doc.MicrosoftID)
	user.SetAppleID(doc.AppleID)
	return
}
