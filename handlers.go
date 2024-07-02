package main

import (
    "embed"
    "html/template"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"
)

//go:embed templates/index.html
var assets embed.FS

const app_data = "tmp"

func upload(w http.ResponseWriter, r *http.Request) string {
    if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
        return ""
    }

    if r.ContentLength > 5*1024*1024 {
        w.Write([]byte("Size " + strconv.FormatInt(r.ContentLength, 10) + "b, too large"))
        return ""
    }

    upload_file_rawdata, upload_file_metadata, _ := r.FormFile("f")
    defer upload_file_rawdata.Close()

    if upload_file_metadata == nil {
        w.Write([]byte("Invalid name (not 'f') or no file"))
        return ""
    }

    upload_file_body, _ := io.ReadAll(upload_file_rawdata)

    if !regexp.MustCompile("^(text|image)/.*").MatchString(http.DetectContentType(upload_file_body)) {
        w.Write([]byte("Invalid file type (not Text/Image)"))
        return ""
    }

    url_path := rand_string() + filepath.Ext(upload_file_metadata.Filename)

    os.WriteFile(filepath.Join(app_data, url_path), upload_file_body, 0o666)
    w.Write([]byte(r.Host + "/" + url_path))

    return url_path
}

func index_page(w http.ResponseWriter) {
    template.Must(template.ParseFS(assets, "templates/index.html")).Execute(w, nil)
}

func get_file(w http.ResponseWriter, r *http.Request) {
    get_file_body, e := os.ReadFile(filepath.Join(app_data, strings.TrimPrefix(r.URL.Path, "/")))

    if e != nil {
        w.Write([]byte("404 not found"))
        return
    }

    if regexp.MustCompile("^(text)/.*").MatchString(http.DetectContentType(get_file_body)) {
        w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
    }

    w.Write(get_file_body)
}
