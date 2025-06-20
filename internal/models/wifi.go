package models

type WiFi struct {
	ID       string `bson:"_id,omitempty" json:"id,omitempty"`
	SSID     string `bson:"ssid" json:"ssid"`
	Password string `bson:"password" json:"password"`
	Location string `bson:"location" json:"location"`
	SharedBy string `bson:"shared_by" json:"shared_by"`
}
