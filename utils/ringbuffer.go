package utils

import misc "github.com/so5dz/utils/misc"

type bufferInUse int

const (
	bufferA bufferInUse = 1
	bufferB bufferInUse = 2
)

type RingBuffer struct {
	windowSize int
	bufferSize int

	bufferA   []float64
	positionA int

	bufferB   []float64
	positionB int
	inUse     bufferInUse
}

func (rb *RingBuffer) Initialize(windowSize int) {
	rb.windowSize = windowSize
	rb.bufferSize = windowSize*2 + 1
	rb.bufferA = make([]float64, rb.bufferSize)
	rb.bufferB = make([]float64, rb.bufferSize)
	rb.positionA = 0
	rb.positionB = 0
	rb.inUse = bufferA
}

func (rb *RingBuffer) Feed(f float64) {
	if rb.inUse == bufferA {
		rb.bufferA[rb.positionA] = f
		rb.positionA++

		if rb.positionA == rb.bufferSize-rb.windowSize {
			rb.positionB = 0
		} else if rb.positionA > rb.bufferSize-rb.windowSize {
			rb.bufferB[rb.positionB] = f
			rb.positionB++
		}
		if rb.positionA == rb.bufferSize {
			rb.inUse = bufferB
		}
	} else {
		rb.bufferB[rb.positionB] = f
		rb.positionB++

		if rb.positionB == rb.bufferSize-rb.windowSize {
			rb.positionA = 0
		} else if rb.positionB > rb.bufferSize-rb.windowSize {
			rb.bufferA[rb.positionA] = f
			rb.positionA++
		}
		if rb.positionB == rb.bufferSize {
			rb.inUse = bufferA
		}
	}
}

func (rb *RingBuffer) Window() []float64 {
	if rb.inUse == bufferA {
		if rb.positionA < rb.windowSize {
			return rb.bufferA[0:rb.windowSize]
		} else {
			return rb.bufferA[rb.positionA-rb.windowSize : rb.positionA]
		}
	} else {
		if rb.positionB < rb.windowSize {
			return rb.bufferB[0:rb.windowSize]
		} else {
			return rb.bufferB[rb.positionB-rb.windowSize : rb.positionB]
		}
	}
}

func (rb *RingBuffer) WindowOf(size int) ([]float64, error) {
	if size > rb.windowSize {
		return nil, misc.NewError("requested window is larger than set at initialization")
	}
	window := rb.Window()
	return window[rb.windowSize-size : size], nil
}

func (rb *RingBuffer) Average() float64 {
	sum := 0.0
	for _, f := range rb.Window() {
		sum += f
	}
	return sum / float64(rb.windowSize)
}

func (rb *RingBuffer) AverageOf(size int) (float64, error) {
	window, err := rb.WindowOf(size)
	if err != nil {
		return 0.0, err
	}
	sum := 0.0
	for _, f := range window {
		sum += f
	}
	return sum / float64(size), nil
}
