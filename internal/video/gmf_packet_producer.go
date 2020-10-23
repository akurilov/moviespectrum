package video

import (
	"github.com/3d0c/gmf"
	"github.com/sirupsen/logrus"
	"io"
	"sync/atomic"
)

type GmfPacketProducer struct {
	log              *logrus.Entry
	inputCtx         *gmf.FmtCtx
	inputStreamIndex int
	out              chan<- *gmf.Packet
	count            uint64
}

func NewGmfPacketProducer(inputCtx *gmf.FmtCtx, inputStreamIndex int, out chan<- *gmf.Packet) *GmfPacketProducer {
	log := logrus.WithFields(logrus.Fields{})
	return &GmfPacketProducer{
		log,
		inputCtx,
		inputStreamIndex,
		out,
		0,
	}
}

func (ctx *GmfPacketProducer) Produce() {
	ctx.log.Infof("Started producing src packets")
	for {
		srcPacket, err := ctx.inputCtx.GetNextPacket()
		if err == nil {
			if ctx.inputStreamIndex == srcPacket.StreamIndex() {
				ctx.out <- srcPacket
				atomic.AddUint64(&ctx.count, 1)
			}
		} else {
			if err == io.EOF {
				break
			} else {
				ctx.log.Warnf("failed to get the next src video packet: %v", err)
			}
		}
	}
	close(ctx.out)
	ctx.log.Infof("Finished producing %d src packets", ctx.count)
}

func (ctx *GmfPacketProducer) Count() uint64 {
	return atomic.LoadUint64(&ctx.count)
}
