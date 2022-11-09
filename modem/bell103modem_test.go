package modem

import (
	"strings"
	"testing"
)

func TestBell103Modem(t *testing.T) {
	testBell103Modem(t, "character", getBell103Modem(600), "a")
	testBell103Modem(t, "word", getBell103Modem(1000), "Word")
	testBell103Modem(t, "sentence", getBell103Modem(2000), "Libero excepturi suscipit deserunt ut consectetur voluptatem. Rerum provident nihil beatae praesentium. Libero aperiam corporis eaque. Hic voluptates accusantium accusamus nostrum")
	testBell103Modem(t, "paragraph", getBell103Modem(1234.5678), "Libero excepturi suscipit deserunt ut consectetur voluptatem. Rerum provident nihil beatae praesentium. Libero aperiam corporis eaque. Hic voluptates accusantium accusamus nostrum. Quae aperiam beatae officiis repellendus suscipit laboriosam voluptatibus. Quod molestias doloremque qui veritatis. Non eum iure dolores. Soluta omnis odit voluptatem libero rem eum. Rerum nulla id numquam dolores rerum totam non voluptatum. Explicabo asperiores enim numquam ut veniam inventore. Vitae delectus qui quos. Praesentium sunt quas quos et. Animi quod placeat magnam est. Eum dicta cupiditate velit consequatur.")
}

func testBell103Modem(t *testing.T, name string, modem Modem, inputString string) {
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

func getBell103Modem(center float64) *Bell103Modem {
	var bellModem Bell103Modem
	bellModem.Initialize(8000, 8000, center)
	return &bellModem
}
