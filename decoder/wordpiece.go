package decoder

import (
	"fmt"
	"strings"
)

// WordPieceDecoder takes care of decoding a list of wordpiece tokens
// back into a readable string.
type WordPieceDecoder struct {
	// The prefix to be used for continuing subwords
	prefix string
	// Whether to cleanup some tokenization artifacts (spaces before punctuation, ...)
	cleanup bool
}

// NewBpeDecoder creates a new BpeDecoder
func NewWordPieceDecoder(prefix string, cleanup bool) *WordPieceDecoder {
	return &WordPieceDecoder{
		prefix:  prefix,
		cleanup: cleanup,
	}
}

// DefaultBpeDecoder create a new BpeDecoder with default suffix (`</w>`)
func DefaultWordpieceDecoder() *WordPieceDecoder {
	return &WordPieceDecoder{
		prefix:  "##",
		cleanup: true,
	}
}

func (wd *WordPieceDecoder) Decode(tokens []string) string {
	output := strings.Join(tokens, " ")
	output = strings.ReplaceAll(output, fmt.Sprintf(" %v", wd.prefix), "")
	if wd.cleanup {
		output = strings.ReplaceAll(output, " .", ".")
		output = strings.ReplaceAll(output, " ?", "?")
		output = strings.ReplaceAll(output, " !", "!")
		output = strings.ReplaceAll(output, " ,", ",")
		output = strings.ReplaceAll(output, " ' ", "'")
		output = strings.ReplaceAll(output, " n't", "n't")
		output = strings.ReplaceAll(output, " 'm", "'m")
		output = strings.ReplaceAll(output, " do not", " don't")
		output = strings.ReplaceAll(output, " 's", "'s")
		output = strings.ReplaceAll(output, " 've", "'ve")
		output = strings.ReplaceAll(output, " 're", "'re")
	}

	return output
}
