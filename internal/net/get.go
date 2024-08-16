package net

import (
	"net/http"
	"strings"

	"github.com/nibazshab/pastebin/internal/db"
)

func getIdx(req *http.Request) string {
	return strings.TrimPrefix(req.URL.Path, "/")
}

func respData(con *[]byte, typ bool, w http.ResponseWriter) {
	if typ {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}

	w.Write(*con)
}

func RespGet(w http.ResponseWriter, req *http.Request) {
	idx := getIdx(req)
	con := new([]byte)
	typ := new(bool)

	err := db.Select(idx, con, typ)
	if err != nil {
		w.Write([]byte("404 not found"))
		return
	}

	respData(con, *typ, w)
}
