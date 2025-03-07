ur-monitor/
├── api/ # Vercel serverless functions
│ ├── line-webhook.go # LINE webhook endpoint
│ ├── monitor.go # Monitoring function trigger endpoint
│ └── health.go # Health check endpoint
│
├── internal/
│ ├── config/ # Configuration
│ │ └── config.go # Config loader
│ ├── handlers/ # Request handlers
│ │ ├── line_handler.go # LINE webhook handler
│ │ └── monitor_handler.go # Monitoring job handler
│ ├── services/ # Business logic
│ │ ├── line_service.go # LINE message processing service
│ │ ├── monitor_service.go # UR property monitoring service
│ │ └── notification_service.go # Notification delivery service
│ ├── models/ # Data models
│ │ ├── user.go # User model
│ │ ├── subscription.go # Subscription model
│ │ └── property.go # Property model
│ ├── repositories/ # Data access
│ │ ├── user_repository.go # User repository
│ │ └── subscription_repository.go # Subscription repository
│ └── clients/ # External API clients
│ ├── line_client.go # LINE API client
│ └── ur_client.go # UR property API client
│
├── pkg/ # Shared utilities
│ ├── db/ # Database utilities
│ │ └── mongodb.go # MongoDB client
│ └── utils/ # General utilities
│ └── utils.go # Utility functions
│
├── vercel.json # Vercel configuration
├── go.mod # Go module file
└── go.sum # Go dependencies file

Build this LINE bot service for monitoring UR apartments.
Suitable for Vercel deployment and to implement the subscription functionality properly.
Break down the approach:

- Structure the service for Vercel serverless deployment
- Implement LINE webhook handling for user interactions
- Set up the monitoring job that checks for available apartments
- Create a database setup for storing user subscriptions
- Connect everything together
