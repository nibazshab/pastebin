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
	portDef = "10002"
	dirDef  = "pastebin_data"

	programName = "Pastebin"
	attDirName  = "attachment"
	embedDir    = "dist/"
)

var (
	//go:embed all:dist
	web embed.FS

	version string
	port    *string
	dir     *string

	attDir string
)

func main() {
	config()
	initAttachment()
	initDb()

	run()
}

func run() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/:uid", requestPaste)
	r.POST("/", maxBodySizeMiddleware(), uploadPaste)

	c := r.Group("/")
	c.Use(cacheControl())
	c.GET("/", indexPage)
	c.GET("/favicon.ico", favicon)

	log.Printf("%s start HTTP server @ 0.0.0.0:%s", programName, *port)

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
	port = flag.String("port", portDef, "PORT")
	dir = flag.String("dir", dirDef, "DIR")
	v := flag.Bool("v", false, "version")

	flag.Parse()

	if *v {
		fmt.Printf("%s %s", programName, version)
		os.Exit(0)
	}
}

func initAttachment() {
	attDir = objectPath(attDirName)
	_, err := os.Stat(attDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			os.Mkdir(attDir, 0o755)
		}
	}
}

func indexPage(c *gin.Context) {
	c.FileFromFS(embedDir, http.FS(web))
}

func favicon(c *gin.Context) {
	c.Data(http.StatusOK, "image/x-icon", []byte{})
}
