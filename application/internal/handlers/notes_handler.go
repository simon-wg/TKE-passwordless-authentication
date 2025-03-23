package handlers

import (
	"chalmers/tkey-group22/application/internal/session_util"
	"chalmers/tkey-group22/application/internal/structs"
	"chalmers/tkey-group22/application/internal/util"
	"encoding/json"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var NotesRepo util.NotesRepository

// GetNotesHandler handles HTTP GET requests to retrieve notes for a signed-in user
// It checks if the request method is GET, retrieves the username from the session,
// fetches the notes for the user from the NotesRepo, converts the notes to JSON,
// and sends the JSON response back to the client
//
// Possible responses:
// - 405 Method Not Allowed: if the request method is not GET
// - 401 Unauthorized: if there is no user signed in
// - 500 Internal Server Error: if there is an error marshalling the notes to JSON
// - 200 OK: if the notes are retrieved and marshalled successfully
func GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username, err := session_util.GetSessionUsername(r)
	if err != nil {
		http.Error(w, "No user signed in", http.StatusUnauthorized)
		return
	}
	notes, _ := NotesRepo.GetNotes(username)

	// Convert notes to JSON
	responseBodyBytes, err := json.Marshal(notes)
	if err != nil {
		http.Error(w, "Unable to marshal notes", http.StatusInternalServerError)
		return
	}

	// Send the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes)
}

// CreateNoteHandler handles HTTP POST requests to create a new note
// It checks if the request method is POST, reads and unmarshals the request body,
// retrieves the username from the session, saves the note using NotesRepo,
// and sends a JSON response with a success message and the ID of the created note
//
// Possible responses:
// - 405 Method Not Allowed: if the request method is not POST
// - 400 Bad Request: if the request body is invalid
// - 401 Unauthorized: if there is no user signed in
// - 500 Internal Server Error: if there is an error saving the note or marshalling the response
// - 200 OK: if the note is created and the response is marshalled successfully
func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	requestBody := structs.SaveNoteRequest{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	name := requestBody.Name
	note := requestBody.Note
	username, err := session_util.GetSessionUsername(r)
	if err != nil {
		http.Error(w, "No user signed in", http.StatusUnauthorized)
		return
	}

	result, err := NotesRepo.CreateNote(username, name, note)
	if result == nil || err != nil {
		http.Error(w, "Failed to save notes", http.StatusInternalServerError)
		return
	}

	responseBody := map[string]interface{}{
		"message": "Notes saved successfully",
		"id":      result.InsertedID.(primitive.ObjectID).Hex(),
	}
	responseBodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
		return
	}

	// Send the response
	sendJSONResponse(w, http.StatusOK, responseBodyBytes)

}

// UpdateNoteHandler handles HTTP POST requests to update an existing note
// It checks if the request method is POST, reads and unmarshals the request body,
// retrieves the username from the session, fetches the current note entry from the repository,
// checks if the current user is the owner of the note, updates the note in the repository,
// and sends a JSON response indicating the success or failure of the update operation
//
// Possible responses:
// - 405 Method Not Allowed: if the request method is not POST
// - 400 Bad Request: if the request body is invalid
// - 401 Unauthorized: if there is no user signed in or the user is not the owner of the note
// - 500 Internal Server Error: if there is an error retrieving or updating the note
// - 200 OK: if the note is updated successfully
func UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	requestBody := structs.UpdateNotesRequest{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username, _ := session_util.GetSessionUsername(r)
	currentEntry, err := NotesRepo.GetNote(requestBody.ID)
	if err != nil {
		http.Error(w, "Error retrieving entry", http.StatusInternalServerError)
	}

	if username != currentEntry.Username {
		http.Error(w, "User not owner of entry", http.StatusUnauthorized)
		return
	}

	result, err := NotesRepo.UpdateNote(requestBody.ID, username, requestBody.Name, requestBody.Note)
	if result == nil || err != nil {
		http.Error(w, "Failed to update note", http.StatusInternalServerError)
		return
	}

	// Send the response
	response := map[string]string{"message": "Note updated successfully"}
	sendJSONResponse(w, http.StatusOK, response)

}

// DeleteNoteHandler handles HTTP DELETE requests to delete a note
// It checks if the request method is DELETE, reads and unmarshals the request body,
// retrieves the username from the session, fetches the note entry from the repository,
// checks if the current user is the owner of the note, deletes the note from the repository,
// and sends a JSON response indicating the success or failure of the delete operation
//
// Possible responses:
// - 405 Method Not Allowed: if the request method is not DELETE
// - 400 Bad Request: if the request body is invalid
// - 401 Unauthorized: if there is no user signed in or the user is not the owner of the note
// - 500 Internal Server Error: if there is an error retrieving or deleting the note
// - 200 OK: if the note is deleted successfully
func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	requestBody := structs.DeleteNoteRequest{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username, err := session_util.GetSessionUsername(r)
	if err != nil {
		http.Error(w, "No user signed in", http.StatusUnauthorized)
		return
	}

	currentEntry, err := NotesRepo.GetNote(requestBody.ID)
	if err != nil {
		http.Error(w, "Error retrieving entry", http.StatusInternalServerError)
		return
	}

	if username != currentEntry.Username {
		http.Error(w, "User not owner of entry", http.StatusUnauthorized)
		return
	}

	result, err := NotesRepo.DeleteNote(requestBody.ID)
	if result == nil || err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	// Send the response
	response := map[string]string{"message": "Note deleted successfully"}
	sendJSONResponse(w, http.StatusOK, response)
}
