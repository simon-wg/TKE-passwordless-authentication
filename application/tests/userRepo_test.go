package tests

import (
	"chalmers/tkey-group22/application/internal/util"
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"strconv"
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
	//testDB.Drop(context.Background())
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
	assert.Equal(t, base64.StdEncoding.EncodeToString(pubkey), user.PublicKeys[0])
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
	assert.Equal(t, base64.StdEncoding.EncodeToString(updatedPubkey), user.PublicKeys[0])
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

func TestAddPublicKey(t *testing.T) {
	setup()
	defer teardown()

	username := "testuser"
	initialPubkey := ed25519.PublicKey([]byte("initialpublickey"))

	// Create the user with the initial public key
	_, err := userRepo.CreateUser(username, initialPubkey)
	assert.NoError(t, err)

	// Add a new public key to the existing user
	newPubkey := ed25519.PublicKey([]byte("newpublickey"))
	result, err := userRepo.AddPublicKey(username, newPubkey)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify the new public key was added
	user, err := userRepo.GetUser(username)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 2, len(user.PublicKeys))
	assert.Equal(t, base64.StdEncoding.EncodeToString(initialPubkey), user.PublicKeys[0])
	assert.Equal(t, base64.StdEncoding.EncodeToString(newPubkey), user.PublicKeys[1])

	// Add more public keys until the maximum limit is reached
	for i := 2; i < util.MaxPublicKeys; i++ {
		pubkey := ed25519.PublicKey([]byte("pubkey" + strconv.Itoa(i)))
		_, err := userRepo.AddPublicKey(username, pubkey)
		assert.NoError(t, err)
	}

	// Try to add another public key beyond the maximum limit
	extraPubkey := ed25519.PublicKey([]byte("extrapubkey"))
	_, err = userRepo.AddPublicKey(username, extraPubkey)
	assert.Error(t, err)
	assert.Equal(t, "user already has the maximum number of public keys", err.Error())
}

func TestRemovePublicKey(t *testing.T) {
	setup()
	defer teardown()

	username := "testuser"
	initialPubkey := ed25519.PublicKey([]byte("initialpublickey"))

	// Create the user with the initial public key
	_, err := userRepo.CreateUser(username, initialPubkey)
	assert.NoError(t, err)

	// Add a new public key to the existing user
	newPubkey := ed25519.PublicKey([]byte("newpublickey"))
	_, err = userRepo.AddPublicKey(username, newPubkey)
	assert.NoError(t, err)

	// Remove the new public key
	result, err := userRepo.RemovePublicKey(username, newPubkey)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify the public key was removed
	user, err := userRepo.GetUser(username)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, len(user.PublicKeys))
	assert.Equal(t, base64.StdEncoding.EncodeToString(initialPubkey), user.PublicKeys[0])

	// Try to remove the last remaining public key
	_, err = userRepo.RemovePublicKey(username, initialPubkey)
	assert.Error(t, err)
	assert.Equal(t, "user must have at least two public keys to remove one", err.Error())
}
