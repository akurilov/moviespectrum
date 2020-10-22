package producer

type Producer interface {
	Produce()
	ConsumedCount() uint64
}
