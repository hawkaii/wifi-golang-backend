package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"` // [longitude, latitude]
	Address     string    `bson:"address" json:"address"`
}

type WiFi struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SSID        string             `json:"ssid"`
	Password    string             `json:"password"`
	Location    Location           `json:"location"`
	Description string             `json:"description"`
}

type GeoJSON struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func (l Location) ToGeoJSON() GeoJSON {
	return GeoJSON{
		Type:        "Point",
		Coordinates: l.Coordinates,
	}
}
