package modem

type Modem interface {
	Modulate(symbols []byte) []float64
	Demodulate(samples []float64) []byte
	DCD() float64
}
