package main

import (
	"net/http"
)

func route(w http.ResponseWriter, r *http.Request) {
	index := r.URL.Path

	if r.Method == http.MethodPost {
		if index == "/" {
			id := upload_data(w, r)
			if id != "" {
				logger(r, id)
			}
		}
	} else {
		if index == "/" || index == "/style.css" || index == "/script.js" || index == "/favicon.ico" {
			index_page(w, index)
		} else {
			show_data(w, r)
		}
	}
}
