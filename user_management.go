package main

import (
    "io"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/sessions"
    "github.com/joho/godotenv"
    "golang.org/x/crypto/bcrypt"
)

var err = godotenv.Load(".dotenv")
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

func hash_password(psw string) (string, error) {
    hashed, err := bcrypt.GenerateFromPassword([]byte(psw), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashed), nil;
}
func check_password(hashed string, psw string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(psw));
}

func http_err(w http.ResponseWriter, err error) {
    http.Error(w, "Internal server error", http.StatusInternalServerError);
    log.Println(err);
}

func register_handler(w http.ResponseWriter, r *http.Request) {
    templ_values := struct {Username_error bool; Password_error bool; Mail_error bool;}{false,false,false};

    if r.Method != http.MethodPost {
        if err := templates.ExecuteTemplate(w, "login", nil); err!=nil {
            http_err(w, err);
        }
        return;
    }

    if err:=r.ParseMultipartForm(10); err!=nil {
        log.Println("Can't parse multiform");
        http_err(w, err);
        return;
    }

    username := r.FormValue("username");
    image, _, err := r.FormFile("image");
    var image_data []byte;
    if err == nil {
        defer image.Close();
        image_data, err = io.ReadAll(image);
        if err != nil {
            log.Println("Can't read data from image");
            http_err(w, err);
            return;
        }
    }

    password := r.FormValue("password");
    mail := r.FormValue("mail");
    hashed, err := hash_password(password);
    if err != nil {
        log.Println("Can't hash psw");
        http_err(w, err);
        return;
    }

    valid_inputs := true;

    row := db.QueryRow("SELECT count(user_id) FROM users WHERE name=$1", username);
    var count int;
    row.Scan(&count);
    if count!=0 {
        templ_values.Username_error = true;
        valid_inputs = false;
    }

    row = db.QueryRow("SELECT count(user_id) FROM users WHERE mail=$1", mail);
    row.Scan(&count);
    if count!=0 {
        templ_values.Mail_error = true;
        valid_inputs = false;
    }

    // TODO aggiungere una validazione serverside della sicurezza della password. No db access required

    if !valid_inputs {
        templates.ExecuteTemplate(w, "login", templ_values);
        return;
    }

    _, err = db.Exec("INSERT INTO users(name, image, psw, mail, confirmed) values ($1, $2, $3, $4, false)", username, image_data, hashed, mail);
    if err != nil {
        http.Error(w, "Error creating user", http.StatusInternalServerError);
        return;
    }

    log_in_user(w,r, &User{username: username, mail: mail});

    templates.ExecuteTemplate(w, "index", nil);
}



func login_handler(w http.ResponseWriter, r *http.Request) {
    templ_values := struct {LoginError bool;}{false};

    if r.Method != http.MethodPost {
        if err := templates.ExecuteTemplate(w, "login", nil); err!=nil {
            http_err(w, err);
        }
        return;
    }

    mail := r.FormValue("mail");
    password := r.FormValue("password");
    hashed, err := hash_password(password);
    if err != nil {
        log.Println("Can't hash psw");
        http_err(w, err);
        return;
    }

    valid_inputs := true;

    rows, err := db.Query("SELECT mail,name FROM users WHERE mail=$1 and psw=$2", mail, hashed);
    if err!=nil {
        log.Println("Can't lookup user");
        http_err(w, err);
        return;
    }
    defer rows.Close();
    var user User;
    if rows.Next() {
        rows.Scan(&user.mail, &user.username);
    }else{
        templ_values.LoginError = true;
        valid_inputs = false;
    }

    if !valid_inputs {
        templates.ExecuteTemplate(w, "login", templ_values);
        return;
    }

    log_in_user(w,r,&user);
    templates.ExecuteTemplate(w, "index", nil);
}

func log_in_user(w http.ResponseWriter, r *http.Request, user *User) {
    // Create a new session
    session, _ := store.Get(r, "x-mines-session")
    session.Values["authenticated"] = true
    session.Values["username"] = user.username
    session.Values["mail"] = user.mail
    session.Save(r, w)
}

type User struct {
    username string;
    mail string;
};
