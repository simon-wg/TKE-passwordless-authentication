package internal

import (
	"encoding/csv"
	"os"
)

/*
	Read returns a map where each username is associated with a publickey.
*/

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

/*
	By calling Write with a filepath, username and a publickey will make a new
	entry in the csv file with the specified properties.
*/

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
