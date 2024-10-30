package main

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
    var err error
    db, err = sql.Open("sqlite3", "./database.db")
    if err != nil {
        panic(err)
    }
}

type User struct {
    Username  string
    Mail      string
    Psw       string
    Confirmed bool
}

