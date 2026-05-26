package normalizer

import (
	"regexp"
	"sync"
	"unicode/utf8"

	"github.com/dlclark/regexp2"
)

var runeToBytePool = sync.Pool{
	New: func() any {
		buf := make([]int, 0, 256)
		return &buf
	},
}

var regexp2MatchPool = sync.Pool{
	New: func() any {
		buf := make([][2]int, 0, 16)
		return &buf
	},
}

func getRuneToByteScratch(minCap int) *[]int {
	buf := runeToBytePool.Get().(*[]int)
	if cap(*buf) < minCap {
		*buf = make([]int, 0, minCap)
	} else {
		*buf = (*buf)[:0]
	}
	return buf
}

func putRuneToByteScratch(buf *[]int) {
	const maxRetainedCap = 1 << 20
	if cap(*buf) > maxRetainedCap {
		*buf = make([]int, 0, 256)
	} else {
		*buf = (*buf)[:0]
	}
	runeToBytePool.Put(buf)
}

func getRegexp2MatchScratch(minCap int) *[][2]int {
	buf := regexp2MatchPool.Get().(*[][2]int)
	if cap(*buf) < minCap {
		*buf = make([][2]int, 0, minCap)
	} else {
		*buf = (*buf)[:0]
	}
	return buf
}

func putRegexp2MatchScratch(buf *[][2]int) {
	const maxRetainedCap = 1 << 14
	if cap(*buf) > maxRetainedCap {
		*buf = make([][2]int, 0, 16)
	} else {
		*buf = (*buf)[:0]
	}
	regexp2MatchPool.Put(buf)
}

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
			nextIdx := byteIdx + utf8.RuneLen(char)
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
	s  string
	re *regexp.Regexp
}

func NewStringPattern(s string) *StringPattern {
	quoted := regexp.QuoteMeta(s)
	return &StringPattern{s: s, re: regexp.MustCompile(quoted)}
}

func (s *StringPattern) FindMatches(inside string) []OffsetsMatch {
	// If we try to find the matches with an empty string, just don't match anything
	if s.s == "" {
		return []OffsetsMatch{
			{
				Offsets: []int{0, len(inside)},
				Match:   false,
			},
		}
	}

	return findMatches(s.re, inside)
}

func buildMatchesFromIndices(matches [][]int, inside string) []OffsetsMatch {

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
		subs          = make([]OffsetsMatch, 0, len(matches)*2+1)
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

func findMatches(re *regexp.Regexp, inside string) []OffsetsMatch {
	matches := re.FindAllStringIndex(inside, -1)
	return buildMatchesFromIndices(matches, inside)
}

// collectRegexp2Indices collects all match [start, end) rune offsets from re
// into dst, reusing the slice to avoid per-call allocation.
func collectRegexp2Indices(re *regexp2.Regexp, s string, dst [][2]int) [][2]int {
	dst = dst[:0]
	m, err := re.FindStringMatch(s)
	for err == nil && m != nil {
		dst = append(dst, [2]int{m.Index, m.Index + m.Length})
		m, err = re.FindNextMatch(m)
	}
	return dst
}

func findMatchesRegexp2(re *regexp2.Regexp, inside string) []OffsetsMatch {
	asciiOnly := true
	for i := 0; i < len(inside); i++ {
		if inside[i] >= utf8.RuneSelf {
			asciiOnly = false
			break
		}
	}

	var (
		runeToByte    []int
		runeToByteBuf *[]int
	)
	if !asciiOnly {
		runeToByteBuf = getRuneToByteScratch(utf8.RuneCountInString(inside) + 1)
		runeToByte = *runeToByteBuf
		for byteIdx := range inside {
			runeToByte = append(runeToByte, byteIdx)
		}
		runeToByte = append(runeToByte, len(inside))
		*runeToByteBuf = runeToByte
	}

	toByte := func(runeIdx int) int {
		if runeIdx < 0 {
			return 0
		}
		if asciiOnly {
			if runeIdx > len(inside) {
				return len(inside)
			}
			return runeIdx
		}
		if runeIdx >= len(runeToByte) {
			return len(inside)
		}
		return runeToByte[runeIdx]
	}

	runeMatchesBuf := getRegexp2MatchScratch(8)
	runeMatches := collectRegexp2Indices(re, inside, *runeMatchesBuf)
	if len(runeMatches) == 0 {
		putRegexp2MatchScratch(runeMatchesBuf)
		if runeToByteBuf != nil {
			putRuneToByteScratch(runeToByteBuf)
		}
		return []OffsetsMatch{{Offsets: []int{0, len(inside)}, Match: false}}
	}

	matches := make([]OffsetsMatch, 0, len(runeMatches)*2+1)
	curr := 0
	for _, rm := range runeMatches {
		start := toByte(rm[0])
		end := toByte(rm[1])
		if start > curr {
			matches = append(matches, OffsetsMatch{Offsets: []int{curr, start}, Match: false})
		}
		matches = append(matches, OffsetsMatch{Offsets: []int{start, end}, Match: true})
		curr = end
	}
	*runeMatchesBuf = runeMatches
	putRegexp2MatchScratch(runeMatchesBuf)

	if runeToByteBuf != nil {
		putRuneToByteScratch(runeToByteBuf)
	}

	if curr < len(inside) {
		matches = append(matches, OffsetsMatch{Offsets: []int{curr, len(inside)}, Match: false})
	}

	return matches
}

// RegexpPattern uses github.com/dlclark/regexp2 for regex matching,
// which supports lookaheads, lookbehinds, and other features not
// available in Go's standard regexp package. This enables compatibility
// with tokenizer patterns used by GPT-4, Qwen, Llama 3, and other
// modern models that rely on .NET/PCRE-style regex syntax.
type RegexpPattern struct {
	re     *regexp2.Regexp
	source string
}

// NewRegexpPattern compiles the given pattern using regexp2 and returns
// a RegexpPattern that implements the Pattern interface. Panics if the
// pattern cannot be compiled.
func NewRegexpPattern(s string) *RegexpPattern {
	re, err := regexp2.Compile(s, 0)
	if err != nil {
		panic(err)
	}
	return &RegexpPattern{re: re, source: s}
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

	return findMatchesRegexp2(rp.re, inside)
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
			nextIdx := byteIdx + utf8.RuneLen(char)
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
	return invert(i.Pattern.FindMatches(inside))
}

func invert(matches []OffsetsMatch) (retVal []OffsetsMatch) {
	res := make([]OffsetsMatch, len(matches))
	for i, m := range matches {
		m.Match = !m.Match
		res[i] = m
	}

	return res
}
