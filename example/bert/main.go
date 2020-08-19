package main

import (
	"fmt"
	"log"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/model/wordpiece"
	"github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/pretokenizer"
	"github.com/sugarme/tokenizer/processor"
)

func main() {
	model, err := wordpiece.NewWordPieceFromFile("../../data/bert/vocab.txt", "[UNK]")
	if err != nil {
		log.Fatal(err)
	}

	tk := tokenizer.NewTokenizer(model)

	fmt.Printf("Vocab size: %v\n", tk.GetVocabSize(false))

	// var id uint32 = 2500
	// val, _ := tk.IdToToken(id)
	// fmt.Printf("Value at Key %v: %v\n", id, val)

	bertPreTokenizer := pretokenizer.NewBertPreTokenizer()
	tk.WithPreTokenizer(bertPreTokenizer)

	bertNormalizer := normalizer.NewBertNormalizer(true, true, true, true)
	tk.WithNormalizer(bertNormalizer)

	sepId, ok := tk.TokenToId("[SEP]")
	if !ok {
		log.Fatalf("Cannot find ID for [SEP] token.\n")
	}
	sep := processor.PostToken{Id: sepId, Value: "[SEP]"}

	clsId, ok := tk.TokenToId("[CLS]")
	if !ok {
		log.Fatalf("Cannot find ID for [CLS] token.\n")
	}
	cls := processor.PostToken{Id: clsId, Value: "[CLS]"}

	postProcess := processor.NewBertProcessing(sep, cls)
	tk.WithPostProcessor(postProcess)

	sentence := `Hello, y'all! How are you 😁 ?`
	// sentence := `a visually stunning rumination on love`

	// en := tk.Encode(sentence)
	en := tk.Encode(tokenizer.NewSingleEncodeInput(sentence))
	// en := tk.Encode(tokenizer.NewDualEncodeInput(sentence, "And its pair."))

	// fmt.Printf("Sentence: '%v'\n", sentence)

	// Output should be:
	// [101, 7592, 1010, 1061, 1005, 2035, 999, 2129, 2024, 2017, 100, 1029, 102]
	// ['[CLS]', 'hello', ',', 'y', "'", 'all', '!', 'how', 'are', 'you', '[UNK]', '?', '[SEP]']
	// [(0, 0), (0, 5), (5, 6), (7, 8), (8, 9), (9, 12), (12, 13), (14, 17), (18, 21), (22, 25), (26, 27),
	// (28, 29), (0, 0)]
	fmt.Printf("Original string: \t'%v'\n", en.Normalized.GetOriginal())
	fmt.Printf("Normalized string: \t'%v'\n", en.Normalized.GetNormalized())
	fmt.Printf("Ids: \t\t\t%v\n", en.GetIds())
	fmt.Printf("Tokens: \t\t%+v\n", en.GetTokens())
	fmt.Printf("Offsets: \t\t%v\n", en.GetOffsets())
	expected := `[{0 0} {0 5} {5 6} {7 8} {8 9} {9 12} {12 13} {14 17} {18 21} {22 25} {26 27} {28 29} {0 0}]`
	fmt.Printf("Expected: \t\t%v\n", expected)

	// for _, tok := range en.GetTokens() {
	// fmt.Printf("'%v'\n", tok)
	// }

}
