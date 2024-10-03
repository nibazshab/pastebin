package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

const _port = "10002"

var port *string

func init() {
	port = flag.String("port", _port, "server port")
	flag.Parse()

	dbInit()
	logInit()
	attachmentInit()
}

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	defer dbClose()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	g := r.Group("/assets")
	{
		g.Use(cacheControl)
		publicFile(g)
	}

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

	go func() {
		if err := r.Run(":" + *port); err != nil {
			log.Fatalln("start error: ", err)
		}
	}()

	<-ch
}

func cacheControl(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=3600")
	c.Next()
}
