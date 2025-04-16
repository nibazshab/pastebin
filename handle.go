package main

import (
	"crypto/rand"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	uidLength = 4
	uidLimit  = 32
	bodyLimit = 100 << 20

	formName    = "f"
	tokenHeader = "token"
	typeHeader  = "type"
	inline      = 1
	attachment  = 2
	link        = 3
	html        = 4
)

var typeLimits = map[string]int64{
	"text":  1 << 20,
	"image": 5 << 20,
}

type resp struct {
	Code  int    `json:",omitempty"`
	Msg   string `json:",omitempty"`
	Link  string `json:",omitempty"`
	Token string `json:",omitempty"`
}

func createPasteHandler(c *gin.Context) {
	fileHeader, err := c.FormFile(formName)
	if err != nil {
		c.IndentedJSON(400, resp{
			Code: 400,
			Msg:  "form name must be: " + formName,
		})
		return
	}

	th, _ := strconv.Atoi(c.GetHeader(typeHeader))
	if th != inline && th != attachment && th != link && th != html {
		th = inline
	}

	var paste Paste
	var e bool

	switch th {
	case html:
		e = paste.htmlCreate(c, fileHeader)
	case link:
		e = paste.linkCreate(c, fileHeader)
	case attachment:
		e = paste.attachmentCreate(c, fileHeader)
	default:
		e = paste.inlineCreate(c, fileHeader)
	}

	if !e {
		c.IndentedJSON(500, resp{
			Code: 500,
			Msg:  "paste fail",
		})
		return
	}

	log.Print(paste.Uid, " ", c.RemoteIP())

	c.IndentedJSON(200, resp{
		Link:  c.Request.Host + "/" + paste.Uid,
		Token: paste.Token,
	})
}

func (p *Paste) inlineCreate(c *gin.Context, fh *multipart.FileHeader) bool {
	f, _ := fh.Open()
	defer f.Close()

	e, t := typeCheck(f, fh.Size)
	if !e {
		return p.attachmentCreate(c, fh)
	}

	bytes, _ := io.ReadAll(f)

	if t == "text" {
		*p = Paste{
			Token: rand.Text(),
			Type:  inline,
			Size:  fh.Size,
			Text:  strings.TrimSpace(string(bytes)),
		}
		return p.create()
	} else {
		*p = Paste{
			Token:    rand.Text(),
			Type:     inline,
			Size:     fh.Size,
			FileName: fh.Filename,
		}

		if !p.create() {
			return false
		}

		if err := c.SaveUploadedFile(fh, filepath.Join(attDir, p.Uid, p.FileName)); err != nil {
			log.Print(err.Error())
			p.delete()

			return false
		}
		return true
	}
}

func (p *Paste) attachmentCreate(c *gin.Context, fh *multipart.FileHeader) bool {
	*p = Paste{
		Token:    rand.Text(),
		Type:     attachment,
		Size:     fh.Size,
		FileName: fh.Filename,
	}

	if !p.create() {
		return false
	}

	if err := c.SaveUploadedFile(fh, filepath.Join(attDir, p.Uid, p.FileName)); err != nil {
		log.Print(err.Error())
		p.delete()

		return false
	}

	return true
}

func (p *Paste) linkCreate(c *gin.Context, fh *multipart.FileHeader) bool {
	f, _ := fh.Open()
	defer f.Close()

	e, t := typeCheck(f, fh.Size)
	if t != "text" || !e {
		return p.attachmentCreate(c, fh)
	}

	bytes, _ := io.ReadAll(f)
	text := strings.TrimSpace(string(bytes))

	if _, err := url.ParseRequestURI(text); err != nil {
		return p.inlineCreate(c, fh)
	}

	*p = Paste{
		Token: rand.Text(),
		Type:  link,
		Size:  fh.Size,
		Text:  text,
	}

	return p.create()
}

func (p *Paste) htmlCreate(c *gin.Context, fh *multipart.FileHeader) bool {
	f, _ := fh.Open()
	defer f.Close()

	e, t := typeCheck(f, fh.Size)
	if t != "text" || !e {
		return p.attachmentCreate(c, fh)
	}

	bytes, _ := io.ReadAll(f)

	*p = Paste{
		Token: rand.Text(),
		Type:  html,
		Size:  fh.Size,
		Text:  strings.TrimSpace(string(bytes)),
	}
	return p.create()
}

func respPasteHandler(c *gin.Context) {
	paste := Paste{
		HashKey: convHash(c.Param("uid")),
	}

	if !paste.get() {
		c.JSON(200, resp{
			Code: 200,
			Msg:  "null",
		})
		return
	}

	switch paste.Type {
	case html:
		c.Data(200, "text/html; charset=utf-8", []byte(paste.Text))
	case link:
		c.Redirect(302, paste.Text)
	case attachment:
		fs := filepath.Join(attDir, paste.Uid, paste.FileName)
		c.FileAttachment(fs, paste.FileName)
	default:
		if paste.Text != "" {
			fs := filepath.Join(attDir, paste.Uid, paste.FileName)
			c.Writer.Header().Set("Content-Disposition", "inline; filename='"+paste.FileName+"'")
			c.File(fs)
		} else {
			c.Data(200, "text/plain; charset=utf-8", []byte(paste.Text))
		}
	}
}

func deletePasteHandler(c *gin.Context) {
	token := c.GetHeader(tokenHeader)

	if token == "" {
		c.JSON(200, resp{
			Code: 200,
			Msg:  "null",
		})
		return
	}

	paste := Paste{
		HashKey: convHash(c.Param("uid")),
		Token:   token,
	}

	paste.delete()

	c.JSON(200, resp{
		Code: 200,
		Msg:  "success",
	})
}
