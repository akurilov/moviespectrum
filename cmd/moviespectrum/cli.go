package main

import (
	"github.com/akurilov/moviespectrum/internal/spectrum"
	"os"
)

func main() {
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
