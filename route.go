package main

import (
    "net/http"
    "os"
    "path/filepath"
)

var (
    db_  string
    log_ string
)

func init() {
    ex, _ := os.Executable()
    datadir := filepath.Join(filepath.Dir(ex), "data")

    if _, err := os.Stat(datadir); os.IsNotExist(err) {
        os.MkdirAll(datadir, os.ModePerm)
    }

    db_ = filepath.Join(datadir, "pastebin.db")
    log_ = filepath.Join(datadir, "log.log")

    db_init()
    log_init()
}

func route(w http.ResponseWriter, r *http.Request) {
    index := r.URL.Path

    if r.Method == http.MethodPost {
        if index == "/" {
            id := upload_data(w, r)
            if id != "" {
                logging(r, id)
            }
        }
    } else {
        if index == "/" || index == "/style.css" || index == "/script.js" {
            index_page(w, index)
        } else {
            show_data(w, r)
        }
    }
}
