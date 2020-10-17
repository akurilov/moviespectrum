package video

import (
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

type converter struct {
	decoderCtx *gmf.CodecCtx
	swsCtx     *gmf.SwsCtx
	encoderCtx *gmf.CodecCtx
}

func convertStream(inputCtx *gmf.FmtCtx, stream *gmf.Stream) (*[]*image.RGBA, error) {

	var frames *[]*image.RGBA = nil

	// decoder context
	var decoderCtx = stream.CodecCtx()
	defer func() {
		gmf.Release(decoderCtx)
		decoderCtx.Free()
	}()
	encoder, err := gmf.FindEncoder(gmf.AV_CODEC_ID_RAWVIDEO)

	width := decoderCtx.Width()
	height := decoderCtx.Height()
	pixFmt := decoderCtx.PixFmt()

	// encoder context
	var encoderCtx *gmf.CodecCtx
	if err == nil {
		encoderCtx = gmf.NewCodecCtx(encoder)
		defer gmf.Release(encoderCtx)
		if encoder.IsExperimental() {
			encoderCtx.SetStrictCompliance(gmf.FF_COMPLIANCE_EXPERIMENTAL)
		}
		encoderCtx.
			SetTimeBase(gmf.AVR{Num: 1, Den: 1}).
			SetPixFmt(gmf.AV_PIX_FMT_RGBA).
			SetWidth(width).
			SetHeight(height)
		err = encoderCtx.Open(nil)
		defer encoderCtx.Free()
	}

	// rescaler
	var swsCtx *gmf.SwsCtx
	if err == nil {
		swsCtx, err = gmf.NewSwsCtx(width, height, pixFmt, width, height, pixFmt, gmf.SWS_BICUBIC)
		defer swsCtx.Free()
	}

	if err == nil {
		frames, err = convert(inputCtx, decoderCtx, swsCtx, encoderCtx)
	}

	return frames, err
}

func convert(
	inputCtx *gmf.FmtCtx, decoderCtx *gmf.CodecCtx, swsCtx *gmf.SwsCtx, encoderCtx *gmf.CodecCtx,
) (
	*[]*image.RGBA, error,
) {
	var frames []*image.RGBA
	var err error
	for {
		inPacket, err := inputCtx.GetNextPacket()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Warn(err)
			}
		}
		var outPackets []*gmf.Packet = nil
		outPackets, err = convertToGmfImages(inPacket, decoderCtx, swsCtx, encoderCtx)
		inPacket.Free()
		if err == nil {
			newFrames := convertToGoImages(outPackets, decoderCtx)
			frames = append(frames, newFrames...)
		} else {
			log.Warn(err)
		}
	}
	return &frames, err
}

func convertToGmfImages(
	inPacket *gmf.Packet, decoderCtx *gmf.CodecCtx, swsCtx *gmf.SwsCtx, encoderCtx *gmf.CodecCtx,
) (
	[]*gmf.Packet, error,
) {
	var outPackets []*gmf.Packet = nil
	gmfFrames, err := decoderCtx.Decode(inPacket)
	if err == nil {
		gmfFrames, err = gmf.DefaultRescaler(swsCtx, gmfFrames)
	}
	if err == nil {
		outPackets, err = encoderCtx.Encode(gmfFrames, -1)
	}
	return outPackets, err
}

func convertToGoImages(packets []*gmf.Packet, decoderCtx *gmf.CodecCtx) []*image.RGBA {
	var frames []*image.RGBA = nil
	for _, packet := range packets {
		width, height := decoderCtx.Width(), decoderCtx.Height()
		frame := new(image.RGBA)
		frame.Pix = packet.Data()
		frame.Stride = 4 * width
		frame.Rect = image.Rect(0, 0, width, height)
		frames = append(frames, frame)
		packet.Free()
	}
	return frames
}
