package main

import (
    "log"
    "net/http"
)

func logger(r *http.Request, id string) {
    ip := r.Header.Get("X-Forwarded-For")
    if ip == "" {
        ip = r.RemoteAddr
    }
    log.Print(id + " - " + ip + " - " + r.Header.Get("user-agent"))
}
