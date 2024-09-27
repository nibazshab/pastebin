package main

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"pastebin/util"
)

var pubDir = getDataFile("attachments")

//go:embed all:dist
var web embed.FS

func handleGet(c *gin.Context) {
	tid := c.Param("id")
	id := convHashId(tid)
	data := dbselect(id)

	if data.Type != "" {
		switch data.Type {
		case "text":
			c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(data.Text))
		case "file":
			pubfiledir := filepath.Join(pubDir, tid)
			file, err := os.Open(filepath.Join(pubfiledir, data.File))
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
			defer file.Close()

			c.Writer.Header().Set("Content-Type", data.Mime)
			c.Status(http.StatusOK)
			io.Copy(c.Writer, file)
		}
		viewTime := getUnixTime()
		upAfterRead(id, data.Count, viewTime)
	} else {
		c.Status(http.StatusNotFound)
	}
}

func handlePost(c *gin.Context) string {
	if !conTypeCheck(c) {
		c.String(http.StatusBadRequest, "ERROR: content-type not multipart/form-data")
		return ""
	}

	size := c.Request.ContentLength
	if size > 104857600 {
		c.String(http.StatusBadRequest, "ERROR: be less than 100mb")
		return ""
	}

	formFile, err := c.FormFile("f")
	if err != nil {
		c.String(http.StatusBadRequest, "ERROR: need name f")
		return ""
	}

	file, _ := formFile.Open()
	defer file.Close()

	buf := make([]byte, 512)
	n, _ := file.Read(buf)

	var data Data
	tid := util.RandStr(4)
	data.ID = convHashId(tid)
	data.Size = size
	data.Create = getUnixTime()
	data.Mime = http.DetectContentType(buf[:n])

	file.Seek(0, 0)
	if strings.HasPrefix(data.Mime, "text") {
		tmptext, _ := io.ReadAll(file)
		data.Text = string(tmptext)
		data.Type = "text"
		dbinsertText(data.ID, data.Text, data.Size, data.Create, data.Type)
	} else {
		pubfiledir := filepath.Join(pubDir, tid)
		os.MkdirAll(pubfiledir, 0o755)
		filename := filepath.Base(formFile.Filename)
		filebody, _ := os.Create(filepath.Join(pubfiledir, filename))
		defer filebody.Close()
		io.Copy(filebody, file)

		data.Type = "file"
		data.File = filename
		data.Size = size
		dbinsertFile(data.ID, data.File, data.Size, data.Mime, data.Create, data.Type)
	}
	c.String(http.StatusOK, c.Request.Host+"/"+tid)
	return tid
}

func indexpage(c *gin.Context) {
	c.FileFromFS("dist/", http.FS(web))
}

func static(g *gin.RouterGroup) {
	public, _ := fs.Sub(web, "dist/assets")
	g.StaticFS("/", http.FS(public))
}

func conTypeCheck(c *gin.Context) bool {
	return strings.HasPrefix(c.Request.Header.Get("Content-Type"), "multipart/form-data")
}

func attachmentsInit() {
	if _, err := os.Stat(pubDir); os.IsNotExist(err) {
		os.MkdirAll(pubDir, 0o755)
	}
}
