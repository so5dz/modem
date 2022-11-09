package stuffing

import "github.com/tunabay/go-bitarray"

type BitStuffer struct {
	Flag       *bitarray.BitArray
	BreakBit   int
	BreakEvery int
	streak     int
	buffer     *bitarray.BitArray
}

func (bs *BitStuffer) Start() {
	bs.streak = 0
	bs.buffer = bitarray.New()
}

func (bs *BitStuffer) PutFlag() {
	bs.streak = 0
	bs.PutBitsWithoutStuffing(bs.Flag)
}

func (bs *BitStuffer) PutMultipleFlags(count int) {
	bs.streak = 0
	bs.PutBitsWithoutStuffing(bs.Flag.Repeat(count))
}

func (bs *BitStuffer) PutBitsWithoutStuffing(bits *bitarray.BitArray) {
	bs.buffer = bs.buffer.Append(bits)
}

func (bs *BitStuffer) PutBytes(bytes []byte) {
	bs.PutBits(bitarray.NewFromBytes(bytes, 0, 8*len(bytes)))
}

func (bs *BitStuffer) PutBits(bits *bitarray.BitArray) {
	stuffedBuffer := bitarray.NewBuffer(bits.Len() + bits.Len()/bs.BreakEvery)

	i := 0
	bits.Iterate(func(_, b int) error {
		stuffedBuffer.PutBitAt(i, byte(b))
		i++

		if b == bs.BreakBit {
			bs.streak++
			if bs.streak == bs.BreakEvery {
				stuffedBuffer.PutBitAt(i, byte(1-bs.BreakBit))
				i++
				bs.streak = 0
			}
		} else {
			bs.streak = 0
		}

		return nil
	})

	bs.buffer = bs.buffer.Append(stuffedBuffer.Slice(0, i))
}

func (bs *BitStuffer) End() *bitarray.BitArray {
	return bs.buffer
}
