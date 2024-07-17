package main

import (
	"embed"
	"net/http"
)

//go:embed assets/*
var assets embed.FS

func index_page(w http.ResponseWriter, index string) {
	if index == "/" {
		index = "assets/index.html"
	} else {
		index = "assets" + index
	}
	data, _ := assets.ReadFile(index)
	w.Write(data)
}
