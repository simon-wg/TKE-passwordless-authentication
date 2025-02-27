package auth

import "fmt"

func Unregister() error {

	username := GetUsername()
	err := VerifyUser(username)

	if err != nil {
		return fmt.Errorf("user verification failed for %s", username)
	}

	return nil
}
