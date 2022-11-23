package stuffing

import (
	"testing"

	"github.com/tunabay/go-bitarray"
)

func TestBitStuffer_every5ones(t *testing.T) {
	bs := &BitStuffer{
		Flag:       bitarray.NewFromByteBits([]byte{0, 1, 1, 1, 1, 1, 1, 0}),
		BreakBit:   1,
		BreakEvery: 5,
	}

	t.Run("nothing", func(t *testing.T) {
		bs.Start()
		result := bs.End()
		expectedResult := bitarray.NewFromByteBits([]byte{})
		assertEqual(t, result, expectedResult)
	})

	t.Run("1 flag", func(t *testing.T) {
		bs.Start()
		bs.PutFlag()
		result := bs.End()
		expectedResult := bitarray.NewFromByteBits([]byte{0, 1, 1, 1, 1, 1, 1, 0})
		assertEqual(t, result, expectedResult)
	})

	t.Run("3 ones", func(t *testing.T) {
		bs.Start()
		bs.PutBits(bitarray.NewFromByteBits([]byte{1, 1, 1}))
		result := bs.End()
		expectedResult := bitarray.NewFromByteBits([]byte{1, 1, 1})
		assertEqual(t, result, expectedResult)
	})

	t.Run("5 ones", func(t *testing.T) {
		bs.Start()
		bs.PutBits(bitarray.NewFromByteBits([]byte{1, 1, 1, 1, 1}))
		result := bs.End()
		expectedResult := bitarray.NewFromByteBits([]byte{1, 1, 1, 1, 1, 0})
		assertEqual(t, result, expectedResult)
	})

	t.Run("4 ones 1 zero 2 ones", func(t *testing.T) {
		bs.Start()
		bs.PutBits(bitarray.NewFromByteBits([]byte{1, 1, 1, 1, 0, 1, 1}))
		result := bs.End()
		expectedResult := bitarray.NewFromByteBits([]byte{1, 1, 1, 1, 0, 1, 1})
		assertEqual(t, result, expectedResult)
	})

	t.Run("5 ones 1 zero", func(t *testing.T) {
		bs.Start()
		bs.PutBits(bitarray.NewFromByteBits([]byte{1, 1, 1, 1, 1, 0}))
		result := bs.End()
		expectedResult := bitarray.NewFromByteBits([]byte{1, 1, 1, 1, 1, 0, 0})
		assertEqual(t, result, expectedResult)
	})

	t.Run("6 ones", func(t *testing.T) {
		bs.Start()
		bs.PutBits(bitarray.NewFromByteBits([]byte{1, 1, 1, 1, 1, 1}))
		result := bs.End()
		expectedResult := bitarray.NewFromByteBits([]byte{1, 1, 1, 1, 1, 0, 1})
		assertEqual(t, result, expectedResult)
	})

	t.Run("1 flag 1 one", func(t *testing.T) {
		bs.Start()
		bs.PutFlag()
		bs.PutBits(bitarray.NewFromByteBits([]byte{1}))
		result := bs.End()
		expectedResult := bitarray.NewFromByteBits([]byte{0, 1, 1, 1, 1, 1, 1, 0, 1})
		assertEqual(t, result, expectedResult)
	})

	t.Run("1 flag 6 ones 1 zero", func(t *testing.T) {
		bs.Start()
		bs.PutFlag()
		bs.PutBits(bitarray.NewFromByteBits([]byte{1, 1, 1, 1, 1, 1}))
		result := bs.End()
		expectedResult := bitarray.NewFromByteBits([]byte{0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1})
		assertEqual(t, result, expectedResult)
	})

	t.Run("1 flag 5 ones 1 flag", func(t *testing.T) {
		bs.Start()
		bs.PutFlag()
		bs.PutBits(bitarray.NewFromByteBits([]byte{1, 1, 1, 1, 1}))
		bs.PutFlag()
		result := bs.End()
		expectedResult := bitarray.NewFromByteBits([]byte{0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0})
		assertEqual(t, result, expectedResult)
	})

	t.Run("4 ones 1 flag 2 ones", func(t *testing.T) {
		bs.Start()
		bs.PutBits(bitarray.NewFromByteBits([]byte{1, 1, 1, 1}))
		bs.PutFlag()
		bs.PutBits(bitarray.NewFromByteBits([]byte{1, 1}))
		result := bs.End()
		expectedResult := bitarray.NewFromByteBits([]byte{1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1})
		assertEqual(t, result, expectedResult)
	})
}

func assertEqual(t *testing.T, got *bitarray.BitArray, expected *bitarray.BitArray) {
	if !got.Equal(expected) {
		t.Errorf("got %v, expected %v", got, expected)
	}
}
