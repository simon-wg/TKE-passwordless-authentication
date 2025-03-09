package util

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PasswordData struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"` // Unique ID set by MongoDB
	Username string             `bson:"username"`      // Username of the user
	Name     string             `bson:"name"`          // Name of company/website for password
	Password string             `bson:"password"`      // Password
}

type PasswordRepository interface {
	CreatePassword(username string, name string, password string) (*mongo.InsertOneResult, error)
	GetUserPasswords(username string) ([]PasswordData, error) // Update return type to match PasswordRepo
	UpdatePassword(username string, updatedUser PasswordData) (*mongo.UpdateResult, error)
	DeletePassword(username string) (*mongo.DeleteResult, error)
}

type PasswordRepo struct {
	db *mongo.Database
}

func NewPasswordRepo(db *mongo.Database) *PasswordRepo {
	return &PasswordRepo{db: db}
}

func (repo *PasswordRepo) CreatePassword(username, name, password string) (*mongo.InsertOneResult, error) {
	collection := repo.db.Collection("companyUsers")

	user := PasswordData{
		ID:       primitive.NewObjectID(),
		Username: username,
		Name:     name,
		Password: password,
	}

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *PasswordRepo) GetUserPasswords(username string) ([]PasswordData, error) {
	collection := repo.db.Collection("companyUsers")

	filter := bson.M{"username": username}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []PasswordData
	for cursor.Next(context.Background()) {
		var user PasswordData
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *PasswordRepo) UpdatePassword(username string, updatedUser PasswordData) (*mongo.UpdateResult, error) {
	collection := repo.db.Collection("companyUsers")

	filter := bson.M{"username": username}
	updatedData := bson.M{
		"$set": bson.M{
			"username": username,
			"name":     updatedUser.Name,
			"password": updatedUser.Password,
		},
	}

	result, err := collection.UpdateOne(context.Background(), filter, updatedData)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *PasswordRepo) DeletePassword(username string) (*mongo.DeleteResult, error) {
	collection := repo.db.Collection("companyUsers")

	filter := bson.M{"username": username}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}
