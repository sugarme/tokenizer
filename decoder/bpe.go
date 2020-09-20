package decoder

import (
	"strings"
)

// Allows decoding Original BPE by joining all the tokens and then replacing
// the suffix used to identify end-of-words by whitespaces
type BpeDecoder struct {
	suffix string
}

// NewBpeDecoder creates a new BpeDecoder
func NewBpeDecoder(suffix string) *BpeDecoder {
	return &BpeDecoder{suffix: suffix}
}

// DefaultBpeDecoder create a new BpeDecoder with default suffix (`</w>`)
func DefaultBpeDecoder() *BpeDecoder {
	return &BpeDecoder{suffix: "</w>"}
}

func (bd *BpeDecoder) Decode(tokens []string) string {
	output := strings.Join(tokens, "")
	output = strings.ReplaceAll(output, bd.suffix, " ")

	return output
}
