package main

import (
	"embed"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed index.html
var static embed.FS

const tmpDir = "tmp"

func handleRaw(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-type", "text/plain; charset=UTF-8")
    c, _ := os.ReadFile(filepath.Join(tmpDir, strings.TrimPrefix(r.URL.Path, "/")))
    w.Write(c)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
    template.Must(template.ParseFS(static, "index.html")).Execute(w, nil)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    if !r.PostForm.Has("t") {
		return
    }

    a := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
    var f string
    for i := 0; i < 7; i++ {
        f += string(a[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(a))])
    }

    os.WriteFile(filepath.Join(tmpDir, f), []byte(r.PostFormValue("t")), 0666)

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
