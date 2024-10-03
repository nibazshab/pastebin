package main

import (
	"hash/fnv"
	"os"
	"path/filepath"
	"time"
)

var dataPath string

func getDataFile(file string) string {
	if dataPath == "" {
		ex, _ := os.Executable()

		dataPath = filepath.Join(filepath.Dir(ex), "pastebin_data")

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
