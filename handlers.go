package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/diego-oniarti/mines1v1/shared"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func indexHandler(c *gin.Context) {
    shared.Render(c, http.StatusOK, "index.html", nil)
}

func loginPageHandler(c *gin.Context) {
    shared.Render(c, http.StatusOK, "login.html", nil)
}

func userPageHandler(c *gin.Context) {
    if !c.GetBool("isAuthenticated") {
        c.Redirect(http.StatusTemporaryRedirect, "/login")
        return;
    }
    shared.Render(c, http.StatusOK, "user.html", nil)
}

func loginHandler(c *gin.Context) {
    mail := c.PostForm("mail")
    password := c.PostForm("password")

    var user User
    err := db.QueryRow("SELECT name, mail, psw FROM users WHERE mail=? AND confirmed=true", mail).
        Scan(&user.Username, &user.Mail, &user.Psw);
    if err == sql.ErrNoRows || bcrypt.CompareHashAndPassword([]byte(user.Psw), []byte(password)) != nil {
        log.Println(err)
        shared.Render(c, http.StatusUnauthorized, "login.html", nil)
        return
    } else if err != nil {
        log.Println("Database error:", err)
        shared.Render(c, http.StatusInternalServerError, "login.html", nil)
        return
    }

    createSession(c, &user);

    c.Redirect(http.StatusTemporaryRedirect, "/")
}

func registerHandler(c *gin.Context) {
    username := strings.TrimSpace(c.PostForm("username"));
    mail := strings.TrimSpace(c.PostForm("mail"));
    password := strings.TrimSpace(c.PostForm("password"));

    validInputs := true;
    if utf8.RuneCountInString(username)<1 { validInputs = false; }
    if utf8.RuneCountInString(password)<8 { validInputs = false; }
    if m, err := regexp.MatchString(".+@.+\\..+", mail); !m || err!=nil { validInputs = false; }
    if m, err := regexp.MatchString("[_!?(){}#$%^&*.,+\\[\\]=+\"']", password); !m || err!=nil { validInputs = false; }
    if m, err := regexp.MatchString("[a-z]", password); !m || err!=nil { validInputs = false; }
    if m, err := regexp.MatchString("[A-Z]", password); !m || err!=nil { validInputs = false; }
    if m, err := regexp.MatchString("[0-9]", password); !m || err!=nil { validInputs = false; }

    if !validInputs {
        log.Println("Someone tried supplying invalid credentials")
        shared.Render(c, http.StatusInternalServerError, "login.html", gin.H{"error": "Error processing registration"})
        return;
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error hashing password:", err)
        shared.Render(c, http.StatusInternalServerError, "login.html", gin.H{"error": "Error processing registration"})
        return
    }

    var usernameExists, mailExists bool
    db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE name=?)", username).Scan(&usernameExists)
    db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE mail=?)", mail).Scan(&mailExists)

    templValues := gin.H{"UsernameError": usernameExists, "MailError": mailExists, "OldUser": username, "OldMail": mail};
    if usernameExists || mailExists {
        shared.Render(c, http.StatusConflict, "login.html", templValues)
        return
    }

    code := shared.RandomString(20, "code_");

    send_mail(mail, code, username);

    _, err = db.Exec("INSERT INTO users (name, psw, mail, confirmed, confirm_key) VALUES (?, ?, ?, false, ?)",
        username, hashedPassword, mail, code)
    if err != nil {
        log.Println("Database error during registration:", err)
        shared.Render(c, http.StatusInternalServerError, "login.html", gin.H{"error": "Error creating user"})
        return
    }

    shared.Render(c, http.StatusOK, "registration_complete.html", nil);
}

func verifyHandler(c *gin.Context) {
    code := c.Request.FormValue("code")
    _, err := db.Exec("UPDATE users SET confirmed=true WHERE confirm_key=?", code);
    if err!=nil {
        log.Println(err)
        shared.Render(c, http.StatusInternalServerError, "/index", nil);
    }
    row := db.QueryRow("SELECT name, mail FROM users WHERE confirm_key=?", code);
    var user User;
    row.Scan(&user.Username, &user.Mail);

    createSession(c, &user);

    c.Redirect(http.StatusPermanentRedirect, "/");
}

func logoutHandler(c *gin.Context) {
    session, _ := shared.Store.Get(c.Request, "session-name");
    session.Values["authenticated"] = false;
    session.Save(c.Request, c.Writer);
    c.Redirect(http.StatusTemporaryRedirect, "/");
}

func deleteAccountHandler(c *gin.Context) {
    session, _ := shared.Store.Get(c.Request, "session-name");
    db.Exec("DELETE FROM users WHERE mail=?", session.Values["email"]);
    session.Values["authenticated"] = false;
    session.Save(c.Request, c.Writer);
    c.Redirect(http.StatusTemporaryRedirect, "/");
}

func createSession(c *gin.Context, user *User) *sessions.Session {
    session, _ := shared.Store.Get(c.Request, "session-name");
    session.Values["authenticated"] = true;
    session.Values["username"] = user.Username;
    session.Values["email"] = user.Mail;
    session.Save(c.Request, c.Writer);
    return session;
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

func lobbyHandle(c *gin.Context) {
    shared.Render(c, http.StatusOK, "lobby.html", nil)
}
