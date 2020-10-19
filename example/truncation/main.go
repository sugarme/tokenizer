package main

import (
	"fmt"
	"log"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
)

func main() {
	tk := pretrained.RobertaBaseSquad2(false, true)

	maxLen := 25

	truncParams := tokenizer.TruncationParams{
		MaxLength: maxLen,
		Strategy:  tokenizer.OnlySecond,
		Stride:    10,
	}
	tk.WithTruncation(&truncParams)

	padToken := "<pad>"
	paddingStrategy := tokenizer.NewPaddingStrategy(tokenizer.WithFixed(maxLen))
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

	pairEn, err := tk.EncodePair(input, pairInput, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Offsets: %+v - length: %v\n\n", pairEn.Offsets, len(pairEn.Offsets))
	fmt.Printf("Words: %+v - length: %v\n\n", pairEn.Words, len(pairEn.Words))
	fmt.Printf("Ids: %+v - length: %v\n\n", pairEn.Ids, len(pairEn.Ids))
	fmt.Printf("TypeIds: %+v - lenght: %v\n\n", pairEn.TypeIds, len(pairEn.TypeIds))
	fmt.Printf("SpecialTokenMask: %+v - length: %v\n\n", pairEn.SpecialTokenMask, len(pairEn.SpecialTokenMask))
	fmt.Printf("AttentionMask: %+v - length: %v\n\n", pairEn.AttentionMask, len(pairEn.AttentionMask))

	fmt.Printf("Tokens: %q - length: %v\n\n", pairEn.Tokens, len(pairEn.Tokens))
	// overflowing encodings with stride
	for i, en := range pairEn.Overflowing {
		fmt.Printf("Overflow %v - Tokens: %q - length: %v\n", i, en.Tokens, len(en.Tokens))
	}
}
