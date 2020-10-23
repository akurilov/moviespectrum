package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/akurilov/moviespectrum/internal/spectrum"
	"github.com/akurilov/moviespectrum/internal/video"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
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
	var img *image.RGBA
	if err == nil {
		img, err = videoFileToSpectrumImage(localTmpFileName)
	}
	var imgBytes []byte
	if err == nil {
		imgBytes = make([]byte, 0)
		imgBuff := bytes.NewBuffer(imgBytes)
		imgWriter := io.Writer(imgBuff)
		err = png.Encode(imgWriter, img)
		imgBytes = imgBuff.Bytes()
	}
	if err == nil {
		ctx.Data(http.StatusOK, "image/png", imgBytes)
	} else {
		err = ctx.Error(err)
	}
	if err != nil {
		log.Error(err)
	}
}

func videoFileToSpectrumImage(videoFileName string) (*image.RGBA, error) {
	var img *image.RGBA
	frameBuff := make(chan *image.RGBA, 100)
	frameProducer, err := video.NewFileRgbaFramesProducer(videoFileName, frameBuff)
	if err == nil {
		spectrumPromise := make(chan *image.RGBA)
		spectrumProducer := spectrum.NewProducerImpl(frameBuff, spectrumPromise)
		go frameProducer.Produce()
		go spectrumProducer.Produce()
		for img == nil {
			select {
			case img = <-spectrumPromise:
				log.Infof("Converted the spectrum to an image")
			case <-time.After(1 * time.Second):
				log.Infof("Processed frames %d", frameProducer.Count())
			}
		}
	}
	return img, err
}
