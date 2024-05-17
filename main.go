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
	"unsafe"
)

//go:embed index.html
var assets embed.FS

const app_data = "tmp"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			if r.Method == http.MethodPost {
				p := RequestPost(w, r)
				if p != "" {
					Log(r, p)
				}
			} else {
				RequestGetWeb(w)
			}
		} else {
			if r.Method == http.MethodGet {
				RequestGetRaw(w, r)
			}
		}
	})
	http.ListenAndServe(":1000", nil)
}

func RequestPost(w http.ResponseWriter, r *http.Request) string {
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
	p := RandString() + filepath.Ext(d.Filename)
	os.WriteFile(filepath.Join(app_data, p), t, 0o666)
	w.Write([]byte(r.Host + "/" + p))
	return p
}

func Log(r *http.Request, p string) {
	xff := r.Header.Get("X-Forwarded-For")
	if xff == "" {
		xff = r.RemoteAddr
	}
	log.Print(xff + " - " + p + " - " + r.Header.Get("user-agent"))
}

func RequestGetWeb(w http.ResponseWriter) {
	template.Must(template.ParseFS(assets, "index.html")).Execute(w, nil)
}

func RequestGetRaw(w http.ResponseWriter, r *http.Request) {
	f, e := os.ReadFile(filepath.Join(app_data, strings.TrimPrefix(r.URL.Path, "/")))
	if e != nil {
		w.Write([]byte("404 not found"))
		return
	}
	if regexp.MustCompile("^(text)/.*").MatchString(http.DetectContentType(f)) {
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	}
	w.Write(f)
}

const letter_bytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letter_idx_bits = 6
	letter_idx_mask = 1<<letter_idx_bits - 1
	letter_idx_max  = 63 / letter_idx_bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandString() string {
	b := make([]byte, 4)
	for i, cache, remain := 3, src.Int63(), letter_idx_max; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letter_idx_max
		}
		if idx := int(cache & letter_idx_mask); idx < len(letter_bytes) {
			b[i] = letter_bytes[idx]
			i--
		}
		cache >>= letter_idx_bits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
