package db

func Insert(idx string, con *[]byte, typ bool) {
	db.Exec("INSERT INTO pastebin_data (id, data, type) VALUES (?, ?, ?)", idx, con, typ)
}

func Select(idx string, con *[]byte, typ *bool) error {
	return db.QueryRow("SELECT data, type FROM pastebin_data WHERE id = ?", idx).Scan(con, typ)
}
