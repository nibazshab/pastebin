package util

import (
	"net/http"
)

func GetUA(r *http.Request) string {
	return r.Header.Get("user-agent")
}
