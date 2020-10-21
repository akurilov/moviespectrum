package video

import (
	"github.com/3d0c/gmf"
	"github.com/sirupsen/logrus"
	"image"
)

type RgbaFrameOutput struct {
	log *logrus.Entry
	out chan<- *image.RGBA
}

func NewRgbaFrameOutput(out chan<- *image.RGBA) *RgbaFrameOutput {
	log := logrus.WithFields(logrus.Fields{})
	return &RgbaFrameOutput{log, out}
}

func (ctx *RgbaFrameOutput) Consume(
	srcPacketInput <-chan *gmf.Packet,
	inputStream *gmf.Stream,
	swsCtx *gmf.SwsCtx,
	encoderCtx *gmf.CodecCtx,
) {
	ctx.log.Infof("Started producing video frame images")
	count := 0
	for srcPacket := range srcPacketInput {
		frames, err := ctx.toFrames(srcPacket, inputStream, swsCtx, encoderCtx)
		if err == nil {
			for _, frame := range frames {
				ctx.out <- frame
			}
			count += len(frames)
		} else {
			ctx.log.Errorf("failed to convert the src packet (%v) to frame: %v", srcPacket, err)
		}
		srcPacket.Free()
	}
	close(ctx.out)
	ctx.log.Infof("Finished producing %d video frame images", count)
}

func (ctx *RgbaFrameOutput) toFrames(
	srcPacket *gmf.Packet,
	inputStream *gmf.Stream,
	swsCtx *gmf.SwsCtx,
	encoderCtx *gmf.CodecCtx,
) ([]*image.RGBA, error) {
	var frames []*image.RGBA = nil
	decodedFrames, err := inputStream.CodecCtx().Decode(srcPacket)
	defer func() {
		for _, decodedFrame := range decodedFrames {
			decodedFrame.Free()
		}
	}()
	if err == nil {
		decodedFrames, err = gmf.DefaultRescaler(swsCtx, decodedFrames)
	} else {
		ctx.log.Errorf("failed to decode the input packet (%v): %v", srcPacket, err)
	}
	var convertedPackets []*gmf.Packet
	if err == nil {
		convertedPackets, err = encoderCtx.Encode(decodedFrames, -1)
	}
	if err == nil {
		width, height := encoderCtx.Width(), encoderCtx.Height()
		for _, convertedPacket := range convertedPackets {
			frame := new(image.RGBA)
			frame.Pix = convertedPacket.Data()
			frame.Stride = 4 * width
			frame.Rect = image.Rect(0, 0, width, height)
			frames = append(frames, frame)
			convertedPacket.Free()
		}
	} else {
		ctx.log.Errorf("faield to encode the gmf frames: %v", err)
	}
	return frames, err
}
