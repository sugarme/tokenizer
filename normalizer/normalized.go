package normalizer

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/sugarme/tokenizer/util"
)

// RangeType is a enum like representing
// which string (original or normalized) then range
// indexes on.
type IndexOn int

const (
	OriginalTarget = iota
	NormalizedTarget
)

// Range is a slice of indexes on either normalized string or original string
// It is INCLUSIVE start and INCLUSIVE end
type Range struct {
	start   int
	end     int
	indexOn IndexOn
}

func NewRange(start int, end int, indexOn IndexOn) (retVal Range) {
	return Range{
		start:   start, // inclusive
		end:     end,   // exclusive
		indexOn: indexOn,
	}
}

func (r Range) Start() (retVal int) {
	return r.start
}

func (r Range) End() (retVal int) {
	return r.end
}

// IntoFullRange convert the current range to cover the case where the
// original provided range was out of bound.
// maxLen is maximal len of string in `chars` (runes)
func (r Range) intoFullRange(maxLen int) (retVal Range) {
	// case: start out of bound (including `None` value)
	if r.start == -1 || r.start > maxLen {
		r.start = 0
	}

	//  case: end out of bound
	if r.end > maxLen {
		r.end = maxLen
	}

	// case: end is None value
	if r.end == -1 {
		// TODO: should we just accept as `None`?
		r.end = maxLen
	}
	return r
}

// NormalizedString keeps both versions of an input string and
// provides methods to access them
type NormalizedString struct {
	original   string
	normalized string
	alignments []Alignment
}

// Alignment maps normalized string to original one using `rune` (Unicode code point)
type Alignment struct {
	Start int
	End   int
}

/* // Normalized is wrapper for a `NormalizedString` and provides
 * // methods to access it.
 * type Normalized struct {
 *   normalizedString NormalizedString
 * } */

// NewNormalizedFrom creates a Normalized instance from string input
func NewNormalizedFrom(s string) (retVal NormalizedString) {
	var alignments []Alignment

	// Break down string to slice of runes
	for i := range []rune(s) {
		alignments = append(alignments, Alignment{
			Start: i,
			End:   i + 1,
		})
	}

	return NormalizedString{
		original:   s,
		normalized: s,
		alignments: alignments,
	}
}

// GetNormalized returns the Normalized struct
func (n NormalizedString) GetNormalized() string {
	return n.normalized
}

// GetOriginal return the original string
func (n NormalizedString) GetOriginal() string {
	return n.original
}

// Alignments returns alignments mapping `chars` from
// normalized string to original string
func (n NormalizedString) Alignments() (retVal []Alignment) {
	return n.alignments
}

// ConvertOffset converts the given offsets range from referential to the the
// other one (`Original` to `Normalized` and vice versa)
func (n NormalizedString) ConvertOffset(inputRange Range) (retVal Range) {

	switch inputRange.indexOn {
	case OriginalTarget:
		var start, end int = -1, -1
		target := inputRange.intoFullRange(n.LenOriginal())

		// If we target before the start of normalized string
		if target.end <= n.alignments[0].Start {
			return NewRange(0, 0, NormalizedTarget)
		}

		// If we target after the end of normalized string
		if target.start > n.alignments[len(n.alignments)-1].End {
			length := n.Len()
			return NewRange(length, length, NormalizedTarget)
		}

		// Otherwise, let find the range
		var alignments []Alignment
		for _, a := range n.alignments {
			if target.end >= a.End {
				alignments = append(alignments, a)
			}
		}
		for i, a := range alignments {
			if a.Start >= target.start && start == -1 {
				// Here we want to keep the first char in the normalized string
				// that is on or *after* the target start.
				start = i
			}
			if a.End <= target.end {
				end = i + 1
			}
		}
		// If we didn't find the start, let's use the end of the normalized string
		if start == -1 {
			start = n.Len()
		}
		// The end must be greater or equal to start, and might be None otherwise
		if end < start {
			end = -1 // None value
		}

		return NewRange(start, end, NormalizedTarget)
	case NormalizedTarget:
		// If we target 0..0 on an empty normalized string, we want to return the
		// entire original one
		fullRange := inputRange.intoFullRange(n.Len())
		if len(n.alignments) == 0 && fullRange.start == 0 && fullRange.end == 0 {
			return NewRange(0, n.LenOriginal(), OriginalTarget)
		} else {
			alignments := n.alignments[fullRange.start:fullRange.end]
			if len(alignments) == 0 {
				return NewRange(-1, -1, OriginalTarget) // range = `None` value
			} else {
				start := alignments[0].Start
				end := alignments[len(alignments)-1].End
				return NewRange(start, end, OriginalTarget)
			}
		}
	}

	return
}

/*
 *
 *
 * func (n NormalizedString) ConvertOffset(inputRange Range) (retVal Range) {
 *   lastAlign := n.alignments[len(n.alignments)-1]
 *   r := inputRange.intoFullRange(lastAlign.End)
 *   start := 0
 *   end := 0
 *   switch inputRange.indexOn {
 *   case OriginalTarget: // convert to normalized
 *     // get all alignments in range
 *     var alignments []Alignment
 *     for _, a := range n.alignments {
 *       if r.end >= a.End {
 *         alignments = append(alignments, a)
 *       }
 *     }
 *     for i, a := range alignments {
 *       if a.Start <= r.start {
 *         start = i
 *       }
 *       if a.End <= r.end {
 *         end = i + 1
 *       }
 *     }
 *
 *     retVal = Range{
 *       start:   start,
 *       end:     end,
 *       indexOn: NormalizedTarget,
 *     }
 *
 *   case NormalizedTarget: // convert to original
 *     alignments := n.alignments[r.start:r.end]
 *     if len(alignments) == 0 {
 *       // log.Fatalf("Cannot convert to original offsets. No alignments are in range.\n")
 *       // NOTE. r.start == r.end -> just switch indexOn and return
 *       r.indexOn = OriginalTarget
 *       return r
 *     }
 *
 *     start = alignments[0].Start
 *     end = alignments[len(alignments)-1].End
 *
 *     retVal = Range{
 *       start:   start,
 *       end:     end,
 *       indexOn: OriginalTarget,
 *     }
 *
 *   default:
 *     log.Fatalf("Invalid 'indexOn' type: %v\n", r.indexOn)
 *   }
 *
 *   return retVal
 * }
 *  */

// RangeOf returns a substring of the given string by indexing chars instead of bytes
// It will return empty string if input range is out of bound
func RangeOf(s string, r []int) (retVal string) {
	runes := []rune(s)
	sLen := len(runes)
	var start, end int
	if len(r) == 0 {
		start = 0
	} else {
		start = r[0]
	}

	if r[len(r)-1] > sLen {
		end = sLen
	} else {
		end = r[len(r)-1]
	}

	// if out of range, return 'empty' string
	if start >= sLen || end > sLen || start >= end {
		return ""
	}

	slicedRunes := runes[start:end]
	return string(slicedRunes)
}

// Range returns a substring of the NORMALIZED string (indexing on character not byte)
func (n NormalizedString) Range(r Range) (retVal string) {
	var nRange Range

	// Convert to NormalizedRange if r is OriginalRange
	switch r.indexOn {
	case OriginalTarget:
		nRange = n.ConvertOffset(r)
	case NormalizedTarget:
		nRange = r
	default:
		log.Fatalf("Invalid Range type: %v\n", r.indexOn)
	}

	return RangeOf(n.GetNormalized(), util.MakeRange(nRange.start, nRange.end))
}

// RangeOriginal returns substring of ORIGINAL string
func (n NormalizedString) RangeOriginal(r Range) string {
	var oRange Range
	switch r.indexOn {
	case NormalizedTarget:
		oRange = n.ConvertOffset(r)
	case OriginalTarget:
		oRange = r
	default:
		log.Fatalf("Invalid Range type: %v\n", r.indexOn)
	}

	rSlice := util.MakeRange(oRange.start, oRange.end)

	return RangeOf(n.GetOriginal(), rSlice)
}

// SliceBytes returns a new NormalizedString that contains only the specified
// range, indexing on BYTES.
//
// Any range that splits a UTF-8 `char` will return `None`.
//
// If we want a slice of the `NormalizedString` based on a `Range::Normalized``,
// the original part of the `NormalizedString` will contain any "additional"
// content on the right, and also on the left. The left will be included
// only if we are retrieving the very beginning of the string, since there
// is no previous part. The right is always included, up to what's covered
// by the next part of the normalized string.  This is important to be able
// to build a new `NormalizedString` from multiple contiguous slices
func (n NormalizedString) SliceBytes(inputRange Range) (retVal *NormalizedString) {
	var (
		r      Range
		s      string
		target IndexOn
	)

	switch inputRange.indexOn {
	case OriginalTarget:
		target = OriginalTarget
		r = inputRange.intoFullRange(len(n.original)) // len in bytes
		s = n.original
	case NormalizedTarget:
		target = NormalizedTarget
		r = inputRange.intoFullRange(len(n.normalized))
		s = n.normalized
	default:
		log.Fatalf("Invalid Range type: %v\n", r.indexOn)
	}

	type runeIdx struct {
		rune    rune
		byteIdx int
		runeIdx int
	}

	var (
		start, end  int = -1, -1 // `None` value
		currRuneIdx int = 0
		chars       []runeIdx
	)

	if r.start == 0 && r.end == 0 {
		start = 0
		end = 0
	}

	// NOTE. range over string is special. It iterates index on byte and unicode
	// code point (rune). See more: https://blog.golang.org/strings
	for i, char := range s {
		// select indexes of bytes in range
		if i < r.end && i >= r.start {
			chars = append(chars, runeIdx{char, i, currRuneIdx})
		}
		currRuneIdx++
	}

	for _, char := range chars {
		if char.byteIdx == r.start {
			start = char.runeIdx
		}
		if char.byteIdx+len(string(char.rune)) == r.end {
			end = char.runeIdx + 1
		}
	}

	if start == -1 || end == -1 { // splitting on `char`
		return nil
	}

	outRange := NewRange(start, end, target)

	return n.Slice(outRange)
}

// Slice returns a new NormalizedString that contains only specified range, indexing
// on `char`
//
// If we want a slice of the `NormalizedString` based on a `Range::Normalized``,
// the original part of the `NormalizedString` will contain any "additional"
// content on the right, and also on the left. The left will be included
// only if we are retrieving the very beginning of the string, since there
// is no previous part. The right is always included, up to what's covered
// by the next part of the normalized string.  This is important to be able
// to build a new `NormalizedString` from multiple contiguous slices
func (n NormalizedString) Slice(inputRange Range) (retVal *NormalizedString) {

	lenOriginal := n.LenOriginal()
	lenNormalized := n.Len()

	var rNormalized, rOriginal Range

	// rNormalized
	switch inputRange.indexOn {
	case OriginalTarget:
		rNormalized = n.ConvertOffset(inputRange)
	case NormalizedTarget:
		rNormalized = inputRange.intoFullRange(lenNormalized)
	}

	// rOriginal
	switch inputRange.indexOn {
	case OriginalTarget:
		rOriginal = inputRange.intoFullRange(lenOriginal)
	case NormalizedTarget:
		endRange := n.ConvertOffset(NewRange(rNormalized.end, n.LenOriginal(), NormalizedTarget))

		rOriginal = n.ConvertOffset(inputRange)

		// If we take the very beginning of the normalized string, we should take
		// all the beginning of the original too
		if rNormalized.start == 0 && rOriginal.start != 0 {
			rOriginal.start = 0
		}

		// If there is a void between the `end` char we target and the next one, we
		// want to include everything in-between from the original string
		if endRange.start >= 0 && endRange.end >= 0 {
			if endRange.start > rOriginal.end {
				rOriginal.end = endRange.start
			}
		}

		// If we target the end of the normalized but the original is longer
		if rNormalized.end == len(n.alignments) && lenOriginal > rOriginal.end {
			rOriginal.end = lenOriginal
		}
	}

	// Shift the alignments according to the part of the original string
	// that will be kept.
	alignmentShift := rOriginal.start

	newAlignments := n.alignments[rNormalized.start:rNormalized.end]
	var shiftAligments []Alignment

	for _, a := range newAlignments {
		shiftAligments = append(shiftAligments, Alignment{
			Start: a.Start - alignmentShift,
			End:   a.End - alignmentShift,
		})
	}

	retVal = &NormalizedString{
		original:   RangeOf(n.GetOriginal(), util.MakeRange(rOriginal.start, rOriginal.end)),
		normalized: RangeOf(n.GetNormalized(), util.MakeRange(rNormalized.start, rNormalized.end)),
		alignments: shiftAligments,
	}

	return retVal
}

type ChangeMap struct {
	RuneVal string
	Changes int
}

// Transform applies transformations to the current normalized version, updating the current
// alignments with the new ones.
// This method expect an Iterator yielding each rune of the new normalized string
// with a `change` interger size equals to:
//   - `1` if this is a new rune
//   - `-N` if the char is right before N removed runes
//   - `0` if this rune represents the old one (even if changed)
// Since it is possible that the normalized string doesn't include some of the `characters` (runes) at
// the beginning of the original one, we need an `initial_offset` which represents the number
// of removed runes at the very beginning.
//
// `change` should never be more than `1`. If multiple runes are added, each of
// them has a `change` of `1`, but more doesn't make any sense.
// We treat any value above `1` as `1`.
//
// E.g. string `élégant`
// Before NFD():  [{233 0} {108 1} {233 2} {103 3} {97 4} {110 5} {116 6}]
// After NFD(): 	[{101 0} {769 1} {108 2} {101 3} {769 4} {103 5} {97 6} {110 7} {116 8}]
// New Alignments:
// {0, 1},
// {0, 1},
// {1, 2},
// {2, 3},
// {2, 3},
// {3, 4},
// {4, 5},
// {5, 6},
// {6, 7},
func (n NormalizedString) Transform(m []ChangeMap, initialOffset int) (retVal NormalizedString) {

	offset := -initialOffset
	var (
		alignments []Alignment
		runeVals   []string
	)

	for i, item := range m {
		// Positive offset means there're added `chars`. This offset needed to be
		// removed from current index to get the previous id.
		idx := i - offset
		offset += item.Changes
		var align Alignment
		if item.Changes > 0 {
			if idx < 1 {
				align = Alignment{Start: 0, End: 0}
			} else { // newly inserted `char`. Hence, use aligment from previous one
				align = n.alignments[idx-1]
			}
		} else {
			align = n.alignments[idx]
		}

		alignments = append(alignments, align)
		runeVals = append(runeVals, item.RuneVal)
	}

	n.alignments = alignments
	n.normalized = strings.Join(runeVals, "")

	return n
}

func (n NormalizedString) NFD() (retVal NormalizedString) {

	s := n.normalized
	var (
		changeMap []ChangeMap
		it        norm.Iter
	)
	// Create slice of (char, changes) to map changing
	// if added (inserted) rune, changes = 1; `-N` if char
	// right before N removed chars
	// changes = 0 if this represents the old one (even if changed)

	// Iterating over string and apply tranformer (NFD). One character at a time
	// A `character` is defined as:
	// - a sequence of runes that starts with a starter,
	// - a rune that does not modify or combine backwards with any other rune,
	// - followed by possibly empty sequence of non-starters, that is, runes that do (typically accents).
	// We will iterate over string and apply transformer to each char
	// If a char composes of one rune, there no changes
	// If more than one rune, first is no change, the rest is 1 changes
	it.InitString(norm.NFD, s)
	for !it.Done() {
		runes := []rune(string(it.Next()))

		for i, r := range runes {

			switch i := i; {
			case i == 0:
				changeMap = append(changeMap, ChangeMap{
					RuneVal: string(r),
					Changes: 0,
				})
			case i > 0:
				changeMap = append(changeMap, ChangeMap{
					RuneVal: string(r),
					Changes: 1,
				})
			}
		}

	}

	return n.Transform(changeMap, 0)
}

func (n NormalizedString) NFC() (retVal NormalizedString) {

	var (
		changeMap []ChangeMap
		it        norm.Iter
	)

	s := n.normalized

	isNFC := norm.Form.IsNormalString(norm.NFC, s)
	if isNFC {
		return
	}

	it.InitString(norm.NFD, s)

	for !it.Done() {
		runes := []rune(string(it.Next()))

		if len(runes) == 1 {
			changeMap = append(changeMap, ChangeMap{
				RuneVal: string(runes),
				Changes: 0,
			})
		} else if len(runes) > 1 {
			changeMap = append(changeMap, ChangeMap{
				RuneVal: string(runes),
				Changes: -1,
			})
		}
	}

	return n.Transform(changeMap, 0)
}

func (n NormalizedString) NFKD() (retVal NormalizedString) {

	s := n.normalized
	isNFKD := norm.Form.IsNormalString(norm.NFKD, s)
	if isNFKD {
		return
	}

	var (
		changeMap []ChangeMap
		it        norm.Iter
	)

	it.InitString(norm.NFKD, s)
	for !it.Done() {
		runes := []rune(string(it.Next()))

		for i, r := range runes {
			switch i := i; {
			case i == 0:
				changeMap = append(changeMap, ChangeMap{
					RuneVal: string(r),
					Changes: 0,
				})
			case i > 0:
				changeMap = append(changeMap, ChangeMap{
					RuneVal: string(r),
					Changes: 1,
				})
			}
		}
	}

	return n.Transform(changeMap, 0)
}

func (n NormalizedString) NFKC() (retVal NormalizedString) {

	var (
		changeMap []ChangeMap
		it        norm.Iter
	)

	s := n.normalized

	isNFKC := norm.Form.IsNormalString(norm.NFKC, s)

	if isNFKC {
		return
	}

	it.InitString(norm.NFKD, n.normalized)

	for !it.Done() {
		runes := []rune(string(it.Next()))

		if len(runes) == 1 {
			changeMap = append(changeMap, ChangeMap{
				RuneVal: string(runes),
				Changes: 0,
			})
		} else if len(runes) > 1 {
			changeMap = append(changeMap, ChangeMap{
				RuneVal: string(runes),
				Changes: -1,
			})
		}
	}

	return n.Transform(changeMap, 0)
}

// Filter applies filtering on NormalizedString
func (n NormalizedString) Filter(fr rune) (retVal NormalizedString) {
	s := n.normalized
	var changeMap []ChangeMap

	var oRunes []rune

	var it norm.Iter
	it.InitString(norm.NFC, s)

	for !it.Done() {
		runes := []rune(string(it.Next()))

		oRunes = append(oRunes, runes...)

	}

	revRunes := make([]rune, 0)
	for i := len(oRunes) - 1; i >= 0; i-- {
		revRunes = append(revRunes, oRunes[i])
	}

	var removed int = 0
	for _, r := range revRunes {
		if r == fr {
			removed += 1
		} else {
			if removed > 0 {
				changeMap = append(changeMap, ChangeMap{
					RuneVal: string(r),
					Changes: -removed,
				})
				removed = 0
			} else if removed == 0 {
				changeMap = append(changeMap, ChangeMap{
					RuneVal: string(r),
					Changes: 0,
				})
			}
		}
	}

	// Flip back changeMap
	var unrevMap []ChangeMap
	for i := len(changeMap) - 1; i >= 0; i-- {
		unrevMap = append(unrevMap, changeMap[i])
	}

	return n.Transform(unrevMap, removed)
}

// Prepend adds given string to the begining of NormalizedString
func (n NormalizedString) Prepend(s string) (retVal NormalizedString) {
	newString := fmt.Sprintf("%s%s", s, n.GetNormalized())
	var newAligments []Alignment
	for i := 0; i < len([]rune(s)); i++ {
		newAligments = append(newAligments, Alignment{Start: 0, End: 0})
	}
	newAligments = append(newAligments, n.alignments...)
	n.normalized = newString
	n.alignments = newAligments

	return n
}

// Append adds given string to the end of NormalizedString
func (n NormalizedString) Append(s string) (retVal NormalizedString) {
	newString := fmt.Sprintf("%s%s", n.GetNormalized(), s)
	var newAligments []Alignment
	lastAlign := n.alignments[len(n.alignments)-1]
	for i := 0; i < len([]rune(s)); i++ {
		newAligments = append(newAligments, Alignment{Start: lastAlign.End, End: lastAlign.End})
	}
	newAligments = append(n.alignments, newAligments...)
	n.normalized = newString
	n.alignments = newAligments

	return n
}

// NormFn is a convenient function type for applying
// on each `char` of normalized string
type NormFn func(rune) rune

// Map maps and applies function to each `char` of normalized string
func (n NormalizedString) Map(nfn NormFn) (retVal NormalizedString) {
	s := n.normalized
	var runes []rune
	for _, r := range []rune(s) {
		runes = append(runes, nfn(r))
	}

	n.normalized = string(runes)

	return n
}

// ForEach applies function on each `char` of normalized string
// Similar to Map???
func (n NormalizedString) ForEach(nfn NormFn) (retVal NormalizedString) {
	s := n.normalized
	var runes []rune
	for _, r := range []rune(s) {
		runes = append(runes, nfn(r))
	}
	n.normalized = string(runes)

	return n
}

// RemoveAccents removes all Unicode Mn group (M non-spacing)
func (n NormalizedString) RemoveAccents() (retVal NormalizedString) {

	s := n.normalized
	b := make([]byte, len(s))

	tf := transform.Chain(transform.RemoveFunc(isMn))

	_, _, err := tf.Transform(b, []byte(s), true)
	if err != nil {
		log.Fatal(err)
	}

	n.normalized = string(b)

	return n
}

// Lowercase transforms string to lowercase
func (n NormalizedString) Lowercase() (retVal NormalizedString) {
	n.normalized = strings.ToLower(n.normalized)

	return n
}

// Uppercase transforms string to uppercase
func (n NormalizedString) Uppercase() (retVal NormalizedString) {
	n.normalized = strings.ToUpper(n.normalized)

	return n
}

// SplitOff truncates string with the range [at, len).
// remaining string will contain the range [0, at).
// The provided `at` indexes on `char` not bytes.
func (n NormalizedString) SplitOff(at int) (retVal NormalizedString) {
	if at < 0 {
		log.Fatal("Split off point must be a positive interger number.")
	}
	s := n.normalized
	if at > len([]rune(s)) {
		n = NewNormalizedFrom("")
	}

	var (
		it       norm.Iter
		runeVals []string
		aligns   []Alignment
	)

	// Split normalized string
	it.InitString(norm.NFC, s)
	for !it.Done() {
		runeVal := string(it.Next())
		runeVals = append(runeVals, runeVal)
	}

	// Alignments
	remainVals := runeVals[0:at]
	for i := range remainVals {
		aligns = append(aligns, Alignment{
			Start: i,
			End:   i + 1,
		})
	}
	n.normalized = strings.Join(remainVals, "")
	n.alignments = aligns

	// Split original string
	originalAt := aligns[len(aligns)].End // changes of last alignment

	var oRuneVals []string
	it.InitString(norm.NFC, n.original)
	for !it.Done() {
		runeVal := string(it.Next())
		oRuneVals = append(oRuneVals, runeVal)
	}

	remainORuneVals := oRuneVals[0:originalAt]
	n.original = strings.Join(remainORuneVals, "")

	return n
}

// MergeWith merges an input string with existing one
func (n NormalizedString) MergeWith(other NormalizedString) (retVal NormalizedString) {
	len := n.Len()
	n.original = strings.Join([]string{n.original, other.original}, "")
	n.normalized = strings.Join([]string{n.normalized, other.normalized}, "")

	var ajustedAligns []Alignment
	for _, a := range other.alignments {
		new := Alignment{
			Start: a.Start + len,
			End:   a.End + len,
		}

		ajustedAligns = append(ajustedAligns, new)
	}

	n.alignments = append(n.alignments, ajustedAligns...)

	return n
}

// LStrip removes leading spaces
func (n NormalizedString) LStrip() (retVal NormalizedString) {
	return n.lrstrip(true, false)
}

// RStrip removes trailing spaces
func (n NormalizedString) RStrip() (retVal NormalizedString) {
	return n.lrstrip(false, true)
}

// Strip remove leading and trailing spaces
func (n NormalizedString) Strip() (retVal NormalizedString) {
	return n.lrstrip(true, true)
}

// lrstrip - Private func to help with exposed strip funcs
func (n NormalizedString) lrstrip(left, right bool) (retVal NormalizedString) {
	var (
		leadingSpaces  int = 0
		trailingSpaces int = 0
		s              string
		changeMap      []ChangeMap
	)

	s = n.normalized

	runes := []rune(s)

	if left {
		for _, r := range runes {
			if !unicode.IsSpace(r) {
				break
			}

			leadingSpaces += 1
		}
	}

	if right {
		for i := len(runes) - 1; i >= 0; i-- {
			if !unicode.IsSpace(runes[i]) {
				break
			}

			trailingSpaces += 1
		}
	}

	if leadingSpaces > 0 || trailingSpaces > 0 {
		for i, r := range runes {
			if i < leadingSpaces || i >= (len(runes)-trailingSpaces) {
				continue
			} else if i == len(runes)-trailingSpaces-1 {
				changeMap = append(changeMap, ChangeMap{
					RuneVal: string(r),
					Changes: -(trailingSpaces),
				})
			} else {
				changeMap = append(changeMap, ChangeMap{
					RuneVal: string(r),
					Changes: 0,
				})
			}
		}

		return n.Transform(changeMap, leadingSpaces)
	}

	return n
}

// Len returns length (number of runes) of normalized string
func (n NormalizedString) Len() int {
	runes := []rune(n.normalized)
	return len(runes)
}

// LenOriginal returns the length of Original string in `char` (rune)
func (n NormalizedString) LenOriginal() int {
	return len([]rune(n.GetOriginal()))
}

// IsEmpty returns whether the normalized string is empty
func (n NormalizedString) IsEmpty() bool {
	return n.Len() == 0
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}
