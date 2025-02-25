package db

import (
	"context"
	"crypto/ed25519"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User struct represents a user in DB
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // Unique ID assigned by MongoDB
	Username  string             `bson:"username"`      // Username
	PublicKey string             `bson:"publicKey"`     // Public key
}

// Interface for UserRepository
type UserRepository interface {
	CreateUser(user *User) (*mongo.InsertOneResult, error)
	GetUser(username string) (*User, error)
	UpdateUser(username string, publicKey string) (*mongo.UpdateResult, error)
	DeleteUser(username string) (*mongo.DeleteResult, error)
}

// CreateUser inserts a new user into the MongoDB database
func CreateUser(db *mongo.Database, userName string, pubkey ed25519.PublicKey) (*mongo.InsertOneResult, error) {

	collection := db.Collection("users")

	user := User{
		ID:        primitive.NewObjectID(),
		Username:  userName,
		PublicKey: string(pubkey), // Store public key as string
	}

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	return result, nil
}
