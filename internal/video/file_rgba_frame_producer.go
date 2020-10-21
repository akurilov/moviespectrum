package video

import (
	"github.com/3d0c/gmf"
	"github.com/sirupsen/logrus"
	"image"
)

var (
	log = logrus.WithFields(logrus.Fields{})
)

type FileRgbaFramesProducer struct {
	inputCtx    *gmf.FmtCtx
	inputStream *gmf.Stream
	decoderCtx  *gmf.CodecCtx
	swsCtx      *gmf.SwsCtx
	encoder     *gmf.Codec
	encoderCtx  *gmf.CodecCtx
	frameOut    *RgbaFrameOutput
}

func NewFileRgbaFramesProducer(inputFileName string, out chan<- *image.RGBA) (*FileRgbaFramesProducer, error) {

	var result *FileRgbaFramesProducer = nil

	inputCtx, err := gmf.NewInputCtx(inputFileName)
	var inputStream *gmf.Stream
	if err == nil {
		inputStream, err = inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	} else {
		log.Errorf("failed to open the video file %s: %v", inputFileName, err)
	}

	var decoderCtx *gmf.CodecCtx
	var swsCtx *gmf.SwsCtx
	var width, height int
	if err == nil {
		decoderCtx = inputStream.CodecCtx()
		swsCtx, width, height, err = initSwsCtx(decoderCtx)
	} else {
		log.Errorf("failed to get the video stream for the input context %v: %v", inputCtx, err)
	}

	var encoder *gmf.Codec
	if err == nil {
		encoder, err = gmf.FindEncoder(gmf.AV_CODEC_ID_RAWVIDEO)
	} else {
		log.Errorf("failed to init the rescaler for the rescaler: %v", err)
	}

	var encoderCtx *gmf.CodecCtx
	if err == nil {
		encoderCtx, err = initEncoderCtx(encoder, width, height)
	} else {
		log.Errorf("failed to find encoder for the raw video format: %v", err)
	}

	var frameOut *RgbaFrameOutput
	if err == nil {
		frameOut = NewRgbaFrameOutput(out)
	}

	if err == nil {
		result = &FileRgbaFramesProducer{
			inputCtx,
			inputStream,
			decoderCtx,
			swsCtx,
			encoder,
			encoderCtx,
			frameOut,
		}
	} else {
		log.Errorf("failed to open the encoder context %v: %v", encoderCtx, err)
	}

	return result, err
}

func initSwsCtx(decoderCtx *gmf.CodecCtx) (*gmf.SwsCtx, int, int, error) {
	width := decoderCtx.Width()
	height := decoderCtx.Height()
	pixFmt := decoderCtx.PixFmt()
	swsCtx, err := gmf.NewSwsCtx(width, height, pixFmt, width, height, gmf.AV_PIX_FMT_RGBA, gmf.SWS_FAST_BILINEAR)
	return swsCtx, width, height, err
}

func initEncoderCtx(encoder *gmf.Codec, width int, height int) (*gmf.CodecCtx, error) {
	encoderCtx := gmf.NewCodecCtx(encoder)
	if encoder.IsExperimental() {
		encoderCtx.SetStrictCompliance(gmf.FF_COMPLIANCE_EXPERIMENTAL)
	}
	encoderCtx.
		SetTimeBase(gmf.AVR{Num: 1, Den: 1}).
		SetPixFmt(gmf.AV_PIX_FMT_RGBA).
		SetWidth(width).
		SetHeight(height)
	return encoderCtx, encoderCtx.Open(nil)
}

func (ctx *FileRgbaFramesProducer) Close() {
	ctx.inputCtx.Free()
	gmf.Release(ctx.inputStream)
	gmf.Release(ctx.decoderCtx)
	ctx.swsCtx.Free()
	gmf.Release(ctx.encoder)
	gmf.Release(ctx.encoderCtx)
}

func (ctx *FileRgbaFramesProducer) Run() {

}
