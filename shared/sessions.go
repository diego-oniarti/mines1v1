package shared

import (
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

func init() {
    Store.Options = &sessions.Options{
    	Path:        "/",
    	Domain:      "localhost",
    	MaxAge:      3600*8,
    	HttpOnly:    true,
    }
}
