package dsp

import "math"

type chebyshevPortionPass struct {
	m  int
	ep float64
	A  []float64
	d1 []float64
	d2 []float64
	w0 []float64
	w1 []float64
	w2 []float64
}

func (b *chebyshevPortionPass) setup(order int) {
	b.m = order / 2
	b.A = make([]float64, b.m)
	b.d1 = make([]float64, b.m)
	b.d2 = make([]float64, b.m)
	b.w0 = make([]float64, b.m)
	b.w1 = make([]float64, b.m)
	b.w2 = make([]float64, b.m)
}

type ChebyshevLowPass struct {
	chebyshevPortionPass
}

type ChebyshevHighPass struct {
	chebyshevPortionPass
}

func (filter *ChebyshevLowPass) Setup(order int, epsilon float64, sampleRate float64, frequency float64) {
	filter.setup(order)

	a := math.Tan(math.Pi * frequency / sampleRate)
	a2 := a * a

	n := float64(order)
	u := math.Log((1 + math.Sqrt(1+epsilon*epsilon)) / epsilon)
	su := math.Sinh(u / n)
	cu := math.Cosh(u / n)
	s := sampleRate

	for i := 0; i < filter.m; i++ {
		b := math.Sin(math.Pi*(2*float64(i)+1)/(2*n)) * su
		c := math.Cos(math.Pi*(2*float64(i)+1)/(2*n)) * cu
		c = b*b + c*c
		s = a2*c + 2*a*b + 1
		filter.A[i] = a2 / (4 * s)
		filter.d1[i] = 2 * (1 - a2*c) / s
		filter.d2[i] = -(a2*c - 2*a*b + 1) / s
	}

	filter.ep = 2 / epsilon
}

func (filter *ChebyshevHighPass) Setup(order int, epsilon float64, sampleRate float64, frequency float64) {
	filter.setup(order)

	a := math.Tan(math.Pi * frequency / sampleRate)
	a2 := a * a
	n := float64(order)
	u := math.Log((1 + math.Sqrt(1+epsilon*epsilon)) / epsilon)
	su := math.Sinh(u / n)
	cu := math.Cosh(u / n)
	s := sampleRate

	for i := 0; i < filter.m; i++ {
		b := math.Sin(math.Pi*(2*float64(i)+1)/(2*n)) * su
		c := math.Cos(math.Pi*(2*float64(i)+1)/(2*n)) * cu
		c = b*b + c*c
		s = a2 + 2*a*b + c
		filter.A[i] = 1 / (4 * s)
		filter.d1[i] = 2 * (c - a2) / s
		filter.d2[i] = -(a2 - 2*a*b + c) / s
	}

	filter.ep = 2 / epsilon
}

func (filter *ChebyshevLowPass) Filter(x float64) float64 {
	for i := 0; i < filter.m; i++ {
		filter.w0[i] = filter.d1[i]*filter.w1[i] + filter.d2[i]*filter.w2[i] + x
		x = filter.A[i] * (filter.w0[i] + 2*filter.w1[i] + filter.w2[i])
		filter.w2[i] = filter.w1[i]
		filter.w1[i] = filter.w0[i]
	}
	return x * filter.ep
}

func (filter *ChebyshevHighPass) Filter(x float64) float64 {
	for i := 0; i < filter.m; i++ {
		filter.w0[i] = filter.d1[i]*filter.w1[i] + filter.d2[i]*filter.w2[i] + x
		x = filter.A[i] * (filter.w0[i] - 2*filter.w1[i] + filter.w2[i])
		filter.w2[i] = filter.w1[i]
		filter.w1[i] = filter.w0[i]
	}
	return x * filter.ep
}
