package pretokenizer

import (
	// "fmt"
	// "unicode"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
)

func isBertPunc(x rune) (retVal bool) {
	// TODO. check whether it covers all?
	// return unicode.IsPunct(x)
	return normalizer.IsBertPunctuation(x)
}

type shouldSplitFn func(x rune) bool

// splitOn splits the given string if satisfying `shouldSplit` function and
// keeps track of the offsets.
func splitOn(s string, shouldSplit shouldSplitFn, includeSplitToken bool) (retVal []tokenizer.PreToken) {

	var (
		words  []tokenizer.PreToken
		offset int    = 0
		word   []rune = make([]rune, 0)
	)

	for _, r := range s {
		if shouldSplit(r) {
			if len(word) > 0 {
				offsets := []int{offset - len(word), offset}
				words = append(words, tokenizer.PreToken{
					Value:   string(word),
					Offsets: offsets,
				})
				word = make([]rune, 0)
			}

			if includeSplitToken {
				offsets := []int{offset, offset + 1}
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
		offsets := []int{offset - len(word), offset}

		words = append(words, tokenizer.PreToken{
			Value:   string(word),
			Offsets: offsets,
		})
		word = make([]rune, 0)
	}

	return words
}

type BertPreTokenizer struct{}

func NewBertPreTokenizer() *BertPreTokenizer {
	return &BertPreTokenizer{}
}

// PreTokenize implements PreTokenizer interface for BertPreTokenizer
func (bt *BertPreTokenizer) PreTokenize(pretokenized *tokenizer.PreTokenizedString) (*tokenizer.PreTokenizedString, error) {
	pretok := pretokenized.Split(func(noop int, sub *normalizer.NormalizedString) []tokenizer.SplitIdx {
		var splits []normalizer.NormalizedString
		whitespace := normalizer.NewRegexpPattern(`\s+`)
		wsSubs := sub.Split(whitespace, normalizer.RemovedBehavior)

		for _, sub := range wsSubs {
			puncSubs := sub.Split(normalizer.NewFnPattern(isBertPunc), normalizer.IsolatedBehavior)
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
