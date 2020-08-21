package pretokenizer

import (
	"unicode"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
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
		offset int    = 0
		word   []rune = make([]rune, 0)
	)

	for _, r := range s {
		if shouldSplit(r) {
			if len(word) > 0 {
				offsets := tokenizer.Offsets{
					Start: offset - len(word),
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
			Start: offset - len(word),
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

type BertPreTokenizer struct{}

func NewBertPreTokenizer() (retVal BertPreTokenizer) {
	return BertPreTokenizer{}
}

// PreTokenize implements PreTokenizer interface for BertPreTokenizer
func (bt BertPreTokenizer) PreTokenize(pretokenized tokenizer.PreTokenizedString) (retVal tokenizer.PreTokenizedString, err error) {

	err = pretokenized.Split(func(noop int, sub *normalizer.NormalizedString) []*normalizer.NormalizedString {

		var res []*normalizer.NormalizedString

		isWhiteSpace := func(r rune) bool {
			return unicode.IsSpace(r)
		}
		p := normalizer.NewFnPattern(isWhiteSpace)
		subs := sub.Split(p, normalizer.RemovedBehavior)

		for _, sub := range subs {
			splits := sub.Split(normalizer.NewFnPattern(isBertPunc), normalizer.IsolatediBehavior)
			res = append(res, splits...)
		}

		return res
	})

	return pretokenized, err
}
