package repositories

import (
	"errors"
	"fmt"

	"github.com/poprih/ur-monitor/lib/models"
	"github.com/poprih/ur-monitor/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository handles database operations for users
type UserRepository struct {
	client *db.MongoClient
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(client *db.MongoClient) *UserRepository {
	return &UserRepository{client: client}
}

// CreateUser adds a new user
func (r *UserRepository) CreateUser(user models.User) error {
	collection := r.client.GetCollection("users")
	_, err := collection.InsertOne(r.client.GetContext(), user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(userID string) (models.User, error) {
	collection := r.client.GetCollection("users")
	filter := bson.M{"line_user_id": userID}

	var user models.User
	err := collection.FindOne(r.client.GetContext(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, fmt.Errorf("user not found")
		}
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUser updates a user
func (r *UserRepository) UpdateUser(user models.User) error {
	collection := r.client.GetCollection("users")
	filter := bson.M{"line_user_id": user.LineUserID}
	update := bson.M{"$set": user}

	_, err := collection.UpdateOne(r.client.GetContext(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// ActivateUser sets a user's active status to true
func (r *UserRepository) ActivateUser(userID string) error {
	collection := r.client.GetCollection("users")
	filter := bson.M{"line_user_id": userID}
	update := bson.M{"$set": bson.M{"active": true}}

	_, err := collection.UpdateOne(r.client.GetContext(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}

	return nil
}

// DeactivateUser sets a user's active status to false
func (r *UserRepository) DeactivateUser(userID string) error {
	collection := r.client.GetCollection("users")
	filter := bson.M{"line_user_id": userID}
	update := bson.M{"$set": bson.M{"active": false}}

	_, err := collection.UpdateOne(r.client.GetContext(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}

// GetAllActiveUsers retrieves all active users
func (r *UserRepository) GetAllActiveUsers() ([]models.User, error) {
	collection := r.client.GetCollection("users")
	filter := bson.M{"active": true}

	cursor, err := collection.Find(r.client.GetContext(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(r.client.GetContext())

	var users []models.User
	if err := cursor.All(r.client.GetContext(), &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return users, nil
}
