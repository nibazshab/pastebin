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
	reText  = regexp.MustCompile("^text/")
	reImage = regexp.MustCompile("^image/")
)

func RespPost(w http.ResponseWriter, req *http.Request) string {
	if !strings.HasPrefix(req.Header.Get("Content-Type"), "multipart/form-data") {
		w.Write([]byte("ERROR: content-type not multipart/form-data"))
		return ""
	}

	if req.ContentLength > 5242880 {
		w.Write([]byte("ERROR: file more than 5mb"))
		return ""
	}

	file_raw, file_meta, err := req.FormFile("f")
	if err != nil {
		w.Write([]byte("ERROR: body name not 'f'"))
		return ""
	}
	defer file_raw.Close()

	con, _ := io.ReadAll(file_raw)
	var mod int

	if reText.MatchString(http.DetectContentType(con)) {
		mod = 1
	} else if reImage.MatchString(http.DetectContentType(con)) {
		mod = 0
	} else {
		w.Write([]byte("ERROR: file not text or image"))
		return ""
	}

	idx := util.RandIdx(4) + filepath.Ext(file_meta.Filename)
	db.Insert(idx, &con, mod)

	w.Write([]byte(req.Host + "/" + idx))
	return idx
}
