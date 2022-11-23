package modem

import "github.com/iskrapw/modem/utils/uart"

type Bell103Modem struct {
	BFSKModemBase
	uartEncoder uart.UARTEncoder
	uartDecoder uart.UARTDecoder
}

func (f *Bell103Modem) Initialize(inputRate float64, outputRate float64, center float64) {
	f.initialize(inputRate, outputRate, center, 200, 300)

	f.uartEncoder = uart.UARTEncoder{
		Data:   8,
		Parity: uart.None,
		Stop:   true,
	}

	f.uartDecoder = uart.UARTDecoder{
		Data:      8,
		Parity:    uart.None,
		Stop:      1,
		BitLength: f.inputSymbolSamples,
	}
}

func (f *Bell103Modem) Modulate(bytes []byte) []float64 {
	output := make([]float64, 0)
	output = append(output, f.getTone(8, 1)...)

	for _, symbol := range f.uartEncoder.Symbolize(bytes) {
		bit := 0
		if symbol == uart.Mark || symbol == uart.Stop {
			bit = 1
		}
		output = append(output, f.getTone(1, bit)...)
	}

	output = append(output, f.getTone(8, 1)...)
	return output
}

func (f *Bell103Modem) Demodulate(samples []float64) []byte {
	bytes := make([]byte, 0, 1+len(samples)/f.inputSymbolSamples)

	for _, s := range samples {
		value := f.getFuzzySymbol(s)

		lineState := uart.LineState_Low
		if value > 0 {
			lineState = uart.LineState_High
		}

		optionalByte := f.uartDecoder.Feed(lineState)
		if len(optionalByte) > 0 {
			bytes = append(bytes, optionalByte...)
		}
	}

	return bytes
}
