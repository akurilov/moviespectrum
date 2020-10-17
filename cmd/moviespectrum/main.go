package main

import (
	"fmt"
	"github.com/akurilov/moviespectrum/internal/util"
	"github.com/akurilov/moviespectrum/internal/video"
	"github.com/akurilov/moviespectrum/internal/youtube"
	"image"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	log                 = logrus.WithFields(logrus.Fields{})
	videoOutFileNameFmt = os.TempDir() + string(os.PathSeparator) + "%s"
)

func main() {

	// get the video input stream
	videoId := "UvbBh9z38TE"
	in, err := youtube.GetVideoContent(videoId)
	if in != nil {
		defer (*in).Close()
	}
	if err == nil {

		var size int64
		videoOutputFileName := fmt.Sprintf(videoOutFileNameFmt, videoId)
		defer os.Remove(videoOutputFileName)
		size, err = util.WriteToFile(in, videoOutputFileName)
		log.Info(
			fmt.Sprintf(
				"Written %d bytes from the input stream to the output file %s", size, videoOutputFileName))

		var frames *[]*image.RGBA
		if err == nil {
			frames, err = video.ConvertToFrames(videoOutputFileName)
			log.Info(fmt.Sprintf("Got %d frames from the video w/ id %s", len(*frames), videoId))
		}

		for _, frame := range *frames {
			bytes := frame.Pix
			pixelCount := len(bytes) / 4 // 4 channels: R, G, B, A
			for i := 0; i < pixelCount; i++ {
				r, g, b := bytes[i], bytes[i+1], bytes[i+2]
				log.Info(fmt.Sprintf("Pixel # %d has color: r(%d), g(%d), b(%d)", i, r, g, b))
			}
		}

	} else {
		log.Errorf("failed to get the video input stream: %s", err)
	}
}
