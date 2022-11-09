package modem

import (
	"strings"
	"testing"
)

func TestRTTYModem_Amateur(t *testing.T) {
	testE2E(t, "one character", getRTTYModem(45.45, 432, 170), "A")
	testE2E(t, "one word", getRTTYModem(45.45, 1234, 170), "WORD")
	testE2E(t, "words with spaces", getRTTYModem(45.45, 2345, 170), "RYRYRY WORD RYRYRY")
	testE2E(t, "words and numbers", getRTTYModem(45.45, 1500, 170), "MSG 1 WORD 12345 POS 12.34 N 56.78 E")
}

func TestRTTYModem_RTTY75N(t *testing.T) {
	getModem := func() Modem {
		return getRTTYModem(75, 1500, 170)
	}
	testE2E(t, "one character", getModem(), "A")
	testE2E(t, "one word", getModem(), "WORD")
	testE2E(t, "words with spaces", getModem(), "RYRYRY WORD RYRYRY")
	testE2E(t, "words and numbers", getModem(), "MSG 1 WORD 12345 POS 12.34 N 56.78 E")
}

func TestRTTYModem_Random(t *testing.T) {
	testE2E(t, "25 baud shift 85", getRTTYModem(25, 1500, 85), "MSG 1 WORD 12345")
	testE2E(t, "50 baud shift 170", getRTTYModem(50, 1500, 170), "MSG 1 WORD 12345 POS 12.34 N 56.78 E")
	testE2E(t, "100 baud shift 300", getRTTYModem(100, 1500, 300), "MSG 1 WORD 12345 POS 12.34 N 56.78 E")
	testE2E(t, "150 baud shift 450 longer msg", getRTTYModem(150, 1500, 450), "MSG 1 WORD 12345 POS 12.34 N 56.78 E A 1 B 2 C 3 D 4 HELLO WORLD THIS IS A TEST RYRYRY 1234567890")
	testE2E(t, "200 baud shift 800 longer msg", getRTTYModem(200, 1600, 800), "MSG 1 WORD 12345 POS 12.34 N 56.78 E A 1 B 2 C 3 D 4 HELLO WORLD THIS IS A TEST RYRYRY 1234567890")
}

func testE2E(t *testing.T, name string, modem Modem, inputString string) {
	t.Run(name, func(t *testing.T) {
		input := []byte(inputString)
		samples := modem.Modulate(input)
		output := modem.Demodulate(samples)
		outputString := string(output)

		if !arraysEqual(input, output) && !strings.Contains(outputString, inputString) {
			t.Errorf("got '%v', want '%v'", outputString, inputString)
		}
	})
}

func getRTTYModem(speed float64, center float64, shift float64) *RTTYModem {
	var rttyModem RTTYModem
	rttyModem.Initialize(speed, 8000, 8000, center, shift, 5, 1.5)
	return &rttyModem
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
