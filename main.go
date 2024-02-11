package main

import (
    "embed"
    "html/template"
    "io"
    "log"
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
                p := handlePost(w, r)
                if p != "" {
                    logRecord(r, p)
                }
            } else {
                handleWeb(w)
            }
        } else {
            if r.Method == http.MethodGet {
                handleRaw(w, r)
            }
        }
    })
    http.ListenAndServe(":10002", nil)
}

func handlePost(w http.ResponseWriter, r *http.Request) string {
    if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
        return ""
    }
    if r.ContentLength > 5*1024*1024 {
        w.Write([]byte("Size " + strconv.FormatInt(r.ContentLength, 10) + "b, too large"))
        return ""
    }
    f, d, _ := r.FormFile("f")
    defer f.Close()
    if d == nil {
        w.Write([]byte("Invalid name (not 'f') or no file"))
        return ""
    }
    t, _ := io.ReadAll(f)
    if !regexp.MustCompile("^(text|image)/.*").MatchString(http.DetectContentType(t)) {
        w.Write([]byte("Invalid file type (not Text/Image)"))
        return ""
    }
    p := randStr() + filepath.Ext(d.Filename)
    os.WriteFile(filepath.Join(AppData, p), t, 0666)
    w.Write([]byte(r.Host + "/" + p))
    return p
}

func logRecord(r *http.Request, p string) {
    xff := r.Header.Get("X-Forwarded-For")
    if xff == "" {
        xff = r.RemoteAddr
    }
    log.Print(xff + " - " + p + " - " + r.Header.Get("user-agent"))
}

func handleWeb(w http.ResponseWriter) {
    template.Must(template.ParseFS(assets, "index.html")).Execute(w, nil)
}

func handleRaw(w http.ResponseWriter, r *http.Request) {
    f, e := os.ReadFile(filepath.Join(AppData, strings.TrimPrefix(r.URL.Path, "/")))
    if e != nil {
        w.Write([]byte("404 not found"))
        return
    }
    if regexp.MustCompile("^(text)/.*").MatchString(http.DetectContentType(f)) {
        w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
    }
    w.Write(f)
}

func randStr() string {
    a := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
    var s string
    for i := 0; i < 4; i++ {
        s += string(a[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(a))])
    }
    return s
}
