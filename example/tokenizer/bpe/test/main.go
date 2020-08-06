package main

import (
	"fmt"
	"log"

	bpe "github.com/sugarme/tokenizer/model/bpe"
	"github.com/sugarme/tokenizer/pretokenizer"
	"github.com/sugarme/tokenizer/tokenizer"
)

func main() {
	model, err := bpe.NewBpeFromFiles("example/tokenizer/bpe/test/model/es-vocab.json", "example/tokenizer/bpe/test/model/es-merges.txt")
	if err != nil {
		log.Fatal(err)
	}

	tk := tokenizer.NewTokenizer(model)

	bl := pretokenizer.NewByteLevel()

	tk.WithPreTokenizer(bl)

	sentence := "Mi estas Julien."

	en := tk.Encode(sentence)

	fmt.Printf("Sentence: '%v'\n", sentence)

	fmt.Printf("Tokens: %+v\n", en.GetTokens())

	// for _, tok := range en.GetTokens() {
	// fmt.Printf("'%v'\n", tok)
	// }

}
