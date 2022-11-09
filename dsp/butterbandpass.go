package dsp

import "math"

type butterworthBandPass struct {
	n  int
	A  []float64
	d1 []float64
	d2 []float64
	d3 []float64
	d4 []float64
	w0 []float64
	w1 []float64
	w2 []float64
	w3 []float64
	w4 []float64
}

func (b *butterworthBandPass) setup(order int) {
	b.n = order / 4
	b.A = make([]float64, b.n)
	b.d1 = make([]float64, b.n)
	b.d2 = make([]float64, b.n)
	b.d3 = make([]float64, b.n)
	b.d4 = make([]float64, b.n)
	b.w0 = make([]float64, b.n)
	b.w1 = make([]float64, b.n)
	b.w2 = make([]float64, b.n)
	b.w3 = make([]float64, b.n)
	b.w4 = make([]float64, b.n)
}

type ButterworthBandPass struct {
	butterworthBandPass
}

func (bpf *ButterworthBandPass) Setup(order int, sampleRate float64, low float64, high float64) {
	bpf.setup(order)

	a := math.Cos(math.Pi*(high+low)/sampleRate) / math.Cos(math.Pi*(high-low)/sampleRate)
	a2 := a * a
	b := math.Tan(math.Pi * (high - low) / sampleRate)
	b2 := b * b
	s := sampleRate

	for i := 0; i < bpf.n; i++ {
		r := math.Sin(math.Pi * (2*float64(i) + 1) / (4 * float64(bpf.n)))
		s = b2 + 2*b*r + 1
		bpf.A[i] = b2 / s
		bpf.d1[i] = 4 * a * (1 + b*r) / s
		bpf.d2[i] = 2 * (b2 - 2*a2 - 1) / s
		bpf.d3[i] = 4 * a * (1 - b*r) / s
		bpf.d4[i] = -(b2 - 2*b*r + 1) / s
	}
}

func (bpf *ButterworthBandPass) Filter(x float64) float64 {
	for i := 0; i < bpf.n; i++ {
		bpf.w0[i] = bpf.d1[i]*bpf.w1[i] + bpf.d2[i]*bpf.w2[i] + bpf.d3[i]*bpf.w3[i] + bpf.d4[i]*bpf.w4[i] + x
		x = bpf.A[i] * (bpf.w0[i] - 2*bpf.w2[i] + bpf.w4[i])
		bpf.w4[i] = bpf.w3[i]
		bpf.w3[i] = bpf.w2[i]
		bpf.w2[i] = bpf.w1[i]
		bpf.w1[i] = bpf.w0[i]
	}
	return x
}
