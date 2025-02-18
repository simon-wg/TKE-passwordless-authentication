package internal

import (
	"encoding/csv"
	"os"
)

// Read reads a CSV file and returns a map where each username is associated
// with a public key.
//
// Parameters:
//   - filePath: The path to the CSV file to be read.
//
// Returns:
//   - map[string]string: A map where the keys are usernames and the values are
//     public keys.
//   - error: An error if one occurred during reading the file.
func Read(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]string)
	for _, record := range records {
		if len(record) < 2 {
			continue
		}
		username := record[0]
		publicKey := record[1]
		userMap[username] = publicKey
	}

	return userMap, nil
}

// Write appends a new entry with the specified username and public key to a
// CSV file.
//
// Parameters:
//   - filePath: The path to the CSV file to be written to.
//   - username: The username to be added.
//   - publicKey: The public key to be associated with the username.
//
// Returns:
//   - error: An error if one occurred during writing to the file.
func Write(filePath, username, publicKey string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if _, err := file.WriteString("\n"); err != nil {
		return err
	}

	record := []string{username, publicKey}
	if err := writer.Write(record); err != nil {
		return err
	}

	return nil
}
