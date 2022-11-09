package stuffing

import (
	"testing"
)

func TestBitDestuffer_every5ones(t *testing.T) {
	t.Run("0000000_", func(t *testing.T) {
		bd := newBitDestuffer()
		bd.Reset()
		result := feedMultiple(bd, 0, 7)
		expectedResult := []byte{}
		assertByteArraysEqual(t, result, expectedResult)
	})

	t.Run("00000000", func(t *testing.T) {
		bd := newBitDestuffer()
		bd.Reset()
		result := feedMultiple(bd, 0, 8)
		expectedResult := []byte{0x00}
		assertByteArraysEqual(t, result, expectedResult)
	})

	t.Run("111110100", func(t *testing.T) {
		bd := newBitDestuffer()
		bd.Reset()
		result := make([]byte, 0)
		result = append(result, feedMultiple(bd, 1, 5)...)
		result = append(result, bd.FeedBit(0)...)
		result = append(result, bd.FeedBit(1)...)
		result = append(result, feedMultiple(bd, 0, 2)...)
		expectedResult := []byte{0b11111100}
		assertByteArraysEqual(t, result, expectedResult)
	})

	t.Run("11111111", func(t *testing.T) {
		bd := newBitDestuffer()
		bd.Reset()
		result := feedMultiple(bd, 1, 8)
		expectedResult := []byte{}
		assertByteArraysEqual(t, result, expectedResult)
	})

	t.Run("010101_01111110_10001000", func(t *testing.T) {
		bd := newBitDestuffer()
		bd.Reset()
		result := make([]byte, 0)
		result = append(result, bd.FeedBit(1)...)
		result = append(result, bd.FeedBit(0)...)
		result = append(result, bd.FeedBit(1)...)
		result = append(result, bd.FeedBit(0)...)
		result = append(result, bd.FeedBit(1)...)

		result = append(result, bd.FeedBit(0)...)
		result = append(result, feedMultiple(bd, 1, 6)...)
		result = append(result, bd.FeedBit(0)...)

		result = append(result, bd.FeedBit(1)...)
		result = append(result, feedMultiple(bd, 0, 3)...)
		result = append(result, bd.FeedBit(1)...)
		result = append(result, feedMultiple(bd, 0, 3)...)

		expectedResult := []byte{0b10101011, 0b10001000}
		assertByteArraysEqual(t, result, expectedResult)
	})
}

func newBitDestuffer() *BitDestuffer {
	return &BitDestuffer{
		BreakBit:   1,
		BreakEvery: 5,
	}
}

func feedMultiple(bd *BitDestuffer, bit int, repetitions int) []byte {
	bytes := make([]byte, 0)
	for i := 0; i < repetitions; i++ {
		bytes = append(bytes, bd.FeedBit(bit)...)
	}
	return bytes
}

func assertByteArraysEqual(t *testing.T, got []byte, expected []byte) {
	if !arraysEqual(got, expected) {
		t.Errorf("got %v, expected %v", got, expected)
	}
}

func arraysEqual(a []byte, b []byte) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i, aElement := range a {
		bElement := b[i]
		if aElement != bElement {
			return false
		}
	}
	return true
}
