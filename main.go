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
	"path/filepath"
	"syscall"

	"github.com/gin-gonic/gin"
)

const (
	version     = "1.2.1"
	programName = "pastebin"

	portDef    = "10002"
	dirDef     = "pastebin_data"
	embedDir   = "dist/"
	attDirName = "attachment"
)

var (
	//go:embed all:dist
	web      embed.FS
	dataPath string
	attDir   string
	port     *string
)

func main() {
	config()
	database()
	run()
}

func run() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.POST("/", limitRequest(), createPasteHandler)
	r.GET("/:uid", respPasteHandler)
	r.DELETE("/:uid", deletePasteHandler)

	g := r.Group("/")
	g.Use(cacheControl())
	g.GET("/favicon.ico", favicon)
	g.GET("/", indexPage)

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
	port = flag.String("port", portDef, "server port")
	dir := flag.String("dir", dirDef, "data directory")
	v := flag.Bool("v", false, "version")

	flag.Parse()

	if *v {
		fmt.Printf("%s %s", programName, version)
		os.Exit(0)
	}

	if filepath.IsAbs(*dir) {
		dataPath = filepath.Clean(*dir)
	} else {
		ex, _ := os.Executable()
		dataPath = filepath.Join(filepath.Dir(ex), *dir)
	}

	info, err := os.Stat(dataPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			os.MkdirAll(dataPath, 0o755)
		} else {
			log.Fatal(err)
		}
	}
	if info != nil && !info.IsDir() {
		log.Fatalf("%s 必须是一个有效的目录", *dir)
	}

	attDir = objectPath(attDirName)
}

func favicon(c *gin.Context) {
	c.Data(http.StatusOK, "image/x-icon", []byte{})
}

func indexPage(c *gin.Context) {
	c.FileFromFS(embedDir, http.FS(web))
}
