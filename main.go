package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

const (
	_port = "10002"
	_dir  = "pastebin_data"
)

var (
	version string
	port    *string
	dir     *string

	attDir string

	//go:embed all:dist
	web embed.FS
)

func main() {
	config()
	attachment()
	initDb()

	run()
}

func run() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/:uid", requestPaste)
	r.POST("/", uploadPaste)

	c := r.Group("/")
	c.Use(cacheControl())
	c.GET("/", indexPage)
	c.GET("/favicon.ico", favicon)

	log.Printf("Pastebin start HTTP server @ 0.0.0.0:%s\n", *port)

	go func() {
		r.Run(":" + *port)
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	_db, _ := db.DB()
	_db.Close()
}

func config() {
	port = flag.String("port", _port, "PORT")
	dir = flag.String("dir", _dir, "DIR")
	v := flag.Bool("v", false, "version")

	flag.Parse()

	if *v {
		fmt.Printf("Pastebin %s", version)
		os.Exit(0)
	}
}

func attachment() {
	attDir = objectPath("attachment")
	_, err := os.Stat(attDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			os.Mkdir(attDir, 0o755)
		}
	}
}

func indexPage(c *gin.Context) {
	c.FileFromFS("dist/", http.FS(web))
}

func favicon(c *gin.Context) {
	c.Data(http.StatusOK, "image/x-icon", []byte{})
}
