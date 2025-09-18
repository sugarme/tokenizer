package pretokenizer

import (
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
)

type CharDelimiterSplit struct {
	Delimiter rune
}

func NewCharDelimiterSplit(delimiter rune) *CharDelimiterSplit {
	return &CharDelimiterSplit{delimiter}
}

// Implement tokenizer.PreTokenizer for CharDelimiterSplit

var _ tokenizer.PreTokenizer = new(CharDelimiterSplit)

func (d *CharDelimiterSplit) PreTokenize(pretokenized *tokenizer.PreTokenizedString) (*tokenizer.PreTokenizedString, error) {
	pretok := pretokenized.Split(func(noop int, normalized *normalizer.NormalizedString) []tokenizer.SplitIdx {
		delimiter := normalizer.NewRunePattern(d.Delimiter)
		splits := normalized.Split(delimiter, normalizer.RemovedBehavior)

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
