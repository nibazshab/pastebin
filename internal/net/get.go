package net

import (
	"net/http"
	"strings"

	"github.com/nibazshab/pastebin/internal/db"
)

func getIdx(req *http.Request) string {
	return strings.TrimPrefix(req.URL.Path, "/")
}

func respData(con *[]byte, mod int, w http.ResponseWriter) {
	if mod == 1 {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}

	w.Write(*con)
}

func RespGet(w http.ResponseWriter, req *http.Request) {
	idx := getIdx(req)
	con := new([]byte)
	mod := new(int)

	err := db.Select(idx, con, mod)
	if err != nil {
		w.Write([]byte("404 not found"))
		return
	}

	respData(con, *mod, w)
}
