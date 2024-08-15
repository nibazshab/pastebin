package log

import "github.com/nibazshab/pastebin/internal/datapath"

func getLogPath() string {
	return datapath.GetDataFile("log.log")
}
