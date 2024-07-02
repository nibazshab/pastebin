package main

import (
    "net/http"
)

func main() {
    http.HandleFunc("/", route)
    http.ListenAndServe(":10002", nil)
}

func route(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/" {
        if r.Method == http.MethodPost {
            url_path := upload(w, r)
            if url_path != "" {
                logging(r, url_path)
            }
        } else {
            index_page(w)
        }
    } else {
        if r.Method == http.MethodGet {
            get_file(w, r)
        }
    }
}
