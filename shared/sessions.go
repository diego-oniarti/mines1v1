package shared

import (
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
