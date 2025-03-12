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
	testDB.Drop(context.Background())
}

func TestCreateUser(t *testing.T) {
	setup()
	defer teardown()

	username := "testuser"
	pubkey := ed25519.PublicKey([]byte("testpublickey"))
	label := "initial key"

	result, err := userRepo.CreateUser(username, pubkey, label)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify user was created
	var user util.User
	err = testDB.Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	assert.NoError(t, err)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, 1, len(user.PublicKeys))
	assert.Equal(t, base64.StdEncoding.EncodeToString(pubkey), user.PublicKeys[0].Key)
	assert.Equal(t, label, user.PublicKeys[0].Label)
}

func TestGetUser(t *testing.T) {
	setup()
	defer teardown()

	username := "testuser"
	pubkey := ed25519.PublicKey([]byte("testpublickey"))
	label := "initial key"

	_, err := userRepo.CreateUser(username, pubkey, label)
	assert.NoError(t, err)

	user, err := userRepo.GetUser(username)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, 1, len(user.PublicKeys))
	assert.Equal(t, base64.StdEncoding.EncodeToString(pubkey), user.PublicKeys[0].Key)
	assert.Equal(t, label, user.PublicKeys[0].Label)
}

func TestUpdateUser(t *testing.T) {
	setup()
	defer teardown()

	username := "testuser"
	pubkey := ed25519.PublicKey([]byte("testpublickey"))
	label := "initial key"

	_, err := userRepo.CreateUser(username, pubkey, label)
	assert.NoError(t, err)

	updatedPubkey := ed25519.PublicKey([]byte("updatedpublickey"))
	updatedLabel := "updated key"
	updatedUser := util.User{
		Username: username,
		PublicKeys: []util.PublicKey{
			{
				Key:   base64.StdEncoding.EncodeToString(updatedPubkey),
				Label: updatedLabel,
			},
		},
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
	assert.Equal(t, base64.StdEncoding.EncodeToString(updatedPubkey), user.PublicKeys[0].Key)
	assert.Equal(t, updatedLabel, user.PublicKeys[0].Label)
}

func TestDeleteUser(t *testing.T) {
	setup()
	defer teardown()

	username := "testuser"
	pubkey := ed25519.PublicKey([]byte("testpublickey"))
	label := "initial key"

	_, err := userRepo.CreateUser(username, pubkey, label)
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
	initialLabel := "initial key"

	// Create the user with the initial public key
	_, err := userRepo.CreateUser(username, initialPubkey, initialLabel)
	assert.NoError(t, err)

	// Add a new public key to the existing user
	newPubkey := ed25519.PublicKey([]byte("newpublickey"))
	newLabel := "new key"
	result, err := userRepo.AddPublicKey(username, newPubkey, newLabel)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify the new public key was added
	user, err := userRepo.GetUser(username)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 2, len(user.PublicKeys))
	assert.Equal(t, base64.StdEncoding.EncodeToString(initialPubkey), user.PublicKeys[0].Key)
	assert.Equal(t, initialLabel, user.PublicKeys[0].Label)
	assert.Equal(t, base64.StdEncoding.EncodeToString(newPubkey), user.PublicKeys[1].Key)
	assert.Equal(t, newLabel, user.PublicKeys[1].Label)

	// Try to add the same public key again
	_, err = userRepo.AddPublicKey(username, newPubkey, newLabel)
	assert.Error(t, err)
	assert.Equal(t, "public key already exists for the user", err.Error())

	// Try to add a new public key with an existing label
	anotherPubkey := ed25519.PublicKey([]byte("anotherpublickey"))
	_, err = userRepo.AddPublicKey(username, anotherPubkey, newLabel)
	assert.Error(t, err)
	assert.Equal(t, "label already exists for the user", err.Error())

	// Add more public keys until the maximum limit is reached
	for i := 2; i < util.MaxPublicKeys; i++ {
		pubkey := ed25519.PublicKey([]byte("pubkey" + strconv.Itoa(i)))
		label := "key" + strconv.Itoa(i)
		_, err := userRepo.AddPublicKey(username, pubkey, label)
		assert.NoError(t, err)
	}

	// Try to add another public key beyond the maximum limit
	extraPubkey := ed25519.PublicKey([]byte("extrapubkey"))
	extraLabel := "extra key"
	_, err = userRepo.AddPublicKey(username, extraPubkey, extraLabel)
	assert.Error(t, err)
	assert.Equal(t, "user already has the maximum number of public keys", err.Error())
}

func TestRemovePublicKey(t *testing.T) {
	setup()
	defer teardown()

	username := "testuser"
	initialPubkey := ed25519.PublicKey([]byte("initialpublickey"))
	initialLabel := "initial key"

	// Create the user with the initial public key
	_, err := userRepo.CreateUser(username, initialPubkey, initialLabel)
	assert.NoError(t, err)

	// Add a new public key to the existing user
	newPubkey := ed25519.PublicKey([]byte("newpublickey"))
	newLabel := "new key"
	_, err = userRepo.AddPublicKey(username, newPubkey, newLabel)
	assert.NoError(t, err)

	// Remove the new public key
	result, err := userRepo.RemovePublicKey(username, newLabel)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify the public key was removed
	user, err := userRepo.GetUser(username)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, len(user.PublicKeys))
	assert.Equal(t, base64.StdEncoding.EncodeToString(initialPubkey), user.PublicKeys[0].Key)
	assert.Equal(t, initialLabel, user.PublicKeys[0].Label)

	// Try to remove the last remaining public key
	_, err = userRepo.RemovePublicKey(username, initialLabel)
	assert.Error(t, err)
	assert.Equal(t, "user must have at least two public keys to remove one", err.Error())
}

func TestGetPublicKeyLabels(t *testing.T) {
	setup()
	defer teardown()

	username := "testuser"
	initialPubkey := ed25519.PublicKey([]byte("initialpublickey"))
	initialLabel := "initial key"

	// Create the user with the initial public key
	_, err := userRepo.CreateUser(username, initialPubkey, initialLabel)
	assert.NoError(t, err)

	// Add a new public key to the existing user
	newPubkey := ed25519.PublicKey([]byte("newpublickey"))
	newLabel := "new key"
	_, err = userRepo.AddPublicKey(username, newPubkey, newLabel)
	assert.NoError(t, err)

	// Retrieve the public key labels
	labels, err := userRepo.GetPublicKeyLabels(username)
	assert.NoError(t, err)
	assert.NotNil(t, labels)
	assert.Equal(t, 2, len(labels))
	assert.Contains(t, labels, initialLabel)
	assert.Contains(t, labels, newLabel)
}
