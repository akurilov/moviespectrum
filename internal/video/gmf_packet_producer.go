package video

import (
	"github.com/3d0c/gmf"
	"github.com/sirupsen/logrus"
	"io"
	"sync/atomic"
)

type GmfPacketProducer struct {
	log         *logrus.Entry
	inputCtx    *gmf.FmtCtx
	inputStream *gmf.Stream
	out         chan<- *gmf.Packet
	count       uint64
}

func NewGmfPacketProducer(inputCtx *gmf.FmtCtx, inputStream *gmf.Stream, out chan<- *gmf.Packet) *GmfPacketProducer {
	log := logrus.WithFields(logrus.Fields{})
	return &GmfPacketProducer{
		log,
		inputCtx,
		inputStream,
		out,
		0,
	}
}

func (ctx *GmfPacketProducer) Produce() {
	log := logrus.WithFields(logrus.Fields{})
	duration := ctx.inputCtx.Duration()
	log.Infof("Video duration: %f [s]", duration)
	approxFrameCount := uint64(ctx.inputStream.GetAvgFrameRate().AVR().Av2qd() * duration)
	log.Infof("Total frame count estimate: %d", approxFrameCount)
	packetStepToPass := approxFrameCount / 1000 // considering that each packet corresponds to a frame
	log.Infof("Pass every %dth frame", packetStepToPass)
	inputStreamIndex := ctx.inputStream.Index()
	ctx.log.Infof("Started producing src packets")
	for {
		srcPacket, err := ctx.inputCtx.GetNextPacket()
		if err == nil {
			if inputStreamIndex == srcPacket.StreamIndex() {
				if packetStepToPass < 2 {
					ctx.out <- srcPacket
				} else if atomic.LoadUint64(&ctx.count)%packetStepToPass == 0 {
					ctx.out <- srcPacket
				} else {
					srcPacket.Free() // discard
				}
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
