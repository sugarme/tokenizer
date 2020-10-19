package main

import (
	"fmt"
	"log"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/decoder"
	"github.com/sugarme/tokenizer/model/bpe"
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

	wpDecoder := decoder.DefaultWordpieceDecoder()
	tk.WithDecoder(wpDecoder)

	return tk
}

func getRoberta() (retVal *tokenizer.Tokenizer) {

	vocabFile := "../../data/roberta-qa-vocab.json"
	mergesFile := "../../data/roberta-qa-merges.txt"
	model, err := bpe.NewBpeFromFiles(vocabFile, mergesFile)
	if err != nil {
		log.Fatal(err)
	}

	tk := tokenizer.NewTokenizer(model)

	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("<s>", true)})
	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("</s>", true)})

	blPreTokenizer := pretokenizer.NewByteLevel()
	tk.WithPreTokenizer(blPreTokenizer)

	postProcess := processor.DefaultRobertaProcessing()
	tk.WithPostProcessor(postProcess)

	bpeDecoder := decoder.NewBpeDecoder("Ġ")
	tk.WithDecoder(bpeDecoder)

	return tk
}

func main() {
	// tk := getBert()
	tk := getRoberta()

	maxLen := 25

	truncParams := tokenizer.TruncationParams{
		MaxLength: maxLen,
		Strategy:  tokenizer.OnlySecond,
		Stride:    10,
	}
	tk.WithTruncation(&truncParams)

	// padToken := "[PAD]"
	padToken := "<pad>"
	paddingStrategy := tokenizer.NewPaddingStrategy(tokenizer.WithFixed(maxLen))
	// paddingStrategy := tokenizer.NewPaddingStrategy(tokenizer.WithBatchLongest())
	padId, ok := tk.TokenToId(padToken)
	if !ok {
		log.Fatalf("'ConvertExampleToFeatures' method call error: cannot find pad token in the vocab.\n")
	}

	paddingParams := tokenizer.PaddingParams{
		Strategy:  *paddingStrategy,
		Direction: tokenizer.Right, // padding right
		PadId:     padId,
		PadTypeId: 1,
		PadToken:  padToken,
	}
	tk.WithPadding(&paddingParams)

	input := "A visually stunning rumination on love."
	pairInput := "This is the long paragraph that I want to put context on it. It is not only about how to deal with anger but also how to maintain being calm at all time."

	// encodeInput := tokenizer.NewDualEncodeInput(tokenizer.NewInputSequence(input), tokenizer.NewInputSequence(pairInput))
	// pairEn, err := tk.Encode(encodeInput, false)
	pairEn, err := tk.EncodePair(input, pairInput, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Offsets: %+v - length: %v\n\n", pairEn.Offsets, len(pairEn.Offsets))
	fmt.Printf("Words: %+v - length: %v\n\n", pairEn.Words, len(pairEn.Words))
	fmt.Printf("Overflow: %+v - length: %v\n\n", pairEn.Overflowing, len(pairEn.Overflowing))
	fmt.Printf("Ids: %+v - length: %v\n\n", pairEn.Ids, len(pairEn.Ids))
	fmt.Printf("TypeIds: %+v - lenght: %v\n\n", pairEn.TypeIds, len(pairEn.TypeIds))
	fmt.Printf("SpecialTokenMask: %+v - length: %v\n\n", pairEn.SpecialTokenMask, len(pairEn.SpecialTokenMask))
	fmt.Printf("AttentionMask: %+v - length: %v\n\n", pairEn.AttentionMask, len(pairEn.AttentionMask))

	fmt.Printf("Tokens: %q - length: %v\n\n", pairEn.Tokens, len(pairEn.Tokens))
	// overflowing encodings with stride
	for i, en := range pairEn.Overflowing {
		fmt.Printf("Overflow %v - Tokens: %q - length: %v\n", i, en.Tokens, len(en.Tokens))
	}

	decodedStr := tk.Decode(pairEn.Ids, true)
	// decodedStr := tk.Decode(pairEn.Ids[5:11], true)

	fmt.Printf("Decoded string: '%v'\n", decodedStr)
}
