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
	notView          = "1"
	imagePreviewSize = 5 << 20
	textPreviewSize  = 1 << 20
)

func getPasteHandler(c *gin.Context) {
	uid := c.Param("uid")

	paste := Paste{
		HashKey: convHash(uid),
	}

	if paste.get() {
		if paste.Text != "" {
			c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(paste.Text))
			return
		}
		_fs := filepath.Join(attDir, paste.Uid, paste.FileName)

		if !paste.Preview {
			c.FileAttachment(_fs, paste.FileName)
			return
		}
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", paste.FileName))
		c.File(_fs)
	}
}

func createPasteHandler(c *gin.Context) {
	fileHeader, err := c.FormFile(formName)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Form name != %s", formName))
		return
	}
	xv := c.GetHeader(previewHeader)

	file, _ := fileHeader.Open()
	defer file.Close()

	p := Paste{
		Uid: randUid(uidLength),
	}

	var _t bool
	if xv == notView {
		_t = false
		p.Preview = false
	} else {
		buf := make([]byte, 512)
		num, _ := file.Read(buf)
		file.Seek(0, io.SeekStart)
		mime := http.DetectContentType(buf[:num])

		_t = strings.HasPrefix(mime, "text") && fileHeader.Size < textPreviewSize
		p.Preview = _t || strings.HasPrefix(mime, "image") && fileHeader.Size < imagePreviewSize
	}

	log.Printf("%s | %s", p.Uid, requestIp(c.Request))

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
		err = c.SaveUploadedFile(fileHeader, _fs)
		if err != nil {
			p.delete()
			c.Status(http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}
	}

	c.String(http.StatusOK, p.Uid)
}

func deletePasteHandler(c *gin.Context) {
	uid := c.Param("uid")
	paste := Paste{
		HashKey: convHash(uid),
	}
	paste.delete()
	c.Status(http.StatusOK)
}

func (p *Paste) inputNewPaste() {
	n := uidLength
	for {
		p.HashKey = convHash(p.Uid)
		if p.create() {
			break
		}
		n++
		p.Uid = randUid(n)
	}
}

func limitRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBodySize {
			c.String(http.StatusRequestEntityTooLarge, fmt.Sprintf("Request > %d", maxBodySize))
			c.Abort()
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
			c.Abort()
			return
		}
		c.Next()
	}
}
