package routes

import (
	"encoding/json"
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
	// Generate PKCE code verifier and challenge (optional, recommended for mobile)
	verifier, err := auth.GenerateCodeVerifier()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to generate code verifier"))
		return
	}
	challenge := auth.GenerateCodeChallenge(verifier)

	// (Optional) Store code_verifier in session/cookie or pass to client for later use in callback
	// For demo, we pass it as a query param (not secure for production)

	// Build the Civic OAuth URL
	params := r.URL.Query()
	params.Set("response_type", "code")
	params.Set("client_id", auth.CivicOauthConfig.ClientID)
	params.Set("redirect_uri", auth.CivicOauthConfig.RedirectURL)
	params.Set("scope", "openid")
	params.Set("code_challenge", challenge)
	params.Set("code_challenge_method", "S256")
	// (Optional) Add state param for CSRF protection
	// params.Set("state", "random_state")

	oauthURL := auth.CivicOauthConfig.Endpoint.AuthURL + "?" + params.Encode()

	// Redirect to Civic OAuth
	http.Redirect(w, r, oauthURL, http.StatusFound)
}

func (h *Handlers) AuthMe(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: Return current user info
}

func (h *Handlers) AuthUpgrade(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: Upgrade verification level
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

func (h *Handlers) CivicCallback(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	code := r.URL.Query().Get("code")
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No code parameter provided"))
		return
	}

	token, err := auth.CivicOauthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}

	idToken, _ := token.Extra("id_token").(string)

	response := map[string]interface{}{
		"access_token": token.AccessToken,
		"id_token":     idToken,
		"expiry":       token.Expiry.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func SetupRouter(cfg *config.Config) http.Handler {
	h := NewHandlers(cfg)
	router := httprouter.New()

	// --- Authentication Endpoints ---
	// router.POST("/api/auth/civic", h.CivicAuth)
	router.GET("/api/auth/me", h.AuthMe)
	router.POST("/api/auth/upgrade", h.AuthUpgrade)
	// router.GET("/api/auth/civic/callback", h.CivicCallback)

	// --- WiFi Management Endpoints ---
	// router.POST("/api/wifi/scan", auth.RequireAuthRouter(h.WiFiScan))
	router.POST("/api/wifi/scan", auth.RequireAuthRouter(h.WiFiScan))
	router.POST("/api/wifi/connect", auth.RequireAuthRouter(h.WiFiConnect))
	router.GET("/api/wifi/nearby", h.WiFiNearby)
	router.GET("/api/wifi/all", h.GetAllWiFi) // New route to get all WiFi networks
	router.GET("/api/wifi/saved", auth.RequireAuthRouter(h.WiFiSaved))

	// --- Statistics Endpoints ---
	router.GET("/api/stats", auth.RequireAuthRouter(h.StatsGet))
	router.PATCH("/api/stats", auth.RequireAuthRouter(h.StatsPatch))

	return router
}
