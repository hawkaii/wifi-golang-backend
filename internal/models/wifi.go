package models

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}

type WiFi struct {
	SSID        string   `json:"ssid"`
	Password    string   `json:"password"`
	Location    Location `json:"location"`
	Description string   `json:"description"`
}
