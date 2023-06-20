package pretokenizer

import (
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
)

type Whitespace struct{}

func NewWhitespace() *Whitespace {
	return new(Whitespace)
}

func DefaultWhitespace() *Whitespace {
	return new(Whitespace)
}

// Implement tokenizer.PreTokenizer for Whitespace

var _ tokenizer.PreTokenizer = new(Whitespace)

func (p *Whitespace) PreTokenize(pretokenized *tokenizer.PreTokenizedString) (*tokenizer.PreTokenizedString, error) {
	pretok := pretokenized.Split(func(noop int, normalized *normalizer.NormalizedString) []tokenizer.SplitIdx {
		s := `\w+|[^\w\s]+`
		rePattern := normalizer.NewRegexpPattern(s)
		invert := normalizer.NewInvertPattern(rePattern)
		splits := normalized.Split(invert, normalizer.RemovedBehavior)

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

type WhitespaceSplit struct{}

func NewWhitespaceSplit() *WhitespaceSplit {
	return new(WhitespaceSplit)
}

// Implement tokenizer.PreTokenizer for WhitespaceSplit

var _ tokenizer.PreTokenizer = new(WhitespaceSplit)

func (p *WhitespaceSplit) PreTokenize(pretokenized *tokenizer.PreTokenizedString) (*tokenizer.PreTokenizedString, error) {
	pretok := pretokenized.Split(func(noop int, normalized *normalizer.NormalizedString) []tokenizer.SplitIdx {
		isWhitespace := normalizer.NewFnPattern(normalizer.IsWhitespace)
		splits := normalized.Split(isWhitespace, normalizer.RemovedBehavior)

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
