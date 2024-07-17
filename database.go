package main

import (
    "database/sql"

    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
    db, _ = sql.Open("sqlite3", "./data.db")

    create_table := `
    CREATE TABLE IF NOT EXISTS pastebin_data (
        id VARCHAR(16) PRIMARY KEY,
        data BLOB,
        mod INTEGER
    );`

    db.Exec(create_table)
}
