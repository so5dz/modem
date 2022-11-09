package modem

import (
	"math"

	"github.com/iskrapw/modem/baudot"
	"github.com/iskrapw/modem/uart"
	"github.com/iskrapw/modem/utils"
)

type RTTYModem struct {
	inputSymbolSamples      int
	outputSymbolSamples     int
	outputStopSymbolSamples int
	space                   float64
	mark                    float64
	toneGenerator           utils.ToneGen
	spaceDetector           utils.ToneDet
	markDetector            utils.ToneDet
	audioBuffer             utils.RingBuffer
	valueAvgBuffer          utils.RingBuffer
	uartEncoder             uart.UARTEncoder
	uartDecoder             uart.UARTDecoder
	baudotSymbolizer        baudot.BaudotSymbolizer
}

func (f *RTTYModem) Initialize(speed float64, inputRate float64, outputRate float64, center float64, shift float64, bps int, stopBits float64) {
	f.inputSymbolSamples = int(math.Round(inputRate / speed))
	f.outputSymbolSamples = int(math.Round(outputRate / speed))
	f.outputStopSymbolSamples = int(math.Round(stopBits * outputRate / speed))
	f.space = center - shift/2
	f.mark = center + shift/2

	f.toneGenerator.Initialize(outputRate)
	f.spaceDetector.Initialize(inputRate, f.space, f.inputSymbolSamples)
	f.markDetector.Initialize(inputRate, f.mark, f.inputSymbolSamples)

	audioBufferSize := f.spaceDetector.WindowSize()
	if f.markDetector.WindowSize() > audioBufferSize {
		audioBufferSize = f.markDetector.WindowSize()
	}
	f.audioBuffer.Initialize(audioBufferSize)
	f.valueAvgBuffer.Initialize(audioBufferSize / 10)

	f.uartEncoder = uart.UARTEncoder{
		Data:   bps,
		Parity: uart.None,
		Stop:   stopBits > 0,
	}

	f.uartDecoder = uart.UARTDecoder{
		BitLength: f.inputSymbolSamples,
		Data:      bps,
		Parity:    uart.None,
		Stop:      stopBits,
	}

	f.baudotSymbolizer = baudot.BaudotSymbolizer{}
}

func (f *RTTYModem) Modulate(bytes []byte) []float64 {
	output := make([]float64, 0)

	byteInts := make([]int, len(bytes))
	for i, b := range bytes {
		byteInts[i] = int(b)
	}

	output = append(output, f.toneGenerator.Tone(f.outputSymbolSamples, f.mark)...)

	baudotWords := f.baudotSymbolizer.Symbolize(bytes)
	uartSymbols := f.uartEncoder.SymbolizeInts(baudotWords)

	for _, symbol := range uartSymbols {
		frequency := f.space
		if symbol == uart.Mark || symbol == uart.Stop {
			frequency = f.mark
		}

		samples := f.outputSymbolSamples
		if symbol == uart.Stop {
			samples = f.outputStopSymbolSamples
		}

		output = append(output, f.toneGenerator.Tone(samples, frequency)...)
	}

	output = append(output, f.toneGenerator.Tone(f.outputSymbolSamples, f.mark)...)

	return output
}

func (f *RTTYModem) Demodulate(samples []float64) []byte {
	words := make([]byte, 0, 1+len(samples)/f.inputSymbolSamples)

	for _, s := range samples {
		value := f.getFuzzySymbol(s)

		lineState := uart.LineState_Low
		if value > 0 {
			lineState = uart.LineState_High
		}

		optionalByte := f.uartDecoder.Feed(lineState)
		if len(optionalByte) > 0 {
			words = append(words, optionalByte...)
		}
	}

	return f.baudotSymbolizer.Desymbolize(asInts(words))
}

func (f *RTTYModem) getFuzzySymbol(sample float64) float64 {
	f.audioBuffer.Feed(sample)
	spaceWindow, _ := f.audioBuffer.WindowOf(f.spaceDetector.WindowSize())
	spaceStrength := f.spaceDetector.Detect(spaceWindow)
	markWindow, _ := f.audioBuffer.WindowOf(f.markDetector.WindowSize())
	markStrength := f.markDetector.Detect(markWindow)

	value := 1 - 2*spaceStrength/(spaceStrength+markStrength)
	f.valueAvgBuffer.Feed(value)
	value = f.valueAvgBuffer.Average()
	return value
}

func asInts(bytes []byte) []int {
	ints := make([]int, len(bytes))
	for i, b := range bytes {
		ints[i] = int(b)
	}
	return ints
}
