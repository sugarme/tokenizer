package main

import (
	"fmt"
	"log"

	"github.com/sugarme/tokenizer/pretrained"
)

func main() {

	tk := pretrained.BertBaseUncased()
	sentence := `Yesterday I saw a [MASK] far away`

	en, err := tk.EncodeSingle(sentence)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tokens: %v\n", en.GetTokens())
	fmt.Printf("offsets: %v\n", en.GetOffsets())

	// Output
	// tokens: [yesterday i saw a [MASK] far away]
	// offsets: [{0 9} {10 11} {12 15} {16 17} {18 24} {25 28} {29 33}]
}
