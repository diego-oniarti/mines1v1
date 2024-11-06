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
            c.Set("email", session.Values["email"])
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

func addUserDataMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Recupera i dati utente, ad esempio da sessione o database
        username := c.GetString("username") // Supponiamo che l'username sia già nel contesto
        email := c.GetString("email")       // Supponiamo che l'email sia già nel contesto

        // Aggiungi i dati utente al contesto per renderli disponibili nei template
        if c.GetBool("isAuthenticated") {
            c.Set("templateData", gin.H{
                "username": username,
                "email":    email,
            })
        }else{
            c.Set("templateData", gin.H{})
        }

        c.Next() // Continua alla richiesta successiva
    }
}
