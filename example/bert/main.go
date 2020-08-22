package main

import (
	"fmt"
	"log"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/model/wordpiece"
	"github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/pretokenizer"
	// "github.com/sugarme/tokenizer/processor"
)

func main() {
	model, err := wordpiece.NewWordPieceFromFile("../../data/bert/vocab.txt", "[UNK]")
	if err != nil {
		log.Fatal(err)
	}

	tk := tokenizer.NewTokenizer(model)

	fmt.Printf("Vocab size: %v\n", tk.GetVocabSize(false))

	bertPreTokenizer := pretokenizer.NewBertPreTokenizer()
	tk.WithPreTokenizer(bertPreTokenizer)

	bertNormalizer := normalizer.NewBertNormalizer(true, true, true, true)
	tk.WithNormalizer(bertNormalizer)

	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("[MASK]", true)})

	// sepId, ok := tk.TokenToId("[SEP]")
	// if !ok {
	// log.Fatalf("Cannot find ID for [SEP] token.\n")
	// }
	// sep := processor.PostToken{Id: sepId, Value: "[SEP]"}
	//
	// clsId, ok := tk.TokenToId("[CLS]")
	// if !ok {
	// log.Fatalf("Cannot find ID for [CLS] token.\n")
	// }
	// cls := processor.PostToken{Id: clsId, Value: "[CLS]"}
	//
	// postProcess := processor.NewBertProcessing(sep, cls)
	// tk.WithPostProcessor(postProcess)

	// sentence := `Hello, y'all! How are you 😁 ?`
	sentence := `Yesterday I saw a [MASK] far away`
	fmt.Printf("Sentence: '%v'\n", sentence)

	input := tokenizer.NewInputSequence(sentence)
	en, err := tk.Encode(tokenizer.NewSingleEncodeInput(input), true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tokens: %v\n", en.GetTokens())
	fmt.Printf("offsets: %v\n", en.GetOffsets())
}
