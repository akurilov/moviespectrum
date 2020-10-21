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
	srcPacketConvertor *GnfPacketToRgbaFrameConvertor,
) {
	ctx.log.Infof("Started producing video frame images")
	count := 0
	for srcPacket := range srcPacketInput {
		frames, err := srcPacketConvertor.Convert(srcPacket)
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
