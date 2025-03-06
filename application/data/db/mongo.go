package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB instance struct
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// ConnectMongoDB establishes a connection to a MongoDB instance using the provided URI and database name
// It returns a MongoDB struct containing the client and database reference, or an error if the connection fails
//
// Parameters:
//   - uri: The connection string URI for the MongoDB instance
//   - dbName: The name of the database to connect to
//
// Returns:
//   - *MongoDB: A struct containing the MongoDB client and database reference
//   - error: An error if the connection to MongoDB fails
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

	// Get a reference to the database
	database := client.Database(dbName)

	return &MongoDB{Client: client, Database: database}, nil
}

// Close disconnects the MongoDB client and logs an error if the disconnection fails
// It prints a message indicating that the disconnection was successful
//
// Parameters:
//   - db: The MongoDB struct containing the client to disconnect
func (db *MongoDB) Close() {
	err := db.Client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
