package main

import (
    "log"
    "net/http"
)

func logging(r *http.Request, url_path string) {
    user_ip := r.Header.Get("X-Forwarded-For")
    if user_ip == "" {
        user_ip = r.RemoteAddr
    }
    log.Print(user_ip + " - " + url_path + " - " + r.Header.Get("user-agent"))
}
