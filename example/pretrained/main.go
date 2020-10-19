package main

import (
	"fmt"
	"log"

	"github.com/sugarme/tokenizer/pretrained"
)

func main() {

	tk := pretrained.BertBaseUncased()

	sentence := `The Gophers craft code using [MASK] language.`
	en, err := tk.EncodeSingle(sentence)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tokens: %q\n", en.Tokens)
	fmt.Printf("offsets: %v\n", en.Offsets)

	// Output
	// tokens: ["the" "go" "##pher" "##s" "craft" "code" "using" "[MASK]" "language" "."]
	// offsets: [[0 3] [4 6] [6 10] [10 11] [12 17] [18 22] [23 28] [29 35] [36 44] [44 45]]
}
