package stream

import (
	"net/http"

	"github.com/nibazshab/pastebin/internal/home"
	"github.com/nibazshab/pastebin/internal/log"
	"github.com/nibazshab/pastebin/internal/net"
)

func Stream(w http.ResponseWriter, r *http.Request) {
	idx := r.URL.Path

	if r.Method == http.MethodPost {
		if idx == "/" {

			idx := net.HttpPost(w, r)

			if idx != "" {
				log.Message(idx, r)
			}
		}
	} else {
		if idx == "/" || idx == "/style.css" || idx == "/script.js" {
			home.HomePage(idx, w)
		} else {
			net.HttpGet(w, r)
		}
	}
}
