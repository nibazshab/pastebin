package net

import (
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/nibazshab/pastebin/internal/db"
	"github.com/nibazshab/pastebin/pkg/util"
)

var (
	re1 = regexp.MustCompile("^(text)/")
	re2 = regexp.MustCompile("^(image)/")
)

func ConTypeCheck(r *http.Request) bool {
	return !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data")
}

func ConLengthCheck(r *http.Request) bool {
	return r.ContentLength > 5242880
}

func HttpPost(w http.ResponseWriter, r *http.Request) string {
	if ConTypeCheck(r) {
		w.Write([]byte("ERROR: content-type not multipart/form-data"))

		return ""
	}

	if ConLengthCheck(r) {
		w.Write([]byte("ERROR: file more than 5mb"))

		return ""
	}

	file_raw, file_meta, err := r.FormFile("f")
	if err != nil {
		w.Write([]byte("ERROR: body name not 'f'"))

		return ""
	}

	defer file_raw.Close()

	con, _ := io.ReadAll(file_raw)

	var mod int

	if re1.MatchString(http.DetectContentType(con)) {
		mod = 1
	} else if re2.MatchString(http.DetectContentType(con)) {
		mod = 0
	} else {
		w.Write([]byte("ERROR: file not text or image"))

		return ""
	}

	db := db.GetDb()

	idx := util.RandString(4) + filepath.Ext(file_meta.Filename)

	db.Exec("INSERT INTO pastebin_data (id, data, mod) VALUES (?, ?, ?)", idx, con, mod)

	w.Write([]byte(r.Host + "/" + idx))

	return idx
}
