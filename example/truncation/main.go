package main

import (
	"fmt"
	"log"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
)

func main() {
	// tk := pretrained.RobertaBaseSquad2(false, true)
	tk := pretrained.BertLargeCasedWholeWordMaskingSquad()

	maxLen := 384

	truncParams := tokenizer.TruncationParams{
		MaxLength: maxLen,
		Strategy:  tokenizer.OnlySecond,
		Stride:    128,
	}
	tk.WithTruncation(&truncParams)

	// padToken := "<pad>"
	padToken := "[PAD]"
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

	// input := "A visually stunning rumination on love."
	// pairInput := "This is the long paragraph that I want to put context on it. It is not only about how to deal with anger but also how to maintain being calm at all time."

	input := "In what country is Normandy located?"
	pairInput := "The Normans (Norman: Nourmands; French: Normands; Latin: Normanni) were the people who in the 10th and 11th centuries gave their name to Normandy, a region in France. They were descended from Norse (\"Norman\" comes from \"Norseman\") raiders and pirates from Denmark, Iceland and Norway who, under their leader Rollo, agreed to swear fealty to King Charles III of West Francia. Through generations of assimilation and mixing with the native Frankish and Roman-Gaulish populations, their descendants would gradually merge with the Carolingian-based cultures of West Francia. The distinct cultural and ethnic identity of the Normans emerged initially in the first half of the 10th century, and it continued to evolve over the succeeding centuries."

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
