package pretokenizer

import (
	"unicode"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
)

type Digits struct {
	IndividualDigits bool
}

func NewDigits(individualDigits bool) *Digits {
	return &Digits{individualDigits}
}

func DefaultDigits() *Digits {
	return NewDigits(false)
}

// PreTokenize implements tokenizer.PreTokenizer.
func (p *Digits) PreTokenize(pretokenized *tokenizer.PreTokenizedString) (*tokenizer.PreTokenizedString, error) {
	isNumeric := normalizer.NewFnPattern(unicode.IsNumber)
	var pretok *tokenizer.PreTokenizedString
	if p.IndividualDigits {
		pretok = pretokenized.Split(func(noop int, normalized *normalizer.NormalizedString) []tokenizer.SplitIdx {
			splits := normalized.Split(isNumeric, normalizer.IsolatedBehavior)
			var splitIdxs []tokenizer.SplitIdx
			for _, s := range splits {
				normalized := s
				splitIdx := tokenizer.SplitIdx{Normalized: &normalized, Tokens: nil}
				splitIdxs = append(splitIdxs, splitIdx)
			}

			return splitIdxs
		})
	} else {
		pretok = pretokenized.Split(func(noop int, normalized *normalizer.NormalizedString) []tokenizer.SplitIdx {
			splits := normalized.Split(isNumeric, normalizer.ContiguousBehavior)
			var splitIdxs []tokenizer.SplitIdx
			for _, s := range splits {
				normalized := s
				splitIdx := tokenizer.SplitIdx{Normalized: &normalized, Tokens: nil}
				splitIdxs = append(splitIdxs, splitIdx)
			}

			return splitIdxs
		})
	}

	return pretok, nil
}
