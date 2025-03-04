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

	formName         = "f"
	previewHeader    = "X-V"
	preview          = "1"
	imagePreviewSize = 5 << 20
	textPreviewSize  = 1 << 20
)

func requestPaste(c *gin.Context) {
	uid := c.Param("uid")

	p := &Paste{
		HashKey: convHash(uid),
	}

	if p.getPaste() {
		if p.Text != "" {
			c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(p.Text))
			return
		}
		_fs := filepath.Join(attDir, p.Uid, p.FileName)

		if !p.Preview {
			c.FileAttachment(_fs, p.FileName)
			return
		}
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", p.FileName))
		c.File(_fs)
	}
}

func uploadPaste(c *gin.Context) {
	fileHeader, err := c.FormFile(formName)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Form name != %s", formName))
		return
	}
	xv := c.GetHeader(previewHeader)

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

		_t = strings.HasPrefix(mime, "text") && fileHeader.Size < textPreviewSize
		p.Preview = _t || strings.HasPrefix(mime, "image") && fileHeader.Size < imagePreviewSize
	} else {
		_t = false
		p.Preview = false
	}

	if _t {
		text, _ := io.ReadAll(file)

		p.Text = string(text)
		p.Size = fileHeader.Size
		p.inputNewPaste()
	} else {
		p.FileName = fileHeader.Filename
		p.Size = fileHeader.Size
		p.inputNewPaste()

		_fs := filepath.Join(attDir, p.Uid, p.FileName)
		c.SaveUploadedFile(fileHeader, _fs)
	}

	c.String(http.StatusOK, p.Uid)
	log.Printf("%s | %s", p.Uid, getRequestIp(c.Request))
}

func (p *Paste) inputNewPaste() {
	n := uidLength
	for {
		p.HashKey = convHash(p.Uid)
		if p.newPaste() {
			break
		}
		n++
		p.Uid = randUid(n)
	}
}

func maxBodySizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBodySize {
			c.String(http.StatusRequestEntityTooLarge, fmt.Sprintf("Request > %d", maxBodySize))
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodySize)
		err := c.Request.ParseMultipartForm(maxBodySize)
		if err != nil {
			var maxBytesError *http.MaxBytesError
			if errors.As(err, &maxBytesError) {
				c.String(http.StatusRequestEntityTooLarge, fmt.Sprintf("Request > %d", maxBodySize))
			} else {
				c.Status(http.StatusBadRequest)
			}
			return
		}
		c.Next()
	}
}
