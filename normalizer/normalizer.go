package normalizer

import (
	"golang.org/x/text/unicode/norm"
)

type Normalizer interface {
	Normalize(normalized *NormalizedString) (*NormalizedString, error)
}

type normalizer struct {
	Normalizer Normalizer
}

func newNormalizer(opts ...DefaultOption) *normalizer {
	n := NewDefaultNormalizer()
	for _, opt := range opts {
		opt(n)
	}
	return &normalizer{
		Normalizer: n,
	}
}

func (n *normalizer) Normalize(normalized *NormalizedString) (*NormalizedString, error) {

	return normalized, nil
}

type Option func(*normalizer)

// WithBertNormalizer creates normalizer with BERT normalization features.
func WithBertNormalizer(cleanText, lowercase, handleChineseChars, stripAccents bool) Option {
	return func(o *normalizer) {
		NewBertNormalizer(cleanText, lowercase, handleChineseChars, stripAccents)
	}
}

// WithUnicodeNormalizer creates normalizer with one of unicode NFD, NFC, NFKD, or NFKC normalization feature.
func WithUnicodeNormalizer(form norm.Form) Option {
	return func(o *normalizer) {
		NewUnicodeNormalizer(form)
	}

}

func NewNormalizer(opts ...Option) Normalizer {

	nml := newNormalizer()

	for _, o := range opts {
		o(nml)
	}

	return nml
}

// Lowercase creates a lowercase normalizer
func Lowercase() Normalizer {
	n := NewDefaultNormalizer()
	n.lower = true
	n.strip = false

	return n
}
