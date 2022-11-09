package modem

import (
	"math"

	"github.com/iskrapw/modem/stuffing"

	"github.com/tunabay/go-bitarray"
)

type Bell202Modem struct {
	BFSKModemBase
	lastValue         float64
	sinceLastCrossing int
	bitStuffer        stuffing.BitStuffer
	bitDestuffer      stuffing.BitDestuffer
}

func (f *Bell202Modem) Initialize(inputRate float64, outputRate float64, center float64) {
	f.initialize(inputRate, outputRate, center, 1000, 1200)
	f.sinceLastCrossing = 0
	f.lastValue = -1

	flag := bitarray.MustParse("01111110")
	f.bitStuffer = stuffing.BitStuffer{Flag: flag, BreakEvery: 5, BreakBit: 1}
	f.bitDestuffer = stuffing.BitDestuffer{BreakEvery: 5, BreakBit: 1}

	// better filter than bfskmodembase
	sampleBPFHalfWidth := 720.0
	f.sampleBPF.Setup(12, 0.5, inputRate, center-sampleBPFHalfWidth, center+sampleBPFHalfWidth)
}

func (f *Bell202Modem) Modulate(bytes []byte) []float64 {
	output := make([]float64, 0)

	f.bitStuffer.Start()
	f.bitStuffer.PutMultipleFlags(4)
	f.bitStuffer.PutBytes(bytes)
	f.bitStuffer.PutMultipleFlags(1)
	stuffedBits := f.bitStuffer.End()

	frequency := 0
	stuffedBits.Iterate(func(_, bit int) error {
		if bit == 0 {
			frequency = 1 - frequency
		}
		output = append(output, f.getTone(1, frequency)...)
		return nil
	})

	return output
}

func (f *Bell202Modem) Demodulate(samples []float64) []byte {
	bytes := make([]byte, 0, 1+len(samples)/f.inputSymbolSamples)

	for _, sample := range samples {
		value := f.getFuzzySymbol(sample)
		if value*f.lastValue < 0 {
			slcInBitLengths := float64(f.sinceLastCrossing) / float64(f.inputSymbolSamples)
			onesBetween := int(math.Round(slcInBitLengths - 1))
			for i := 0; i < onesBetween; i++ {
				bytes = append(bytes, f.bitDestuffer.FeedBit(1)...)
			}
			bytes = append(bytes, f.bitDestuffer.FeedBit(0)...)
			f.sinceLastCrossing = 0
		} else {
			f.sinceLastCrossing++
		}
		f.lastValue = value
	}

	return bytes
}
