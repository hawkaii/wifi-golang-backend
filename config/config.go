package config

import (
	"os"
)

type Config struct {
	MongoURI    string
	OAuthClientID     string
	OAuthClientSecret string
}

func Load() *Config {
	return &Config{
		MongoURI:    os.Getenv("MONGO_URI"),
		OAuthClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		OAuthClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
	}
}
