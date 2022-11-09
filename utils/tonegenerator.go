package utils

import (
	"math"
)

type ToneGen struct {
	sampleRate float64
	phase      float64
}

func (tg *ToneGen) Initialize(sampleRate float64) {
	tg.sampleRate = sampleRate
	tg.phase = 0.0
}

func (tg *ToneGen) Sample(frequency float64) float64 {
	sample := math.Sin(tg.phase)
	tg.phase += math.Pi * 2.0 * frequency / tg.sampleRate
	return sample
}

func (tg *ToneGen) Tone(samples int, frequency float64) []float64 {
	output := make([]float64, samples)
	omega := math.Pi * 2.0 * frequency / tg.sampleRate
	for i := 0; i < samples; i++ {
		output[i] = math.Sin(float64(i)*omega + tg.phase)
	}
	tg.phase += float64(samples) * omega
	return output
}

func (tg *ToneGen) Chirp(samples int, startFrequency float64, endFrequency float64) []float64 {
	output := make([]float64, samples)
	c := (endFrequency - startFrequency) / (float64(samples) / tg.sampleRate)
	for i := 0; i < samples; i++ {
		t := float64(i) / tg.sampleRate
		output[i] = math.Sin(tg.phase + math.Pi*2.0*(0.5*c*t*t+startFrequency*t))
	}
	t := float64(samples) / tg.sampleRate
	tg.phase += math.Pi * 2.0 * (0.5*c*t*t + startFrequency*t)
	return output
}

func (tg *ToneGen) PhaseShift(offset float64) {
	tg.phase += offset
}
