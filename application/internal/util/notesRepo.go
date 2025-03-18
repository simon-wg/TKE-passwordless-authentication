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
	Name     string             `bson:"name"`          // Name of company/website for note
	Note     string             `bson:"note"`          // Note as a string
}

type NotesRepository interface {
	CreateNote(username string, name string, note string) (*mongo.InsertOneResult, error)
	GetNotes(username string) ([]NoteData, error)
	GetNote(id string) (NoteData, error)
	UpdateNote(id string, username string, name string, note string) (*mongo.UpdateResult, error)
	DeleteNote(id string) (*mongo.DeleteResult, error)
}

type NotesRepo struct {
	db *mongo.Database
}

const repoName = "user_notes"

func NewNotesRepo(db *mongo.Database) *NotesRepo {
	return &NotesRepo{db: db}
}

func (repo *NotesRepo) CreateNote(username, name, note string) (*mongo.InsertOneResult, error) {
	collection := repo.db.Collection(repoName)

	user := NoteData{
		ID:       primitive.NewObjectID(),
		Username: username,
		Name:     name,
		Note:     note,
	}

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *NotesRepo) GetNotes(username string) ([]NoteData, error) {
	collection := repo.db.Collection(repoName)

	filter := bson.M{"username": username}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []NoteData
	for cursor.Next(context.Background()) {
		var note NoteData
		if err := cursor.Decode(&note); err != nil {
			return nil, err
		}
		users = append(users, note)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *NotesRepo) GetNote(id string) (NoteData, error) {
	collection := repo.db.Collection(repoName)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return NoteData{}, err
	}

	filter := bson.M{"_id": objectID}
	var note NoteData
	err = collection.FindOne(context.Background(), filter).Decode(&note)
	if err != nil {
		return NoteData{}, err
	}

	return note, nil
}

func (repo *NotesRepo) UpdateNote(id string, username string, name string, note string) (*mongo.UpdateResult, error) {
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
			"note":     note,
		},
	}

	result, err := collection.UpdateOne(context.Background(), filter, updatedData)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *NotesRepo) DeleteNote(id string) (*mongo.DeleteResult, error) {
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
