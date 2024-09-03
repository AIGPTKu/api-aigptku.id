package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoConn(uri, dbName string) *mongo.Database {
	// Connect to MongoDB.
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic("Failed to connect MongoDB: " + err.Error())
	}

	// Check the connection.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic("Failed to ping MongoDB: " + err.Error())
	}

	return client.Database(dbName)
}