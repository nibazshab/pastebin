package db

func initSQL() string {
	return "CREATE TABLE IF NOT EXISTS pastebin_data (id VARCHAR(16) PRIMARY KEY, data BLOB, type TINYINT(1));"
}
