package pretokenizer

import (
	"unicode"

	"github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/tokenizer"
)

func isBertPunc(x rune) (retVal bool) {
	// TODO. check whether it covers all?
	return unicode.IsPunct(x)
}

type shouldSplitFn func(x rune) bool

// splitOn splits the given string if satisfying `shouldSplit` function and
// keeps track of the offsets.
func splitOn(s string, shouldSplit shouldSplitFn, includeSplitToken bool) (retVal []tokenizer.PreToken) {

	var (
		words  []tokenizer.PreToken
		offset uint   = 0
		word   []rune = make([]rune, 0)
	)

	for _, r := range s {
		if shouldSplit(r) {
			if len(word) > 0 {
				offsets := tokenizer.Offsets{
					Start: offset - uint(len(word)),
					End:   offset,
				}
				words = append(words, tokenizer.PreToken{
					Value:   string(word),
					Offsets: offsets,
				})
				word = make([]rune, 0)
			}

			if includeSplitToken {
				offsets := tokenizer.Offsets{
					Start: offset,
					End:   offset + 1,
				}
				words = append(words, tokenizer.PreToken{
					Value:   string([]rune{r}),
					Offsets: offsets,
				})
			}
		} else {
			word = append(word, r)
		}

		offset += 1
	}

	// Potential last word
	if len(word) > 0 {
		offsets := tokenizer.Offsets{
			Start: offset - uint(len(word)),
			End:   offset,
		}

		words = append(words, tokenizer.PreToken{
			Value:   string(word),
			Offsets: offsets,
		})
		word = make([]rune, 0)
	}

	return words
}

type BertTokenizer struct{}

// Implement PreTokenizer interface for BertTokenizer:
// ===================================================

func (bt *BertTokenizer) PreTokenize(normalized *normalizer.Normalized) (retVal []tokenizer.PreToken) {

	var splitTokens []tokenizer.PreToken

	shouldSplit := func(r rune) bool {
		return unicode.IsSpace(r)
	}

	preToks := splitOn(normalized.GetNormalized(), shouldSplit, false)
	for _, preTok := range preToks {
		token := preTok.Value
		offsets := preTok.Offsets

		splitToksTmp := splitOn(token, isBertPunc, true)

		var splitToks []tokenizer.PreToken
		for _, splitTok := range splitToksTmp {
			tok := splitTok.Value
			offStart := splitTok.Offsets.Start + offsets.Start
			offEnd := splitTok.Offsets.End + offsets.Start

			splitToks = append(splitToks, tokenizer.PreToken{
				Value:   tok,
				Offsets: tokenizer.Offsets{Start: offStart, End: offEnd},
			})
		}

		splitTokens = append(splitTokens, splitToks...)
	}

	return splitTokens
}
