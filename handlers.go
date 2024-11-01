package main

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "io"
    "log"
    "net/http"
)

// RenderIndexPage renders the index page for logged-in users
func indexHandler(c *gin.Context) {
    username := c.GetString("username")
    c.HTML(http.StatusOK, "index.html", gin.H{"username": username})
}

// Render the login page
func loginPageHandler(c *gin.Context) {
    c.HTML(http.StatusOK, "login.html", nil)
}

// Render the profile page
func profilePageHandler(c *gin.Context) {
    isAuthenticated := c.GetBool("isAuthenticated")
    if !isAuthenticated {
        c.Redirect(http.StatusTemporaryRedirect, "/login")
        return;
    }
    c.HTML(http.StatusOK, "user.html", nil)
}


// Handle login with password verification
func loginHandler(c *gin.Context) {
    mail := c.PostForm("mail")
    password := c.PostForm("password")

    var user User
    err := db.QueryRow("SELECT username, psw, confirmed FROM users WHERE mail=?", mail).
        Scan(&user.Username, &user.Psw, &user.Confirmed)
    if err == sql.ErrNoRows || bcrypt.CompareHashAndPassword([]byte(user.Psw), []byte(password)) != nil {
        c.HTML(http.StatusUnauthorized, "login.html", gin.H{"LoginError": "Invalid email or password"})
        return
    } else if err != nil {
        log.Println("Database error:", err)
        c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Internal server error"})
        return
    }

    // Set session for authenticated user
    session, _ := store.Get(c.Request, "session-name")
    session.Values["authenticated"] = true
    session.Values["username"] = user.Username
    session.Save(c.Request, c.Writer)

    c.Redirect(http.StatusTemporaryRedirect, "/")
}

// Handle user registration with image upload, unique username and email checks
func registerHandler(c *gin.Context) {
    username := c.PostForm("username")
    mail := c.PostForm("mail")
    password := c.PostForm("password")

    // Image processing
    var imageData []byte
    image, _, err := c.Request.FormFile("image")
    if err == nil {
        defer image.Close()
        imageData, err = io.ReadAll(image)
        if err != nil {
            c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Error reading image file"})
            return
        }
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error hashing password:", err)
        c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Error processing registration"})
        return
    }

    // Input validation
    var usernameExists, mailExists bool
    db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE name=?)", username).Scan(&usernameExists)
    db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE mail=?)", mail).Scan(&mailExists)

    templValues := gin.H{"UsernameError": usernameExists, "MailError": mailExists, "PasswordError": false}
    if usernameExists || mailExists {
        c.HTML(http.StatusConflict, "register.html", templValues)
        return
    }

    // Insert the user into the database
    _, err = db.Exec("INSERT INTO users (name, image, psw, mail, confirmed) VALUES (?, ?, ?, ?, false)",
        username, imageData, hashedPassword, mail)
    if err != nil {
        log.Println("Database error during registration:", err)
        c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Error creating user"})
        return
    }

    // Automatically log in user after registration
    session, _ := store.Get(c.Request, "session-name")
    session.Values["authenticated"] = true
    session.Values["username"] = username
    session.Save(c.Request, c.Writer)

    c.HTML(http.StatusOK, "registration_complete.html", nil);
}
