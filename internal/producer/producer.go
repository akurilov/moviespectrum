package producer

type Producer interface {
	Produce()
	Count() uint64
}
