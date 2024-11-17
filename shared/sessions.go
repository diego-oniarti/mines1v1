package shared

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var Store *sessions.CookieStore 

func init() {
    err := godotenv.Load(".env")
    if err != nil { panic("Error loading .env file") }
    Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
    log.Println(os.Getenv("SESSION_SECRET"))
    Store.Options = &sessions.Options{
    	Path:        "/",
    	Domain:      "localhost",
    	MaxAge:      3600*8,
    	HttpOnly:    true,
    }
}
