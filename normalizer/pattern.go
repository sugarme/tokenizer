package normalizer

import (
	"log"
	"reflect"
	"regexp"
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
	Offsets []int // slice of 2 elements (start, end)
	Match   bool
}

// RunePattern is a wrapper of primitive rune
// so that it can implement `Pattern` interface
type RunePattern struct {
	rune
}

func NewRunePattern(r rune) RunePattern {
	return RunePattern{r}
}

// FindMaches implements Pattern interface for RunePattern
func (r RunePattern) FindMatches(inside string) (retVal []OffsetsMatch) {

	if len(inside) == 0 {
		return []OffsetsMatch{
			{
				Offsets: []int{0, 0},
				Match:   false,
			},
		}
	}

	var subs []OffsetsMatch
	var prevStart int = 0
	var hasPrevious bool = false

	for byteIdx, char := range inside {
		if char == r.rune {
			nextIdx := byteIdx + len(string(char))
			// 1. Add previous unmatched if any
			if hasPrevious {
				prev := OffsetsMatch{Offsets: []int{prevStart, byteIdx}}
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

func NewStringPattern(s string) StringPattern {
	return StringPattern{s}
}

func (s StringPattern) FindMatches(inside string) (retVal []OffsetsMatch) {
	// If we try to find the matches with an empty string, just don't match anything
	if len(s.string) == 0 {
		return []OffsetsMatch{
			{
				Offsets: []int{0, len([]rune(s.string))},
				Match:   false,
			},
		}
	}

	quoted := regexp.QuoteMeta(s.string)
	re := regexp.MustCompile(quoted)

	return findMatches(re, inside)
}

func findMatches(re *regexp.Regexp, inside string) (retVal []OffsetsMatch) {
	matches := re.FindAllStringIndex(inside, -1)
	var subs []OffsetsMatch
	for i, m := range matches {
		// 1. First unmatched substring if first match is not start at 0
		if i == 0 && m[0] > 0 {
			first := OffsetsMatch{
				Offsets: []int{0, m[0]},
				Match:   false,
			}
			subs = append(subs, first)
		}

		// 2. Matched sub itself
		matched := OffsetsMatch{
			Offsets: []int{m[0], m[1]},
			Match:   true,
		}
		subs = append(subs, matched)

		// 3. Unmatched sub between matched sub if any
		if i+1 < len(matches) {
			next := matches[i+1]

			if next[0] > m[1] {
				between := OffsetsMatch{
					Offsets: []int{m[1], next[0]},
					Match:   false,
				}
				subs = append(subs, between)
			}
		}
	}

	return subs
}

type RegexpPattern struct {
	re *regexp.Regexp
}

func NewRegexpPattern(s string) RegexpPattern {
	re := regexp.MustCompile(s)
	return RegexpPattern{
		re: re,
	}
}

// FindMatches implements Pattern interface for RegexpPattern
func (rp RegexpPattern) FindMatches(inside string) (retVal []OffsetsMatch) {
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

/*
 *
 * func (rp RegexpPattern) FindMatches(inside string) (retVal []OffsetsMatch) {
 *   if len(inside) == 0 {
 *     return []OffsetsMatch{
 *       {
 *         offsets: tokenizer.Offsets{Start: 0, End: 0},
 *         match:   false,
 *       },
 *     }
 *   }
 *
 *   // charIndices contains position of rune in string
 *   var charIndices [][]int
 *   var currIdx int = 0
 *   for charIdx := range []rune(inside) {
 *     charIndices = append(charIndices, []int{currIdx, charIdx})
 *     currIdx = charIdx
 *   }
 *   var (
 *     charIdx int = 0
 *     prev    int = 0
 *     splits  []OffsetsMatch
 *   )
 *
 *   matches := rp.re.FindAllStringIndex(inside, -1)
 *   for _, m := range matches {
 *     prevIdx := charIdx
 *     startIdx := charIdx
 *     for charIdx < len(charIndices) && charIndices[charIdx][0] < m[1] {
 *       if charIndices[charIdx][0] == m[0] {
 *         startIdx = charIdx
 *       }
 *       charIdx++
 *     }
 *
 *     if prev != m[0] {
 *       splits = append(splits, OffsetsMatch{
 *         offsets: tokenizer.Offsets{Start: prevIdx, End: startIdx},
 *         match:   false,
 *       })
 *     }
 *
 *     splits = append(splits, OffsetsMatch{
 *       offsets: tokenizer.Offsets{Start: startIdx, End: charIdx},
 *       match:   true,
 *     })
 *     prev = m[1]
 *   }
 *
 *   if prev != len(inside) {
 *     splits = append(splits, OffsetsMatch{
 *       offsets: tokenizer.Offsets{Start: charIdx, End: len(charIndices)},
 *       match:   false,
 *     })
 *   }
 *
 *   return splits
 * } */

// PatternFn is a func type to apply pattern
type PatternFn func(rune) bool

type FnPattern struct {
	fn PatternFn
}

func NewFnPattern(fn PatternFn) FnPattern {
	return FnPattern{fn}
}

// FindMatches implements Pattern interface for FnPattern
func (fp FnPattern) FindMatches(inside string) (retVal []OffsetsMatch) {
	if len(inside) == 0 {
		return []OffsetsMatch{
			{
				Offsets: []int{0, 0},
				Match:   false,
			},
		}
	}

	var (
		currIdx int = 0
		matches []OffsetsMatch
	)

	for byteIdx, char := range inside {
		m := OffsetsMatch{
			Offsets: []int{currIdx, byteIdx},
			Match:   fp.fn(char),
		}
		matches = append(matches, m)
	}

	return
}

// Invert the `is_match` flags for the wrapped Pattern. This is usefull
// for example when we use a regex that matches words instead of a delimiter,
// and we want to match the delimiter.
type Invert struct {
	Pattern Pattern
}

// FindMatches implement Pattern interface for Invert
func (i Invert) FindMatches(inside string) (retVal []OffsetsMatch) {

	var matches []OffsetsMatch

	switch reflect.TypeOf(i.Pattern).Name() {
	case "StringPattern":
		matches = i.Pattern.(StringPattern).FindMatches(inside)
	case "RunePattern":
		matches = i.Pattern.(RunePattern).FindMatches(inside)
	case "FnPattern":
		matches = i.Pattern.(FnPattern).FindMatches(inside)

	default:
		log.Fatalf("Unsupported type - %v\n", reflect.TypeOf(i.Pattern).Name())
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
