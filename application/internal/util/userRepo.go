package util

import (
	"context"
	"crypto/ed25519"

	"go.mongodb.org/mongo-driver/bson"
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

// Deletes user from database provided a userName
func DeleteUser(db *mongo.Database, userName string) (*mongo.DeleteResult, error) {

	collection := db.Collection("users")

	filter := bson.M{"username": userName}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Given username and public key, updates user in database with provided user struct.
func UpdateUser(db *mongo.Database, userName string, publicKey string, updatedUser User) (*mongo.UpdateResult, error) {

	collection := db.Collection("users")

	filter := bson.M{"username": userName}
	updatedData := bson.M{
		"$set": bson.M{
			"username":  updatedUser.Username,
			"publicKey": updatedUser.PublicKey,
		},
	}

	result, err := collection.UpdateOne(context.Background(), filter, updatedData)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Given a username, retrieves user from database and returns user struct.
func GetUser(db *mongo.Database, userName string) (*User, error) {

	collection := db.Collection("users")

	filter := bson.M{"username": userName}
	var user User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
