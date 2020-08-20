package normalizer

import (
	"regexp"

	"github.com/sugarme/tokenizer"
)

// Pattern is used to split a NormalizedString
type Pattern interface {
	// FindMatches slices the given string in a list of pattern match positions, with
	// a boolean indicating whether this is a match or not.
	//
	// NOTE. This method *must* cover the whole string in its outputs, with
	// contiguous ordered slices.
	FindMatches(inside string) (retVal []OffsetsMatch)
}

// OfsetsMatch contains a combination of Offsets position
// and a boolean indicates whether this is a match or not.
type OffsetsMatch struct {
	offsets tokenizer.Offsets
	match   bool
}

// String is a wrapper of primitive string
// so that it can implement `Pattern` interface
type StringPattern struct {
	string
}

func NewStringPattern(s string) String {
	return StringPattern{s}
}

func (s StringPattern) FindMatches(inside string) (retVal OffsetsMatch) {
	// If we try to find the matches with an empty string, just don't match anything
	if len(s.string) == 0 {
		return OffsetsMatch{
			offsets: tokenizer.Offsets{0, len([]rune(s.string))},
			match:   false,
		}
	}

	quoted := regexp.QuoteMeta(s.string)
	re := regexp.MustCompile(quoted)
	matches := re.FindAllStringIndex(inside, -1)
	var subs []OffsetsMatch
	for i, m := range matches {
		// 1. First unmatched substring if first match is not start at 0
		if i == 0 && m[0] > 0 {
			first := OffsetMatch{
				offsets: Offsets{0, m[0]},
				match:   false,
			}
			subs = append(subs, first)
		}

		// 2. Matched sub itself
		matched := OffsetsMatch{
			offsets: Offsets{m[0], m[1]},
			match:   true,
		}
		subs = append(subs, matched)

		// 3. Unmatched sub between matched sub if any
		if i+1 < len(matches) {
			next = matches[i+1]

			if next[0] > m[1] {
				between := OffsetsMatch{
					offsets: Offsets{m[1], next[0]},
					match:   false,
				}
				subs = append(subs, between)
			}
		}
	}

	return subs
}

// RunePattern is a wrapper of primitive rune
// so that it can implement `Pattern` interface
type RunePattern struct {
	rune
}

func NewRunePattern(r rune) RunePattern {
	return RunePattern{r}
}

func (r RunePattern) FindMatches(inside string) (retVal []OffsetsMatch) {

	// TODO. implement
	return
}
