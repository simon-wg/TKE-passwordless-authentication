package util

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User struct represents a user in DB
// It contains the user's unique ID, username and public key
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // Unique ID set by MongoDB
	Username  string             `bson:"username"`      // Username of the user
	PublicKey string             `bson:"publicKey"`     // Public key of the user
}

// Interface for UserRepository
// This interface defines the methods that a UserRepository should implement
type UserRepository interface {
	CreateUser(userName string, pubkey ed25519.PublicKey) (*mongo.InsertOneResult, error)
	GetUser(username string) (*User, error)
	UpdateUser(userName string, updatedUser User) (*mongo.UpdateResult, error)
	DeleteUser(userName string) (*mongo.DeleteResult, error)
}

// UserRepo holds the database reference
type UserRepo struct {
	db *mongo.Database
}

// NewUserRepo initializes a new UserRepositoryImpl with a given database
// It returns a pointer to the new UserRepo
//
// Parameters:
//   - db: The MongoDB database reference
//
// Returns:
//   - *UserRepo: A pointer to the new UserRepo
func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{db: db}
}

// CreateUser inserts a new user with the specified username and public key into the MongoDB collection
// The public key is encoded to base64 before storing
//
// Parameters:
//   - userName: The username of the new user
//   - pubkey: The ed25519 public key of the new user
//
// Returns:
//   - *mongo.InsertOneResult: The result of the insert operation
//   - error: An error if the insert operation fails
func (repo *UserRepo) CreateUser(userName string, pubkey ed25519.PublicKey) (*mongo.InsertOneResult, error) {
	collection := repo.db.Collection("users")

	// Encodes public key to base64 to allow storing in MongoDB
	encodedPubKey := base64.StdEncoding.EncodeToString(pubkey)
	user := User{
		ID:        primitive.NewObjectID(),
		Username:  userName,
		PublicKey: string(encodedPubKey),
	}

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetUser retrieves a user from the database by their username
// It returns a pointer to a User struct and an error if the retrieval fails
//
// Parameters:
//   - userName: The username of the user to retrieve
//
// Returns:
//   - *User: A pointer to the User struct containing the user's information
//   - error: An error if the retrieval fails
func (repo *UserRepo) GetUser(userName string) (*User, error) {
	collection := repo.db.Collection("users")

	filter := bson.M{"username": userName}
	var user User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	decodedKey, err := DecodePublicKey(user.PublicKey)

	if err != nil {
		return nil, err
	}

	user.PublicKey = string(decodedKey)

	return &user, nil
}

// UpdateUser updates the user document in the MongoDB collection with the given username
//
// Parameters:
//   - userName: The username of the user to be updated
//   - updatedUser: A User struct containing the new values for the username and public key
//
// Returns:
//   - *mongo.UpdateResult: The result of the update operation
//   - error: An error if the update operation fails
func (repo *UserRepo) UpdateUser(userName string, updatedUser User) (*mongo.UpdateResult, error) {
	collection := repo.db.Collection("users")

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

// DeleteUser deletes a user from the "users" collection in the MongoDB database
// It returns a pointer to mongo.DeleteResult and an error if the deletion fails
//
// Parameters:
//   - userName: The username of the user to be deleted
//
// Returns:
//   - *mongo.DeleteResult: The result of the delete operation
//   - error: An error if the deletion fails, otherwise nil
func (repo *UserRepo) DeleteUser(userName string) (*mongo.DeleteResult, error) {
	collection := repo.db.Collection("users")

	filter := bson.M{"username": userName}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DecudePublicKey decodes a public key from base64 format to ed25519.PublicKey
// It returns the decoded public key and an error if the decoding fails
//
// Parameters:
//   - encodedKey: The public key encoded in base64 format
//
// Returns:
//   - ed25519.PublicKey: The decoded public key
//   - error: An error if the decoding fails
func DecodePublicKey(encodedKey string) (ed25519.PublicKey, error) {
	data, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, err
	}
	return ed25519.PublicKey(data), nil
}
