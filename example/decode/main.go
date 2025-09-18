package main

import (
	"fmt"

	"github.com/sugarme/tokenizer/pretrained"
)

func main() {
	tk := pretrained.BertBaseUncased()

	en, err := tk.EncodeSingle("Goodmorning, how are you today?", true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("tokens: %+v\n", en.Tokens)
	fmt.Printf("decoded string: '%v'\n", tk.Decode(en.Ids, true))
}
