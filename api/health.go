package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Database  string    `json:"database"`
	Config    bool      `json:"config"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow GET method
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Check database connection
	// dbStatus := "disconnected"
	// if client, err := db.GetMongoClient(); err == nil {
	// 	if err := client.Ping(r.Context(), nil); err == nil {
	// 		dbStatus = "connected"
	// 	}
	// }

	// Check if config loaded successfully
	// configLoaded := false
	// if _, err := config.GetConfig(); err == nil {
	// 	configLoaded = true
	// }

	// Prepare health status response
	status := HealthStatus{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   os.Getenv("VERSION"),
		// Database:  dbStatus,
		// Config:    configLoaded,
	}

	// If database is not connected, set status to warning
	// if dbStatus != "connected" {
	// 	status.Status = "warning"
	// }

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}
