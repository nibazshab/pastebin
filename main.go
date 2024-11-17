package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

const (
	_port = "10002"
	_path = "pastebin_data"
)

var (
	version string
	port    *string
	path    *string
)

func main() {
	flagInit()
	dbInit()
	logInit()
	attachmentInit()

	run()
}

func run() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	defer dbClose()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	//g := r.Group("/assets")
	//{
	//	g.Use(cacheControl)
	//	publicFile(g)
	//}

	r.GET("/favicon.ico", cacheControl, func(c *gin.Context) {
		c.Data(http.StatusOK, "image/x-icon", []byte{})
	})

	r.GET("/", cacheControl, indexPage)
	r.GET("/:id", handleReqData)

	r.POST("/", func(c *gin.Context) {
		pathId, is := handleUploadData(c)
		if is {
			logWrite(c, pathId)
		}
	})

	fmt.Printf("pastebin %s\n", version)
	log.Printf("start HTTP server @ 0.0.0.0:%s\n", *port)

	go func() {
		if err := r.Run(":" + *port); err != nil {
			log.Fatalln("start error: ", err)
		}
	}()

	<-ch
}

func flagInit() {
	port = flag.String("port", _port, "server port")
	path = flag.String("path", _path, "data directory")
	_v := flag.Bool("v", false, "version")

	flag.Parse()

	if *_v {
		fmt.Printf("Version %s\nVisit github.com/nibazshab/pastebin", version)
		os.Exit(0)
	}
}

func cacheControl(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=3600")
	c.Next()
}
