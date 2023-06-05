package decoder

import (
	"fmt"
	"strings"

	"github.com/sugarme/tokenizer"
)

// WordPieceDecoder takes care of decoding a list of wordpiece tokens
// back into a readable string.
type WordPieceDecoder struct {
	*DecoderBase
	// The prefix to be used for continuing subwords
	prefix string
	// Whether to cleanup some tokenization artifacts (spaces before punctuation, ...)
	cleanup bool
}

// NewBpeDecoder creates a new BpeDecoder
func NewWordPieceDecoder(prefix string, cleanup bool) *WordPieceDecoder {
	base := new(DecoderBase)
	d := &WordPieceDecoder{
		DecoderBase: base,
		prefix:      prefix,
		cleanup:     cleanup,
	}

	d.DecoderBase.Decoder = interface{}(d).(tokenizer.Decoder)

	return d
}

// DefaultBpeDecoder create a new BpeDecoder with default suffix (`</w>`)
func DefaultWordpieceDecoder() *WordPieceDecoder {
	return &WordPieceDecoder{
		prefix:  "##",
		cleanup: true,
	}
}

/*
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
*/

func (wd *WordPieceDecoder) Cleanup(tok string) string {
	output := tok
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

	return output
}

func (wd *WordPieceDecoder) DecodeChain(tokens []string) []string {
	var toks []string
	for i, token := range tokens {
		var tok string
		if i != 0 {
			if strings.HasPrefix(token, wd.prefix) {
				tok = strings.Replace(token, wd.prefix, "", 1)
			} else {
				tok = fmt.Sprintf(" %s", token)
			}
		}

		if wd.cleanup {
			tok = wd.Cleanup(tok)
		}

		toks = append(toks, tok)
	}

	return toks
}
