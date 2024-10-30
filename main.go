package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

var templates = template.Must(template.ParseGlob("./templates/*.html"));
func session_template(w http.ResponseWriter, name string, r *http.Request, vals any) error {
    if vals == nil {vals = struct{}{}}

    b, _ := json.Marshal(&vals);
    var c map[string]interface{};
    json.Unmarshal(b, &c);

    session, _ := store.Get(r, "x-mines-session")
    if auth, ok := session.Values["authenticated"].(bool); ok && auth {
        c["SessionName"] = session.Values["username"];
    }
    fmt.Println("porcoiddio",name)

    err := templates.ExecuteTemplate(w, name, c);
    return err;
}

var db *sql.DB;
func init_db() error {
    var err error;
    db, err = sql.Open("sqlite3", "./database.db");
    if err != nil {
        return err;
    }
    _, err = db.Exec("PRAGMA foreign_keys = ON;");
    if err != nil {
        return err
    }
    return nil
}


func index_handler(w http.ResponseWriter, r *http.Request) {
    err := session_template(w, "index", r, nil);
    if err!=nil {
        http.Error(w, "Internal error", http.StatusInternalServerError);
        log.Println(err);
        return;
    }
}

var upgrader = websocket.Upgrader {
    CheckOrigin: func(r *http.Request) bool {return true;},
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w,r,nil);
    if err != nil {                     
        log.Println("Error upgrading to websocker:", err);
        return;
    }
    defer conn.Close();

    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error reading message:", err);
            break;
        }
        log.Printf("Received: %s\nType: %d", message, messageType);
        if err := conn.WriteMessage(messageType, message); err != nil {
            log.Println("Error writing message:", err);
            break;
        }
    }
}

func setCookieHeadersMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set security-related headers
        w.Header().Set("Content-Type", "text/html; charset=UTF-8")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")

        // Proceed to the next handler
        next.ServeHTTP(w, r)
    })
}


func main() {
    fmt.Println("Start")
    init_db();
    defer db.Close()
    
    mux := http.NewServeMux()

    file_server := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")));
    mux.Handle("/static/", file_server);

    mux.HandleFunc("/", index_handler);

    mux.HandleFunc("/ws", websocketHandler);

    mux.HandleFunc("/login", login_handler);
    mux.HandleFunc("/register", register_handler);

    log.Fatal(http.ListenAndServe(":2357", mux));
}
