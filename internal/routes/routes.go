package routes

import (
	"net/http"

	"wifi-go-backend/config"
	"wifi-go-backend/internal/auth"

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

func (h *Handlers) WiFiConnect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: Connect to network
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
	router.GET("/api/wifi/connect", auth.RequireAuthRouter(h.WiFiConnect))
	router.GET("/api/wifi/nearby", h.WiFiNearby)
	router.GET("/api/wifi/saved", auth.RequireAuthRouter(h.WiFiSaved))

	// --- Statistics Endpoints ---
	router.GET("/api/stats", auth.RequireAuthRouter(h.StatsGet))
	router.PATCH("/api/stats", auth.RequireAuthRouter(h.StatsPatch))

	return router
}
