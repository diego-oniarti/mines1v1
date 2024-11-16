package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func indexHandler(c *gin.Context) {
    render(c, http.StatusOK, "index.html", nil)
}

func loginPageHandler(c *gin.Context) {
    render(c, http.StatusOK, "login.html", nil)
}

func userPageHandler(c *gin.Context) {
    if !c.GetBool("isAuthenticated") {
        c.Redirect(http.StatusTemporaryRedirect, "/login")
        return;
    }
    render(c, http.StatusOK, "user.html", nil)
}

func loginHandler(c *gin.Context) {
    mail := c.PostForm("mail")
    password := c.PostForm("password")

    var user User
    err := db.QueryRow("SELECT name, mail, psw FROM users WHERE mail=? AND confirmed=true", mail).
        Scan(&user.Username, &user.Mail, &user.Psw);
    if err == sql.ErrNoRows || bcrypt.CompareHashAndPassword([]byte(user.Psw), []byte(password)) != nil {
        log.Println(err)
        render(c, http.StatusUnauthorized, "login.html", nil)
        return
    } else if err != nil {
        log.Println("Database error:", err)
        render(c, http.StatusInternalServerError, "login.html", nil)
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
        render(c, http.StatusInternalServerError, "login.html", gin.H{"error": "Error processing registration"})
        return;
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error hashing password:", err)
        render(c, http.StatusInternalServerError, "login.html", gin.H{"error": "Error processing registration"})
        return
    }

    var usernameExists, mailExists bool
    db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE name=?)", username).Scan(&usernameExists)
    db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE mail=?)", mail).Scan(&mailExists)

    templValues := gin.H{"UsernameError": usernameExists, "MailError": mailExists, "OldUser": username, "OldMail": mail};
    if usernameExists || mailExists {
        render(c, http.StatusConflict, "login.html", templValues)
        return
    }

    code := get_code();

    send_mail(mail, code, username);

    _, err = db.Exec("INSERT INTO users (name, psw, mail, confirmed, confirm_key) VALUES (?, ?, ?, false, ?)",
        username, hashedPassword, mail, code)
    if err != nil {
        log.Println("Database error during registration:", err)
        render(c, http.StatusInternalServerError, "login.html", gin.H{"error": "Error creating user"})
        return
    }

    render(c, http.StatusOK, "registration_complete.html", nil);
}

func verifyHandler(c *gin.Context) {
    code := c.Request.FormValue("code")
    _, err := db.Exec("UPDATE users SET confirmed=true WHERE confirm_key=?", code);
    if err!=nil {
        log.Println(err)
        render(c, http.StatusInternalServerError, "/index", nil);
    }
    row := db.QueryRow("SELECT name, mail FROM users WHERE confirm_key=?", code);
    var user User;
    row.Scan(&user.Username, &user.Mail);

    createSession(c, &user);

    c.Redirect(http.StatusPermanentRedirect, "/");
}

func logoutHandler(c *gin.Context) {
    session, _ := store.Get(c.Request, "session-name");
    session.Values["authenticated"] = false;
    session.Save(c.Request, c.Writer);
    c.Redirect(http.StatusTemporaryRedirect, "/");
}

func deleteAccountHandler(c *gin.Context) {
    session, _ := store.Get(c.Request, "session-name");
    db.Exec("DELETE FROM users WHERE mail=?", session.Values["email"]);
    session.Values["authenticated"] = false;
    session.Save(c.Request, c.Writer);
    c.Redirect(http.StatusTemporaryRedirect, "/");
}

func createSession(c *gin.Context, user *User) *sessions.Session {
    session, _ := store.Get(c.Request, "session-name");
    session.Values["authenticated"] = true;
    session.Values["username"] = user.Username;
    session.Values["email"] = user.Mail;
    session.Save(c.Request, c.Writer);
    return session;
}

func get_code() string{
    const dict = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM-_0123456789";
    var b strings.Builder;
    fmt.Fprint(&b, "code_")
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

func render(c *gin.Context, code int, templateName string, data gin.H) {
    // Unisci i dati globali con quelli specifici dell'handler
    if data==nil {
        data = gin.H{};
    }
    globalData := c.MustGet("templateData").(gin.H)
    for k, v := range globalData {
        data[k] = v
        fmt.Println(k,v)
    }
    c.HTML(code, templateName, data)
}

func lobbyHandle(c *gin.Context) {
    render(c, http.StatusOK, "lobby.html", nil)
}
func singlePlayerHandler(c *gin.Context) {
    var width, height, bombs, tempo int;
    valid := true;
    tmp := func(name, def string) (int) {
        v,e := strconv.Atoi(c.DefaultQuery(name, def))
        if e!=nil || v<=0 { valid=false; return -1; }
        return v;
    }
    width = tmp("width" , "18")
    height = tmp("height", "14")
    bombs = tmp("bombs" , "40")
    timed := c.DefaultQuery("timed" , "off")
    tempo = tmp("tempo" , "3000")
    if !valid {
        c.Status(400);
        return;
    }
    
    render(c, http.StatusOK, "singlePlayer.html", nil);
}
