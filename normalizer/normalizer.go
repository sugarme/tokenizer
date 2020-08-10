package normalizer

import (
	"golang.org/x/text/unicode/norm"
)

type Normalizer interface {
	Normalize(normalized NormalizedString) (NormalizedString, error)
}

type normalizer struct {
	Normalizer Normalizer
}

func newNormalizer() normalizer {
	return normalizer{
		Normalizer: NewDefaultNormalizer(),
	}
}

func (n normalizer) Normalize(normalized NormalizedString) (NormalizedString, error) {

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
		o(&nml)
	}

	return nml
}
