package stream

import (
	"net/http"

	"github.com/nibazshab/pastebin/internal/home"
	"github.com/nibazshab/pastebin/internal/log"
	"github.com/nibazshab/pastebin/internal/net"
)

func Stream(w http.ResponseWriter, req *http.Request) {
	idx := req.URL.Path

	if req.Method == http.MethodPost {
		if idx == "/" {
			idx := net.RespPost(w, req)
			if idx != "" {
				log.Message(idx, req)
			}
		}
	} else {
		if idx == "/" || idx == "/style.css" || idx == "/script.js" {
			home.HomePage(idx, w)
		} else {
			net.RespGet(w, req)
		}
	}
}
