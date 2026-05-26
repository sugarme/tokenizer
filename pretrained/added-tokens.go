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

// CreateAddedTokensWithIds preserves the explicit IDs from tokenizer.json
// instead of letting AddedVocabulary recompute them. This is required for
// tokenizers with compacted vocabularies where added token IDs are not
// simply model.GetVocabSize() + offset.
func CreateAddedTokensWithIds(data []tokenizer.TokenConfig) []tokenizer.AddedTokenWithId {
	result := make([]tokenizer.AddedTokenWithId, 0, len(data))
	for _, d := range data {
		tok := tokenizer.DefaultAddedToken()
		tok.Content = d.Content
		tok.LStrip = d.Lstrip
		tok.Normalized = d.Normalized
		tok.RStrip = d.Rstrip
		tok.SingleWord = d.SingleWord

		result = append(result, tokenizer.AddedTokenWithId{
			Id:      int(d.Id),
			Special: d.Special,
			Token:   tok,
		})
	}
	return result
}
