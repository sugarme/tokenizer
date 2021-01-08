package main

import (
	"fmt"
	"log"
	"unicode"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/decoder"
	"github.com/sugarme/tokenizer/model/wordlevel"
	"github.com/sugarme/tokenizer/normalizer"
	// "github.com/sugarme/tokenizer/pretokenizer"
)

type customNormalizer struct{}

// implement Normalizer interface
// - Lowercase
// - Remove accents
func (n *customNormalizer) Normalize(input *normalizer.NormalizedString) (*normalizer.NormalizedString, error) {
	return input.Lowercase().RemoveAccents(), nil
}

type customPreTokenizer struct{}

// implement PreTokenizer interface
// - Split on whitespace
// - Split on punctuation
func (pt *customPreTokenizer) PreTokenize(pretokenized *tokenizer.PreTokenizedString) (*tokenizer.PreTokenizedString, error) {
	pretok := pretokenized.Split(func(noop int, sub *normalizer.NormalizedString) []tokenizer.SplitIdx {
		var splits []normalizer.NormalizedString

		// split on whitespace
		whitespace := normalizer.NewRegexpPattern(`\s+`)
		wsSubs := sub.Split(whitespace, normalizer.RemovedBehavior)

		// split on punctuation
		for _, sub := range wsSubs {
			// puncSubs := sub.Split(normalizer.NewFnPattern(normalizer.IsPunctuation), normalizer.IsolatediBehavior)
			puncSubs := sub.Split(normalizer.NewFnPattern(unicode.IsPunct), normalizer.IsolatediBehavior)
			splits = append(splits, puncSubs...)
		}

		var splitIdxs []tokenizer.SplitIdx
		for _, s := range splits {
			normalized := s
			splitIdx := tokenizer.SplitIdx{Normalized: &normalized, Tokens: nil}
			splitIdxs = append(splitIdxs, splitIdx)
		}

		return splitIdxs
	})

	return pretok, nil
}

func getWordTokenizer() *tokenizer.Tokenizer {
	wlbuilder := wordlevel.NewWordLevelBuilder()
	// wlbuilder.UnkToken("[UNK]")
	model := wlbuilder.Build()
	tk := tokenizer.NewTokenizer(model)

	// custom normalizer
	n := new(customNormalizer)
	tk.WithNormalizer(n)

	// custom pretokenizer
	pt := new(customPreTokenizer)
	tk.WithPreTokenizer(pt)

	// Added tokens to vocab
	tk.AddTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("goodnight", false)})
	tk.AddTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("goodmorning", false)})
	tk.AddTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("hello", false)})

	// Added decoder
	wlDecoder := decoder.DefaultWordpieceDecoder()
	tk.WithDecoder(wlDecoder)

	return tk
}

func runWordLevel() {
	tk := getWordTokenizer()

	toks, err := tk.Tokenize("Hello world!")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q\n", toks)

	enc, err := tk.EncodeSingle("Hello World!")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", enc)

	decodeStr := tk.Decode(enc.Ids, true)
	fmt.Printf("decoded string: %q\n", decodeStr)

}
