package main

import (
	"github.com/akurilov/moviespectrum/internal/spectrum"
	"github.com/akurilov/moviespectrum/internal/video"
	"github.com/sirupsen/logrus"
	"image"
	"image/png"
	"os"
	"strconv"
)

func main() {
	for _, videoFileName := range os.Args[1:] {
		log := logrus.WithFields(logrus.Fields{"videoFileName": videoFileName})
		converter, err := video.NewRgbaFramesConverter(videoFileName)
		if err == nil {
			rawSpectrum := spectrum.NewSpectrum(100, 100)
			frameBuff := converter.ProduceFrameOutput()
			count := 0
			for frame := range frameBuff {
				frameFile, _ := os.Create(videoFileName + "_frame" + strconv.Itoa(count) + ".png")
				png.Encode(frameFile, frame)
				frameFile.Close()
				bytes := frame.Pix
				pixelCount := len(bytes) / 4 // 4 channels: R, G, B, A
				for i := 0; i < pixelCount; i++ {
					r, g, b := bytes[i], bytes[i+1], bytes[i+2]
					var cw *spectrum.ColorWeight
					cw, err = spectrum.NewColorWeight(r, g, b)
					if err == nil {
						log.Debugf("Pixel # %d has color weight: h(%f), w(%f)", i, cw.Color(), cw.Weight())
						err = rawSpectrum.AddMeasurement(cw)
						if err != nil {
							log.Errorf("failed to add the spectrum measurement: %v", cw)
						}
					} else {
						log.Errorf(
							"Failed to calculate the color weight for the color: r(%d), g(%d), b(%d)", r, g, b)
					}
				}
			}
			log.Info("Finished the conversion to raw spectrum")
			converter.Close()
			normalizedSpectrum := rawSpectrum.Normalize()
			log.Info("Normalized the spectrum")
			var img *image.RGBA
			img, err = normalizedSpectrum.ToImage()
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
		}
	}
}
