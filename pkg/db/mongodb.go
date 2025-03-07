package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient wraps a MongoDB client
type MongoClient struct {
	client *mongo.Client
	db     *mongo.Database
	ctx    context.Context
	cancel context.CancelFunc
}

// NewMongoClient creates and initializes a MongoDB client
func NewMongoClient() (*MongoClient, error) {
	// Get MongoDB connection string from environment
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, fmt.Errorf("MONGODB_URI environment variable is not set")
	}

	// Get database name from environment or use default
	dbName := os.Getenv("MONGODB_DB")
	if dbName == "" {
		dbName = "ur_monitor"
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Get database
	db := client.Database(dbName)

	return &MongoClient{
		client: client,
		db:     db,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Disconnect closes the MongoDB connection
func (mc *MongoClient) Disconnect() {
	defer mc.cancel()
	if mc.client != nil {
		_ = mc.client.Disconnect(mc.ctx)
	}
}

// GetCollection returns a collection from the database
func (mc *MongoClient) GetCollection(name string) *mongo.Collection {
	return mc.db.Collection(name)
}

// GetContext returns the context
func (mc *MongoClient) GetContext() context.Context {
	return mc.ctx
}
