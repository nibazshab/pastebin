package main

import (
	"errors"
	"hash/fnv"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func typeCheck(f multipart.File, n int64) (bool, string) {
	buf := make([]byte, 512)
	num, _ := f.Read(buf)
	f.Seek(0, io.SeekStart)
	mime := http.DetectContentType(buf[:num])

	for k, v := range typeLimits {
		if strings.HasPrefix(mime, k) && n < v {
			return true, k
		}
	}
	return false, ""
}

func limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > bodyLimit {
			c.JSON(413, resp{
				Code: 413,
				Msg:  "size must less: " + strconv.Itoa(bodyLimit),
			})
			c.Abort()
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, bodyLimit)
		err := c.Request.ParseMultipartForm(bodyLimit)
		if err != nil {
			var maxBytesError *http.MaxBytesError
			if errors.As(err, &maxBytesError) {
				c.JSON(413, resp{
					Code: 413,
					Msg:  "size must less: " + strconv.Itoa(bodyLimit),
				})
			} else {
				c.Status(400)
			}
			c.Abort()
			return
		}
		c.Next()
	}
}

func cacheControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=86400")
		c.Next()
	}
}

func objectPath(obj string) string {
	return filepath.Join(dataPath, obj)
}

func convHash(uid string) int64 {
	hash := fnv.New64a()
	hash.Write([]byte(uid))
	return int64(hash.Sum64())
}
