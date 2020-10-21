package main

import (
	"github.com/akurilov/moviespectrum/internal/spectrum"
	"github.com/akurilov/moviespectrum/internal/video"
	"github.com/sirupsen/logrus"
	"image"
	"image/png"
	"os"
)

const (
	FrameBuffSize = 100
)

func main() {
	for _, videoFileName := range os.Args[1:] {
		log := logrus.WithFields(logrus.Fields{"videoFileName": videoFileName})
		frameBuff := make(chan *image.RGBA, FrameBuffSize)
		frameProducer, err := video.NewFileRgbaFramesProducer(videoFileName, frameBuff)
		if err == nil {
			spectrumPromise := make(chan *image.RGBA)
			spectrumProducer := spectrum.NewProducer(frameBuff, spectrumPromise)
			go frameProducer.Produce()
			go spectrumProducer.Produce()
			img := <-spectrumPromise
			log.Infof("Converted the spectrum to an image")
			if err == nil {
				outImgFileName := videoFileName + ".png"
				outImgFile, err := os.Create(outImgFileName)
				defer outImgFile.Close()
				if err == nil {
					err = png.Encode(outImgFile, img)
					if err == nil {
						log.Infof(
							"Processing done, spectrum saved to the corresponding PNG file %s", outImgFileName)
					} else {
						log.Errorf("failed to encode the PNG image: %s", err)
					}
				} else {
					log.Fatalf("failed to open the output image file: %s", err)
				}
			} else {
				log.Fatalf("failed to generate the spectrum image: %s", err)
			}
		} else {
			log.Fatalf("failed to init the frame producer: %s", err)
		}
	}
}
