package main

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "os"
)

var db *sql.DB

func init() {
    _, err := os.Stat("./database.db")
    db_exists := !os.IsNotExist(err)

    db, err = sql.Open("sqlite3", "./database.db")
    if err != nil {
        panic(err)
    }
    
    if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err!=nil {
        panic(err);
    }

    if db_exists {
        return;
    }

    users := `
    CREATE TABLE users (
        user_id integer primary key autoincrement, 
        name string, 
        image bytea, 
        psw string, 
        mail string, 
        confirmed boolean,
        confirm_key string
    );`;
    friends:= `
    CREATE TABLE friends (
        user_id1 integer,
        user_id2 integer,
        foreign key (user_id1) references users(user_id),
        foreign key (user_id2) references users(user_id),
        primary key (user_id1, user_id2)
    );`;
    scores := `
    CREATE TABLE scores (
        score_id integer primary key autoincrement,
        user_id integer,
        mode integer,
        time interval,
        seed integer,
        date timestamp,
        foreign key (user_id) references users(user_id) 
    );`;

    if _, err := db.Exec(users); err!=nil {
        panic(err);
    }
    if _, err := db.Exec(friends); err!=nil {
        panic(err);
    }
    if _, err := db.Exec(scores); err!=nil {
        panic(err);
    }
}

type User struct {
    Username  string
    Mail      string
    Psw       string
    Confirmed bool
}

