package main

import (
	"fmt"
	"log"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
)

func main() {
	// any model with file `tokenizer.json` available. Eg. `"tiiuae/falcon-7b"`
	configFile, err := tokenizer.CachedPath("bert-base-uncased", "tokenizer.json")
	if err != nil {
		panic(err)
	}

	tk, err := pretrained.FromFile(configFile)
	if err != nil {
		panic(err)
	}

	sentence := `The Gophers craft code using [MASK] language.`
	en, err := tk.EncodeSingle(sentence)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tokens: %q\n", en.Tokens)
	fmt.Printf("offsets: %v\n", en.Offsets)
}

// tokens: ["the" "go" "##pher" "##s" "craft" "code" "using" "[MASK]" "language" "."]
// offsets: [[0 3] [4 6] [6 10] [10 11] [12 17] [18 22] [23 28] [29 35] [36 44] [44 45]]
