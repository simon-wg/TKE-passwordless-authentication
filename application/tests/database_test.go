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
