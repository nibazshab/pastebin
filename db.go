package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Paste struct {
	HashKey  int64  `gorm:"primaryKey"`
	Uid      string `gorm:"unique"`
	Text     string
	FileName string
	Size     int64
	Preview  bool
}

func (*Paste) TableName() string {
	return "pastebin"
}

func initDb() {
	dbFile := objectPath("pastebin.db3")
	db, _ = gorm.Open(sqlite.Open(dbFile + "?_journal=WAL&_vacuum=incremental"))

	db.AutoMigrate(&Paste{})
}

func (p *Paste) getPaste() bool {
	return db.Model(p).First(p).Error == nil
}

func (p *Paste) newPaste() bool {
	return db.Create(p).Error == nil
}
