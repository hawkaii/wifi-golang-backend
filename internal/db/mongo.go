package db

import (
	"context"
	"fmt"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
)

func GetMongoClient() (*mongo.Client, error) {
	var err error
	clientOnce.Do(func() {
		uri := os.Getenv("MONGO_URI")
		fmt.Println("Connecting to MongoDB at:", uri)
		clientInstance, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	})
	return clientInstance, err
}

func GetWiFiCollection() (*mongo.Collection, error) {
	client, err := GetMongoClient()
	if err != nil {
		return nil, err
	}
	return client.Database("wifi_db").Collection("wifi"), nil
}
