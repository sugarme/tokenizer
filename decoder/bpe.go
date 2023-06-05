package decoder

import (
	"strings"

	"github.com/sugarme/tokenizer"
)

// Allows decoding Original BPE by joining all the tokens and then replacing
// the suffix used to identify end-of-words by whitespaces
type BpeDecoder struct {
	*DecoderBase

	suffix string
}

// NewBpeDecoder creates a new BpeDecoder
func NewBpeDecoder(suffix string) *BpeDecoder {
	base := new(DecoderBase)
	d := &BpeDecoder{
		DecoderBase: base,
		suffix:      suffix,
	}

	d.DecoderBase.Decoder = interface{}(d).(tokenizer.Decoder)

	return d
}

// DefaultBpeDecoder create a new BpeDecoder with default suffix (`</w>`)
func DefaultBpeDecoder() *BpeDecoder {
	return &BpeDecoder{suffix: "</w>"}
}

/*
func (bd *BpeDecoder) Decode(tokens []string) string {
	output := strings.Join(tokens, "")
	output = strings.ReplaceAll(output, bd.suffix, " ")

	return output
}
*/

func (bd *BpeDecoder) DecodeChain(tokens []string) []string {
	var toks []string
	for i, token := range tokens {
		replacement := " "
		if i == len(tokens)-1 {
			replacement = ""
		}

		tok := strings.ReplaceAll(token, bd.suffix, replacement)
		toks = append(toks, tok)
	}

	return toks
}
