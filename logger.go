package main

import (
    "io"
    "log"
    "net/http"
    "os"
)

func log_init() {
    f, err := os.Create(log_)
    if err != nil {
        log.Fatalf("log.log error: %v", err)
    }
    defer f.Close()
}

func logging(r *http.Request, id string) {
    ip := r.Header.Get("X-Forwarded-For")
    if ip == "" {
        ip = r.RemoteAddr
    }

    f, _ := os.OpenFile(log_, os.O_APPEND|os.O_RDWR, os.ModePerm)
    defer f.Close()

    multiWriter := io.MultiWriter(os.Stdout, f)
    log.SetOutput(multiWriter)

    log.Print(id + " - " + ip + " - " + r.Header.Get("user-agent"))
}
