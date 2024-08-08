package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func SQLCreateTable() string {
	sql := "CREATE TABLE IF NOT EXISTS pastebin_data (id VARCHAR(16) PRIMARY KEY, data BLOB, mod INTEGER);"

	return sql
}

func Init() {
	db, _ = sql.Open("sqlite3", GetDbFile())

	if err := db.Ping(); err != nil {
		log.Fatalf("db connect error: %v", err)
	}

	db.Exec(SQLCreateTable())
}

func Close() {
	err := db.Close()
	if err != nil {
		log.Fatalf("db close error: %v", err)
	}
}

func GetDb() *sql.DB {
	return db
}
