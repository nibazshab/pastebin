package net

import (
	"net/http"
	"strings"

	"github.com/nibazshab/pastebin/internal/db"
)

func HttpGetIdx(r *http.Request) string {
	return strings.TrimPrefix(r.URL.Path, "/")
}

func HttpGetMod(con []byte, mod int, w http.ResponseWriter) {
	if mod == 1 {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}

	w.Write(con)
}

func HttpGet(w http.ResponseWriter, r *http.Request) {
	db := db.GetDb()

	var con []byte
	var mod int

	err := db.QueryRow("SELECT data, mod FROM pastebin_data WHERE id = ?", HttpGetIdx(r)).Scan(&con, &mod)
	if err != nil {
		w.Write([]byte("404 not found"))

		return
	}

	HttpGetMod(con, mod, w)
}
