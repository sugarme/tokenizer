package main

import (
	"fmt"
	"log"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/decoder"
	"github.com/sugarme/tokenizer/model/wordpiece"
	"github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/pretokenizer"
	"github.com/sugarme/tokenizer/processor"
	"github.com/sugarme/tokenizer/util"
)

func runBERT() {
	// Example 1.
	splitOnAddedToken()

	// Example 2.
	bertTokenize()
}

func getBert() (retVal *tokenizer.Tokenizer) {

	util.CdToThis()
	vocabFile := "../../data/bert-base-uncased-vocab.txt"

	model, err := wordpiece.NewWordPieceFromFile(vocabFile, "[UNK]")
	if err != nil {
		log.Fatal(err)
	}

	tk := tokenizer.NewTokenizer(model)
	fmt.Printf("Vocab size: %v\n", tk.GetVocabSize(false))

	bertNormalizer := normalizer.NewBertNormalizer(true, true, true, true)
	tk.WithNormalizer(bertNormalizer)

	bertPreTokenizer := pretokenizer.NewBertPreTokenizer()
	tk.WithPreTokenizer(bertPreTokenizer)

	wpDecoder := decoder.DefaultWordpieceDecoder()

	tk.WithDecoder(wpDecoder)

	return tk
}

func splitOnAddedToken() {

	tk := getBert()
	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("[MASK]", true)})

	sentence := `Yesterday I saw a [MASK] far away`
	// sentence := `Hello, y'all! How are you 😁 ?`
	// Output:
	// tokens: [yesterday i saw a [MASK] far away]
	// offsets: [{0 9} {10 11} {12 15} {16 17} {18 24} {25 28} {29 33}]

	// sentence := `Looks like one [MASK] is missing`
	// Output:
	// tokens: [looks like one [MASK] is missing]
	// offsets: [{0 5} {6 10} {11 14} {15 21} {22 24} {25 32}]

	input := tokenizer.NewInputSequence(sentence)
	en, err := tk.Encode(tokenizer.NewSingleEncodeInput(input), false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tokens: %v\n", en.GetTokens())
	fmt.Printf("offsets: %v\n", en.GetOffsets())
	fmt.Printf("word idx: %v\n", en.GetWords())
}

func bertTokenize() {
	tk := getBert()

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

	input := tokenizer.NewInputSequence(sentence)
	en, err := tk.Encode(tokenizer.NewSingleEncodeInput(input), true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("original: '%v'\n", sentence)
	fmt.Printf("tokens: %v\n", en.GetTokens())
	fmt.Printf("offsets: %v\n", en.GetOffsets())
	fmt.Printf("word Ids: %v\n", en.GetWords())

	decodedStr := tk.Decode(en.Ids, true)
	fmt.Printf("decodedStr: '%v'\n", decodedStr)

	// Output:
	// original: 'Hello, y'all! How are you 😁 ?'
	// tokens: [[CLS] hello , y ' all ! how are you [UNK] ? [SEP]]
	// offsets: [{0 0} {0 5} {5 6} {7 8} {8 9} {9 12} {12 13} {14 17} {18 21} {22 25} {26 27} {28 29} {0 0}]
	// word Ids: [-1 0 1 2 3 4 5 6 7 8 9 10 -1]
}
