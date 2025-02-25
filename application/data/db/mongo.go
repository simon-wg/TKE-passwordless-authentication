package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB instance struct
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// Initializes MongoDB connection
func ConnectMongoDB(uri, dbName string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping MongoDB to confirm connection
	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB")

	// Get a reference to the database
	database := client.Database(dbName)

	return &MongoDB{Client: client, Database: database}, nil
}

// Close terminates the MongoDB connection
func (db *MongoDB) Close() {
	err := db.Client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Disconnected from MongoDB.")
}
