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

type RgbaFramesConverter struct {
	inputFileName   string
	inputCtx        *gmf.FmtCtx
	inputStream     *gmf.Stream
	decoderCtx      *gmf.CodecCtx
	swsCtx          *gmf.SwsCtx
	encoder         *gmf.Codec
	encoderCtx      *gmf.CodecCtx
	srcPacketOutput chan *gmf.Packet
	rawPacketOutput chan *gmf.Packet
	frameOutput     chan *image.RGBA
}

func NewRgbaFramesConverter(inputFileName string) (*RgbaFramesConverter, error) {

	var result *RgbaFramesConverter = nil

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
		decoderCtx, swsCtx, width, height, err = initDecoderCtx(inputStream)
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

	srcPacketOutput := make(chan *gmf.Packet, 1000)
	rawPacketOutput := make(chan *gmf.Packet, 1000)
	frameOutput := make(chan *image.RGBA, 1000)

	if err == nil {
		result = &RgbaFramesConverter{
			inputFileName,
			inputCtx,
			inputStream,
			decoderCtx,
			swsCtx,
			encoder,
			encoderCtx,
			srcPacketOutput,
			rawPacketOutput,
			frameOutput,
		}
	} else {
		log.Errorf("failed to open the encoder context %v: %v", encoderCtx, err)
	}

	return result, err
}

func initDecoderCtx(stream *gmf.Stream) (*gmf.CodecCtx, *gmf.SwsCtx, int, int, error) {
	decoderCtx := stream.CodecCtx()
	width := decoderCtx.Width()
	height := decoderCtx.Height()
	pixFmt := decoderCtx.PixFmt()
	swsCtx, err := gmf.NewSwsCtx(width, height, pixFmt, width, height, pixFmt, gmf.SWS_FAST_BILINEAR)
	return decoderCtx, swsCtx, width, height, err
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

func (ctx *RgbaFramesConverter) Close() {
	ctx.inputCtx.Free()
	gmf.Release(ctx.inputStream)
	gmf.Release(ctx.decoderCtx)
	ctx.swsCtx.Free()
	gmf.Release(ctx.encoder)
	gmf.Release(ctx.encoderCtx)
}

func (ctx *RgbaFramesConverter) ProduceFrameOutput() <-chan *image.RGBA {
	go ctx.produceSrcPacketOutput()
	go ctx.produceRawPacketOutput()
	go ctx.produceFrameOutput()
	return ctx.frameOutput
}

func (ctx *RgbaFramesConverter) produceSrcPacketOutput() {
	log.Infof("Started producing src packets")
	count := 0
	srcIdx := ctx.inputStream.Index()
	for {
		srcPacket, err := ctx.inputCtx.GetNextPacket()
		if err == nil {
			if srcIdx == srcPacket.StreamIndex() {
				ctx.srcPacketOutput <- srcPacket
				count++
			}
		} else {
			if err == io.EOF {
				break
			} else {
				log.Warnf("failed to get the next src video packet: %v", err)
			}
		}
	}
	close(ctx.srcPacketOutput)
	log.Infof("Finished producing %d src packets", count)
}

func (ctx *RgbaFramesConverter) produceRawPacketOutput() {
	log.Infof("Started producing raw packets")
	count := 0
	var rawPackets []*gmf.Packet
	stream, err := ctx.inputCtx.GetStream(ctx.inputStream.Index())
	defer gmf.Release(stream)
	if err == nil {
		for srcPacket := range ctx.srcPacketOutput {
			rawPackets, err = ctx.toRawPackets(stream, srcPacket)
			if err == nil {
				for _, rawPacket := range rawPackets {
					ctx.rawPacketOutput <- rawPacket
				}
				count += len(rawPackets)
			} else {
				log.Errorf("failed to convert the src packet (%v) to raw packet: %v", srcPacket, err)
			}
			srcPacket.Free()
		}
	} else {
		log.Errorf("failed to get the stream by index %d: %v", ctx.inputStream.Index(), err)
	}
	close(ctx.rawPacketOutput)
	log.Infof("Finished producing %d raw packets", count)
}

func (ctx *RgbaFramesConverter) toRawPackets(stream *gmf.Stream, inputPacket *gmf.Packet) ([]*gmf.Packet, error) {
	var outPackets []*gmf.Packet = nil
	gmfFrames, err := stream.CodecCtx().Decode(inputPacket)
	if err == nil {
		gmfFrames, err = gmf.DefaultRescaler(ctx.swsCtx, gmfFrames)
	} else {
		log.Errorf("failed to decode the input packet (%v): %v", inputPacket, err)
	}
	if err == nil {
		outPackets, err = ctx.encoderCtx.Encode(gmfFrames, -1)
	}
	if err != nil {
		log.Errorf("faield to encode the gmf frames: %v", err)
	}
	for _, gmfFrame := range gmfFrames {
		gmfFrame.Free()
	}
	return outPackets, err
}

func (ctx *RgbaFramesConverter) produceFrameOutput() {
	log.Infof("Started producing frames")
	count := 0
	width, height := ctx.decoderCtx.Width(), ctx.decoderCtx.Height()
	for rawPacket := range ctx.rawPacketOutput {
		frame := new(image.RGBA)
		frame.Pix = rawPacket.Data()
		frame.Stride = 4 * width
		frame.Rect = image.Rect(0, 0, width, height)
		ctx.frameOutput <- frame
		count++
		rawPacket.Free()
	}
	close(ctx.frameOutput)
	log.Infof("Finished producing %d frames", count)
}
