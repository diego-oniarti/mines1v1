package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

var templates = template.Must(template.ParseGlob("./templates/*.html"));

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
    err := templates.ExecuteTemplate(w, "index", nil);
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

func main() {
    fmt.Println("Start")
    init_db();
    defer db.Close()
    
    file_server := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")));
    http.Handle("/static/", file_server);

    http.HandleFunc("/", index_handler)

    http.HandleFunc("/ws", websocketHandler);

    http.HandleFunc("/login", login_handler);
    http.HandleFunc("/register", register_handler);

    log.Fatal(http.ListenAndServe(":2357", nil));
}
