package routes

import (
	"encoding/json"
	"net/http"

	"wifi-go-backend/config"
	"wifi-go-backend/internal/auth"
	"wifi-go-backend/internal/db"
	"wifi-go-backend/internal/models"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (h *Handlers) WiFiConnect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse request body for wifi_id and location
	type ConnectRequest struct {
		WiFiID    string  `json:"wifi_id"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	var req ConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}

	if req.WiFiID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wifi_id is required"))
		return
	}

	coll, err := db.GetWiFiCollection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Database connection error"))
		return
	}

	// Find WiFi by ID
	var wifi models.WiFi
	objID, err := primitive.ObjectIDFromHex(req.WiFiID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid wifi_id"))
		return
	}
	err = coll.FindOne(r.Context(), map[string]interface{}{"_id": objID}).Decode(&wifi)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("WiFi not found"))
		return
	}

	// Check if the provided location is within 100 meters of the WiFi
	if len(wifi.Location.Coordinates) != 2 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("WiFi location data invalid"))
		return
	}
	wifiLng := wifi.Location.Coordinates[0]
	wifiLat := wifi.Location.Coordinates[1]
	dist := haversine(req.Latitude, req.Longitude, wifiLat, wifiLng)
	if dist > 0.1 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("You are too far from this WiFi to connect"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ssid":        wifi.SSID,
		"password":    wifi.Password,
		"location":    wifi.Location,
		"description": wifi.Description,
	})
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
	router.GET("/api/wifi/nearby", h.WiFiNearby)
	router.GET("/api/wifi/saved", auth.RequireAuthRouter(h.WiFiSaved))

	// --- Statistics Endpoints ---
	router.GET("/api/stats", auth.RequireAuthRouter(h.StatsGet))
	router.PATCH("/api/stats", auth.RequireAuthRouter(h.StatsPatch))

	return router
}
