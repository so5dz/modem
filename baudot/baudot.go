package baudot

import (
	"strings"
)

type BaudotTable string

const (
	Letters BaudotTable = "\000E\nA SIU\rDRJNFCKTZLWHYPQOBG\000MXV\000"
	Figures BaudotTable = "\0003\n- \a87\r$4',!:(5\")2#6019?&\000./;\000"
)

type BaudotMode int

const (
	BaudotMode_Letters BaudotMode = 1
	BaudotMode_Figures BaudotMode = 2
)

type BaudotSymbolizer struct {
	txMode BaudotMode
	rxMode BaudotMode
}

const (
	LettersShift   int = 31
	FiguresShift   int = 27
	LineFeed       int = 2
	CarriageReturn int = 8
)

func (s *BaudotSymbolizer) Symbolize(bytes []byte) []int {
	words := make([]int, 0)

	words = append(words, LettersShift)
	words = append(words, LettersShift)
	s.txMode = BaudotMode_Letters

	for _, c := range bytes {
		if c == '\n' {
			words = append(words, LineFeed)

		} else if c == '\r' {
			words = append(words, CarriageReturn)

		} else if isLowercase(c) || isUppercase(c) || isSpace(c) {
			if isLowercase(c) {
				c -= 32
			}

			if s.txMode != BaudotMode_Letters {
				words = append(words, LettersShift)
				s.txMode = BaudotMode_Letters
			}

			words = append(words, toBaudot(c, Letters))

		} else {
			b := toBaudot(c, Figures)

			if b > 0 && s.txMode != BaudotMode_Figures {
				words = append(words, FiguresShift)
				s.txMode = BaudotMode_Figures
			}

			words = append(words, b)
		}
	}

	words = append(words, 0)
	return words
}

func (s *BaudotSymbolizer) Desymbolize(words []int) []byte {
	bytes := make([]byte, 0, len(words))

	for _, symbol := range words {
		if symbol == LettersShift {
			s.rxMode = BaudotMode_Letters
		} else if symbol == FiguresShift {
			s.rxMode = BaudotMode_Figures
		} else if s.rxMode == BaudotMode_Letters {
			bytes = append(bytes, byte(toChar(symbol, Letters)))
		} else if s.rxMode == BaudotMode_Figures {
			bytes = append(bytes, byte(toChar(symbol, Figures)))
		}
	}

	return bytes
}

func isLowercase(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isUppercase(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isSpace(c byte) bool {
	return c == ' '
}

func toBaudot(c byte, table BaudotTable) int {
	return strings.Index(string(table), string(c))
}

func toChar(s int, table BaudotTable) int {
	return int(string(table)[s])
}
