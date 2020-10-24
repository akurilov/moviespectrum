package spectrum

import (
	"github.com/sirupsen/logrus"
	"image"
	"sync/atomic"
)

const (
	Width  = 360
	Height = 100
)

type ProducerImpl struct {
	log        *logrus.Entry
	frameInput <-chan *image.RGBA
	out        chan<- *Spectrum
	count      uint64
}

func NewProducerImpl(frameInput <-chan *image.RGBA, out chan<- *Spectrum) *ProducerImpl {
	log := logrus.WithFields(logrus.Fields{})
	return &ProducerImpl{log, frameInput, out, 0}
}

func (ctx *ProducerImpl) Produce() {
	accumulator := NewSpectrum(Width, Height)
	for frame := range ctx.frameInput {
		bytes := frame.Pix
		pixelCount := len(bytes) / 4 // 4 channels: R, G, B, A
		for i := 0; i < pixelCount; i += 4 {
			r, g, b := bytes[i], bytes[i+1], bytes[i+2]
			cw, err := NewColorWeight(r, g, b)
			if err == nil {
				ctx.log.Debugf("Pixel # %d has color weight: h(%f), w(%f)", i, cw.Color(), cw.Weight())
				err = accumulator.AddMeasurement(cw)
				if err != nil {
					ctx.log.Errorf("failed to add the spectrum measurement: %v", cw)
				}
			} else {
				ctx.log.Errorf(
					"Failed to calculate the color weight for the color: r(%d), g(%d), b(%d)", r, g, b)
			}
		}
	}
	atomic.AddUint64(&ctx.count, 1)
	logrus.Info("Finished the spectrum accumulation, converting to the image")
	ctx.out <- accumulator
	close(ctx.out)
}

func (ctx *ProducerImpl) Count() uint64 {
	return atomic.LoadUint64(&ctx.count)
}
