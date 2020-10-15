package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/3d0c/gmf"
	yt "github.com/kkdai/youtube/v2"
	"github.com/sirupsen/logrus"
)

var (
	log      = logrus.WithFields(logrus.Fields{})
	ytClient = yt.Client{}
)

func main() {

	// download the specified youtube video
	ytVideo, err := ytClient.GetVideo("wEf6lVAuYQ0")
	if err != nil {
		log.Panic(err)
	}
	videoFormat := &ytVideo.Formats[0]
	videoStream, err := ytClient.GetStream(ytVideo, videoFormat)
	videoBody := videoStream.Body
	defer videoBody.Close()
	if err != nil {
		log.Panic(err)
	}
	tmpVideoFile, err := ioutil.TempFile("", "*")
	tmpVideoFileName := tmpVideoFile.Name()
	defer os.Remove(tmpVideoFileName)
	if err != nil {
		log.Panic(err)
	}
	tmpVideoFileSize, err := io.Copy(tmpVideoFile, videoBody)
	if err != nil {
		log.Panic(err)
	}
	tmpVideoFile.Close()
	log.Info(fmt.Sprintf("Video temporary file ready '%s', %d bytes", tmpVideoFileName, tmpVideoFileSize))

	// Using GMF to convert the
	inputCtx, err := gmf.NewInputCtx(tmpVideoFileName)
	defer inputCtx.Free()
	if err != nil {
		log.Panic(err)
	}
	srcVideoStream, err := inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	if err != nil {
		log.Panic(err)
	}
	ist, err := inputCtx.GetStream(srcVideoStream.Index())
	if err != nil {
		log.Panic(err)
	}
	frameNumber := 0
	for {
		packet, err := inputCtx.GetNextPacket()
		defer packet.Free()
		if err != nil {
			if err != io.EOF {
				log.Panic(err)
			}
			break
		}
		frames, err := ist.CodecCtx().Decode(packet)
		if err != nil {
			log.Warn(err)
			continue
		}
		for _, frame := range frames {
			log.Info(fmt.Sprintf("Frame #%d: %s", frameNumber, frame))
			frameNumber++
		}
	}
}
