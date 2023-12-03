package mongodb

import (
	"context"
	"encoding/json"
	"github.com/startcodextech/goauth/internal/infrastructure/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	// MongoRepository represents the persistence for the entity.
	// It implements the Repository interface.
	MongoRepository[T persistence.Entity] struct {
		collection *mongo.Collection
		operations []func(ctx mongo.SessionContext) error
		session    mongo.Session
	}
)

// NewMongoRepository implements the Repository interface.
// If the interface changes and the implementation doesn't implement the
func NewMongoRepository[T persistence.Entity](db *mongo.Database, name string) MongoRepository[T] {
	return MongoRepository[T]{
		collection: db.Collection(name),
		operations: []func(mongo.SessionContext) error{},
		session:    nil,
	}
}

var _ persistence.Repository[persistence.Entity] = (*MongoRepository[persistence.Entity])(nil)

// StartTx starts a transaction and returns the session
// that should be used by the repository.
func (r *MongoRepository[T]) StartTx(ctx context.Context) (interface{}, error) {
	session, err := r.collection.Database().Client().StartSession()
	if err != nil {
		return nil, err
	}

	r.session = session
	return r.session, err
}

// SetTx sets the transaction to be used by the repository.
// This is used when the repository is created inside a transaction.
func (r *MongoRepository[T]) SetTx(session interface{}) {
	r.session = session.(mongo.Session)
}

// Save persists the user.
// If the user doesn't exist in the database it's inserted.
// If the user already exists in the database it's updated.
func (r *MongoRepository[T]) Save(ctx context.Context, model persistence.Model) {
	opts := options.Update().SetUpsert(true)
	op := func(ctx mongo.SessionContext) error {

		data, err := r.marshal(model)
		if err != nil {
			return err
		}

		filter := bson.M{"_id": model.ID()}

		update := bson.M{
			"$set": data,
		}

		_, err = r.collection.UpdateOne(ctx, filter, update, opts)
		return err
	}
	r.operations = append(r.operations, op)
}

// Delete removes the user from the persistence.
// If the user doesn't exist in the database the function returns an error.
func (r *MongoRepository[T]) Delete(ctx context.Context, ID string) {
	op := func(ctx mongo.SessionContext) error {
		_, err := r.collection.DeleteOne(ctx, bson.M{"_id": ID})
		return err
	}

	r.operations = append(r.operations, op)
}

// Find retrieves users from the persistence.
// It can retrieve all the users or filter them by a criteria.
func (r *MongoRepository[T]) Find(ctx context.Context, filter map[string]interface{}, factoryModel func() persistence.Model) ([]persistence.Model, error) {

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = cursor.Close(ctx); err != nil {
			return
		}
	}()

	result := []persistence.Model{}

	for cursor.Next(ctx) {
		var data bson.M
		if err = cursor.Decode(&data); err != nil {
			return nil, err
		}

		model := factoryModel()
		if err = model.UnmarshalFromMap(data); err != nil {
			return nil, err
		}

		result = append(result, model)
	}

	if cursor.Err() != nil {
		return nil, cursor.Err()
	}

	return result, nil
}

// Commit commits the transaction.
// If the transaction fails the function returns an error.
func (r *MongoRepository[T]) Commit(ctx context.Context) error {

	if r.session == nil {
		_, err := r.StartTx(ctx)
		if err != nil {
			return err
		}
	}

	defer r.session.EndSession(ctx)

	_, err := r.session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		for _, op := range r.operations {
			if err := op(sessCtx); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		if r.session.AbortTransaction(ctx) != nil {
			return err
		}
		return err
	}

	err = r.session.CommitTransaction(ctx)
	if err != nil {
		return err
	}

	r.operations = []func(mongo.SessionContext) error{}

	r.session = nil

	return nil
}

// marshal converts the user domain entity to a mongodb document.
// This is needed because the bson package doesn't support
func (r *MongoRepository[T]) marshal(model persistence.Model) (T, error) {
	var entity T

	err := json.Unmarshal(model.Marshal(), &entity)
	if err != nil {
		return entity, err
	}

	return entity, nil
}
