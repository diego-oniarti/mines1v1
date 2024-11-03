package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gorilla/sessions"
    "github.com/joho/godotenv"
    "os"
)

var store *sessions.CookieStore

func main() {
    err := godotenv.Load(".env")
    if err != nil {
        panic("Error loading .env file")
    }

    defer db.Close();

    sessionKey := []byte(os.Getenv("SESSION_SECRET"))
    store = sessions.NewCookieStore(sessionKey)

    r := gin.Default()
    r.LoadHTMLGlob("templates/*")
    r.Use(SessionMiddleware())

    r.Static("/static", "./static")

    r.GET("/", indexHandler)
    r.POST("/", indexHandler)
    r.GET("/login", loginPageHandler)
    r.POST("/login", loginHandler)
    r.POST("/register", registerHandler)
    r.POST("/verify", verifyHandler)
    r.GET("/user", profilePageHandler)

    r.Run(":2357")
}
