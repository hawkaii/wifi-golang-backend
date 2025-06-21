package routes

import (
	"encoding/json"
	"net/http"

	"wifi-go-backend/internal/db"
	"wifi-go-backend/internal/models"

	"github.com/julienschmidt/httprouter"
)

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

	// Ensure location is GeoJSON format
	wifi.Location.Type = "Point"
	if len(wifi.Location.Coordinates) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Coordinates must be [longitude, latitude]"))
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
