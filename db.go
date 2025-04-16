package main

import (
	"errors"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"pastebin/randstrings"
)

const (
	dbName    = "pastebin.db3"
	tableName = "pastebin"
)

var db *gorm.DB

type Paste struct {
	HashKey  int64  `gorm:"primaryKey"`
	Uid      string `gorm:"unique"`
	Token    string
	Type     int
	Size     int64
	Text     string
	FileName string
}

func (*Paste) TableName() string {
	return tableName
}

func database() {
	dbFile := objectPath(dbName)
	db, _ = gorm.Open(sqlite.Open(dbFile+"?_journal=WAL&_vacuum=incremental"), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Silent),
	})

	db.AutoMigrate(&Paste{})
}

func (p *Paste) create() bool {
	n := uidLength
	for {
		p.Uid = randstrings.RandStringBytesMaskImprSrcUnsafe(n)
		p.HashKey = convHash(p.Uid)
		if err := db.Create(p).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				if n < uidLimit {
					n++
					continue
				}
			}
			log.Print(err.Error())
			return false
		}
		return true
	}
}

func (p *Paste) get() bool {
	return db.First(p).Error == nil
}

func (p *Paste) delete() {
	db.Where(p).Delete(&Paste{})
}
