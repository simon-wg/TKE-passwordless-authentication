package main

import (
	"chalmers/tkey-group22/application/internal"
	"fmt"
	"net/http"
)

func main() {
	challengeService := internal.NewED25519ChallengeService()

	http.HandleFunc("/api/register", internal.RegisterHandler)
	http.HandleFunc("/api/login", internal.LoginHandler(challengeService))
	http.HandleFunc("/api/verify", internal.VerifyHandler(challengeService))

	fmt.Println("Mock application running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
