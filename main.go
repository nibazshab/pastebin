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
    "strconv"
    "strings"
    "time"
)

//go:embed index.html
var assets embed.FS
const AppData = "tmp"

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/" {
            if r.Method == http.MethodPost {
                handlePost(w, r)
            } else {
                handleWeb(w, r)
            }
        } else {
            handleRaw(w, r)
        }
    })
    http.ListenAndServe(":10002", nil)
}

func handleWeb(w http.ResponseWriter, r *http.Request) {
    template.Must(template.ParseFS(assets, "index.html")).Execute(w, nil)
}

func handleRaw(w http.ResponseWriter, r *http.Request) {
    f, e := os.ReadFile(filepath.Join(AppData, strings.TrimPrefix(r.URL.Path, "/")))
    if e != nil {
        w.Write([]byte("404 not found"))
        return
    }
    if regexp.MustCompile("^(text)/.*").MatchString(http.DetectContentType(f)) {
        w.Header().Set("Content-Type", "text/plain")
    }
    w.Write(f)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
    if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
        return
    }
    if r.ContentLength > 5*1024*1024 {
        w.Write([]byte("Size " + strconv.FormatInt(r.ContentLength, 10) + "b, too large"))
        return
    }
    f, d, _ := r.FormFile("f")
    defer f.Close()
    if d == nil {
        w.Write([]byte("Invalid name (not 'f') or no file"))
        return
    }
    t, _ := io.ReadAll(f)
    if !regexp.MustCompile("^(text|image)/.*").MatchString(http.DetectContentType(t)) {
        w.Write([]byte("Invalid file type (not Text/Image)"))
        return
    }
    p := randStr() + filepath.Ext(d.Filename)
    os.WriteFile(filepath.Join(AppData, p), t, 0666)
    w.Write([]byte(r.Host + "/" + p))
}

func randStr() string {
    a := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
    var s string
    for i := 0; i < 4; i++ {
        s += string(a[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(a))])
    }
    return s
}
