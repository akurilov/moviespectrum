package video

import (
	"github.com/3d0c/gmf"
	"github.com/sirupsen/logrus"
	"io"
)

type GmfPacketOutput struct {
	log *logrus.Entry
	out chan<- *gmf.Packet
}

func NewGmfPacketOutput(out chan<- *gmf.Packet) *GmfPacketOutput {
	log := logrus.WithFields(logrus.Fields{})
	return &GmfPacketOutput{
		log,
		out,
	}
}

func (ctx *GmfPacketOutput) ConsumeGmfInput(inputCtx *gmf.FmtCtx, inputStreamIndex int) {
	ctx.log.Infof("Started producing src packets")
	count := 0
	for {
		srcPacket, err := inputCtx.GetNextPacket()
		if err == nil {
			if inputStreamIndex == srcPacket.StreamIndex() {
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
	}
	close(ctx.out)
	ctx.log.Infof("Finished producing %d src packets", count)
}
