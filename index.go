package main

import (
    "embed"
    "html/template"
    "net/http"
)

//go:embed templates/index.html
var assets embed.FS

func index_page(w http.ResponseWriter) {
    template.Must(template.ParseFS(assets, "templates/index.html")).Execute(w, nil)
}
