package main

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var logFile string

func logInit() {
	logFile = getDataFile("log.log")

	logObj, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalln("log open error: ", err)
	}
	defer func(logObj *os.File) {
		_ = logObj.Close()
	}(logObj)
}

func logWrite(c *gin.Context, pathId string) {
	logObj, _ := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0o644)
	defer func(logObj *os.File) {
		_ = logObj.Close()
	}(logObj)

	multiWriter := io.MultiWriter(os.Stdout, logObj)
	log.SetOutput(multiWriter)
	log.Printf("%s | %s | %s", pathId, GetUserIP(c.Request), GetUserUA(c.Request))
}
