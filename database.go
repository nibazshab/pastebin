package main

import (
    "database/sql"
    "log"
    "os"
    "path/filepath"

    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
    ex, _ := os.Executable()
    db_dir := filepath.Join(filepath.Dir(ex), "data")

    if _, err := os.Stat(db_dir); os.IsNotExist(err) {
        os.MkdirAll(db_dir, os.ModePerm)
    }

    db, _ = sql.Open("sqlite3", filepath.Join(db_dir, "pastebin.db"))

    if err := db.Ping(); err != nil {
        log.Fatalf("database error: %v", err)
    }

    create_table := `
    CREATE TABLE IF NOT EXISTS pastebin_data (
        id VARCHAR(16) PRIMARY KEY,
        data BLOB,
        mod INTEGER
    );`

    db.Exec(create_table)
}
