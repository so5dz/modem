package utils

import "math"

type AutoGain struct {
	Attack float64
	Decay  float64
	Level  float64
}

func (a *AutoGain) Feed(x float64) float64 {
	if math.Abs(x) > a.Level {
		a.Level *= (1 + a.Attack)
	} else {
		a.Level *= (1 - a.Decay)
	}
	return x / a.Level
}
