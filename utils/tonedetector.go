package utils

import "math"

type ToneDet struct {
	sampleRate   float64
	frequency    float64
	windowSize   int
	wCoefficient float64
}

func (td *ToneDet) Initialize(sampleRate float64, frequency float64, windowSize int) {
	td.sampleRate = sampleRate
	td.frequency = math.Round(float64(windowSize)*frequency/sampleRate) / float64(windowSize)
	td.windowSize = int(math.Round(math.Ceil(float64(windowSize)*frequency/sampleRate) * sampleRate / frequency))
	td.wCoefficient = 2.0 * math.Cos(2.0*math.Pi*td.frequency)
}

func (td *ToneDet) Detect(window []float64) float64 {
	d1 := 0.0
	d2 := 0.0
	for _, x := range window {
		y := x + td.wCoefficient*d1 - d2
		d2 = d1
		d1 = y
	}
	return (d2*d2 + d1*d1 - td.wCoefficient*d1*d2) / float64(td.windowSize)
}

func (td *ToneDet) WindowSize() int {
	return td.windowSize
}
