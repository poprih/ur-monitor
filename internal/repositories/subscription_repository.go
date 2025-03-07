package repositories

import (
	"fmt"

	"github.com/poprih/ur-monitor/internal/models"
	"github.com/poprih/ur-monitor/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
)

// SubscriptionRepository handles database operations for subscriptions
type SubscriptionRepository struct {
	client *db.MongoClient
}

// NewSubscriptionRepository creates a new SubscriptionRepository
func NewSubscriptionRepository(client *db.MongoClient) *SubscriptionRepository {
	return &SubscriptionRepository{client: client}
}

// AddSubscription adds a new subscription
func (r *SubscriptionRepository) AddSubscription(subscription models.Subscription) error {
	collection := r.client.GetCollection("subscriptions")
	_, err := collection.InsertOne(r.client.GetContext(), subscription)
	if err != nil {
		return fmt.Errorf("failed to add subscription: %w", err)
	}
	return nil
}

// RemoveSubscription removes a subscription
func (r *SubscriptionRepository) RemoveSubscription(userID, danchiID string) error {
	collection := r.client.GetCollection("subscriptions")
	filter := bson.M{"user_id": userID, "danchi_id": danchiID}

	_, err := collection.DeleteOne(r.client.GetContext(), filter)
	if err != nil {
		return fmt.Errorf("failed to remove subscription: %w", err)
	}
	return nil
}

// GetSubscriptionsByUserID retrieves all subscriptions for a user
func (r *SubscriptionRepository) GetSubscriptionsByUserID(userID string) ([]models.Subscription, error) {
	collection := r.client.GetCollection("subscriptions")
	filter := bson.M{"user_id": userID}

	cursor, err := collection.Find(r.client.GetContext(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find subscriptions: %w", err)
	}
	defer cursor.Close(r.client.GetContext())

	var subscriptions []models.Subscription
	if err := cursor.All(r.client.GetContext(), &subscriptions); err != nil {
		return nil, fmt.Errorf("failed to decode subscriptions: %w", err)
	}

	return subscriptions, nil
}

// GetAllSubscriptions retrieves all subscriptions
func (r *SubscriptionRepository) GetAllSubscriptions() ([]models.Subscription, error) {
	collection := r.client.GetCollection("subscriptions")

	cursor, err := collection.Find(r.client.GetContext(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find subscriptions: %w", err)
	}
	defer cursor.Close(r.client.GetContext())

	var subscriptions []models.Subscription
	if err := cursor.All(r.client.GetContext(), &subscriptions); err != nil {
		return nil, fmt.Errorf("failed to decode subscriptions: %w", err)
	}

	return subscriptions, nil
}

// GetSubscriptionsByDanchiID retrieves all subscriptions for a danchi
func (r *SubscriptionRepository) GetSubscriptionsByDanchiID(danchiID string) ([]models.Subscription, error) {
	collection := r.client.GetCollection("subscriptions")
	filter := bson.M{"danchi_id": danchiID}

	cursor, err := collection.Find(r.client.GetContext(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find subscriptions: %w", err)
	}
	defer cursor.Close(r.client.GetContext())

	var subscriptions []models.Subscription
	if err := cursor.All(r.client.GetContext(), &subscriptions); err != nil {
		return nil, fmt.Errorf("failed to decode subscriptions: %w", err)
	}

	return subscriptions, nil
}

// SubscriptionExists checks if a subscription already exists
func (r *SubscriptionRepository) SubscriptionExists(userID, danchiID string) (bool, error) {
	collection := r.client.GetCollection("subscriptions")
	filter := bson.M{"user_id": userID, "danchi_id": danchiID}

	count, err := collection.CountDocuments(r.client.GetContext(), filter)
	if err != nil {
		return false, fmt.Errorf("failed to check subscription existence: %w", err)
	}

	return count > 0, nil
}
