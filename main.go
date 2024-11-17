package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
    "github.com/diego-oniarti/mines1v1/gamemodes"
)


func main() {
    err := godotenv.Load(".env")
    if err != nil {
        panic("Error loading .env file")
    }

    defer db.Close()

    r := gin.Default()
    r.LoadHTMLGlob("templates/*")
    r.Use(SessionMiddleware())
    r.Use(addUserDataMiddleware())

    r.Static("/static", "./static")

    r.GET("/", indexHandler)
    r.POST("/", indexHandler)
    r.GET("/login", loginPageHandler)
    r.POST("/login", loginHandler)
    r.POST("/register", registerHandler)
    r.POST("/verify", verifyHandler)
    r.GET("/user", userPageHandler)
    r.POST("/logout", logoutHandler)
    r.POST("/deleteAccount", deleteAccountHandler)
    r.GET("/lobby", lobbyHandle)

    r.GET("/singlePlayer", gamemodes.SinglePlayerPage)
    r.GET("/wsSinglePlayer", gamemodes.SinglePlayerWs)

    r.Run(":2357")
}
