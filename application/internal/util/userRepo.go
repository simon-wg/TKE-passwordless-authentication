package util

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User struct represents a user in DB
// It contains the user's unique ID, username and public key
type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"` // Unique ID set by MongoDB
	Username   string             `bson:"username"`      // Username of the user
	PublicKeys []PublicKey        `bson:"publicKeys"`    // Public key of the user
}

type PublicKey struct {
	Label string `bson:"label"` // Label for the public key
	Key   string `bson:"key"`   // Public key encoded in base64
}

// Interface for UserRepository
// This interface defines the methods that a UserRepository should implement
type UserRepository interface {
	CreateUser(userName string, pubkey ed25519.PublicKey, label string) (*mongo.InsertOneResult, error)
	GetUser(username string) (*User, error)
	UpdateUser(userName string, updatedUser User) (*mongo.UpdateResult, error)
	DeleteUser(userName string) (*mongo.DeleteResult, error)
	AddPublicKey(userName string, newPubKey ed25519.PublicKey, label string) (*mongo.UpdateResult, error)
	RemovePublicKey(userName string, label string) (*mongo.UpdateResult, error)
	GetPublicKeyLabels(userName string) ([]string, error)
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

// Max num of keys a single user can have
const MaxPublicKeys = 5

// CreateUser inserts a new user with the specified username and public key and label into the MongoDB collection.
// The public key is encoded to base64 before storing.
//
// Parameters:
//   - userName: The username of the new user.
//   - pubkey: The ed25519 public key of the new user.
//   - label: The label for the public key.
//
// Returns:
//   - *mongo.InsertOneResult: The result of the insert operation.
//   - error: An error if the insert operation fails.
func (repo *UserRepo) CreateUser(userName string, pubkey ed25519.PublicKey, label string) (*mongo.InsertOneResult, error) {
	collection := repo.db.Collection("users")

	// Encodes public key to base64 to allow storing in MongoDB
	encodedPubKey := base64.StdEncoding.EncodeToString(pubkey)
	println(pubkey)
	user := User{
		ID:       primitive.NewObjectID(),
		Username: userName,
		PublicKeys: []PublicKey{
			{
				Key:   encodedPubKey,
				Label: label,
			},
		},
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
			"username":   updatedUser.Username,
			"publicKeys": updatedUser.PublicKeys,
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

// GetPublicKeyLabels retrieves all the labels of the public keys associated with the given user
//
// Parameters:
//   - userName: The username of the user
//
// Returns:
//   - []string: A slice of labels for the user's public keys
//   - error: An error if the retrieval fails
func (repo *UserRepo) GetPublicKeyLabels(userName string) ([]string, error) {
	user, err := repo.GetUser(userName)
	if err != nil {
		return nil, err
	}

	labels := make([]string, len(user.PublicKeys))
	for i, pubkey := range user.PublicKeys {
		labels[i] = pubkey.Label
	}

	return labels, nil
}

// AddPublicKey adds a new public key to the user's list of public keys.
// It encodes the new public key to base64 and updates the user's document in the MongoDB collection.
//
// Parameters:
//   - userName: The username of the user to be updated.
//   - newPubKey: The new ed25519 public key to be added.
//   - label: The label for the new public key.
//
// Returns:
//   - *mongo.UpdateResult: The result of the update operation.
//   - error: An error if the update operation fails.
func (repo *UserRepo) AddPublicKey(userName string, newPubKey ed25519.PublicKey, label string) (*mongo.UpdateResult, error) {
	user, err := repo.GetUser(userName)
	if err != nil {
		return nil, err
	}

	if len(user.PublicKeys) >= MaxPublicKeys {
		return nil, errors.New("user already has the maximum number of public keys")
	}

	encodedPubKey := base64.StdEncoding.EncodeToString(newPubKey)

	for _, pubkey := range user.PublicKeys {
		if pubkey.Key == encodedPubKey {
			return nil, errors.New("public key already exists for the user")
		}
		if pubkey.Label == label {
			return nil, errors.New("label already exists for the user")
		}

	}

	user.PublicKeys = append(user.PublicKeys, PublicKey{
		Key:   encodedPubKey,
		Label: label,
	})

	result, err := repo.UpdateUser(userName, *user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// RemovePublicKey removes an existing public key from the user's list of public keys.
// It encodes the public key to be removed to base64 and updates the user's document in the MongoDB collection.
//
// Parameters:
//   - userName: The username of the user to be updated.
//   - label: The label of the public key to be removed.
//
// Returns:
//   - *mongo.UpdateResult: The result of the update operation.
//   - error: An error if the update operation fails.
func (repo *UserRepo) RemovePublicKey(userName string, label string) (*mongo.UpdateResult, error) {
	user, err := repo.GetUser(userName)
	if err != nil {
		return nil, err
	}

	if len(user.PublicKeys) <= 1 {
		return nil, errors.New("user must have at least two public keys to remove one")
	}

	keyFound := false
	for i, pubkey := range user.PublicKeys {
		if pubkey.Label == label {
			user.PublicKeys = append(user.PublicKeys[:i], user.PublicKeys[i+1:]...)
			keyFound = true
			break
		}
	}

	if !keyFound {
		return nil, errors.New("specified public key to be removed is not found")
	}

	result, err := repo.UpdateUser(userName, *user)
	if err != nil {
		return nil, err
	}

	return result, nil
}
