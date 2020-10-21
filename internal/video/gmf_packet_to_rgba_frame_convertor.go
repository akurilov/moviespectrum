package video

import (
	"github.com/3d0c/gmf"
	"image"
)

type GnfPacketToRgbaFrameConvertor struct {
	inputStream *gmf.Stream
	swsCtx *gmf.SwsCtx
	encoderCtx *gmf.CodecCtx
}

func (ctx *RgbaFrameOutput) toFrames(srcPacket *gmf.Packet
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
