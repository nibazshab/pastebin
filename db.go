package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	tablePrefix = "pastebin_"
	dbFile      = getDataFile("database.sqlite")
)

var db *gorm.DB

type Data struct {
	ID     uint32 `gorm:"primaryKey"`
	Text   string
	File   string
	Size   int64
	Mime   string
	Create int64
	View   int64
	Count  int `gorm:"default:0"`
	Type   string
}

func dbInit() {
	var err error

	db, err = gorm.Open(sqlite.Open(dbFile+"?_journal=WAL&_vacuum=incremental"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: tablePrefix,
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}

	err = db.AutoMigrate(&Data{})
	if err != nil {
		log.Fatalf("db init error: %v", err)
	}
}

func dbclose() {
	dB, _ := db.DB()
	err := dB.Close()
	if err != nil {
		log.Fatalf("db close error: %v", err)
	}
}

func dbselect(id uint32) *Data {
	c := Data{ID: id}
	db.Where(c).First(&c)
	return &c
}

func upAfterRead(id uint32, count int, view int64) {
	count++
	db.Model(&Data{ID: id}).Updates(Data{Count: count, View: view})
}

func dbinsertText(id uint32, text string, size int64, create int64, types string) {
	data := &Data{
		ID:     id,
		Text:   text,
		Size:   size,
		Create: create,
		Type:   types,
	}
	db.Create(data)
}

func dbinsertFile(id uint32, file string, size int64, mime string, create int64, types string) {
	data := &Data{
		ID:     id,
		File:   file,
		Size:   size,
		Mime:   mime,
		Create: create,
		Type:   types,
	}
	db.Create(data)
}
