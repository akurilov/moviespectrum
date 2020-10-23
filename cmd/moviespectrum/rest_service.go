package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/textproto"
)

var (
	log = logrus.WithFields(logrus.Fields{})
)

func main() {
	router := gin.Default()
	router.POST("/upload", handleUpload)
	_ = router.Run(":8080")
}

func handleUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	var fileMimeHeader textproto.MIMEHeader
	var fileSize int64
	if err == nil {
		fileMimeHeader = file.Header
		fileSize = file.Size
	} else {
		log.Errorf("failed to parse the upload request %v: %s", c, err)
	}
	log.Infof("Upload file type: %s, size: %d", fileMimeHeader, fileSize)
	c.String(http.StatusOK, "OK")
}
