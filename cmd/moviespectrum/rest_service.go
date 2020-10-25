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

const (
	VideoFileSizeLimit = 128 * 1024 * 1024
)

var (
	log = logrus.WithFields(logrus.Fields{})
)

func main() {
	router := gin.Default()
	router.Static("/", "assets")
	router.POST("/upload", handleUpload)
	_ = router.Run(":8080")
}

func handleUpload(ctx *gin.Context) {
	respStatus := http.StatusInternalServerError
	fileHeader, err := ctx.FormFile("file")
	var contentType string
	if err == nil {
		contentType = fileHeader.Header.Get("Content-Type")
	}
	if !strings.HasPrefix(contentType, "video/") {
		err = errors.New(fmt.Sprintf("Expected video content-type, got: %s", contentType))
		respStatus = http.StatusBadRequest
	}
	if err == nil {
		if fileHeader.Size > VideoFileSizeLimit {
			err = errors.New(
				fmt.Sprintf(
					"Expected video file size not more than %d bytes, got %d",
					VideoFileSizeLimit,
					fileHeader.Size,
				),
			)
			respStatus = http.StatusRequestEntityTooLarge
		}
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
		msg := fmt.Sprintf("Error: %s", err)
		log.Errorf("%s, responding %d", msg, respStatus)
		ctx.String(respStatus, msg)
	}
}
