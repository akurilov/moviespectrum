package main

import (
	"fmt"
	"github.com/akurilov/moviespectrum/internal/util"
	"github.com/akurilov/moviespectrum/internal/video"
	"github.com/akurilov/moviespectrum/internal/youtube"
	"os"

	yt "github.com/kkdai/youtube/v2"
	"github.com/sirupsen/logrus"
)

var (
	log                 = logrus.WithFields(logrus.Fields{})
	ytClient            = &yt.Client{}
	videoOutFileNameFmt = os.TempDir() + string(os.PathSeparator) + "%s"
)

func main() {

	// get the video input stream
	videoId := "wEf6lVAuYQ0"
	in, err := youtube.GetVideoContent(ytClient, videoId)
	if in != nil {
		defer (*in).Close()
	}
	if err == nil {
		var size int64
		videoOutputFileName := fmt.Sprint(videoOutFileNameFmt, videoId)
		defer os.Remove(videoOutputFileName)
		size, err = util.WriteToFile(in, videoOutputFileName)
		log.Info(fmt.Sprintf("Written %d bytes from the input stream to the output file", size))
		if err == nil {
			err = video.ConvertToFrames(videoOutputFileName)
		}
	} else {
		log.Errorf("failed to get the video input stream: %s", err)
	}
}
