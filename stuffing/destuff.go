package stuffing

type BitDestufferState int

const (
	BitDestufferState_Normal       BitDestufferState = 0
	BitDestufferState_IgnoreZero   BitDestufferState = 10
	BitDestufferState_PossibleFlag BitDestufferState = 20
)

type BitDestuffer struct {
	BreakBit   int
	BreakEvery int
	streak     int
	buffer     byte
	bufferSize int
	state      BitDestufferState
}

func (bd *BitDestuffer) Reset() {
	bd.streak = 0
	bd.buffer = 0
	bd.bufferSize = 0
	bd.state = BitDestufferState_Normal
}

func (bd *BitDestuffer) FeedBit(bit int) []byte {
	optionalByte := make([]byte, 0, 1)

	switch bd.state {
	case BitDestufferState_Normal:
		bd.buffer <<= 1
		bd.buffer += byte(bit & 1)
		bd.bufferSize++

		if bd.bufferSize == 8 {
			bd.bufferSize = 0
			optionalByte = append(optionalByte, bd.buffer)
			bd.buffer = 0
		}

		if bit == bd.BreakBit {
			bd.streak++
		} else {
			bd.streak = 0
		}

		if bd.streak == bd.BreakEvery {
			bd.state = BitDestufferState_IgnoreZero
		}

	case BitDestufferState_IgnoreZero:
		if bit == bd.BreakBit {
			bd.streak++
			bd.state = BitDestufferState_PossibleFlag
		} else {
			bd.streak = 0
			bd.state = BitDestufferState_Normal
		}

	case BitDestufferState_PossibleFlag:
		if bit == bd.BreakBit {
			bd.streak++
		} else {
			bd.streak = 0
			bd.bufferSize = 0
			bd.state = BitDestufferState_Normal
		}

	}

	return optionalByte
}
