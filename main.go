package main

import (
    "fmt"
    "log"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
    // http.HandleFunc("/", handler)
    file_server := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")));
    http.Handle("/static/", file_server);
    log.Fatal(http.ListenAndServe(":8080", nil));
}
