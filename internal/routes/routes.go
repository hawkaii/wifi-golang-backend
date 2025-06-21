package routes

import (
	"encoding/json"
	"net/http"

	"wifi-go-backend/config"
	"wifi-go-backend/internal/auth"
	"wifi-go-backend/internal/db"
	"wifi-go-backend/internal/models"

	"github.com/julienschmidt/httprouter"
)

// Handlers struct for dependency injection
// (following "Let's Go Further" by Alex Edwards)
type Handlers struct {
	Cfg *config.Config
}

func NewHandlers(cfg *config.Config) *Handlers {
	return &Handlers{Cfg: cfg}
}

func (h *Handlers) CivicAuth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: Implement Civic authentication
}

func (h *Handlers) AuthMe(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: Return current user info
}

func (h *Handlers) AuthUpgrade(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: Upgrade verification level
}

func (h *Handlers) WiFiScan(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var wifi models.WiFi
	if err := json.NewDecoder(r.Body).Decode(&wifi); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}

	// Ensure description is present
	if wifi.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Description is required"))
		return
	}

	coll, err := db.GetWiFiCollection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Database connection error"))
		return
	}

	// Only add if there is no WiFi with the same SSID at the same address
	filter := map[string]interface{}{
		"ssid":             wifi.SSID,
		"location.address": wifi.Location.Address,
	}
	count, err := coll.CountDocuments(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to check existing WiFi"))
		return
	}
	if count > 0 {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("WiFi with this SSID already exists at this address"))
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
}

func (h *Handlers) WiFiConnect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: Connect to network
}

func (h *Handlers) WiFiNearby(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: Get nearby networks
}

func (h *Handlers) WiFiSaved(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: Get saved networks
}

func (h *Handlers) StatsGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: Get user statistics
}

func (h *Handlers) StatsPatch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: Update statistics
}

func SetupRouter(cfg *config.Config) http.Handler {
	h := NewHandlers(cfg)
	router := httprouter.New()

	// --- Authentication Endpoints ---
	router.POST("/api/auth/civic", h.CivicAuth)
	router.GET("/api/auth/me", h.AuthMe)
	router.POST("/api/auth/upgrade", h.AuthUpgrade)

	// --- WiFi Management Endpoints ---
	router.POST("/api/wifi/scan", auth.RequireAuthRouter(h.WiFiScan))
	router.POST("/api/wifi/connect", auth.RequireAuthRouter(h.WiFiConnect))
	router.POST("/api/wifi/nearby", auth.RequireAuthRouter(h.WiFiNearby))
	router.GET("/api/wifi/saved", auth.RequireAuthRouter(h.WiFiSaved))

	// --- Statistics Endpoints ---
	router.GET("/api/stats", auth.RequireAuthRouter(h.StatsGet))
	router.PATCH("/api/stats", auth.RequireAuthRouter(h.StatsPatch))

	return router
}
