package db

func Insert(idx string, con *[]byte, mod int) {
	db.Exec("INSERT INTO pastebin_data (id, data, mod) VALUES (?, ?, ?)", idx, con, mod)
}

func Select(idx string, con *[]byte, mod *int) error {
	return db.QueryRow("SELECT data, mod FROM pastebin_data WHERE id = ?", idx).Scan(con, mod)
}
