package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const tablePrefix = "pastebin_"

var db *gorm.DB

type Data struct {
	ID       uint32 `gorm:"primaryKey"`
	Text     string
	FileName string
	Size     int64
	Create   int64
	LastView int64
	Count    int `gorm:"default:0"`
	Type     string
	Preview  bool
}

func dbInit() {
	var err error
	dbFile := getDataFile("database.sqlite")

	db, err = gorm.Open(sqlite.Open(dbFile+"?_journal=WAL&_vacuum=incremental"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: tablePrefix,
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalln("db connect error: ", err)
	}

	err = db.AutoMigrate(&Data{})
	if err != nil {
		log.Fatalln("db init error: ", err)
	}
}

func dbClose() {
	_db, _ := db.DB()
	err := _db.Close()
	if err != nil {
		log.Fatalln("db close error: ", err)
	}
}

func dbGetDataByID(data *Data) *Data {
	db.Where(data).First(data)
	return data
}

func dbUpdateDataInfo(data *Data) {
	db.Model(&Data{ID: data.ID}).Updates(data)
}

func dbWriteData(data *Data) bool {
	err := db.Create(data).Error
	return err == nil
}
