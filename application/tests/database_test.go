package tests

import (
	dbconnect "chalmers/tkey-group22/application/data/db"
	"chalmers/tkey-group22/application/internal/util"
	"context"
	"crypto/ed25519"

	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

const testDBName = "testdb"
const testUser = "testuser"

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
	result, err := repo.CreateUser(testUser, pubkey)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Check that the user was created and is stored in the database
	user, err := repo.GetUser(testUser)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUser, user.Username)
	assert.Equal(t, string(pubkey), user.PublicKey)
}

func TestDeleteUser(t *testing.T) {

	_, repo := setupTestDB(t)

	// Generate a new key pair
	pubkey, _, err := ed25519.GenerateKey(nil)
	assert.NoError(t, err)

	// Create a new user
	result, err := repo.CreateUser(testUser, pubkey)
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
	result, err := repo.CreateUser(testUser, pubkey)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	newUser := util.User{
		Username:  "newTestUser",
		PublicKey: "123",
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
	assert.Equal(t, newUser.PublicKey, user.PublicKey)

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
	result, err := repo.CreateUser(testUser, pubkey)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// GetUser should return correct user
	user, err = repo.GetUser(testUser)
	assert.NoError(t, err)
	assert.NotNil(t, user)

}
