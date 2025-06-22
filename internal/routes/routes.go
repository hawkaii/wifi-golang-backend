package routes

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

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

// GeminiRecommendHandler handles /api/gemini/recommend requests
func GeminiRecommendHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	startLatStr := r.URL.Query().Get("start_lat")
	startLngStr := r.URL.Query().Get("start_lng")
	endLatStr := r.URL.Query().Get("end_lat")
	endLngStr := r.URL.Query().Get("end_lng")
	if startLatStr == "" || startLngStr == "" || endLatStr == "" || endLngStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing required query parameters: start_lat, start_lng, end_lat, end_lng"))
		return
	}
	startLat, err1 := strconv.ParseFloat(startLatStr, 64)
	startLng, err2 := strconv.ParseFloat(startLngStr, 64)
	endLat, err3 := strconv.ParseFloat(endLatStr, 64)
	endLng, err4 := strconv.ParseFloat(endLngStr, 64)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid coordinate values"))
		return
	}
	apiKey := os.Getenv("GEMINI_API_KEY")
	recommender, err := NewLocationRecommender(apiKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create Gemini recommender"))
		return
	}
	defer recommender.Close()
	ctx := r.Context()
	jsonStops, err := recommender.FindStopsBetween(ctx, startLat, startLng, endLat, endLng)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Gemini recommendation failed: " + err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonStops)
}

// GeminiRecommendStopsWiFiHandler handles /api/gemini/recommendstopswifi requests
// It calls the recommender, then for each stop, lists all nearby WiFi networks
func GeminiRecommendStopsWiFiHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	startLatStr := r.URL.Query().Get("start_lat")
	startLngStr := r.URL.Query().Get("start_lng")
	endLatStr := r.URL.Query().Get("end_lat")
	endLngStr := r.URL.Query().Get("end_lng")
	if startLatStr == "" || startLngStr == "" || endLatStr == "" || endLngStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing required query parameters: start_lat, start_lng, end_lat, end_lng"))
		return
	}
	startLat, err1 := strconv.ParseFloat(startLatStr, 64)
	startLng, err2 := strconv.ParseFloat(startLngStr, 64)
	endLat, err3 := strconv.ParseFloat(endLatStr, 64)
	endLng, err4 := strconv.ParseFloat(endLngStr, 64)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid coordinate values"))
		return
	}
	apiKey := os.Getenv("GEMINI_API_KEY")
	recommender, err := NewLocationRecommender(apiKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create Gemini recommender"))
		return
	}
	defer recommender.Close()
	ctx := r.Context()
	stopsResp, err := recommender.FindStops(ctx, RecommendationRequest{
		StartCoordinate: Coordinate{Latitude: startLat, Longitude: startLng},
		EndCoordinate:   Coordinate{Latitude: endLat, Longitude: endLng},
		StopType:        "any",
		MaxStops:        5,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Gemini recommendation failed: " + err.Error()))
		return
	}

	type WiFiWithStop struct {
		Stop  Coordinate               `json:"stop"`
		WiFis []map[string]interface{} `json:"wifis"`
	}
	var results []WiFiWithStop
	for _, stop := range stopsResp.Stops {
		coll, err := db.GetWiFiCollection()
		if err != nil {
			continue
		}
		const radiusKm = 1.0
		const earthRadiusKm = 6371.0
		filter := map[string]interface{}{
			"location": map[string]interface{}{
				"$geoWithin": map[string]interface{}{
					"$centerSphere": []interface{}{
						[]float64{stop.Longitude, stop.Latitude},
						radiusKm / earthRadiusKm,
					},
				},
			},
		}
		cur, err := coll.Find(ctx, filter)
		if err != nil {
			continue
		}
		var wifis []map[string]interface{}
		for cur.Next(ctx) {
			var wifi models.WiFi
			if err := cur.Decode(&wifi); err == nil {
				wifis = append(wifis, map[string]interface{}{
					"id":          wifi.ID,
					"ssid":        wifi.SSID,
					"location":    wifi.Location,
					"description": wifi.Description,
				})
			}
		}
		cur.Close(ctx)
		results = append(results, WiFiWithStop{
			Stop:  stop,
			WiFis: wifis,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"stops_with_wifi":   results,
		"route_description": stopsResp.Route,
	})
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
	router.GET("/api/wifi/saved", auth.RequireAuthRouter(h.WiFiSaved))

	// --- Statistics Endpoints ---
	router.GET("/api/stats", auth.RequireAuthRouter(h.StatsGet))
	router.PATCH("/api/stats", auth.RequireAuthRouter(h.StatsPatch))
	router.POST("/api/wifi/nearby/stops", h.NearbyWiFiForStopsHandler)

	// --- Gemini Recommender Endpoint ---
	router.GET("/api/gemini/recommendstops", GeminiRecommendHandler)
	router.GET("/api/gemini/recommendstopswifi", GeminiRecommendStopsWiFiHandler)

	return router
}
