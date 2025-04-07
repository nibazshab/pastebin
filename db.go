package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbName    = "pastebin.db3"
	tableName = "pastebin"
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
	return tableName
}

func database() {
	dbFile := objectPath(dbName)
	db, _ = gorm.Open(sqlite.Open(dbFile + "?_journal=WAL&_vacuum=incremental"))

	db.AutoMigrate(&Paste{})
}

func (p *Paste) get() bool {
	return db.Model(p).First(p).Error == nil
}

func (p *Paste) create() bool {
	return db.Create(p).Error == nil
}

func (p *Paste) delete() {
	db.Delete(p)
}
