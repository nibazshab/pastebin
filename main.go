package main

import (
  "io/ioutil"
  "net/http"
  "html/template"
  "path/filepath"
  "embed"
  "strings"
  "math/rand"
  "time"
)

//go:embed index.html
var static embed.FS

const tmpDir = "tmp"

func handleRaw(w http.ResponseWriter, r *http.Request) {
  f := filepath.Join(tmpDir, strings.TrimPrefix(r.URL.Path, "/"))
  w.Header().Set("Content-type", "text/plain; charset=UTF-8")
  c, _ := ioutil.ReadFile(f)
  w.Write(c)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
  t := template.Must(template.ParseFS(static, "index.html"))
  t.Execute(w, nil)
}

func handlePost(w http.ResponseWriter, r *http.Request){
  r.ParseForm()
  if !r.PostForm.Has("t") {
    return
  }

  rand.Seed(time.Now().UnixNano())
  a := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
  var f string
  for i := 0; i < 7; i++ {
    char := a[rand.Intn(len(a))]
    f += string(char)
  }

  t := r.PostFormValue("t")
  ioutil.WriteFile(filepath.Join(tmpDir, f), []byte(t), 0666)

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
