package routes

import (
	"encoding/json"
	"net/http"
	"wifi-go-backend/config"
	"wifi-go-backend/internal/db"
	"wifi-go-backend/internal/models"
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
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var wifi models.WiFi
		err := json.NewDecoder(r.Body).Decode(&wifi)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid request body"))
			return
		}
		coll, err := db.GetWiFiCollection()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Database connection error"))
			return
		}
		_, err = coll.InsertOne(r.Context(), wifi)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to save WiFi details"))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("WiFi details saved"))
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
