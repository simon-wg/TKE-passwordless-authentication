package tests

import (
	"chalmers/tkey-group22/application/internal/util"
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testDB *mongo.Database
var userRepo *util.UserRepo

func setup() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	testDB = client.Database("testDB")
	userRepo = util.NewUserRepo(testDB)
}

func teardown() {
	testDB.Drop(context.Background())
}

func TestCreateUser(t *testing.T) {
	setup()
	defer teardown()

	username := "testuser"
	pubkey := ed25519.PublicKey([]byte("testpublickey"))

	result, err := userRepo.CreateUser(username, pubkey)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify user was created
	var user util.User
	err = testDB.Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	assert.NoError(t, err)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, 1, len(user.PublicKeys))
	assert.Equal(t, base64.StdEncoding.EncodeToString(pubkey), user.PublicKeys[0])
}

func TestGetUser(t *testing.T) {
	setup()
	defer teardown()

	username := "testuser"
	pubkey := ed25519.PublicKey([]byte("testpublickey"))

	_, err := userRepo.CreateUser(username, pubkey)
	assert.NoError(t, err)

	user, err := userRepo.GetUser(username)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, 1, len(user.PublicKeys))
	assert.Equal(t, string(pubkey), user.PublicKeys[0])
}

func TestUpdateUser(t *testing.T) {
	setup()
	defer teardown()

	username := "testuser"
	pubkey := ed25519.PublicKey([]byte("testpublickey"))

	_, err := userRepo.CreateUser(username, pubkey)
	assert.NoError(t, err)

	updatedPubkey := ed25519.PublicKey([]byte("updatedpublickey"))
	updatedUser := util.User{
		Username:   username,
		PublicKeys: []string{base64.StdEncoding.EncodeToString(updatedPubkey)},
	}

	result, err := userRepo.UpdateUser(username, updatedUser)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify user was updated
	user, err := userRepo.GetUser(username)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, 1, len(user.PublicKeys))
	assert.Equal(t, string(updatedPubkey), user.PublicKeys[0])
}

func TestDeleteUser(t *testing.T) {
	setup()
	defer teardown()

	username := "testuser"
	pubkey := ed25519.PublicKey([]byte("testpublickey"))

	_, err := userRepo.CreateUser(username, pubkey)
	assert.NoError(t, err)

	result, err := userRepo.DeleteUser(username)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify user was deleted
	var user util.User
	err = testDB.Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	assert.Error(t, err)
	assert.Equal(t, mongo.ErrNoDocuments, err)
}
