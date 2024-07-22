package main

import (
    "database/sql"
    "os"
    "path/filepath"

    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
    bin_path, _ := os.Executable()
    db_path := filepath.Join(filepath.Dir(bin_path), "data.db")

    db, _ = sql.Open("sqlite3", db_path)

    create_table := `
    CREATE TABLE IF NOT EXISTS pastebin_data (
        id VARCHAR(16) PRIMARY KEY,
        data BLOB,
        mod INTEGER
    );`

    db.Exec(create_table)
}
