package main

import (
    "io"
    "net/http"
    "path/filepath"
    "regexp"
    "strings"
)

func upload_data(w http.ResponseWriter, r *http.Request) string {
    if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
        w.Write([]byte("ERROR: content-type not multipart/form-data"))
        return ""
    }

    if r.ContentLength > 5242880 {
        w.Write([]byte("ERROR: file more than 5mb"))
        return ""
    }

    file_raw, file_meta, err := r.FormFile("f")
    if err != nil {
        w.Write([]byte("ERROR: body name not 'f'"))
        return ""
    }

    defer file_raw.Close()

    con, _ := io.ReadAll(file_raw)
    var mod int

    if regexp.MustCompile("^(text)/.*").MatchString(http.DetectContentType(con)) {
        mod = 1
    } else if regexp.MustCompile("^(image)/.*").MatchString(http.DetectContentType(con)) {
        mod = 0
    } else {
        w.Write([]byte("ERROR: file not text or image"))
        return ""
    }

    id := rand_string() + filepath.Ext(file_meta.Filename)

    db.Exec("INSERT INTO pastebin_data (id, data, mod) VALUES (?, ?, ?)", id, con, mod)
    w.Write([]byte(r.Host + "/" + id))

    return id
}

func show_data(w http.ResponseWriter, r *http.Request) {
    var con []byte
    var mod int

    err := db.QueryRow("SELECT data, mod FROM pastebin_data WHERE id = ?", strings.TrimPrefix(r.URL.Path, "/")).Scan(&con, &mod)
    if err != nil {
        w.Write([]byte("404 not found"))
        return
    }

    if mod == 1 {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    }

    w.Write(con)
}
