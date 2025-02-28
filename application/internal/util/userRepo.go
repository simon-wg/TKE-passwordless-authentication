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
	ID        primitive.ObjectID `bson:"_id,omitempty"` // Unique ID set by MongoDB
	Username  string             `bson:"username"`      // Username of the user
	PublicKey string             `bson:"publicKey"`     // Public key of the user
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

func (repo *UserRepo) DeleteUser(userName string) (*mongo.DeleteResult, error) {
	collection := repo.db.Collection("users")

	filter := bson.M{"username": userName}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Decodes Public Key from stored base64 format to ed25519.PublicKey
func DecodePublicKey(encodedKey string) (ed25519.PublicKey, error) {
	data, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, err
	}
	return ed25519.PublicKey(data), nil
}
