package main

import (
	"github.com/akurilov/moviespectrum/internal/spectrum"
	"github.com/akurilov/moviespectrum/internal/video"
	"github.com/sirupsen/logrus"
	"image"
	"os"
	"time"
)

const (
	FrameBuffSize = 100
)

func main() {
	videoFileName := os.Args[1]
	log := logrus.WithFields(logrus.Fields{})
	frameBuff := make(chan *image.RGBA, FrameBuffSize)
	frameProducer, err := video.NewFileRgbaFramesProducer(videoFileName, frameBuff)
	if err == nil {
		spectrumPromise := make(chan *spectrum.Spectrum)
		spectrumProducer := spectrum.NewProducerImpl(frameBuff, spectrumPromise)
		go frameProducer.Produce()
		go spectrumProducer.Produce()
		var s *spectrum.Spectrum = nil
		for s == nil {
			select {
			case s = <-spectrumPromise:
				log.Infof("Got the resulting spectrum")
			case <-time.After(1 * time.Second):
				log.Infof("Processed frames %d", frameProducer.Count())
			}
		}
		outImgFileName := videoFileName + ".svg"
		outImgFile, err := os.Create(outImgFileName)
		defer func() { _ = outImgFile.Close() }()
		if err == nil {
			s.ToSvgImage(outImgFile)
			log.Infof("Processing done, spectrum saved to the corresponding file \"%s\"", outImgFileName)
		} else {
			log.Fatalf("failed to open the output image file: %s", err)
		}
	} else {
		log.Fatalf("failed to init the frame producer: %s", err)
	}
}
