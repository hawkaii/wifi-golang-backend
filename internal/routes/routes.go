package routes

import (
	"net/http"
	"wifi-go-backend/config"
)

func SetupRouter(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// --- Authentication Endpoints ---
	mux.HandleFunc("/api/auth/civic", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement Civic authentication
	})
	mux.HandleFunc("/api/auth/me", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Return current user info
	})
	mux.HandleFunc("/api/auth/upgrade", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Upgrade verification level
	})

	// --- WiFi Management Endpoints ---
	mux.HandleFunc("/api/wifi/scan", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Scan for networks
	})
	mux.HandleFunc("/api/wifi/connect", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Connect to network
	})
	mux.HandleFunc("/api/wifi/nearby", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Get nearby networks
	})
	mux.HandleFunc("/api/wifi/saved", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Get saved networks
	})

	// --- Statistics Endpoints ---
	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// TODO: Get user statistics
		} else if r.Method == http.MethodPatch {
			// TODO: Update statistics
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return mux
}
