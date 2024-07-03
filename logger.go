package main

import (
    "log"
    "net/http"
)

func logging(r *http.Request, url_path string) {
    client_ip := r.Header.Get("X-Forwarded-For")
    if client_ip == "" {
        client_ip = r.RemoteAddr
    }
    log.Print(client_ip + " - " + url_path + " - " + r.Header.Get("user-agent"))
}
