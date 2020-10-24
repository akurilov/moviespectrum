package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/akurilov/moviespectrum/internal/spectrum"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	log = logrus.WithFields(logrus.Fields{})
)

func main() {
	router := gin.Default()
	router.Static("/", "assets")
	router.POST("/upload", handleUpload)
	//router.POST("/youtube", handleYouTubeLink)
	_ = router.Run(":8080")
}

func handleUpload(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("file")
	var contentType string
	if err == nil {
		contentType = fileHeader.Header.Get("Content-Type")
	}
	if !strings.HasPrefix(contentType, "video/") {
		err = errors.New(fmt.Sprintf("Expected video content-type, got: %s", contentType))
	}
	var localTmpFile *os.File
	if err == nil {
		localTmpFile, err = ioutil.TempFile("", "*")
	}
	var localTmpFileName string
	if err == nil {
		localTmpFileName = localTmpFile.Name()
		defer func() { _ = os.Remove(localTmpFileName) }()
		err = ctx.SaveUploadedFile(fileHeader, localTmpFileName)
	}
	var imgBytes []byte
	if err == nil {
		imgBytes = make([]byte, 0)
		imgBuff := bytes.NewBuffer(imgBytes)
		imgWriter := io.Writer(imgBuff)
		err = spectrum.WriteVideoFileSpectrumSvg(localTmpFileName, imgWriter)
		imgBytes = imgBuff.Bytes()
	}
	if err == nil {
		ctx.Data(http.StatusOK, "image/svg+xml", imgBytes)
	} else {
		err = ctx.Error(err)
	}
	if err != nil {
		log.Error(err)
	}
}
