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

func getBert() (retVal *tokenizer.Tokenizer) {

	vocabFile := "../../data/bert-base-uncased-vocab.txt"
	model, err := wordpiece.NewWordPieceFromFile(vocabFile, "[UNK]")
	if err != nil {
		log.Fatal(err)
	}

	tk := tokenizer.NewTokenizer(model)

	bertNormalizer := normalizer.NewBertNormalizer(true, true, true, true)
	tk.WithNormalizer(bertNormalizer)

	bertPreTokenizer := pretokenizer.NewBertPreTokenizer()
	tk.WithPreTokenizer(bertPreTokenizer)

	truncParams := tokenizer.TruncationParams{
		MaxLength: 25,
		Strategy:  tokenizer.OnlySecond,
		Stride:    0,
	}
	tk.WithTruncation(&truncParams)

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

	return tk
}

func main() {
	tk := getBert()

	input := "A visually stunning rumination on love."
	pairInput := "This is the long paragraph that I want to put context on it. It is not only about how to deal with anger but also how to maintain being calm at all time."

	// encodeInput := tokenizer.NewDualEncodeInput(tokenizer.NewInputSequence(input), tokenizer.NewInputSequence(pairInput))
	// pairEn, err := tk.Encode(encodeInput, false)
	pairEn, err := tk.EncodePair(input, pairInput)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Ids: %+v - length: %v\n\n", pairEn.Ids, len(pairEn.Ids))
	fmt.Printf("Tokens: %+q - length: %v\n\n", pairEn.Tokens, len(pairEn.Tokens))
	fmt.Printf("Offsets: %+v\n\n", pairEn.Offsets)
	fmt.Printf("Words: %+v\n\n", pairEn.Words)
	fmt.Printf("Overflow: %+v\n\n", pairEn.Overflowing)
	fmt.Printf("TypeIds: %+v\n\n", pairEn.TypeIds)
	fmt.Printf("SpecialTokenMask: %+v\n\n", pairEn.SpecialTokenMask)
	fmt.Printf("AttentionMask: %+v\n\n", pairEn.AttentionMask)
}
