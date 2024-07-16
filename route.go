package main

import (
    "net/http"
)

func route(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/" {
        if r.Method == http.MethodPost {
            id := upload_data(w, r)
            if id != "" {
                logger(r, id)
            }
        } else {
            index_page(w)
        }
    } else {
        if r.Method == http.MethodGet {
            show_data(w, r)
        }
    }
}
