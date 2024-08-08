package db

import (
	"path/filepath"

	"github.com/nibazshab/pastebin/internal/dir"
)

func GetDbFile() string {
	return filepath.Join(dir.Init(), "pastebin.db")
}
