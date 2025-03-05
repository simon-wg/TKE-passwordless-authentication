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

func (repo *UserRepo) DeleteUser(userName string) (*mongo.DeleteResult, error) {
	collection := repo.db.Collection("users")

	filter := bson.M{"username": userName}
	result, err := collection.DeleteOne(context.Background(), filter)
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
		publicKeys[i] = string(key)
	}
	return publicKeys
}
