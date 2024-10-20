package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
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
        log.Printf("Received: %s", message);
        if err := conn.WriteMessage(messageType, message); err != nil {
            log.Println("Error writing message:", err);
            break;
        }
    }
}

func main() {
    // http.HandleFunc("/", handler)
    file_server := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")));
    http.Handle("/static/", file_server);

    http.HandleFunc("/ws", websocketHandler);

    log.Fatal(http.ListenAndServe(":8080", nil));
}
