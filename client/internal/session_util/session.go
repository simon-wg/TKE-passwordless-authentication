package session_util

import (
	"github.com/gorilla/sessions"
)

// Global session store
var Store = sessions.NewCookieStore([]byte("your-secret-key"))
