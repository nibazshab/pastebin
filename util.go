package main

import (
	"errors"
	"hash/fnv"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
)

func cacheControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=86400")
		c.Next()
	}
}

func getRequestIp(req *http.Request) string {
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = req.RemoteAddr
	}
	return ip
}

var dataPath string

func objectPath(objName string) string {
	if dataPath == "" {
		if filepath.IsAbs(*dir) {
			dataPath = filepath.Clean(*dir)
		} else {
			ex, _ := os.Executable()
			dataPath = filepath.Join(filepath.Dir(ex), *dir)
		}

		_, err := os.Stat(dataPath)
		if errors.Is(err, fs.ErrNotExist) {
			os.MkdirAll(dataPath, 0o755)
		}
	}

	return filepath.Join(dataPath, objName)
}

func convHash(uid string) int64 {
	hasher := fnv.New64a()
	hasher.Write([]byte(uid))
	return int64(hasher.Sum64())
}

// http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

var src = rand.NewSource(time.Now().UnixNano())

func randUid(n int) string { // RandStringBytesMaskImprSrcUnsafe
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
