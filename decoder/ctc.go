package decoder

import (
	"strings"

	"github.com/sugarme/tokenizer"
)

type CTC struct {
	*DecoderBase

	PadToken           string // the pad token used by CTC to delimit a new token
	WordDelimiterToken string // the word delimiter token. It will be replace by a `<space>`
	Cleanup            bool   // whether to cleanup some tokenization artifacts, mainly spaces before punctuation and some abbreviated english forms
}

func NewCTC(padToken string, wordDelimiterToken string, cleanup bool) *CTC {
	base := new(DecoderBase)

	d := &CTC{
		DecoderBase:        base,
		PadToken:           padToken,
		WordDelimiterToken: wordDelimiterToken,
		Cleanup:            cleanup,
	}

	d.DecoderBase.Decoder = interface{}(d).(tokenizer.Decoder)

	return d
}

func DefaultCTC() *CTC {
	base := new(DecoderBase)
	return &CTC{
		DecoderBase:        base,
		PadToken:           "<pad>",
		WordDelimiterToken: "|",
		Cleanup:            true,
	}
}

// dedup deduplicates consecutive elements.
func dedup(tokens []string) []string {
	var toks []string

	var previous string

	for _, tok := range tokens {
		if tok != previous {
			toks = append(toks, tok)
			previous = tok
		}
	}

	return toks
}

func contains(items []string, item string) bool {
	for _, itm := range items {
		if itm == item {
			return true
		}
	}

	return false
}

func cleanup(tok string) string {
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

func (d *CTC) DecodeChain(tokens []string) []string {
	var toks []string

	uniqueTokens := dedup(tokens)
	for _, token := range uniqueTokens {
		replaced := strings.ReplaceAll(token, d.PadToken, "")
		if d.Cleanup {
			replaced = cleanup(replaced)
			replaced = strings.ReplaceAll(replaced, d.WordDelimiterToken, " ")
		}

		if len(replaced) > 0 {
			toks = append(toks, replaced)
		}
	}

	return toks
}
