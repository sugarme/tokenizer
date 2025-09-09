package pretokenizer

import (
	// "log"
	"strings"

	"github.com/gengzongjie/tokenizer"
	"github.com/gengzongjie/tokenizer/normalizer"
)

// PrependScheme defines how the meta character should be prepended
type PrependScheme int

const (
	// Never specifies that the space should not be prepended
	Never PrependScheme = iota
	// First specifies that the scheme should be prepended only once, on the first split
	First
	// Always specifies that the scheme should always be prepended
	Always
)

// Metaspace constructs a Metaspace struct.
// It replaces all the whitespaces by the provided meta character
// and then splits on this character.
type Metaspace struct {
	Replacement    string
	PrependScheme  PrependScheme
	AddPrefixSpace bool // Deprecated: use PrependScheme instead
	StrRep         string
}

func NewMetaspace(replacement string, addPrefixSpace bool) *Metaspace {
	// Convert the boolean to PrependScheme for backward compatibility
	scheme := Never
	if addPrefixSpace {
		scheme = Always
	}

	return &Metaspace{
		Replacement:    replacement,
		PrependScheme:  scheme,
		AddPrefixSpace: addPrefixSpace, // Keep for backward compatibility
		StrRep:         replacement,
	}
}

// NewMetaspaceWithScheme creates a new Metaspace with a specific prepend scheme
func NewMetaspaceWithScheme(replacement string, scheme PrependScheme) *Metaspace {
	// Set AddPrefixSpace for backward compatibility
	addPrefixSpace := scheme != Never

	return &Metaspace{
		Replacement:    replacement,
		PrependScheme:  scheme,
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
	splitFn := func(idx int, normalized *normalizer.NormalizedString) []tokenizer.SplitIdx {
		var splits []normalizer.NormalizedString
		whitespace := normalizer.NewRegexpPattern(`\s`)

		normalized = normalized.Replace(whitespace, m.StrRep)

		// log.Printf("normalized: %+v\n", normalized)

		// Apply the prepend scheme
		switch m.PrependScheme {
		case Always:
			// Always prepend the replacement if it's not already there
			if !strings.HasPrefix(normalized.GetNormalized(), m.Replacement) {
				normalized = normalized.Prepend(m.StrRep)
			}
		case First:
			// Only prepend on the first split (idx == 0)
			if idx == 0 && !strings.HasPrefix(normalized.GetNormalized(), m.Replacement) {
				normalized = normalized.Prepend(m.StrRep)
			}
		case Never:
			// Never prepend
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
