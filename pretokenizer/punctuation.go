package pretokenizer

import (
	"unicode"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
)

// bpunc is the BERT extension of the Punctuation character range
var bpunc = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0021, 0x002f, 1}, // 33-47
		{0x003a, 0x0040, 1}, // 58-64
		{0x005b, 0x0060, 1}, // 91-96
		{0x007b, 0x007e, 1}, // 123-126
	},
	LatinOffset: 4, // All less than 0x00FF
}

// IsPunctuation checks whether rune c is a BERT punctuation character
/*
 * fn is_punc(x: char) -> bool {
 *     char::is_ascii_punctuation(&x) || x.is_punctuation()
 * } */
func isPunctuation(c rune) bool {
	return unicode.In(c, bpunc, unicode.P)
}

type Punctuation struct {
	Behavior normalizer.SplitDelimiterBehavior
}

func DefaultSplit() normalizer.SplitDelimiterBehavior {
	return normalizer.IsolatedBehavior
}

func NewPunctuation(behavior normalizer.SplitDelimiterBehavior) *Punctuation {
	return &Punctuation{behavior}
}

func DefaultPunctuation() *Punctuation {
	behavior := DefaultSplit()
	return NewPunctuation(behavior)
}

// Implement tokenizer.PreTokenizer

var _ tokenizer.PreTokenizer = new(Punctuation)

/*
impl PreTokenizer for Punctuation {
    fn pre_tokenize(&self, pretokenized: &mut PreTokenizedString) -> Result<()> {
        pretokenized.split(|_, s| s.split(is_punc, self.behavior))
    }
}
*/
// PreTokenize implements tokenizer.PreTokenizer.
func (p *Punctuation) PreTokenize(pretokenized *tokenizer.PreTokenizedString) (*tokenizer.PreTokenizedString, error) {
	pretok := pretokenized.Split(func(noop int, sub *normalizer.NormalizedString) []tokenizer.SplitIdx {
		isPunc := normalizer.NewFnPattern(isPunctuation)
		splits := sub.Split(isPunc, p.Behavior)

		var splitIdxs []tokenizer.SplitIdx
		for _, s := range splits {
			normalized := s
			splitIdx := tokenizer.SplitIdx{Normalized: &normalized, Tokens: nil}
			splitIdxs = append(splitIdxs, splitIdx)
		}

		return splitIdxs
	})

	return pretok, nil
}
