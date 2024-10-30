package main

import (
    "github.com/gin-gonic/gin"
)

// SessionMiddleware checks for an active session and applies it only where required
func SessionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        session, _ := store.Get(c.Request, "session-name")

        // Set authentication status in the context
        if auth, ok := session.Values["authenticated"].(bool); ok && auth {
            c.Set("isAuthenticated", true)
            c.Set("username", session.Values["username"])
        } else {
            c.Set("isAuthenticated", false)
        }

        // Only protect specific routes
        // if c.FullPath() == "/" && !c.GetBool("isAuthenticated") {
        //     c.Redirect(http.StatusTemporaryRedirect, "/login")
        //     c.Abort()
        //     return
        // }

        // Proceed with the request
        c.Next()
    }
}
