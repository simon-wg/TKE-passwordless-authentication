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

// GetNotesHandler handles HTTP GET requests to retrieve notes for a signed-in user.
// It checks if the request method is GET, retrieves the username from the session,
// fetches the notes for the user from the NotesRepo, converts the notes to JSON,
// and sends the JSON response back to the client.
//
// If the request method is not GET, it responds with "Invalid request method" and
// a 405 Method Not Allowed status code.
//
// If there is no user signed in, it responds with "No user signed in" and a 401
// Unauthorized status code.
//
// If there is an error marshalling the notes to JSON, it responds with "Unable to
// marshal notes" and a 500 Internal Server Error status code.
//
// The response content type is set to "application/json".
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

// CreateNoteHandler handles the creation of a new note.
// It expects a POST request with a JSON body containing the note details.
// The request body should be in the format:
//
//	{
//	  "name": "note name",
//	  "note": "note content"
//	}
//
// The function performs the following steps:
// 1. Validates that the request method is POST.
// 2. Reads and unmarshals the request body into a SaveNoteRequest struct.
// 3. Retrieves the username from the session.
// 4. Calls the NotesRepo.CreateNote function to save the note.
// 5. Returns a JSON response with a success message and the ID of the created note.
//
// If any step fails, an appropriate HTTP error response is returned.
//
// Parameters:
// - w: http.ResponseWriter to write the response.
// - r: *http.Request containing the request details.
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes)
}

// UpdateNoteHandler handles the HTTP request for updating a note.
// It expects a POST request with a JSON body containing the note details to be updated.
// The request body should match the structs.UpdateNotesRequest structure.
// The handler performs the following steps:
// 1. Validates the request method is POST.
// 2. Reads and unmarshals the request body into a structs.UpdateNotesRequest object.
// 3. Retrieves the username from the session.
// 4. Fetches the current note entry from the repository using the provided note ID.
// 5. Checks if the current user is the owner of the note.
// 6. Updates the note in the repository with the new details.
// 7. Returns a JSON response indicating the success or failure of the update operation.
//
// If any step fails, an appropriate HTTP error response is returned.
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

	responseBody := map[string]string{"message": "Note updated successfully"}
	responseBodyBytes, _ := json.Marshal(responseBody)

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes)
}

// DeleteNoteHandler handles the deletion of a note.
// It expects a DELETE request with a JSON body containing the note ID to be deleted.
// The handler performs the following steps:
// 1. Verifies that the request method is DELETE.
// 2. Reads and unmarshals the request body into a DeleteNoteRequest struct.
// 3. Retrieves the username from the session.
// 4. Fetches the note entry from the repository using the provided note ID.
// 5. Checks if the authenticated user is the owner of the note.
// 6. Deletes the note from the repository if the user is the owner.
// 7. Returns a success message in JSON format if the note is deleted successfully.
//
// If any of the steps fail, an appropriate HTTP error response is returned.
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

	responseBody := map[string]string{"message": "Note deleted successfully"}
	responseBodyBytes, _ := json.Marshal(responseBody)

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes)
}
