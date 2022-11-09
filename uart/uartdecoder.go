package uart

type UARTState int

const (
	UartState_Idle  UARTState = 0
	UartState_Start UARTState = 10
	UartState_Data  UARTState = 20
	UartState_Stop  UARTState = 30
)

type LineState bool

const (
	LineState_Low   LineState = false
	LineState_High  LineState = true
	LineState_Start           = LineState_Low
	LineState_Stop            = LineState_High
)

type UARTDecoder struct {
	BitLength int
	Data      int
	Parity    Parity
	Stop      float64

	state            UARTState
	uartStateCounter int
	lastLineState    LineState
	lowStateCounter  int
	highStateCounter int
	dataBuffer       int
}

func (r *UARTDecoder) Feed(lineState LineState) []byte {
	output := make([]byte, 0, 1)

	if lineState == LineState_High {
		r.highStateCounter++
	} else if lineState == LineState_Low {
		r.lowStateCounter++
	}

	r.uartStateCounter--

	switch r.state {
	case UartState_Idle:
		if (r.lastLineState != LineState_Start) && (lineState == LineState_Start) {
			r.state = UartState_Start
			r.resetLineStateCounters()
			r.uartStateCounter = r.BitLength
		}

	case UartState_Start:
		if r.uartStateCounter == 0 {
			if r.lowStateCounter > r.BitLength/2 {
				r.state = UartState_Data
				r.uartStateCounter = r.Data * r.BitLength
				r.dataBuffer = 0
			} else {
				r.state = UartState_Idle
			}
		}

	case UartState_Data:
		if r.uartStateCounter%r.BitLength == r.BitLength/2 {
			r.dataBuffer <<= 1
			if lineState == LineState_High {
				r.dataBuffer |= 1
			}
		}
		if r.uartStateCounter == 0 {
			r.state = UartState_Stop
			r.uartStateCounter = int(0.8 * r.Stop * float64(r.BitLength))
			r.resetLineStateCounters()
		}

	case UartState_Stop:
		if r.uartStateCounter == 0 {
			if r.highStateCounter > int(0.5*r.Stop*float64(r.BitLength)) {
				output = append(output, byte(r.dataBuffer))
			}
			r.state = UartState_Idle
		}
	}

	r.lastLineState = lineState
	return output
}

func (r *UARTDecoder) resetLineStateCounters() {
	r.lowStateCounter = 0
	r.highStateCounter = 0
}
