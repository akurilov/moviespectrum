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
	consumedCount    uint64
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
	count := 0
	for {
		srcPacket, err := ctx.inputCtx.GetNextPacket()
		if err == nil {
			if ctx.inputStreamIndex == srcPacket.StreamIndex() {
				ctx.out <- srcPacket
				count++
			}
		} else {
			if err == io.EOF {
				break
			} else {
				ctx.log.Warnf("failed to get the next src video packet: %v", err)
			}
		}
		atomic.AddUint64(&ctx.consumedCount, 1)
	}
	close(ctx.out)
	ctx.log.Infof("Finished producing %d src packets", count)
}

func (ctx *GmfPacketProducer) ConsumedCount() uint64 {
	return atomic.LoadUint64(&ctx.consumedCount)
}
