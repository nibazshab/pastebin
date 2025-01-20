package main

import (
	"hash/fnv"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"unsafe"
)

var dataPath string

func getDataFile(file string) string {
	if dataPath == "" {
		if filepath.IsAbs(*path) {
			dataPath = filepath.Clean(*path)
		} else {
			ex, _ := os.Executable()
			dataPath = filepath.Join(filepath.Dir(ex), *path)
		}

		if _, err := os.Stat(dataPath); os.IsNotExist(err) {
			_ = os.MkdirAll(dataPath, 0o755)
		}
	}

	dataFile := filepath.Join(dataPath, file)
	return dataFile
}

func getUnixTime() int64 {
	return time.Now().Unix()
}

func convHashId(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

func GetUserIP(req *http.Request) string {
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = req.RemoteAddr
	}

	return ip
}

func GetUserUA(req *http.Request) string {
	return req.Header.Get("user-agent")
}

// http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStr(n int) string { // RandStringBytesMaskImprSrcUnsafe
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
