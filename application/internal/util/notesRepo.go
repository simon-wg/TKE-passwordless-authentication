package util

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NoteData struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"` // Unique ID set by MongoDB
	Username string             `bson:"username"`      // Username of the user
	Name     string             `bson:"name"`          // Name of company/website for password
	Password string             `bson:"password"`      // Password
}

type NotesRepository interface {
	CreatePassword(username string, name string, password string) (*mongo.InsertOneResult, error)
	GetUserPasswords(username string) ([]NoteData, error)
	GetPassword(id string) (NoteData, error)
	UpdatePassword(id string, username string, name string, password string) (*mongo.UpdateResult, error)
	DeletePassword(id string) (*mongo.DeleteResult, error)
}

type PasswordRepo struct {
	db *mongo.Database
}

const repoName = "userPasswords"

func NewPasswordRepo(db *mongo.Database) *PasswordRepo {
	return &PasswordRepo{db: db}
}

func (repo *PasswordRepo) CreatePassword(username, name, password string) (*mongo.InsertOneResult, error) {
	collection := repo.db.Collection(repoName)

	user := NoteData{
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

func (repo *PasswordRepo) GetUserPasswords(username string) ([]NoteData, error) {
	collection := repo.db.Collection(repoName)

	filter := bson.M{"username": username}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []NoteData
	for cursor.Next(context.Background()) {
		var password NoteData
		if err := cursor.Decode(&password); err != nil {
			return nil, err
		}
		users = append(users, password)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *PasswordRepo) GetPassword(id string) (NoteData, error) {
	collection := repo.db.Collection(repoName)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return NoteData{}, err
	}

	filter := bson.M{"_id": objectID}
	var password NoteData
	err = collection.FindOne(context.Background(), filter).Decode(&password)
	if err != nil {
		return NoteData{}, err
	}

	return password, nil
}

func (repo *PasswordRepo) UpdatePassword(id string, username string, name string, password string) (*mongo.UpdateResult, error) {
	collection := repo.db.Collection(repoName)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	updatedData := bson.M{
		"$set": bson.M{
			"username": username,
			"name":     name,
			"password": password,
		},
	}

	result, err := collection.UpdateOne(context.Background(), filter, updatedData)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *PasswordRepo) DeletePassword(id string) (*mongo.DeleteResult, error) {
	collection := repo.db.Collection(repoName)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}
