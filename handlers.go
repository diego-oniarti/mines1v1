package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
            c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Error reading image file"})
            return
        }
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error hashing password:", err)
        c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Error processing registration"})
        return
    }

    // Input validation
    var usernameExists, mailExists bool
    db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE name=?)", username).Scan(&usernameExists)
    db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE mail=?)", mail).Scan(&mailExists)

    templValues := gin.H{"UsernameError": usernameExists, "MailError": mailExists, "PasswordError": false}
    if usernameExists || mailExists {
        c.HTML(http.StatusConflict, "login.html", templValues)
        return
    }

    code := get_code();

    send_mail(mail, code, username);

    // Insert the user into the database
    _, err = db.Exec("INSERT INTO users (name, image, psw, mail, confirmed, confirm_key) VALUES (?, ?, ?, ?, false, ?)",
        username, imageData, hashedPassword, mail, code)
    if err != nil {
        log.Println("Database error during registration:", err)
        c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Error creating user"})
        return
    }

    c.HTML(http.StatusOK, "registration_complete.html", nil);
}

func verifyHandler(c *gin.Context) {
    code := c.Request.FormValue("code")
    _, err := db.Exec("UPDATE users SET confirmed=true WHERE confirm_key=?", code);
    if err!=nil {
        log.Println(err)
        c.HTML(http.StatusInternalServerError, "/index", nil);
    }
    row := db.QueryRow("SELECT name FROM users WHERE confirm_key=?", code);
    var username string;
    row.Scan(&username)

    session, _ := store.Get(c.Request, "session-name")
    session.Values["authenticated"] = true
    session.Values["username"] = username
    session.Save(c.Request, c.Writer)

    c.Redirect(http.StatusPermanentRedirect, "/");
}

func get_code() string{
    const dict = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM-_0123456789";
    var b strings.Builder;
    rng := rand.New(rand.NewSource(time.Now().UnixNano()));
    for i:=0; i<20; i++ {
        fmt.Fprint(&b,dict[rng.Int()%len(dict)]);
    }
    return b.String();
}

func send_mail(mail, code, name string) {
    from := os.Getenv("EMAIL")
    pass := os.Getenv("MAIL_PASSWORD")

    tmpl, err := template.ParseFiles("templates/mail.html");
    if err!=nil {
        log.Println(err);
        return;
    }
    var builder strings.Builder;
    err = tmpl.Execute(&builder, struct{Name string; Code string}{name, code})
    if err!=nil {
        log.Println(err);
        return;
    }

    msg := "From: " + from + "\n" +
        "To: " + mail + "\n" +
        "Subject: Miens 1v1 Verification Email\n" +
        "MIME-Version: 1.0\n" +
        "Content-Type: text/html; charset=\"UTF-8\"\n\n" +
        builder.String()

    err = smtp.SendMail("smtp.gmail.com:587",
        smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
        from, []string{mail}, []byte(msg))

    if err != nil {
        log.Printf("smtp error: %s", err)
        return
    }

    log.Print("Email sent")
}
