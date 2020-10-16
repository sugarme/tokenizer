package tokenizer_test

import (
	"fmt"
	"log"

	"github.com/sugarme/tokenizer/pretrained"
)

func ExampleTokenizer_Encode() {

	tk := pretrained.BertBaseUncased()
	sentence := `Yesterday I saw a [MASK] far away`

	en, err := tk.EncodeSingle(sentence)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tokens: %v\n", en.GetTokens())
	fmt.Printf("offsets: %v\n", en.GetOffsets())

	// Output:
	// tokens: [yesterday i saw a [MASK] far away]
	// offsets: [[0 9] [10 11] [12 15] [16 17] [18 24] [25 28] [29 33]]
}

func ExamplePreTokenizer_Split() {

	tk := pretrained.BertBaseUncased()

	sentence := `Hello, y'all! How are you 😁 ?`

	en, err := tk.EncodeSingle(sentence, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tokens: %v\n", en.GetTokens())
	fmt.Printf("offsets: %v\n", en.GetOffsets())
	fmt.Printf("word Ids: %v\n", en.GetWords())

	// Output:
	// tokens: [[CLS] hello , y ' all ! how are you [UNK] ? [SEP]]
	// offsets: [[0 0] [0 5] [5 6] [7 8] [8 9] [9 12] [12 13] [14 17] [18 21] [22 25] [26 30] [31 32] [0 0]]
	// word Ids: [-1 0 1 2 3 4 5 6 7 8 9 10 -1]
}
