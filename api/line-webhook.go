// package api

// import (
// 	"io"
// 	"log"
// 	"net/http"

// 	"github.com/poprih/ur-monitor/internal/handlers"
// )

// // Handler is the entry point for the LINE webhook API endpoint
// func Handler(w http.ResponseWriter, r *http.Request) {
// 	// Verify request method
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Read request body
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		log.Printf("Error reading request body: %v", err)
// 		http.Error(w, "Unable to read request body", http.StatusBadRequest)
// 		return
// 	}
// 	defer r.Body.Close()

// 	// Validate request body
// 	if len(body) == 0 {
// 		http.Error(w, "Empty request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Process the webhook
// 	err = handlers.HandleLineWebhook(w, r, body)
// 	if err != nil {
// 		log.Printf("Error handling webhook: %v", err)
// 		http.Error(w, "Error processing webhook", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }
