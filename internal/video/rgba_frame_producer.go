package video

import (
	"github.com/3d0c/gmf"
	"github.com/sirupsen/logrus"
	"image"
	"sync/atomic"
)

type RgbaFrameProducer struct {
	log            *logrus.Entry
	srcPacketInput <-chan *gmf.Packet
	convertor      *GnfPacketToRgbaFrameConvertor
	out            chan<- *image.RGBA
	count          uint64
}

func NewRgbaFrameProducer(
	srcPacketInput <-chan *gmf.Packet,
	convertor *GnfPacketToRgbaFrameConvertor,
	out chan<- *image.RGBA,
) *RgbaFrameProducer {
	log := logrus.WithFields(logrus.Fields{})
	return &RgbaFrameProducer{log, srcPacketInput, convertor, out, 0}
}

func (ctx *RgbaFrameProducer) Produce() {
	ctx.log.Infof("Started producing video frame images")
	for srcPacket := range ctx.srcPacketInput {
		frames, err := ctx.convertor.Convert(srcPacket)
		if err == nil {
			for _, frame := range frames {
				ctx.out <- frame
			}
		} else {
			ctx.log.Errorf("failed to convert the src packet (%v) to frame: %v", srcPacket, err)
		}
		srcPacket.Free()
		atomic.AddUint64(&ctx.count, uint64(len(frames)))
	}
	close(ctx.out)
	ctx.log.Infof("Finished producing %d video frame images", ctx.count)
}

func (ctx *RgbaFrameProducer) Count() uint64 {
	return atomic.LoadUint64(&ctx.count)
}
