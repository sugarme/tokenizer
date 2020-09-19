package main

import (
	"fmt"
	"log"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/model/bpe"
	"github.com/sugarme/tokenizer/pretokenizer"
)

func main() {

	model, err := bpe.DefaultBPE()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("model vocab: %+v\n", model.GetVocab())

	// var vocab map[string]int = make(map[string]int)
	// vocab["<s>"] = 0
	// vocab["<pad>"] = 1
	// vocab["</s>"] = 2
	// vocab["<unk>"] = 3
	// vocab["<mask>"] = 4
	//
	// var merges bpe.Merges = make(map[bpe.Pair]bpe.PairVal)
	//
	// model := bpe.NewBPE(vocab, merges)

	tk := tokenizer.NewTokenizer(model)

	unkToken := "<unk>"
	model.UnkToken = &unkToken

	specialToks := []tokenizer.AddedToken{
		tokenizer.NewAddedToken("<s>", true),
		tokenizer.NewAddedToken("<pad>", true),
		tokenizer.NewAddedToken("</s>", true),
		tokenizer.NewAddedToken("<unk>", true),
		tokenizer.NewAddedToken("<mask>", true),
	}

	tk.AddSpecialTokens(specialToks)

	fmt.Printf("vocab size: %v\n", tk.GetVocabSize(true))
	// fmt.Printf("vocab: %+v\n", tk.GetModel().GetVocab())
	fmt.Printf("vocab: %+v\n", tk.GetVocab(true))

	bytelevel := pretokenizer.NewByteLevel()
	tk.WithPreTokenizer(bytelevel)

	input := "Hello world!"
	inputSeq := tokenizer.NewInputSequence(input)

	encode, err := tk.Encode(tokenizer.NewSingleEncodeInput(inputSeq), true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("encoding: %+v\n", encode)

}
