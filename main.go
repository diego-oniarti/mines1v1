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

    // Initialize Gin router
    r := gin.Default()
    r.LoadHTMLGlob("templates/*")
    r.Use(SessionMiddleware())

    // Serve static files from the "static" directory
    r.Static("/static", "./static")

    // Define routes
    r.GET("/", indexHandler)
    r.POST("/", indexHandler)
    r.GET("/login", loginPageHandler)
    r.POST("/login", loginHandler)
    r.POST("/register", registerHandler)
    r.POST("/verify", verifyHandler)

    // Start server on port 8080
    r.Run(":2357")
}
