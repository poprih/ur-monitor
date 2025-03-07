package services

import (
	"github.com/poprih/ur-monitor/lib/repositories"
)

// MonitorService handles UR property monitoring
type MonitorService struct {
	userRepo *repositories.UserRepository
	subRepo  *repositories.SubscriptionRepository
}

// NewMonitorService creates a new monitor service
func NewMonitorService(userRepo *repositories.UserRepository, subRepo *repositories.SubscriptionRepository) *MonitorService {
	return &MonitorService{
		userRepo: userRepo,
		subRepo:  subRepo,
	}
}

// FetchPropertiesForDanchi fetches available properties for a specific danchi
// func (s *
