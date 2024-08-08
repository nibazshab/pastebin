package log

import (
	"path/filepath"

	"github.com/nibazshab/pastebin/internal/dir"
)

func GetLogFile() string {
	return filepath.Join(dir.Init(), "log.log")
}
