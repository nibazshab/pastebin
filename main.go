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

	port_     = "10002"
	dataPath_ = "pastebin_data"
	attDir_   = "attachment"

	webPath = "web"
)

var (
	port     *string
	dataPath string
	attDir   string

	//go:embed all:web
	web    embed.FS
	webMap = make(map[string]bool)
)

func main() {
	config()
	database()
	run()
}

func run() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(cacheControl())
	r.POST("/", limit(), createPasteHandler)
	r.GET("/*src", static, respPasteHandler)
	r.DELETE("/:uid", deletePasteHandler)

	log.Print(programName, " start HTTP server @ 0.0.0.0:", *port)
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
	v := flag.Bool("v", false, "version")
	port = flag.String("port", port_, "server port")
	dir := flag.String("dir", dataPath_, "data directory")

	flag.Parse()

	if *v {
		fmt.Print(programName, " ", version)
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
		log.Fatal(*dir, " 必须是一个有效的目录")
	}

	attDir = objectPath(attDir_)

	_fs, _ := fs.Sub(web, webPath)
	fs.WalkDir(_fs, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			webMap["/"+path] = true
		}
		return nil
	})
}

func static(c *gin.Context) {
	src := c.Param("src")

	switch src {
	case "/favicon.ico":
		c.Data(http.StatusOK, "image/x-icon", []byte{})
	case "/":
		c.FileFromFS(webPath+src, http.FS(web))
	default:
		if !webMap[src] {
			c.Next()
			return
		}
		c.FileFromFS(webPath+src, http.FS(web))
	}
	c.Abort()
}
