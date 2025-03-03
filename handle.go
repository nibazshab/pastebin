package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	uidLength   = 4
	maxBodySize = 100 * 1024 * 1024
	preview     = "1"
)

func requestPaste(c *gin.Context) {
	uid := c.Param("uid")
	hashKey := convHash(uid)

	p := &Paste{
		HashKey: hashKey,
	}

	if p.getPaste() {
		if p.Text != "" {
			c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(p.Text))
		} else {
			_fs := filepath.Join(attDir, p.Uid, p.FileName)

			if !p.Preview {
				c.FileAttachment(_fs, p.FileName)
			}
			c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", p.FileName))
			c.File(_fs)
		}
	}
}

func uploadPaste(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodySize)
	err := c.Request.ParseMultipartForm(maxBodySize)
	if err != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			c.String(http.StatusBadRequest, fmt.Sprintf("Request > %d", maxBodySize))
		} else {
			c.String(http.StatusBadRequest, "Request error")
		}
		return
	}

	fileHeader, err := c.FormFile("f")
	if err != nil {
		c.String(http.StatusBadRequest, "Form name != f")
		return
	}

	xv := c.GetHeader("X-V")

	file, _ := fileHeader.Open()
	defer file.Close()

	p := &Paste{
		Uid: randUid(uidLength),
	}

	var _t bool
	if xv == preview {
		buf := make([]byte, 512)
		num, _ := file.Read(buf)
		file.Seek(0, io.SeekStart)
		mime := http.DetectContentType(buf[:num])

		_t = strings.HasPrefix(mime, "text")
		p.Preview = _t || strings.HasPrefix(mime, "image")
	} else {
		_t = false
		p.Preview = false
	}

	if _t {
		text, _ := io.ReadAll(file)
		p.Text = string(text)
		p.Size = fileHeader.Size
		p.input()
	} else {
		p.FileName = fileHeader.Filename
		p.Size = fileHeader.Size
		p.input()

		_fs := filepath.Join(attDir, p.Uid, p.FileName)
		c.SaveUploadedFile(fileHeader, _fs)
	}

	c.String(http.StatusOK, p.Uid)
	log.Printf("%s | %s", p.Uid, getRequestIp(c.Request))
}

func (p *Paste) input() {
	n := uidLength
	for {
		p.HashKey = convHash(p.Uid)
		if p.newPaste() {
			break
		} else {
			n++
			p.Uid = randUid(n)
		}
	}
}
