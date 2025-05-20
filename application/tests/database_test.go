package tests

import (
	dbconnect "chalmers/tkey-group22/application/data/db"
	"chalmers/tkey-group22/application/internal/util"
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"strconv"

	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

const testDBName = "testdb"
const testUser = "testuser"
const testLabel = "label"

func setupTestDB(t *testing.T) (*mongo.Client, *util.UserRepo) {
	// Connect to the test database
	inst, err := dbconnect.ConnectMongoDB("mongodb://localhost:27017", testDBName)

	if err != nil {
		log.Fatal(err)
	}

	client := inst.Client
	database := inst.Database

	testdb := client.Database(testDBName)
	repo := util.NewUserRepo(testdb)

	// Clean test database after test

	t.Cleanup(func() {
		database.Drop(context.Background())
		client.Disconnect(context.Background())
	})

	return client, repo
}

func TestCreateUser(t *testing.T) {

	_, repo := setupTestDB(t)

	// Generate a new key pair
	pubkey, _, err := ed25519.GenerateKey(nil)
	assert.NoError(t, err)

	// Create a new user
	result, err := repo.CreateUser(testUser, pubkey, testLabel)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	// Check that the user was created and is stored correctly in the database
	user, err := repo.GetUser(testUser)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUser, user.Username)
	assert.Equal(t, base64.StdEncoding.EncodeToString(pubkey), user.PublicKeys[0].Key)
	assert.Equal(t, 1, len(user.PublicKeys))

}

func TestDeleteUser(t *testing.T) {

	_, repo := setupTestDB(t)

	// Generate a new key pair
	pubkey, _, err := ed25519.GenerateKey(nil)
	assert.NoError(t, err)

	// Create a new user
	result, err := repo.CreateUser(testUser, pubkey, testLabel)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Remove user
	repo.DeleteUser(testUser)

	// Check that user is not in database
	user, err := repo.GetUser(testUser)
	assert.Error(t, err)
	assert.Nil(t, user)

}

func TestUpdateUser(t *testing.T) {

	_, repo := setupTestDB(t)

	// Generate a new key pair
	pubkey, _, err := ed25519.GenerateKey(nil)
	assert.NoError(t, err)

	// Create a new user
	result, err := repo.CreateUser(testUser, pubkey, testLabel)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	newPubKey := ed25519.PublicKey([]byte("updatedpublickey"))

	newUser := util.User{
		Username: "newUser",
		PublicKeys: []util.PublicKey{
			{
				Key:   base64.StdEncoding.EncodeToString(newPubKey),
				Label: "newLabel",
			},
		},
	}

	// Update user data
	repo.UpdateUser(testUser, newUser)

	// Check that the old user data is no longer in database
	user, err := repo.GetUser(testUser)
	assert.Nil(t, user)
	assert.Error(t, err)

	// Check that new user data is in the database
	user, err = repo.GetUser(newUser.Username)
	assert.NotNil(t, user)
	assert.NoError(t, err)
	assert.Equal(t, newUser.Username, user.Username)
	assert.Equal(t, 1, len(user.PublicKeys))
	assert.Equal(t, newUser.PublicKeys[0], user.PublicKeys[0])
	assert.Equal(t, newUser.PublicKeys[0].Label, user.PublicKeys[0].Label)

}

func TestGetUser(t *testing.T) {

	_, repo := setupTestDB(t)

	// Database should return null if requesting non existing user
	user, err := repo.GetUser("DONOTEXIST")
	assert.Error(t, err)
	assert.Nil(t, user)

	// Generate a new key pair
	pubkey, _, err := ed25519.GenerateKey(nil)
	assert.NoError(t, err)

	// Create a new user
	result, err := repo.CreateUser(testUser, pubkey, testLabel)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// GetUser should return correct user
	user, err = repo.GetUser(testUser)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUser, user.Username)
	assert.Equal(t, testLabel, user.PublicKeys[0].Label)
	assert.Equal(t, 1, len(user.PublicKeys))

}

func TestAddPublicKey(t *testing.T) {
	_, repo := setupTestDB(t)

	username := "testuser"
	initialPubkey := ed25519.PublicKey([]byte("initialpublickey"))
	initialLabel := "initialKey"

	// Create the user with the initial public key
	_, err := repo.CreateUser(username, initialPubkey, initialLabel)
	assert.NoError(t, err)

	// Add a new public key to the existing user
	newPubkey := ed25519.PublicKey([]byte("newpublickey"))
	newLabel := "newKey"
	result, err := repo.AddPublicKey(username, newPubkey, newLabel)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify the new public key was added
	user, err := repo.GetUser(username)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 2, len(user.PublicKeys))
	assert.Equal(t, base64.StdEncoding.EncodeToString(initialPubkey), user.PublicKeys[0].Key)
	assert.Equal(t, initialLabel, user.PublicKeys[0].Label)
	assert.Equal(t, base64.StdEncoding.EncodeToString(newPubkey), user.PublicKeys[1].Key)
	assert.Equal(t, newLabel, user.PublicKeys[1].Label)

	// Try to add the same public key again
	_, err = repo.AddPublicKey(username, newPubkey, newLabel)
	assert.Error(t, err)
	assert.Equal(t, "public key already exists for the user", err.Error())

	// Try to add a new public key with an existing label
	anotherPubkey := ed25519.PublicKey([]byte("anotherpublickey"))
	_, err = repo.AddPublicKey(username, anotherPubkey, newLabel)
	assert.Error(t, err)
	assert.Equal(t, "label already exists for the user", err.Error())

	// Add more public keys until the maximum limit is reached
	for i := 2; i < util.MaxPublicKeys; i++ {
		pubkey := ed25519.PublicKey([]byte("pubkey" + strconv.Itoa(i)))
		label := "key" + strconv.Itoa(i)
		_, err := repo.AddPublicKey(username, pubkey, label)
		assert.NoError(t, err)
	}

	// Try to add another public key beyond the maximum limit
	extraPubkey := ed25519.PublicKey([]byte("extrapubkey"))
	extraLabel := "extraKey"
	_, err = repo.AddPublicKey(username, extraPubkey, extraLabel)
	assert.Error(t, err)
	assert.Equal(t, "user already has the maximum number of public keys", err.Error())
}

func TestRemovePublicKey(t *testing.T) {
	_, repo := setupTestDB(t)

	username := "testuser"
	initialPubkey := ed25519.PublicKey([]byte("initialpublickey"))
	initialLabel := "initiaKey"

	// Create the user with the initial public key
	_, err := repo.CreateUser(username, initialPubkey, initialLabel)
	assert.NoError(t, err)

	// Add a new public key to the existing user
	newPubkey := ed25519.PublicKey([]byte("newpublickey"))
	newLabel := "neKey"
	_, err = repo.AddPublicKey(username, newPubkey, newLabel)
	assert.NoError(t, err)

	// Remove the new public key
	result, err := repo.RemovePublicKey(username, newLabel)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify the public key was removed
	user, err := repo.GetUser(username)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, len(user.PublicKeys))
	assert.Equal(t, base64.StdEncoding.EncodeToString(initialPubkey), user.PublicKeys[0].Key)
	assert.Equal(t, initialLabel, user.PublicKeys[0].Label)

	// Try to remove the last remaining public key
	_, err = repo.RemovePublicKey(username, initialLabel)
	assert.Error(t, err)
	assert.Equal(t, "user must have at least two public keys to remove one", err.Error())
}

func TestGetPublicKeyLabels(t *testing.T) {
	_, repo := setupTestDB(t)

	username := "testuser"
	initialPubkey := ed25519.PublicKey([]byte("initialpublickey"))
	initialLabel := "initiaKey"

	// Create the user with the initial public key
	_, err := repo.CreateUser(username, initialPubkey, initialLabel)
	assert.NoError(t, err)

	// Add a new public key to the existing user
	newPubkey := ed25519.PublicKey([]byte("newpublickey"))
	newLabel := "newKey"
	_, err = repo.AddPublicKey(username, newPubkey, newLabel)
	assert.NoError(t, err)

	// Retrieve the public key labels
	labels, err := repo.GetPublicKeyLabels(username)
	assert.NoError(t, err)
	assert.NotNil(t, labels)
	assert.Equal(t, 2, len(labels))
	assert.Contains(t, labels, initialLabel)
	assert.Contains(t, labels, newLabel)
}
