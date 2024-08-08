package home

import (
	"net/http"

	"github.com/nibazshab/pastebin/web"
)

func HomePage(idx string, w http.ResponseWriter) {
	if idx == "/" {
		idx = "dist/index.html"
	} else {
		idx = "dist" + idx
	}

	data, _ := web.Web.ReadFile(idx)

	w.Write(data)
}
