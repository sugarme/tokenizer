package pretrained

import (
	"github.com/sugarme/tokenizer"
)

func CreateAddedTokens(data []tokenizer.TokenConfig) (specialToks, toks []tokenizer.AddedToken) {
	for _, d := range data {
		tok := tokenizer.DefaultAddedToken()
		tok.Content = d.Content
		tok.LStrip = d.Lstrip
		tok.Normalized = d.Normalized
		tok.RStrip = d.Rstrip
		tok.SingleWord = d.SingleWord

		if d.Special {
			specialToks = append(specialToks, tok)
		} else {
			toks = append(toks, tok)
		}
	}

	return specialToks, toks
}
