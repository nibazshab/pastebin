package main

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"pastebin/util"
)

var logFile = getDataFile("log.log")

func logInit() {
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalf("log open error: %v", err)
	}
	defer f.Close()
}

func logging(c *gin.Context, id string) {
	f, _ := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0o644)
	defer f.Close()

	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)
	log.Print(id + " | " + util.GetUserIP(c.Request) + " | " + util.GetUserUA(c.Request))
}
