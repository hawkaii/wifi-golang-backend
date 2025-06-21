package routes

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"wifi-go-backend/internal/db"
	"wifi-go-backend/internal/models"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

func (h *Handlers) WiFiNearby(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse query params for latitude and longitude
	latStr := r.URL.Query().Get("latitude")
	lngStr := r.URL.Query().Get("longitude")
	if latStr == "" || lngStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("latitude and longitude query parameters are required"))
		return
	}
	lat, err1 := strconv.ParseFloat(latStr, 64)
	lng, err2 := strconv.ParseFloat(lngStr, 64)
	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid latitude or longitude"))
		return
	}

	coll, err := db.GetWiFiCollection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Database connection error"))
		return
	}

	// 1km radius in degrees (approx)
	const radiusKm = 1.0
	const earthRadiusKm = 6371.0
	// MongoDB $geoWithin with $centerSphere expects [lng, lat]
	filter := map[string]interface{}{
		"location": map[string]interface{}{
			"$geoWithin": map[string]interface{}{
				"$centerSphere": []interface{}{
					[]float64{lng, lat},
					radiusKm / earthRadiusKm,
				},
			},
		},
	}
	cur, err := coll.Find(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to search WiFi"))
		return
	}
	defer cur.Close(r.Context())

	var results []map[string]interface{}
	for cur.Next(r.Context()) {
		var wifi models.WiFi
		if err := cur.Decode(&wifi); err != nil {
			continue
		}
		// Coordinates: [lng, lat]
		var wifiLat, wifiLng float64
		if len(wifi.Location.Coordinates) == 2 {
			wifiLng = wifi.Location.Coordinates[0]
			wifiLat = wifi.Location.Coordinates[1]
		}
		dist := haversine(lat, lng, wifiLat, wifiLng)
		results = append(results, map[string]interface{}{
			"id":          wifi.ID, // Add the ID field
			"ssid":        wifi.SSID,
			"location":    wifi.Location,
			"distance":    dist,
			"description": wifi.Description,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// GetAllWiFi returns all WiFi networks in the database
func (h *Handlers) GetAllWiFi(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	coll, err := db.GetWiFiCollection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Database connection error"))
		return
	}
	cur, err := coll.Find(r.Context(), map[string]interface{}{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to fetch WiFi list"))
		return
	}
	defer cur.Close(r.Context())

	var results []map[string]interface{}
	for cur.Next(r.Context()) {
		var wifi models.WiFi
		if err := cur.Decode(&wifi); err != nil {
			continue
		}
		results = append(results, map[string]interface{}{
			"id":          wifi.ID,
			"ssid":        wifi.SSID,
			"location":    wifi.Location,
			"description": wifi.Description,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// Haversine formula to calculate distance between two lat/lng points in km
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371 // Earth radius in km
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLng := (lng2 - lng1) * math.Pi / 180.0
	lat1R := lat1 * math.Pi / 180.0
	lat2R := lat2 * math.Pi / 180.0
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLng/2)*math.Sin(dLng/2)*math.Cos(lat1R)*math.Cos(lat2R)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
