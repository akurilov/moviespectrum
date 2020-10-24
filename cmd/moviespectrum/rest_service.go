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
	var imgBytes []byte
	if err == nil {
		imgBytes = make([]byte, 0)
		imgBuff := bytes.NewBuffer(imgBytes)
		imgWriter := io.Writer(imgBuff)
		err = writeVideoFileSpectrumSvg(localTmpFileName, imgWriter)
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

func writeVideoFileSpectrumSvg(videoFileName string, out io.Writer) error {
	var s *spectrum.Spectrum
	frameBuff := make(chan *image.RGBA, 100)
	frameProducer, err := video.NewFileRgbaFramesProducer(videoFileName, frameBuff)
	if err == nil {
		spectrumPromise := make(chan *spectrum.Spectrum)
		spectrumProducer := spectrum.NewProducerImpl(frameBuff, spectrumPromise)
		go frameProducer.Produce()
		go spectrumProducer.Produce()
		for s == nil {
			select {
			case s = <-spectrumPromise:
				log.Infof("Converted the spectrum to an image")
			case <-time.After(1 * time.Second):
				log.Infof("Processed frames %d", frameProducer.Count())
			}
		}
	}
	if s != nil {
		s.ToSvgImage(out)
	}
	return err
}
