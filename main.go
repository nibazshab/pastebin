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

var _port = "10002"

var port *string

func init() {
	port = flag.String("port", _port, "server port")
	flag.Parse()

	dbInit()
	logInit()
	attachmentsInit()
}

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	defer dbclose()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	g := r.Group("/assets")
	{
		g.Use(cacheControl)
		static(g)
	}

	r.GET("/favicon.ico", cacheControl, func(c *gin.Context) {
		c.Data(http.StatusOK, "image/x-icon", []byte{})
	})

	r.GET("/", cacheControl, indexpage)
	r.GET("/:id", handleGet)
	r.POST("/", func(c *gin.Context) {
		id := handlePost(c)
		if id != "" {
			logging(c, id)
		}
	})

	go func() {
		if err := r.Run(":" + *port); err != nil {
			log.Fatalf("start error: %v", err)
		}
	}()

	<-ch
}

func cacheControl(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=3600")
	c.Next()
}
