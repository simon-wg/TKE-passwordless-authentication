package session_util

import (
	"github.com/gorilla/sessions"
)

// Global session store
// TODO: Move key to .env file
var Store = sessions.NewCookieStore([]byte("your-secret-key"))
