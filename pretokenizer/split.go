package pretokenizer

import (
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
)

type Split struct {
	Pattern  normalizer.Pattern
	Behavior normalizer.SplitDelimiterBehavior
	Invert   bool
}

func NewSplit(pattern normalizer.Pattern, behavior normalizer.SplitDelimiterBehavior, invert bool) *Split {

	return &Split{
		Pattern:  pattern,
		Behavior: behavior,
		Invert:   invert,
	}
}

// Implement tokenizer.PreTokenizer for Split
var _ tokenizer.PreTokenizer = new(Split)

func (s *Split) PreTokenize(pretokenized *tokenizer.PreTokenizedString) (*tokenizer.PreTokenizedString, error) {
	if s.Invert {
		pretok := pretokenized.Split(func(noop int, normalized *normalizer.NormalizedString) []tokenizer.SplitIdx {
			invert := normalizer.NewInvertPattern(s.Pattern)
			splits := normalized.Split(invert, s.Behavior)

			var splitIdxs []tokenizer.SplitIdx
			for _, s := range splits {
				normalized := s
				splitIdx := tokenizer.SplitIdx{Normalized: &normalized, Tokens: nil}
				splitIdxs = append(splitIdxs, splitIdx)
			}

			return splitIdxs
		})

		return pretok, nil

	} else {
		pretok := pretokenized.Split(func(noop int, normalized *normalizer.NormalizedString) []tokenizer.SplitIdx {
			splits := normalized.Split(s.Pattern, s.Behavior)

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
}
