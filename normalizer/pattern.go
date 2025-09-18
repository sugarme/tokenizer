package normalizer

import (
	"log"
	// "reflect"
	"regexp"

	"github.com/sugarme/tokenizer/util"
)

// Pattern is used to split a NormalizedString
type Pattern interface {
	// FindMatches slices the given string in a list of pattern match positions, with
	// a boolean indicating whether this is a match or not.
	//
	// NOTE. This method *must* cover the whole string in its outputs, with
	// contiguous ordered slices.
	FindMatches(inside string) []OffsetsMatch
}

// OfsetsMatch contains a combination of Offsets position
// and a boolean indicates whether this is a match or not.
type OffsetsMatch struct {
	Offsets []int // slice of 2 elements (start, end)
	Match   bool
}

// RunePattern is a wrapper of primitive rune
// so that it can implement `Pattern` interface
type RunePattern struct {
	rune
}

func NewRunePattern(r rune) *RunePattern {
	return &RunePattern{r}
}

// FindMaches implements Pattern interface for RunePattern
func (r *RunePattern) FindMatches(inside string) []OffsetsMatch {

	if len(inside) == 0 {
		return []OffsetsMatch{
			{
				Offsets: []int{0, 0},
				Match:   false,
			},
		}
	}

	var (
		subs        []OffsetsMatch
		prevStart   int  = 0
		hasPrevious bool = false
	)

	for byteIdx, char := range inside {
		if char == r.rune {
			nextIdx := byteIdx + len(string(char))
			// 1. Add previous unmatched if any
			if hasPrevious {
				prev := OffsetsMatch{Offsets: []int{prevStart, byteIdx}, Match: false}
				subs = append(subs, prev)
			}

			// 2. Add matched one
			matched := OffsetsMatch{
				Offsets: []int{byteIdx, nextIdx},
				Match:   char == r.rune,
			}
			subs = append(subs, matched)

			// 3. update prevStart
			prevStart = nextIdx
			hasPrevious = false
		} else {
			hasPrevious = true
		}
	}

	// 4. Add last unmatched if any
	if hasPrevious {
		prev := OffsetsMatch{Offsets: []int{prevStart, len(inside)}}
		subs = append(subs, prev)
	}

	return subs
}

// String is a wrapper of primitive string
// so that it can implement `Pattern` interface
type StringPattern struct {
	string
}

func NewStringPattern(s string) *StringPattern {
	return &StringPattern{s}
}

func (s *StringPattern) FindMatches(inside string) []OffsetsMatch {
	// If we try to find the matches with an empty string, just don't match anything
	if s.string == "" {
		return []OffsetsMatch{
			{
				Offsets: []int{0, len(inside)},
				Match:   false,
			},
		}
	}

	quoted := regexp.QuoteMeta(s.string)

	re := regexp.MustCompile(quoted)

	return findMatches(re, inside)
}

func findMatches(re *regexp.Regexp, inside string) []OffsetsMatch {

	matches := re.FindAllStringIndex(inside, -1)

	// 0. If no matches, just return
	if len(matches) == 0 {
		return []OffsetsMatch{
			{
				Offsets: []int{0, len(inside)},
				Match:   false,
			},
		}
	}

	var (
		currIndex int = 0
		subs      []OffsetsMatch
	)

	// 1. Sub before matched if any
	if matches[0][0] > 0 {
		offsets := []int{0, matches[0][0]}
		first := OffsetsMatch{
			Offsets: offsets,
			Match:   false,
		}
		subs = append(subs, first)
		currIndex += matches[0][0]
	}

	for i, m := range matches {

		// 2. matched itself
		sub := OffsetsMatch{
			Offsets: m,
			Match:   true,
		}
		subs = append(subs, sub)
		currIndex += m[1] - m[0]

		// 3. unmatched in between if any (will not if 2 continuous matched)
		if i+1 < len(matches) {
			next := matches[i+1]
			current := matches[i]
			if current[1] != next[0] { // not continuous matches
				offsets := []int{m[1], next[0]}
				between := OffsetsMatch{
					Offsets: offsets,
					Match:   false,
				}
				subs = append(subs, between)
				currIndex += offsets[1] - offsets[0]
			}
		}
	}

	// 4. Last unmatched if any
	if currIndex < len(inside) {
		offsets := []int{currIndex, len(inside)}
		last := OffsetsMatch{
			Offsets: offsets,
			Match:   false,
		}

		subs = append(subs, last)
	}

	return subs
}

type RegexpPattern struct {
	re *regexp.Regexp
}

func NewRegexpPattern(s string) *RegexpPattern {
	re := regexp.MustCompile(s)
	return &RegexpPattern{
		re: re,
	}
}

// FindMatches implements Pattern interface for RegexpPattern
func (rp *RegexpPattern) FindMatches(inside string) []OffsetsMatch {
	if len(inside) == 0 {
		return []OffsetsMatch{
			{
				Offsets: []int{0, 0},
				Match:   false,
			},
		}
	}

	return findMatches(rp.re, inside)
}

// PatternFn is a func type to apply pattern
type PatternFn func(rune) bool

type FnPattern struct {
	fn PatternFn
}

func NewFnPattern(fn PatternFn) *FnPattern {
	return &FnPattern{fn}
}

// FindMatches implements Pattern interface for FnPattern
func (fp *FnPattern) FindMatches(inside string) []OffsetsMatch {
	if len(inside) == 0 {
		return []OffsetsMatch{
			{
				Offsets: []int{0, 0},
				Match:   false,
			},
		}
	}

	var (
		subs        []OffsetsMatch
		prevStart   int  = 0
		hasPrevious bool = false
	)

	for byteIdx, char := range inside {
		if fp.fn(char) {
			nextIdx := byteIdx + len(string(char))
			// 1. Add previous unmatched if any
			if hasPrevious {
				prev := OffsetsMatch{Offsets: []int{prevStart, byteIdx}, Match: false}
				subs = append(subs, prev)
			}

			// 2. Add matched one
			matched := OffsetsMatch{
				Offsets: []int{byteIdx, nextIdx},
				Match:   true,
			}
			subs = append(subs, matched)

			// 3. update prevStart
			prevStart = nextIdx
			hasPrevious = false
		} else {
			hasPrevious = true
		}
	}

	// 4. Add last unmatched if any
	if hasPrevious {
		prev := OffsetsMatch{Offsets: []int{prevStart, len(inside)}}
		subs = append(subs, prev)
	}

	return subs

}

// Invert the `is_match` flags for the wrapped Pattern. This is usefull
// for example when we use a regex that matches words instead of a delimiter,
// and we want to match the delimiter.
type Invert struct {
	Pattern Pattern
}

func NewInvertPattern(p Pattern) *Invert {
	return &Invert{p}
}

// FindMatches implement Pattern interface for Invert
func (i *Invert) FindMatches(inside string) []OffsetsMatch {
	var matches []OffsetsMatch
	typ := util.GetType(i.Pattern)
	switch typ {
	case "*StringPattern":
		matches = i.Pattern.(*StringPattern).FindMatches(inside)
	case "*RunePattern":
		matches = i.Pattern.(*RunePattern).FindMatches(inside)
	case "*FnPattern":
		matches = i.Pattern.(*FnPattern).FindMatches(inside)
	case "*RegexpPattern":
		matches = i.Pattern.(*RegexpPattern).FindMatches(inside)

	default:
		log.Fatalf("Unsupported type - %q\n", typ)
	}

	return invert(matches)
}

func invert(matches []OffsetsMatch) (retVal []OffsetsMatch) {
	var res []OffsetsMatch
	for _, m := range matches {
		m.Match = !m.Match
		res = append(res, m)
	}

	return res
}
