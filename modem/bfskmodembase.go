package modem

import (
	// "encoding/binary"

	"encoding/binary"
	"math"
	"os"

	dspbw "github.com/iskrapw/dsp/filter/butterworth"
	dspch "github.com/iskrapw/dsp/filter/chebyshev"
	"github.com/iskrapw/modem/utils"
)

type BFSKModemBase struct {
	inputSymbolSamples  int
	outputSymbolSamples int
	space               float64
	center              float64
	mark                float64
	lastPhase           float64
	lastValue           float64
	toneGenerator       utils.ToneGen
	iGenerator          utils.ToneGen
	qGenerator          utils.ToneGen
	iLPF                dspbw.LowPass
	qLPF                dspbw.LowPass
	valueLPF            dspbw.LowPass
	valueLPF2           dspbw.LowPass
	sampleBPF           dspch.BandPass
	phaseDeltaMax       float64
	todoDebug           *os.File
}

func (f *BFSKModemBase) initialize(inputRate float64, outputRate float64, center float64, shift float64, baud float64) {
	f.inputSymbolSamples = int(math.Round(inputRate / baud))
	f.outputSymbolSamples = int(math.Round(outputRate / baud))
	f.toneGenerator.Initialize(outputRate)

	f.space = center - shift/2
	f.center = center
	f.mark = center + shift/2

	f.iGenerator.Initialize(inputRate)
	f.qGenerator.Initialize(inputRate)
	f.qGenerator.PhaseShift(math.Pi / 2)

	iqLPFCutoff := (baud + center) / 2
	f.iLPF.Setup(4, inputRate, iqLPFCutoff)
	f.qLPF.Setup(4, inputRate, iqLPFCutoff)

	f.valueLPF.Setup(2, inputRate, baud/2)
	f.valueLPF2.Setup(4, inputRate, 3*baud/2)

	sampleBPFHalfWidth := 0.95 * shift
	f.sampleBPF.Setup(8, 0.25, inputRate, center-sampleBPFHalfWidth, center+sampleBPFHalfWidth)

	f.phaseDeltaMax = 2 * baud / (0.5 * inputRate)

	// Set non-zero values into every filter, buffer etc
	f.getFuzzySymbol(1e-3)

	f.todoDebug, _ = os.Create("debug.raw")
}

func (f *BFSKModemBase) getTone(time float64, bit int) []float64 {
	frequency := f.space
	if bit == 1 {
		frequency = f.mark
	}

	length := int(math.Round(time * float64(f.outputSymbolSamples)))
	return f.toneGenerator.Tone(length, frequency)
}

func (f *BFSKModemBase) getFuzzySymbol(sample float64) float64 {
	sample = f.sampleBPF.Filter(sample)
	f.debugWrite(sample)

	I := sample * f.iGenerator.Sample(f.center)
	I = f.iLPF.Filter(I)

	Q := sample * f.qGenerator.Sample(f.center)
	Q = f.qLPF.Filter(Q)

	phase := math.Atan2(Q, I)
	phaseDelta := (phase - f.lastPhase)
	f.lastPhase = phase

	var value float64
	if math.Abs(phaseDelta) < 1 {
		value = phaseDelta / f.phaseDeltaMax
	} else {
		value = f.lastValue
	}
	f.lastValue = value

	value = f.valueLPF.Filter(value)
	value = f.valueLPF2.Filter(value)
	f.debugWrite(value)

	return value
}

func (f *BFSKModemBase) DCD() float64 {
	return 0.0 // todo
}

// todo remove
func (f *BFSKModemBase) debugWrite(value float64) {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], math.Float32bits(float32(value)))
	f.todoDebug.Write(buf[:])
}
