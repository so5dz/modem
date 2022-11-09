package dsp

import "math"

type chebyshevBandPass struct {
	m  int
	ep float64
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

func (filter *chebyshevBandPass) setup(order int) {
	filter.m = order / 4
	filter.A = make([]float64, filter.m)
	filter.d1 = make([]float64, filter.m)
	filter.d2 = make([]float64, filter.m)
	filter.d3 = make([]float64, filter.m)
	filter.d4 = make([]float64, filter.m)
	filter.w0 = make([]float64, filter.m)
	filter.w1 = make([]float64, filter.m)
	filter.w2 = make([]float64, filter.m)
	filter.w3 = make([]float64, filter.m)
	filter.w4 = make([]float64, filter.m)
}

type ChebyshevBandPass struct {
	chebyshevBandPass
}

func (filter *ChebyshevBandPass) Setup(order int, epsilon float64, sampleRate float64, low float64, high float64) {
	filter.setup(order)

	a := math.Cos(math.Pi*(low+high)/sampleRate) / math.Cos(math.Pi*(high-low)/sampleRate)
	a2 := a * a
	b := math.Tan(math.Pi * (high - low) / sampleRate)
	b2 := b * b
	n := float64(order)
	u := math.Log((1 + math.Sqrt(1+epsilon*epsilon)) / epsilon)
	su := math.Sinh(2 * u / n)
	cu := math.Cosh(2 * u / n)
	s := sampleRate

	for i := 0; i < filter.m; i++ {
		r := math.Sin(math.Pi*(2*float64(i)+1)/n) * su
		c := math.Cos(math.Pi*(2*float64(i)+1)/n) * cu
		c = r*r + c*c
		s = b2*c + 2*b*r + 1
		filter.A[i] = b2 / (4 * s)
		filter.d1[i] = 4 * a * (1 + b*r) / s
		filter.d2[i] = 2 * (b2*c - 2*a2 - 1) / s
		filter.d3[i] = 4 * a * (1 - b*r) / s
		filter.d4[i] = -(b2*c - 2*b*r + 1) / s
	}

	filter.ep = 2 / epsilon
}

func (filter *ChebyshevBandPass) Filter(x float64) float64 {
	for i := 0; i < filter.m; i++ {
		filter.w0[i] = filter.d1[i]*filter.w1[i] + filter.d2[i]*filter.w2[i] + filter.d3[i]*filter.w3[i] + filter.d4[i]*filter.w4[i] + x
		x = filter.A[i] * (filter.w0[i] - 2*filter.w2[i] + filter.w4[i])
		filter.w4[i] = filter.w3[i]
		filter.w3[i] = filter.w2[i]
		filter.w2[i] = filter.w1[i]
		filter.w1[i] = filter.w0[i]
	}
	return x * filter.ep
}
