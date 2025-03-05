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
type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"` // Unique ID set by MongoDB
	Username   string             `bson:"username"`      // Username of the user
	PublicKeys []string           `bson:"publicKeys"`    // Public key of the user
}

// Interface for UserRepository
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
func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{db: db}
}

// Max num of keys a single user can have
const MaxPublicKeys = 5

// CreateUser inserts a new user with the specified username and public key into the MongoDB collection.
// The public key is encoded to base64 before storing.
//
// Parameters:
//   - userName: The username of the new user.
//   - pubkey: The ed25519 public key of the new user.
//
// Returns:
//   - *mongo.InsertOneResult: The result of the insert operation.
//   - error: An error if the insert operation fails.
func (repo *UserRepo) CreateUser(userName string, pubkey ed25519.PublicKey) (*mongo.InsertOneResult, error) {
	collection := repo.db.Collection("users")

	// Encodes public key to base64 to allow storing in MongoDB
	encodedPubKey := base64.StdEncoding.EncodeToString(pubkey)
	println(pubkey)
	user := User{
		ID:         primitive.NewObjectID(),
		Username:   userName,
		PublicKeys: []string{encodedPubKey},
	}

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetUser retrieves a user from the database by their username.
// It takes a username as a string and returns a pointer to a User object and an error.
// If the user is found, their public key is decoded and assigned to the PublicKey field of the User object.
// If any error occurs during the process, it returns nil and the error.
func (repo *UserRepo) GetUser(userName string) (*User, error) {
	collection := repo.db.Collection("users")

	filter := bson.M{"username": userName}
	var user User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	decodedKeys, err := decodePublicKeys(user.PublicKeys)

	if err != nil {
		return nil, err
	}

	user.PublicKeys = convertDecodedKeysToStrings(decodedKeys)

	return &user, nil
}

// UpdateUser updates the user document in the MongoDB collection with the given username.
// It replaces the username and public key fields with the values from the updatedUser parameter.
//
// Parameters:
//   - userName: The username of the user to be updated.
//   - updatedUser: A User struct containing the new values for the username and public key.
//
// Returns:
//   - *mongo.UpdateResult: The result of the update operation.
//   - error: An error if the update operation fails.
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

// DeleteUser deletes a user from the "users" collection in the MongoDB database.
// It takes a userName as a parameter, which specifies the username of the user to be deleted.
// It returns a pointer to mongo.DeleteResult and an error if the deletion fails.
//
// Parameters:
//   - userName: The username of the user to be deleted.
//
// Returns:
//   - *mongo.DeleteResult: The result of the delete operation.
//   - error: An error if the deletion fails, otherwise nil.
func (repo *UserRepo) DeleteUser(userName string) (*mongo.DeleteResult, error) {
	collection := repo.db.Collection("users")

	filter := bson.M{"username": userName}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// AddPublicKey adds a new public key to the user's list of public keys.
// It encodes the new public key to base64 and updates the user's document in the MongoDB collection.
//
// Parameters:
//   - userName: The username of the user to be updated.
//   - newPubKey: The new ed25519 public key to be added.
//
// Returns:
//   - *mongo.UpdateResult: The result of the update operation.
//   - error: An error if the update operation fails.
func (repo *UserRepo) AddPublicKey(userName string, newPubKey ed25519.PublicKey) (*mongo.UpdateResult, error) {
	user, err := repo.GetUser(userName)
	if err != nil {
		return nil, err
	}

	if len(user.PublicKeys) >= MaxPublicKeys {
		return nil, errors.New("user already has the maximum number of public keys")
	}

	encodedPubKey := base64.StdEncoding.EncodeToString(newPubKey)

	user.PublicKeys = append(user.PublicKeys, encodedPubKey)

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
//   - pubKeyToRemove: The ed25519 public key to be removed.
//
// Returns:
//   - *mongo.UpdateResult: The result of the update operation.
//   - error: An error if the update operation fails.
func (repo *UserRepo) RemovePublicKey(userName string, pubKeyToRemove ed25519.PublicKey) (*mongo.UpdateResult, error) {
	user, err := repo.GetUser(userName)
	if err != nil {
		return nil, err
	}

	if len(user.PublicKeys) <= 1 {
		return nil, errors.New("user must have at least two public keys to remove one")
	}

	encodedPubKey := base64.StdEncoding.EncodeToString(pubKeyToRemove)

	keyFound := false
	for i, key := range user.PublicKeys {
		if key == encodedPubKey {
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

// Decodes Public Keys from stored base64 format to a slice of ed25519.PublicKey
func decodePublicKeys(encodedKeys []string) ([]ed25519.PublicKey, error) {
	var publicKeys []ed25519.PublicKey
	for _, encodedKey := range encodedKeys {
		data, err := base64.StdEncoding.DecodeString(encodedKey)
		if err != nil {
			return nil, err
		}
		publicKeys = append(publicKeys, ed25519.PublicKey(data))
	}
	return publicKeys, nil
}

// convertDecodedKeysToStrings converts a slice of ed25519.PublicKey to a slice of strings
func convertDecodedKeysToStrings(decodedKeys []ed25519.PublicKey) []string {
	publicKeys := make([]string, len(decodedKeys))
	for i, key := range decodedKeys {
		publicKeys[i] = base64.StdEncoding.EncodeToString(key)
	}
	return publicKeys
}
