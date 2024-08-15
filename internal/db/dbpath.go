package db

import "github.com/nibazshab/pastebin/internal/datapath"

func getDbPath() string {
	return datapath.GetDataFile("pastebin.db")
}
