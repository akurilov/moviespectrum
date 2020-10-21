package spectrum

import (
	"github.com/sirupsen/logrus"
	"image"
)

const (
	Width  = 360
	Height = 100
)

type Producer struct {
	log        *logrus.Entry
	frameInput <-chan *image.RGBA
	out        chan<- *image.RGBA
}

func NewProducer(frameInput <-chan *image.RGBA, out chan<- *image.RGBA) *Producer {
	log := logrus.WithFields(logrus.Fields{})
	return &Producer{log, frameInput, out}
}

func (ctx *Producer) Produce() {
	accumulator := NewSpectrum(Width, Height)
	for frame := range ctx.frameInput {
		bytes := frame.Pix
		pixelCount := len(bytes) / 4 // 4 channels: R, G, B, A
		for i := 0; i < pixelCount; i++ {
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
	logrus.Info("Finished the spectrum accumulatiom, converting to the image")
	spectrumImg, err := accumulator.ToImage()
	if err == nil {
		ctx.out <- spectrumImg
	} else {
		logrus.Errorf("failed to generate the spectrum image: %v", err)
	}
	close(ctx.out)
}
