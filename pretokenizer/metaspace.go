package pretokenizer

import (
	// "log"
	"strings"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
)

// Metaspace constructs a Metaspace struct.
// It replaces all the whitespaces by the provided meta character
// and then splits on this character.
type Metaspace struct {
	Replacement    string
	AddPrefixSpace bool
	StrRep         string
}

func NewMetaspace(replacement string, addPrefixSpace bool) *Metaspace {
	return &Metaspace{
		Replacement:    replacement,
		AddPrefixSpace: addPrefixSpace,
		StrRep:         replacement,
	}
}

func (m *Metaspace) GetReplacement() string {
	return m.Replacement
}

func (m *Metaspace) SetReplacement(replacement string) {
	m.Replacement = replacement
	m.StrRep = replacement
}

func DefaultMetaspace() *Metaspace {
	return NewMetaspace("▁", true) // NOTE. `▁`  != `_`
}

// PreTokenize implements PreTokenizer interface
func (m *Metaspace) PreTokenize(pretokenized *tokenizer.PreTokenizedString) (*tokenizer.PreTokenizedString, error) {
	// func(int, *normalizer.NormalizedString) []SplitIdx
	splitFn := func(_ int, normalized *normalizer.NormalizedString) []tokenizer.SplitIdx {
		var splits []normalizer.NormalizedString
		whitespace := normalizer.NewRegexpPattern(`\s`)

		normalized = normalized.Replace(whitespace, m.StrRep)

		// log.Printf("normalized: %+v\n", normalized)

		if m.AddPrefixSpace && !strings.HasPrefix(normalized.GetNormalized(), m.Replacement) {
			normalized = normalized.Prepend(m.StrRep)
		}

		replacement := normalizer.NewRegexpPattern(m.Replacement)
		splits = normalized.Split(replacement, normalizer.MergedWithNextBehavior)

		// log.Printf("splits: %+v\n", splits)

		var splitIdxs []tokenizer.SplitIdx
		for _, s := range splits {
			normalized := s
			splitIdx := tokenizer.SplitIdx{Normalized: &normalized, Tokens: nil}
			splitIdxs = append(splitIdxs, splitIdx)
		}

		return splitIdxs
	}

	return pretokenized.Split(splitFn), nil
}

// DecodeChain implements Decoder interface.
func (m *Metaspace) DecodeChain(tokens []string) []string {
	var toks []string
	for i, token := range tokens {
		chars := strings.Split(token, "")

		var newChars []string
		for _, c := range chars {
			if c == m.Replacement {
				if i == 0 && m.AddPrefixSpace {
					// nil
				} else {
					newChars = append(newChars, " ")
				}
			} else {
				newChars = append(newChars, c)
			}
		}

		newTok := strings.Join(newChars, "")
		toks = append(toks, newTok)
	}

	return toks
}

func (m *Metaspace) Decode(tokens []string) string {
	out := m.DecodeChain(tokens)

	return strings.Join(out, "")
}
