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
	encoder, err := gmf.FindEncoder(gmf.AV_CODEC_ID_RAWVIDEO)
	if err == nil {
		for {
			inPacket, err := inputCtx.GetNextPacket()
			if err != nil {
				if err != io.EOF {
					log.Error(err)
				}
				break
			}
			if inPacket == nil {
				break
			} else {
				defer inPacket.Free()
				var decoderCtx = stream.CodecCtx()
				gmfFrames, err := decoderCtx.Decode(inPacket)
				if err != nil {
					log.Warn(err)
					continue
				}
				encoderCtx := gmf.NewCodecCtx(encoder)
				defer gmf.Release(encoderCtx)
				encoderCtx.
					SetTimeBase(gmf.AVR{Num: 1, Den: 1}).
					SetPixFmt(gmf.AV_PIX_FMT_RGBA).
					SetWidth(decoderCtx.Width()).
					SetHeight(decoderCtx.Height())
				var outPackets []*gmf.Packet = nil
				outPackets, err = encoderCtx.Encode(gmfFrames, drain)
			}
		}
	} else {
		log.Errorf("failed to find the raw video encoder")
	}
	return &frames, err
}
