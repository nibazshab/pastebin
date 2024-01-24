package main

import (
    "embed"
    "html/template"
    "io"
    "math/rand"
    "net/http"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "time"
)

//go:embed index.html
var static embed.FS
const AppData = "tmp"

func handleIndex(w http.ResponseWriter, r *http.Request) {
    template.Must(template.ParseFS(static, "index.html")).Execute(w, nil)
}

func handleRaw(w http.ResponseWriter, r *http.Request) {
    f, e := os.Open(filepath.Join(AppData, strings.TrimPrefix(r.URL.Path, "/")))
    if e != nil {
        w.Write([]byte("404 not found"))
        return
    }
    defer f.Close()
    d, _ := f.Stat()
    http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
    if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
        return
    }
    if r.ContentLength > 5*1024*1024 {
        w.Write([]byte("Size too large (max 5M)"))
        return
    }
    c, d, _ := r.FormFile("f")
    if d == nil {
        w.Write([]byte("Invalid name (not 'f') or no file"))
        return
    }
    defer c.Close()
    t, _ := io.ReadAll(c)
    if !regexp.MustCompile("^(text|image)/.*").MatchString(http.DetectContentType(t)) {
        w.Write([]byte("Invalid file type (not TEXT/IMAGE)"))
        return
    }
    f := randStr() + filepath.Ext(d.Filename)
    os.WriteFile(filepath.Join(AppData, f), t, 0666)
    w.Write([]byte(r.Host + "/" + f))
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/" {
            if r.Method == http.MethodPost {
                handlePost(w, r)
            } else {
                handleIndex(w, r)
            }
        } else {
            handleRaw(w, r)
        }
    })
    http.ListenAndServe(":10002", nil)
}

func randStr() string {
    a := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
    var s string
    for i := 0; i < 4; i++ {
        s += string(a[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(a))])
    }
    return s
}
