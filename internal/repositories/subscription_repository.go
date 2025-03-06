package repositories

import (
	"database/sql"
	"fmt"
)

// Subscription represents a user's subscription to a UR property
type Subscription struct {
	ID       int
	UserID   string
	DanchiID string
}

// SubscriptionRepository handles database operations for subscriptions
type SubscriptionRepository struct {
	DB *sql.DB
}

// NewSubscriptionRepository creates a new SubscriptionRepository
func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

// AddSubscription adds a new subscription to the database
func (repo *SubscriptionRepository) AddSubscription(subscription Subscription) error {
	query := "INSERT INTO subscriptions (user_id, danchi_id) VALUES (?, ?)"
	_, err := repo.DB.Exec(query, subscription.UserID, subscription.DanchiID)
	if err != nil {
		return fmt.Errorf("failed to add subscription: %w", err)
	}
	return nil
}

// GetSubscriptionsByUserID retrieves subscriptions for a given user
func (repo *SubscriptionRepository) GetSubscriptionsByUserID(userID string) ([]Subscription, error) {
	query := "SELECT id, user_id, danchi_id FROM subscriptions WHERE user_id = ?"
	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []Subscription
	for rows.Next() {
		var subscription Subscription
		if err := rows.Scan(&subscription.ID, &subscription.UserID, &subscription.DanchiID); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}
