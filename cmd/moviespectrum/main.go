package main

import (
	"fmt"
	"github.com/akurilov/moviespectrum/internal/spectrum"
	"github.com/akurilov/moviespectrum/internal/util"
	"github.com/akurilov/moviespectrum/internal/video"
	"github.com/akurilov/moviespectrum/internal/youtube"
	"image"
	"image/png"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	log                 = logrus.WithFields(logrus.Fields{})
	tmpFilePrefix       = os.TempDir() + string(os.PathSeparator)
	tmpFileNameFmt      = tmpFilePrefix + "%s"
	spectrumFileNameFmt = tmpFileNameFmt + ".png"
)

func main() {

	for _, videoId := range os.Args[1:] {

		// get the video input stream
		in, err := youtube.GetVideoContent(videoId)
		if in != nil {
			defer (*in).Close()
		}
		if err == nil {

			var size int64
			videoOutputFileName := fmt.Sprintf(tmpFileNameFmt, videoId)
			defer os.Remove(videoOutputFileName)
			size, err = util.WriteToFile(in, videoOutputFileName)
			log.Debugf("Written %d bytes from the input stream to the output file %s", size, videoOutputFileName)

			var frames *[]*image.RGBA
			if err == nil {
				frames, err = video.ConvertToFrames(videoOutputFileName)
				log.Debugf("Got %d frames from the video w/ id %s", len(*frames), videoId)
			}

			rawSpectrum := spectrum.NewSpectrum(100, 100)
			for _, frame := range *frames {
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
			normalizedSpectrum := rawSpectrum.Normalize()
			log.Debugf("The resulting normalized spectrum: %v", normalizedSpectrum)
			var img *image.RGBA
			img, err = normalizedSpectrum.ToImage()
			if err == nil {
				outImgFileName := fmt.Sprintf(spectrumFileNameFmt, videoId)
				outImgFile, err := os.Create(outImgFileName)
				defer outImgFile.Close()
				if err == nil {
					err = png.Encode(outImgFile, img)
					if err == nil {
						log.Infof("Done, spectrum saved to the file %s", outImgFileName)
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
			log.Fatalf("failed to get the video input stream: %s", err)
		}
	}
}
