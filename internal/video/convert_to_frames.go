package video

import (
	"fmt"
	"github.com/3d0c/gmf"
	"github.com/sirupsen/logrus"
	"image"
	"io"
)

var (
	log = logrus.WithFields(logrus.Fields{})
)

func ConvertToFrames(inFileName string) (*[]*image.RGBA, error) {
	var frames *[]*image.RGBA = nil
	var err error = nil
	var inputCtx *gmf.FmtCtx = nil
	inputCtx, err = gmf.NewInputCtx(inFileName)
	if inputCtx != nil {
		defer inputCtx.Free()
	}
	if err == nil {
		var stream *gmf.Stream = nil
		stream, err = inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
		if err == nil {
			frames, err = convertStream(inputCtx, stream)
		}
	}
	return frames, err
}

func convertStream(inputCtx *gmf.FmtCtx, stream *gmf.Stream) (*[]*image.RGBA, error) {
	var frames []*image.RGBA = nil
	var err error = nil
	frameNumber := 0
	for {
		packet, err := inputCtx.GetNextPacket()
		if packet != nil {
			defer packet.Free()
		}
		if err != nil {
			if err != io.EOF {
				log.Error(err)
			}
			break
		}
		gmfFrames, err := stream.CodecCtx().Decode(packet)
		if err != nil {
			log.Warn(err)
			continue
		}
		for _, gmfFrame := range gmfFrames {
			log.Debug(fmt.Sprintf("Frame #%d: %s", frameNumber, gmfFrame))
			var frame *image.RGBA = nil
			frame, err = convertGmfFrame(gmfFrame)
			if err == nil {
				frames = append(frames, frame)
			} else {
				log.Errorf("failed to convert the GMF frame #%d", frameNumber)
			}
			frameNumber++
		}
	}
	return &frames, err
}

func convertGmfFrame(gmfFrame *gmf.Frame) (*image.RGBA, error) {
	var frame *image.RGBA = nil
	var err error = nil
	return frame, err
}
