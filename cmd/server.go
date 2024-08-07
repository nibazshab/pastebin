package cmd

import (
	"net/http"

	"github.com/nibazshab/pastebin/internal/db"
	"github.com/nibazshab/pastebin/internal/log"
	"github.com/nibazshab/pastebin/internal/stream"
)

func init() {
	db.Init()
	log.Init()
}

func Start() {
	defer db.Close()

	http.HandleFunc("/", stream.Stream)

	if err := http.ListenAndServe(":10002", nil); err != nil {
		log.Fatalf("http start error: %v", err)
	}
}
