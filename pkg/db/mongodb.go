package db

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/poprih/ur-monitor/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

var (
	client *mongo.Client
	once   sync.Once
	err    error
)

// GetMongoClient returns a MongoDB client, creating it if necessary
func GetMongoClient() (*mongo.Client, error) {
	once.Do(func() {
		// Get configuration
		cfg, configErr := config.GetConfig()
		if configErr != nil {
			err = configErr
			return
		}

		// Create a context with timeout for the connection
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Create client options
		clientOptions := options.Client().ApplyURI(cfg.MongoDBURI)

		// Connect to MongoDB
		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			return
		}

		// Ping the database to verify connection
		err = client.Ping(ctx, readpref.Primary())
	})

	return client, err
}

// GetDatabase returns a specific MongoDB database
func GetDatabase() (*mongo.Database, error) {
	client, err := GetMongoClient()
	if err != nil {
		return nil, err
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	return client.Database(cfg.MongoDBDatabase), nil
}

// CloseConnection closes the MongoDB connection
func CloseConnection(ctx context.Context) error {
	if client != nil {
		return client.Disconnect(ctx)
	}
	return nil
}
