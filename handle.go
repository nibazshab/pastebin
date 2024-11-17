package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	pathIdNum     = 4
	attachmentDir string

	//go:embed all:dist
	web embed.FS
)

func handleReqData(c *gin.Context) {
	reqPathId := c.Param("id")
	dbHashId := convHashId(reqPathId)

	dataSelectId := &Data{
		ID: dbHashId,
	}
	reqData := dbGetDataByID(dataSelectId)

	if reqData.Type != "" {
		if reqData.Type == "text" {
			c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(reqData.Text))
		} else {
			attachmentFileDir := filepath.Join(attachmentDir, reqPathId)
			fileObj, err := os.Open(filepath.Join(attachmentFileDir, reqData.FileName))
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
			defer func(fileObj *os.File) {
				_ = fileObj.Close()
			}(fileObj)

			if reqData.Preview == false {
				c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", reqData.FileName))
			}

			c.Status(http.StatusOK)
			_, _ = io.Copy(c.Writer, fileObj)
		}

		infoData := &Data{
			ID:       dbHashId,
			Count:    reqData.Count + 1,
			LastView: getUnixTime(),
		}
		dbUpdateDataInfo(infoData)
	} else {
		c.Status(http.StatusNotFound)
	}
}

func handleUploadData(c *gin.Context) (string, bool) {
	// check con-type
	if !conTypeCheck(c) {
		c.String(http.StatusBadRequest, "ERROR: content-type not multipart/form-data")
		return "", false
	}

	// check file-size
	const maxBodySize = 100 * 1024 * 1024
	fileSize := c.Request.ContentLength
	if fileSize > maxBodySize {
		c.String(http.StatusBadRequest, "ERROR: be less than 100mb")
		return "", false
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodySize)

	// get file
	bodyObj, err := c.FormFile("f")
	if err != nil {
		c.String(http.StatusBadRequest, "ERROR: need name f")
		return "", false
	}

	fileBody, _ := bodyObj.Open()
	defer func(fileObj multipart.File) {
		_ = fileObj.Close()
	}(fileBody)

	// set file-id
	respPathId := RandStr(pathIdNum)

	// set file-type
	fileType := c.PostForm("t")
	if fileType != "file" {
		buf := make([]byte, 512)
		num, _ := fileBody.Read(buf)
		fileMime := http.DetectContentType(buf[:num])
		_, _ = fileBody.Seek(0, 0)

		if strings.HasPrefix(fileMime, "text") {
			fileType = "text"
		} else {
			fileType = "file"
		}
	}

	// set preview
	fileView, _ := strconv.ParseBool(c.PostForm("v"))
	if fileView {
		buf := make([]byte, 512)
		num, _ := fileBody.Read(buf)
		fileMime := http.DetectContentType(buf[:num])
		_, _ = fileBody.Seek(0, 0)

		if strings.HasPrefix(fileMime, "text") {
			fileView = true
		} else {
			if strings.HasPrefix(fileMime, "image") {
				fileView = true
			} else {
				fileView = false
			}
		}
	}

	// write to database
	if fileType == "text" { // text
		fileText, _ := io.ReadAll(fileBody)

		textData := &Data{
			Text:    string(fileText),
			Size:    fileSize,
			Create:  getUnixTime(),
			Type:    fileType,
			Preview: fileView,
		}
		respPathId = writeData(textData, respPathId)
	} else { // file
		fileName := filepath.Base(bodyObj.Filename)

		fileData := &Data{
			FileName: fileName,
			Size:     fileSize,
			Create:   getUnixTime(),
			Type:     fileType,
			Preview:  fileView,
		}
		respPathId = writeData(fileData, respPathId)

		// write to filesystem
		attachmentFileDir := filepath.Join(attachmentDir, respPathId)
		_ = os.MkdirAll(attachmentFileDir, 0o755)
		fileObj, _ := os.Create(filepath.Join(attachmentFileDir, fileName))
		defer func(filebody *os.File) {
			_ = filebody.Close()
		}(fileObj)

		_, _ = io.Copy(fileObj, fileBody)
	}

	// return file-id
	c.String(http.StatusOK, fmt.Sprintf("%s/%s", c.Request.Host, respPathId))
	return respPathId, true
}

func writeData(data *Data, pathId string) string {
	i := 1
	for {
		data.ID = convHashId(pathId)
		is := dbWriteData(data)
		if is {
			break
		} else {
			if i < 10 {
				i++
				pathId = RandStr(pathIdNum)
			} else {
				i = 1
				pathIdNum++
				pathId = RandStr(pathIdNum)
			}
		}
	}
	return pathId
}

func attachmentInit() {
	attachmentDir = getDataFile("attachments")

	fileInfo, err := os.Stat(attachmentDir)
	if err != nil {
		if os.IsNotExist(err) {
			_ = os.MkdirAll(attachmentDir, 0o755)
		} else {
			log.Fatalln("attachments init error")
		}
	} else {
		if !fileInfo.IsDir() {
			log.Fatalln("attachments init error")
		}
	}
}

func indexPage(c *gin.Context) {
	c.FileFromFS("dist/", http.FS(web))
}

//func publicFile(g *gin.RouterGroup) {
//	public, _ := fs.Sub(web, "dist/assets")
//	g.StaticFS("/", http.FS(public))
//}

func conTypeCheck(c *gin.Context) bool {
	return strings.HasPrefix(c.Request.Header.Get("Content-Type"), "multipart/form-data")
}
