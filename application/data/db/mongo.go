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
	Client     *mongo.Client
	Collection *mongo.Collection
}

// Initializes MongoDB connection
func ConnectMongoDB(uri, dbName, collectionName string) (*MongoDB, error) {

	// Sets URI for DB to connect to.
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

	collection := client.Database(dbName).Collection(collectionName)
	return &MongoDB{Client: client, Collection: collection}, nil
}

func (db *MongoDB) Close() {
	err := db.Client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Disconnected from MongoDB.")
}
