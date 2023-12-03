package mongodb

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

// MongoDriver represents the MongoDB driver.
type MongoDriver struct {
	client   *mongo.Client
	database *mongo.Database
	logger   watermill.LoggerAdapter
}

// New creates a new MongoDB client.
// param dbName: The name of the database to use.
func New(ctx context.Context, dbName string, logger watermill.LoggerAdapter) MongoDriver {
	severApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("DB_URI")).SetServerAPIOptions(severApi)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", err, nil)
		panic(err)
	}

	return MongoDriver{
		client:   client,
		database: client.Database(dbName),
		logger:   logger,
	}
}

// Client returns the client instance.
func (d MongoDriver) Client() *mongo.Client {
	return d.client
}

// Database returns the database instance.
func (d MongoDriver) Database() *mongo.Database {
	return d.database
}

// Disconnect disconnects the client.
// example: defer driver.Disconnect()()
func (d MongoDriver) Disconnect(ctx context.Context) func() {
	return func() {
		if err := d.client.Disconnect(ctx); err != nil {
			d.logger.Error("Failed to disconnect from MongoDB", err, nil)
			panic(err)
		}
	}
}

func (d MongoDriver) Ping(ctx context.Context) {
	if err := d.database.RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		d.logger.Error("Failed to ping MongoDB", err, nil)
		panic(err)
	}
}
