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

	fileRaw, fileMeta, err := req.FormFile("f")
	if err != nil {
		w.Write([]byte("ERROR: body name not 'f'"))
		return ""
	}
	defer fileRaw.Close()

	fileExt := filepath.Ext(fileMeta.Filename)
	if len(fileExt) > 12 {
		w.Write([]byte("ERROR: file extension too long"))
		return ""
	}

	con, _ := io.ReadAll(fileRaw)
	if len(con) > 5242880 {
		w.Write([]byte("ERROR: file more than 5mb"))
		return ""
	}

	var typ bool
	if reText.MatchString(http.DetectContentType(con)) {
		typ = true
	} else if reImage.MatchString(http.DetectContentType(con)) {
		typ = false
	} else {
		w.Write([]byte("ERROR: file not text or image"))
		return ""
	}

	idx := util.RandIdx(4) + fileExt
	db.Insert(idx, &con, typ)

	w.Write([]byte(req.Host + "/" + idx))
	return idx
}
