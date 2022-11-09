package dsp

import "math"

type butterworthPortionPass struct {
	n  int
	A  []float64
	d1 []float64
	d2 []float64
	w0 []float64
	w1 []float64
	w2 []float64
}

func (b *butterworthPortionPass) setup(order int) {
	b.n = order / 2
	b.A = make([]float64, b.n)
	b.d1 = make([]float64, b.n)
	b.d2 = make([]float64, b.n)
	b.w0 = make([]float64, b.n)
	b.w1 = make([]float64, b.n)
	b.w2 = make([]float64, b.n)
}

type ButterworthLowPass struct {
	butterworthPortionPass
}

type ButterworthHighPass struct {
	butterworthPortionPass
}

func (b *ButterworthLowPass) Setup(order int, sampleRate float64, frequency float64) {
	b.setup(order)

	a := math.Tan(math.Pi * frequency / sampleRate)
	a2 := a * a
	s := sampleRate

	for i := 0; i < b.n; i++ {
		r := math.Sin(math.Pi * (2*float64(i) + 1) / (4 * float64(b.n)))
		s = a2 + 2*a*r + 1
		b.A[i] = a2 / s
		b.d1[i] = 2 * (1 - a2) / s
		b.d2[i] = -(a2 - 2*a*r + 1) / s
	}
}

func (b *ButterworthHighPass) Setup(order int, sampleRate float64, frequency float64) {
	b.setup(order)

	a := math.Tan(math.Pi * frequency / sampleRate)
	a2 := a * a
	s := sampleRate

	for i := 0; i < b.n; i++ {
		r := math.Sin(math.Pi * (2*float64(i) + 1) / (4 * float64(b.n)))
		s = a2 + 2*a*r + 1
		b.A[i] = 1 / s
		b.d1[i] = 2 * (1 - a2) / s
		b.d2[i] = -(a2 - 2*a*r + 1) / s
	}
}

func (b *ButterworthLowPass) Filter(x float64) float64 {
	for i := 0; i < b.n; i++ {
		b.w0[i] = b.d1[i]*b.w1[i] + b.d2[i]*b.w2[i] + x
		x = b.A[i] * (b.w0[i] + 2*b.w1[i] + b.w2[i])
		b.w2[i] = b.w1[i]
		b.w1[i] = b.w0[i]
	}
	return x
}

func (b *ButterworthHighPass) Filter(x float64) float64 {
	for i := 0; i < b.n; i++ {
		b.w0[i] = b.d1[i]*b.w1[i] + b.d2[i]*b.w2[i] + x
		x = b.A[i] * (b.w0[i] - 2.0*b.w1[i] + b.w2[i])
		b.w2[i] = b.w1[i]
		b.w1[i] = b.w0[i]
	}
	return x
}
