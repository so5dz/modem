package modem

import (
	"strings"
	"testing"
)

func TestBell202Modem(t *testing.T) {
	testBell202Modem(t, "character", getBell202Modem(1000), "a")
	testBell202Modem(t, "word", getBell202Modem(1234.5), "Word")
	testBell202Modem(t, "sentence", getBell202Modem(2345.6), "Libero excepturi suscipit deserunt ut consectetur voluptatem. Rerum provident nihil beatae praesentium. Libero aperiam corporis eaque. Hic voluptates accusantium accusamus nostrum")
	testBell202Modem(t, "paragraph", getBell202Modem(2000.1234), "Libero excepturi suscipit deserunt ut consectetur voluptatem. Rerum provident nihil beatae praesentium. Libero aperiam corporis eaque. Hic voluptates accusantium accusamus nostrum. Quae aperiam beatae officiis repellendus suscipit laboriosam voluptatibus. Quod molestias doloremque qui veritatis. Non eum iure dolores. Soluta omnis odit voluptatem libero rem eum. Rerum nulla id numquam dolores rerum totam non voluptatum. Explicabo asperiores enim numquam ut veniam inventore. Vitae delectus qui quos. Praesentium sunt quas quos et. Animi quod placeat magnam est. Eum dicta cupiditate velit consequatur.")
}

func testBell202Modem(t *testing.T, name string, modem Modem, inputString string) {
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

func getBell202Modem(center float64) *Bell202Modem {
	var bellModem Bell202Modem
	bellModem.Initialize(8000, 8000, center)
	return &bellModem
}
