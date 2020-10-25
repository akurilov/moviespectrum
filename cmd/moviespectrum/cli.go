package main

import (
	"github.com/akurilov/moviespectrum/internal/spectrum"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	log := logrus.WithFields(logrus.Fields{})
	videoFileName := os.Args[1]
	outImgFileName := videoFileName + ".svg"
	outImgFile, err := os.Create(outImgFileName)
	if err == nil {
		defer func() { _ = outImgFile.Close() }()
		err = spectrum.WriteVideoFileSpectrumSvg(videoFileName, outImgFile)
	}
	if err != nil {
		log.Error(err)
	}
}
