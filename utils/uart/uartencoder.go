package uart

type Parity int

const (
	None Parity = 0
	Odd  Parity = 1
	Even Parity = 2
)

type Symbol int

const (
	Start Symbol = -1
	Stop  Symbol = -2
	Space Symbol = 0
	Mark  Symbol = 1
)

type UARTEncoder struct {
	Data   int
	Parity Parity
	Stop   bool
}

func (s *UARTEncoder) SymbolizeInts(ints []int) []Symbol {
	bytes := make([]byte, len(ints))
	for i, b := range ints {
		bytes[i] = byte(b)
	}
	return s.Symbolize(bytes)
}

func (s *UARTEncoder) Symbolize(bytes []byte) []Symbol {
	output := make([]Symbol, 0, len(bytes)*s.byteSymbols())
	for _, b := range bytes {
		output = append(output, Start)

		shift := s.Data - 1
		for mask := 1 << shift; mask > 0; mask >>= 1 {
			if int(b)&mask > 0 {
				output = append(output, Mark)
			} else {
				output = append(output, Space)
			}
			shift -= 1
		}

		if s.Parity != None {
			// todo
		}

		if s.Stop {
			output = append(output, Stop)
		}
	}
	return output
}

func (s *UARTEncoder) byteSymbols() int {
	bs := 1 + s.Data
	if s.Parity != None {
		bs += 1
	}
	if s.Stop {
		bs += 1
	}
	return bs
}
