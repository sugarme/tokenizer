package normalizer

import (
	"bytes"
	"log"
	"reflect"
	"strings"
	"unicode"

	"github.com/sugarme/tokenizer/util"
	slice "github.com/sugarme/tokenizer/util/slice"

	// "golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// SplitDelimiterBehavior is a enum-like type . It defines the expected behavior
// for the delimiter of a Split Pattern
// When splitting on `'-'` for example, with input `the-final--countdown`:
//   - RemovedBehavior => `[ "the", "final", "countdown" ]`
//   - IsolatedBehavior => `[ "the", "-", "final", "-", "-", "countdown" ]`
//   - MergedWithPreviousBehavior => `[ "the-", "final-", "-", "countdown" ]`
//   - MergedWithNextBehavior => `[ "the", "-final", "-", "-countdown" ]`
//   - Contiguous => `[ "the", "-", "final", "--", "countdown" ]`
type SplitDelimiterBehavior int

const (
	RemovedBehavior = iota
	IsolatedBehavior
	MergedWithPreviousBehavior
	MergedWithNextBehavior
	ContiguousBehavior
)

type OffsetsRemove struct {
	Offsets      []int
	ShouldRemove bool
}

// RangeType is a enum like representing
// which string (original or normalized) then range
// indexes on.
type IndexOn int

const (
	OriginalTarget = iota
	NormalizedTarget
)

// Range is a slice of indexes on either normalized string or original string
// It is INCLUSIVE start and EXCLUSIVE end
type Range struct {
	start   int
	end     int
	indexOn IndexOn
}

func NewRange(start int, end int, indexOn IndexOn) (retVal *Range) {
	return &Range{
		start:   start, // inclusive
		end:     end,   // exclusive
		indexOn: indexOn,
	}
}

func (r *Range) Start() (retVal int) {
	return r.start
}

func (r *Range) End() (retVal int) {
	return r.end
}

// Len returns the length of the current Range if not unbounded
func (r *Range) Len() int {
	end := r.end // can be unbounded if == -1

	if r.start < 0 { // unbounded
		return end
	} else {
		return end - r.start
	}
}

// Values returns range values (start, end)
func (r *Range) Values() []int {
	if r == nil {
		return nil
	}

	return []int{r.start, r.end}
}

// IndexOn returns the target where range index on
func (r *Range) On() IndexOn {
	return r.indexOn
}

// IntoFullRange convert the current range to cover the case where the
// original provided range was out of bound.
// maxLen is maximal len of string in `chars` (runes)
func (r *Range) IntoFullRange(maxLen int) (retVal *Range) {
	// case: start out of bound (including `None` value)
	// if r.start == -1 || r.start > maxLen {
	if r.start == -1 {
		r.start = 0
	}

	/*
	 *   if r.start > r.end {
	 *     r.start = r.end - 1
	 *   }
	 *  */
	//  case: end out of bound
	if r.end > maxLen {
		r.end = maxLen
	}

	/*
	 *   // case: end is None value
	 *   if r.end == -1 {
	 *     // TODO: should we just accept as `None`?
	 *     r.end = maxLen
	 *   }
	 *  */

	return r
}

// A `NormalizedString` takes care of processing an "original" string to modify
// it and obtain a "normalized" string. It keeps both version of the string,
// alignments information between both and provides an interface to retrieve
// ranges of each string, using offsets from any of them.
//
// It is possible to retrieve a part of the original string, by indexing it with
// offsets from the normalized one, and the other way around too. It is also
// possible to convert offsets from one referential to the other one easily.
type NormalizedString struct {
	// The original version of the string, before any modification
	original string
	// The normalized version of the string, after all modifications
	normalized string
	// Mapping from normalized string to original one: (start, end) for each
	// byte of the normalized string
	alignments [][]int
	// Mapping from original string to normalized one: (start, end) for each
	// byte of the original string
	alignmentsOriginal [][]int
	// If this NormalizedString is a slice of a bigger one, we keep the track
	// of the missing part, so that we can still give offsets from this original
	// string.
	originalShift int
}

// NewNormalizedFrom creates a Normalized instance from string input
func NewNormalizedFrom(s string) (retVal *NormalizedString) {
	/*
	 *   // NOTE. Really need to make a deep copy, otherwise
	 *   // `aligments` and `alignmentsOriginal` updates each other later on!
	 *   alignmentsOriginal := make([][]int, len(alignments))
	 *   copy(alignmentsOriginal, alignments)
	 *  */
	return &NormalizedString{
		original:           s,
		normalized:         s,
		alignments:         createAligns(s),
		alignmentsOriginal: createAligns(s),
		originalShift:      0,
	}
}

// createALigns creates alignments from input string.
// NOTE:It is used in `NewNormalizedFrom` to create 2 slices
// (alignments and alignmentsOriginal) without sharing data
// (data elements are in different memory locations).
func createAligns(s string) [][]int {
	var alignments [][]int
	currIdx := 0
	for i, r := range []rune(s) {
		charLen := len([]byte(string(r)))
		var align []int
		if i == 0 {
			align = []int{0, charLen}
		} else {
			align = []int{currIdx, currIdx + charLen}
		}
		for byteIdx := 0; byteIdx < charLen; byteIdx++ {
			alignments = append(alignments, align)
		}
		currIdx += charLen
	}

	return alignments
}

func NewNormalizedString(original, normalized string, alignments, alignmentsOriginal [][]int, originalShift int) *NormalizedString {
	return &NormalizedString{
		original:           original,
		normalized:         normalized,
		alignments:         alignments,
		alignmentsOriginal: alignmentsOriginal,
		originalShift:      originalShift,
	}
}

// GetNormalized returns the Normalized struct
func (n *NormalizedString) GetNormalized() string {
	return n.normalized
}

// GetOriginal return the original string
func (n *NormalizedString) GetOriginal() string {
	return n.original
}

// Alignments returns alignments mapping normalized string to original string
func (n *NormalizedString) Alignments() (retVal [][]int) {
	return n.alignments
}

// AlignmentsOriginal returns original alignments mapping to original string
func (n *NormalizedString) AlignmentsOriginal() (retVal [][]int) {
	return n.alignmentsOriginal
}

// OffsetsOriginal returns the original offsets
func (n *NormalizedString) OffsetsOriginal() []int {
	return []int{n.originalShift, n.originalShift + n.LenOriginal()}
}

// Shift returns original shift
func (n *NormalizedString) Shift() int {
	return n.originalShift
}

// ConvertOffsets converts the given offsets range from one referential to the other one:
// `Original => Normalized` or `Normalized => Original`
//
// Returns `nil` when targeting something that is outside range
func (n *NormalizedString) ConvertOffset(inputRange *Range) (retVal *Range) {
	lenOriginal := n.LenOriginal()
	lenNormalized := n.Len()
	var (
		isOriginal bool
		target     *Range
		indexOn    IndexOn
		alignments [][]int
	)

	switch inputRange.indexOn {
	case OriginalTarget:
		isOriginal = true
		indexOn = NormalizedTarget
		target = inputRange.IntoFullRange(lenOriginal)
		alignments = n.alignmentsOriginal
	case NormalizedTarget:
		isOriginal = false
		indexOn = OriginalTarget
		target = inputRange.IntoFullRange(lenNormalized)
		alignments = n.alignments
	}

	// If we target an empty range, let's return the same
	if target.start == target.end {
		return target
	}

	// If we target 0..0 on an empty string, we want to expand to the entire equivalent
	if isOriginal && len(n.alignmentsOriginal) == 0 && reflect.DeepEqual(target.Values(), []int{0, 0}) {
		return NewRange(0, lenNormalized, indexOn)
	}

	if !isOriginal && len(n.alignments) == 0 && reflect.DeepEqual(target.Values, []int{0, 0}) {
		return NewRange(0, lenOriginal, indexOn)
	}

	// Otherwise, just convert them.
	// NOTE. if out of bound, return nil
	var lowerB, upperB int = 0, len(alignments)
	if target.start < lowerB || target.start > upperB || target.end < lowerB || target.end > upperB {
		return nil
	}

	newAlignments := alignments[target.start:target.end]
	newRange := expandAlignments(newAlignments)
	return NewRange(newRange[0], newRange[1], indexOn)
}

// Range returns a substring of the NORMALIZED string
func (n *NormalizedString) Range(r *Range) (retVal string) {

	bytes := []byte(n.normalized)
	switch r.indexOn {
	case OriginalTarget:
		nRange := n.ConvertOffset(r)
		retVal = string(bytes[nRange.start:nRange.end])
	case NormalizedTarget:
		nRange := r.IntoFullRange(n.Len())
		retVal = string(bytes[nRange.start:nRange.end])
	}

	return retVal
}

// RangeOriginal returns substring of ORIGINAL string
func (n *NormalizedString) RangeOriginal(r *Range) (retVal string) {

	bytes := []byte(n.original)
	switch r.indexOn {
	case NormalizedTarget:
		oRange := n.ConvertOffset(r)
		retVal = string(bytes[oRange.start:oRange.end])
	case OriginalTarget:
		oRange := r.IntoFullRange(n.LenOriginal())
		retVal = string(bytes[oRange.start:oRange.end])
	default:
		log.Fatalf("Invalid Range type: %v\n", r.indexOn)
	}

	return retVal
}

// validateRange validates the given range, to make sure it is on char boundaries
// if range on boundaries, return `nil`, otherwise, just return the input.
func (n *NormalizedString) validateRange(inputRange *Range) (retVal *Range) {
	var (
		r *Range
		s string
	)

	switch inputRange.indexOn {
	case OriginalTarget:
		r = inputRange.IntoFullRange(len(n.original)) // len in bytes
		s = n.original
	case NormalizedTarget:
		r = inputRange.IntoFullRange(len(n.normalized))
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

	// all good, just return the inputRange
	return inputRange
}

// Slice returns a slice of the current NormalizedString
// If the range is not on char boundaries, return `nil`
func (n *NormalizedString) Slice(inputRange *Range) (retVal *NormalizedString) {
	fullRange := n.validateRange(inputRange)
	if fullRange == nil {
		return nil
	}

	// 1. Find range on normalized (nRange) and original string (oRange)
	var nRange, oRange *Range
	switch fullRange.indexOn {
	case OriginalTarget:
		nRange = n.ConvertOffset(fullRange)
		oRange = fullRange
	case NormalizedTarget:
		nRange = fullRange
		oRange = n.ConvertOffset(fullRange)
	}

	// 2. `nShift` is normalized shift on original string
	nShift := oRange.start

	// 3. `oShift` is original shift on original string
	var oShift int
	if r := expandAlignments(n.alignmentsOriginal[0:nShift]); r == nil {
		oShift = 0
	} else {
		oShift = r[1]
	}

	var (
		sOriginal, sNormalized           string
		sAlignments, sAlignmentsOriginal [][]int
		sOriginalShift                   int
	)

	sOriginal = n.RangeOriginal(fullRange)
	sNormalized = n.Range(fullRange)
	for _, a := range n.alignments[nRange.start:nRange.end] {
		sAlignment := []int{a[0] - nShift, a[1] - nShift}
		sAlignments = append(sAlignments, sAlignment)
	}
	for _, a := range n.alignmentsOriginal[oRange.start:oRange.end] {
		sAlignmentOriginal := []int{a[0] - oShift, a[1] - oShift}
		sAlignmentsOriginal = append(sAlignmentsOriginal, sAlignmentOriginal)
	}
	sOriginalShift = n.originalShift + oRange.start

	return &NormalizedString{
		original:           sOriginal,
		normalized:         sNormalized,
		alignments:         sAlignments,
		alignmentsOriginal: sAlignmentsOriginal,
		originalShift:      sOriginalShift,
	}
}

type ChangeMap struct {
	RuneVal string
	Changes int
}

// This method expect an iterator yielding each `char` of the new normalized string
// with a `change` of int type equals to:
//   - `1` if this is a new char
//   - `-N` if the char is right before N removed chars
//   - `0` if the char is replacing the existing one
//
// Since it is possible that the normalized string doesn't include some of the characters at
// the beginning of the original one, we need an `initialOffset` which represents the number
// of removed chars at the very beginning.
func (n *NormalizedString) TransformRange(inputRange *Range, changeMap []ChangeMap, initialOffset int) (retVal *NormalizedString) {

	// fmt.Printf("normalized: %+v\n", n)
	// fmt.Printf("inputRange: %v\n", inputRange)

	// I. Update `original` alignments based on `initialOffset`.
	// If `initialOffset` > 0, there are some chars being replaced.

	// I.1. Determine range on `original` string and `normalized` string
	// based on `inputRange`
	var nRange *Range
	switch inputRange.indexOn {
	case NormalizedTarget:
		nRange = inputRange.IntoFullRange(n.Len())
	case OriginalTarget:
		nRange = n.ConvertOffset(inputRange)
	}

	// NOTE: temp fix out of range occured at bpe training example
	if nRange.end > len(n.alignments) {
		nRange.end = len(n.alignments)
	}

	// fmt.Printf("nRange: %v\n", nRange)

	// I.2. Get removed chars to a slice (`removedChars`); number of removed
	// bytes (`initialRemoved`); and update original alignments (`n.alignmensOriginal`)

	// log.Printf("===== TransformRange call (initial offset: %v) ====\n", initialOffset)
	// Retrieve the original characters that are being replaced. This let us
	// compute the change in byte sizes along the way.
	nBytes := []byte(n.normalized)[nRange.start:nRange.end]
	replacedNormalized := util.NewRuneIter(bytes.Runes(nBytes))
	/*
	 *   fmt.Printf("replaceddNoramlized len: %v\n", replacedNormalized.Len())
	 *   fmt.Printf("replacedNormalized: %+q\n", func(it *util.RuneIter) []string {
	 *     var chars []string
	 *     for {
	 *       r, ok := it.Next()
	 *       if !ok {
	 *         break
	 *       }
	 *       chars = append(chars, string(r))
	 *     }
	 *     return chars
	 *   }(replacedNormalized))
	 *  */
	replacedNormalized.Reset()

	// Handle the initial offset in the original alignment. All the characters
	// that were removed from the normalized one should have their width reduced
	// by the number of bytes we remove
	endShiftStart := nRange.end
	initialRemoved := 0
	if initialOffset > 0 {
		// log.Printf("=> Clearing alignment for %v chars\n", initialOffset)
		removedBytes := 0
		var removedChars []rune
		for i := 0; i < initialOffset; i++ {
			c, ok := replacedNormalized.Next()
			if !ok {
				// We want to panic here, because the NormalizedString is in
				// a bad state if this happens. We already modified a lot of things
				log.Fatalf("1. Expected to remove %v characters but couldn't find them ...\n", initialOffset)
			}
			removedBytes += len([]byte(string(c)))
			removedChars = append(removedChars, c)
		}
		initialRemoved = removedBytes
		offset := nRange.start
		oShift := 0
		// Then we remove all these chars, updating the alignments along the way
		for _, c := range removedChars {
			removedORange := expandAlignments(n.alignments[offset : offset+len(string(c))])
			offset += len(string(c))
			oShift += len(string(c))
			alignments := n.alignmentsOriginal[removedORange[0]:removedORange[1]]
			// log.Printf("Clearing alignments for char: %q - alignments: %+v\n", c, alignments)
			// 1. Get new alignments
			var newAlignments [][]int
			for _, offsets := range alignments {
				// At the very end we will apply the global shift to the remaining
				// original offsets. We should start after these to avoid doing it twice
				if offsets[1] > endShiftStart {
					endShiftStart = offsets[1]
				}

				offsets[1] = applySign(offsets[1], -(oShift))
				// Make sure the starting offset is always smaller or equal to the end
				if offsets[0] > offsets[1] {
					offsets[0] = offsets[1]
				}

				newAlignments = append(newAlignments, offsets)
			}
			// 2. Update alignmentOriginal with new alignments in `removedORange` range
			var aligns [][]int = append(n.alignmentsOriginal[:removedORange[0]], newAlignments...)
			aligns = append(aligns, n.alignmentsOriginal[removedORange[1]:]...)
			n.alignmentsOriginal = aligns

			// log.Printf("Cleared: %+v\n", alignments)
		}
	} // end `If` block

	// II.	Do the transformation based on input `changeMap`.
	// 1. 	Update alignments on normalized string(`n.alignments`) and original
	// 			string (`n.alignmentsOriginal`);
	// 2. 	collect transformed `chars` to a string
	// 			variable (`normalized`)

	// oShift is the shift to be applied to all original alignments along the way
	// NOTE. `oShift` and `offset` here are different from ones inside previous block
	oShift := -(initialRemoved)
	offset := initialRemoved + nRange.start
	var normalizedAlignments [][]int

	// log.Printf("Applying transformations...\n")
	var (
		normalizedRunes []rune
	)

	replacedNormalized.Reset() // make sure iterating from begining
	for _, item := range changeMap {
		/*
		 *     var changeType string
		 *     switch {
		 *     case item.Changes == 0:
		 *       changeType = "Replacing"
		 *     case item.Changes > 0:
		 *       changeType = "Adding"
		 *     case item.Changes < 0:
		 *       changeType = fmt.Sprintf("Replacing + Removing %v following chars", item.Changes)
		 *     default:
		 *       changeType = "Undefined"
		 *     }
		 *     // log.Printf("### %+q with size %v : %v with offset %v ###\n", item.RuneVal, len(item.RuneVal), changeType, offset)
		 *     // log.Printf("### '%v' with size %v : %v with offset %v ###\n", item.RuneVal, len(item.RuneVal), changeType, offset)
		 *     fmt.Printf("'%v' - changes: %v\n", item.RuneVal, item.Changes)
		 *  */
		idx := offset
		// fmt.Printf("idx: %v\n", idx)
		var align []int
		if item.Changes > 0 {
			if idx < 1 {
				align = []int{0, 0}
			} else {
				// This is a newly inserted character, so it shares the same alignment
				// as the previous one
				align = n.alignments[idx-1]
			}
		} else {
			align = n.alignments[idx]
		}

		// If we are replacing a character, find it and compute the change in size
		var replacedChar rune
		if item.Changes <= 0 {
			replacedChar, _ = replacedNormalized.Next()
		} else {
			replacedChar = 0 // nil rune value
		}
		replacedCharSize := 0
		if replacedChar > 0 {
			replacedCharSize = len([]byte(string(replacedChar)))
		}
		replacedCharSizeChange := len(item.RuneVal) - replacedCharSize
		if replacedChar > 0 {
			// log.Printf("Replacing char: '%v' - with a change in size: %v", string(replacedChar), replacedCharSizeChange)
		}

		// If we are removing some characters, find them too
		var nChanges int = 0
		if item.Changes < 0 {
			nChanges = -item.Changes
		}

		var (
			removedChars       []rune
			totalBytesToRemove int
		)
		for i := 0; i < nChanges; i++ {
			c, ok := replacedNormalized.Next()
			if !ok {
				// We want to panic here, because the NormalizedString is in
				// a bad state if this happens. We already modified a lot of things
				// log.Fatalf("2. Expected to remove %v characters but couldn't find them ...\n", nChanges)
			}
			removedChars = append(removedChars, c)
			totalBytesToRemove += len(string(c))
		}
		// fmt.Printf("removedChars: %+q\n", string(removedChars))
		// log.Printf("Total bytes to remove: %v\n", totalBytesToRemove)

		// If we are removing characters, there are two possible scenarios:
		//   1. We remove characters that are part of the original string (most likely)
		//   2. We remove characters that were previously added (with NFD for example)
		//      and are just part of the normalized string

		var removingFromOriginal, removingFromNormalized int = 0, 0
		if totalBytesToRemove > 0 {
			start := n.alignments[idx][1]
			end := n.alignments[idx+totalBytesToRemove][1]
			originalRange := util.MakeRange(start, end)
			// fmt.Printf("start: %v - end: %v; range: (%+v)\n", start, end, originalRange)
			removingFromOriginal = len(originalRange)
			removingFromNormalized = totalBytesToRemove - len(originalRange)

			// log.Printf("Bytes to remove from original alignments: %v\n", removingFromOriginal)
			// log.Printf("Bytes to remove from normalized alignments: %v\n", removingFromNormalized)
		}

		// Update the original alignments for the **current** `char`
		// fmt.Printf("align: %v\n", align)
		// fmt.Printf("current alignmentsOriginal: %+v\n", n.alignmentsOriginal)
		alignments := n.alignmentsOriginal[align[0]:align[1]]
		// fmt.Printf("alignments: %+v\n", alignments)
		if len(alignments) > 0 {
			// log.Printf("Updating original alignments: %+v\n", alignments)

			// Let's compute the actual change in size for the original alignment.
			// This value may be different than `replacedCharSizeChange` because the
			// char might not exist in the original string. It might have been added.
			originalRange := expandAlignments(alignments)

			replacedCharSizeChangeOriginal := len(item.RuneVal) - (originalRange[1] - originalRange[0])
			// log.Printf("Change in size for the current char: %v\n", replacedCharSizeChange)
			// log.Printf("Change in size in the original alignment for the current char: %v\n", replacedCharSizeChangeOriginal)

			var newAlignments [][]int
			for _, offsets := range alignments {
				var newAlign []int = make([]int, 2)

				// A new char is being added, we need to extend the current
				// alignment. No need to apply the shift in this case as it
				// should have been applied already (with a previous changes == 0).
				if item.Changes > 0 {
					newAlign[0] = offsets[0]
					newAlign[1] = offsets[1] + len(item.RuneVal)
				} else {
					// Otherwise we just apply the shift
					newAlign[0] = offsets[0]
					newAlign[1] = applySign(offsets[1], replacedCharSizeChange)
					newAlign[1] = applySign(newAlign[1], -(removingFromNormalized))

					// NOTE: We only apply oShift if we are:
					// - Removing characters
					// - Replacing one, that has the same size in both normalized and
					//   original alignments.
					// Otherwise it means we are modifying multiple times the same
					// original character. This happens when normalized characters
					// are added (not part of the original), and then replaced.
					if item.Changes < 0 || replacedCharSizeChange == replacedCharSizeChangeOriginal {
						newAlign[0] = applySign(newAlign[0], oShift)
						newAlign[1] = applySign(newAlign[1], oShift)
					}
				}

				newAlignments = append(newAlignments, newAlign)
			}

			// Now, update original alignments for current `char` with new alignments
			// fmt.Printf("newAlignments: %+v\n", newAlignments)
			var aligns [][]int
			if align[0] == 0 {
				aligns = append(newAlignments, n.alignmentsOriginal[align[1]:]...)
			} else {
				aligns = append(n.alignmentsOriginal[:align[0]], newAlignments...)
				aligns = append(aligns, n.alignmentsOriginal[align[1]:]...)
			}
			n.alignmentsOriginal = aligns
			// log.Printf("Updated to: %+v\n", newAlignments)
			// fmt.Printf("current n.alignmentsOriginal: %+v\n", n.alignmentsOriginal)
		}

		// If some were removed, we need to zero them out in the original alignments
		if removingFromOriginal > 0 {
			start := n.alignments[idx][1]
			end := n.alignments[idx+totalBytesToRemove][1]
			// They should use the original alignment of the current character
			newIdx := n.alignmentsOriginal[align[0]][1]
			alignments := n.alignmentsOriginal[start:end]
			var newAlignments [][]int
			if len(alignments) > 0 {
				// log.Printf("Removing original alignments: %v\n", alignments)
				for _, offsets := range alignments {
					offsets[0] = newIdx
					offsets[1] = newIdx
					newAlignments = append(newAlignments, offsets)
				}

				aligns := append(n.alignmentsOriginal[:start], newAlignments...)
				aligns = append(aligns, n.alignmentsOriginal[end:]...)
				n.alignmentsOriginal = aligns
			}
		}

		// Keep track of the changes for next offsets
		offset += replacedCharSize
		offset += totalBytesToRemove
		// For the original only the real modifications count
		oShift += replacedCharSizeChange
		oShift -= totalBytesToRemove
		// log.Printf("New normalized alignment: %vx %v\n", len(item.RuneVal), align)

		for i := 0; i < len(item.RuneVal); i++ {
			normalizedAlignments = append(normalizedAlignments, align)
		}

		// Then we keep only the char for string reconstruction
		normalizedRunes = append(normalizedRunes, []rune(item.RuneVal)...)

		// fmt.Printf("replacedCharSize: %v\n", replacedCharSize)
		// fmt.Printf("totalBytesToRemove: %v\n", totalBytesToRemove)
		// fmt.Printf("next offset: %v\n", offset)

	} // End of `For` block

	// Apply the changes to the remaining original alignments
	if oShift != 0 {
		// log.Printf("Shifting the end  from %v using shift: %v\n", endShiftStart, oShift)

		var endShift [][]int
		for _, item := range n.alignments[endShiftStart:] {
			endShift = append(endShift, item)
		}
		endRange := expandAlignments(endShift)

		// log.Printf("End range: %+v\n", endRange)

		if endRange != nil {
			alignments := n.alignmentsOriginal[endRange[0]:endRange[1]]
			var newAlignments [][]int
			if len(alignments) > 0 {
				// log.Printf("Alignments before shifting: %+v\n", alignments)
				for _, offsets := range alignments {
					newOffset0 := applySign(offsets[0], oShift)
					newOffset1 := applySign(offsets[1], oShift)
					newAlignments = append(newAlignments, []int{newOffset0, newOffset1})
				}
				aligns := append(n.alignmentsOriginal[:endRange[0]], newAlignments...)
				aligns = append(aligns, n.alignmentsOriginal[endRange[1]:]...)
				n.alignmentsOriginal = aligns
				// log.Printf("After: %+v\n", newAlignments)
			}
		}
	}

	// replace alignments with new ones in range
	var newAlignments [][]int
	// fmt.Printf("nRange: %v\n", nRange)
	// fmt.Printf("normalizedAlignments: %+v\n", normalizedAlignments)
	// fmt.Printf("alignments to be replaced: %v\n", n.alignments[nRange.start:nRange.end])
	// fmt.Printf("normalized alignments before action: %v\n", n.alignments)
	if len(n.alignments[:nRange.start]) == 0 { // at the beginning
		newAlignments = append(normalizedAlignments, n.alignments[nRange.end:]...)
	} else { // in the middle
		var beforeAligns, afterAligns [][]int
		for i, a := range n.alignments {
			if i < nRange.start {
				beforeAligns = append(beforeAligns, a)
			}
			if i >= nRange.end {
				afterAligns = append(afterAligns, a)
			}
		}
		newAlignments = append(beforeAligns, normalizedAlignments...)
		newAlignments = append(newAlignments, afterAligns...)
	}
	n.alignments = newAlignments

	// Finally, change the `normalized` string
	newNormalized := n.normalized[:nRange.start] + string(normalizedRunes) + n.normalized[nRange.end:]
	n.normalized = newNormalized

	// log.Printf("New normalized alignments: %+v\nNew original alignments: %+v\n", n.alignments, n.alignmentsOriginal)
	// log.Printf("New normalized string: %q\n", n.normalized)

	return n
}

// Transform applies transformations to the current normalized version, updating the current
// alignments with the new ones.
// This method expect an Iterator yielding each rune of the new normalized string
// with a `change` interger size equals to:
//   - `1` if this is a new rune
//   - `-N` if the char is right before N removed runes
//   - `0` if this rune represents the old one (even if changed)
//
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
func (n *NormalizedString) Transform(m []ChangeMap, initialOffset int) (retVal *NormalizedString) {
	start := 0
	end := len(n.original)
	wholeRange := NewRange(start, end, OriginalTarget)
	return n.TransformRange(wholeRange, m, initialOffset)
}

func (n *NormalizedString) NFD() (retVal *NormalizedString) {

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
		// runes := []rune(string(it.Next()))
		runes := bytes.Runes(it.Next())

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

	// for i, c := range changeMap {
	// fmt.Printf("%v - Char: %+q - changes: %v\n", i, c.RuneVal, c.Changes)
	// }

	return n.Transform(changeMap, 0)
}

func (n *NormalizedString) NFC() (retVal *NormalizedString) {
	var (
		changeMap []ChangeMap
		it        norm.Iter
	)

	s := n.normalized

	isNFC := norm.Form.IsNormalString(norm.NFC, s)
	if isNFC {
		return n
	}

	it.InitString(norm.NFD, s)

	for !it.Done() {
		// runes := []rune(string(it.Next()))
		runes := bytes.Runes(it.Next())

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

func (n *NormalizedString) NFKD() (retVal *NormalizedString) {

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

func (n *NormalizedString) NFKC() (retVal *NormalizedString) {

	var (
		changeMap []ChangeMap
		it        norm.Iter
	)

	s := n.normalized

	isNFKC := norm.Form.IsNormalString(norm.NFKC, s)

	if isNFKC {
		return n
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
func (n *NormalizedString) Filter(fn func(rune) bool) (retVal *NormalizedString) {

	var (
		removed   int = 0
		runes     []rune
		changeMap []ChangeMap
	)

	for _, r := range []rune(n.normalized) {
		runes = append(runes, r)
	}

	revRunes := slice.Reverse(runes).([]rune)

	for _, r := range revRunes {
		if fn(r) {
			if removed > 0 {
				changeMap = append(changeMap, ChangeMap{
					RuneVal: string(r),
					Changes: -(removed),
				})
				removed = 0
			} else {
				changeMap = append(changeMap, ChangeMap{
					RuneVal: string(r),
					Changes: 0,
				})
			}
		} else {
			removed += 1
		}
	}

	revChangeMap := slice.Reverse(changeMap).([]ChangeMap)

	// fmt.Printf("Alignments: %+v\n", n.alignments)
	// fmt.Printf("changeMap: %+v\n", revChangeMap)

	// for _, item := range revChangeMap {
	// fmt.Printf("item: %+v\n", item)
	// }

	return n.Transform(revChangeMap, removed)
	// return n.Transform(revChangeMap, 0)

}

// Prepend adds given string to the begining of NormalizedString
func (n *NormalizedString) Prepend(s string) (retVal *NormalizedString) {
	chars := []rune(n.normalized)
	var changeMap []ChangeMap
	if len(chars) == 0 {
		return n
	}

	next := chars[0]
	for i, r := range []rune(s) {
		var c ChangeMap
		if i == 0 {
			c = ChangeMap{string(r), 0}
		} else {
			c = ChangeMap{string(r), 1}
		}
		changeMap = append(changeMap, c)
	}
	changeMap = append(changeMap, ChangeMap{string(next), 1})
	inputRange := NewRange(0, len([]byte(string(next))), NormalizedTarget)

	return n.TransformRange(inputRange, changeMap, 0)
}

// Append adds given string to the end of NormalizedString
func (n *NormalizedString) Append(s string) (retVal *NormalizedString) {

	if n.normalized == "" {
		return n
	} else {
		var lastRuneIdx int
		var lastRune rune
		for i, r := range n.normalized {
			lastRuneIdx = i
			lastRune = r
		}

		inputRange := NewRange(lastRuneIdx, len(n.normalized), NormalizedTarget)
		var changeMap []ChangeMap = []ChangeMap{{string(lastRune), 0}}
		for _, r := range []rune(s) {
			changeMap = append(changeMap, ChangeMap{string(r), 1})
		}
		return n.TransformRange(inputRange, changeMap, 0)
	}
}

// NormFn is a convenient function type for applying
// on each `char` of normalized string
type NormFn func(rune) rune

// Map maps and applies function to each `char` of normalized string
func (n *NormalizedString) Map(nfn NormFn) (retVal *NormalizedString) {
	s := n.normalized
	var changeMap []ChangeMap
	for _, r := range []rune(s) {
		changeMap = append(changeMap, ChangeMap{string(r), 0})
	}

	return n.Transform(changeMap, 0)
}

// ForEach applies function on each `char` of normalized string
// Similar to Map???
func (n *NormalizedString) ForEach(nfn NormFn) (retVal *NormalizedString) {
	s := n.normalized
	var runes []rune
	for _, r := range []rune(s) {
		runes = append(runes, nfn(r))
	}
	n.normalized = string(runes)

	return n
}

// RemoveAccents removes all Unicode Mn group (M non-spacing)
func (n *NormalizedString) RemoveAccents() (retVal *NormalizedString) {
	return n.Filter(func(r rune) bool {
		return !unicode.Is(unicode.Mn, r)
	})
}

// Lowercase transforms string to lowercase
func (n *NormalizedString) Lowercase() (retVal *NormalizedString) {
	n.normalized = strings.ToLower(n.normalized)

	return n
}

// Uppercase transforms string to uppercase
func (n *NormalizedString) Uppercase() (retVal *NormalizedString) {
	n.normalized = strings.ToUpper(n.normalized)

	return n
}

// Clear clears the normalized part of the string
func (n *NormalizedString) Clear() {
	length := n.Len()
	n.Transform([]ChangeMap{}, length)
}

// Split the current string in many subparts. Specify what to do with the
// delimiter.
//
// This method will always ensure that the entire `NOrmalizedString` is covered in the
// produced subparts. This means that the delimiter parts will also be included,
// and will appear empty if we don't want to include them (their `original`
// part will still be present). It should always be possible to merge all the
// subparts back to the original `NormalizedString`
//
// ## Splitting Behavior for the delimiter
//
// The behavior can be one of the followings:
// When splitting on `'-'` for example, with input `the-final--countdown`:
//   - RemovedBehavior => `[ "the", "", "final", "", "", "countdown" ]`
//   - IsolatedBehavior => `[ "the", "-", "final", "-", "-", "countdown" ]`
//   - MergedWithPreviousBehavior => `[ "the-", "final-", "-", "countdown" ]`
//   - MergedWithNextBehavior => `[ "the", "-final", "-", "-countdown" ]`
//   - Contiguous => `[ "the", "-", "final", "--", "countdown" ]`
func (n *NormalizedString) Split(pattern Pattern, behavior SplitDelimiterBehavior) (retVal []NormalizedString) {

	// fmt.Printf("input normalized: %v\n", n)

	matches := pattern.FindMatches(n.GetNormalized())

	// fmt.Printf("length of matches: %v\n", len(matches))
	// for i, m := range matches {
	// fmt.Printf("%v match: %v\n", i, m)
	// }

	// Process the matches according to the selected behavior: []OfssetsMatch
	// where `Match` field is `shouldRemove`
	var splits []OffsetsMatch
	switch behavior {
	case IsolatedBehavior:
		for _, m := range matches {
			m.Match = false
			splits = append(splits, m)
		}
	case RemovedBehavior:
		splits = matches
	case MergedWithPreviousBehavior:
		previousMatch := false
		var acc []OffsetsMatch
		for _, m := range matches {
			if m.Match && !previousMatch {
				if len(acc) > 0 {
					// update last item of acc
					acc[len(acc)-1].Offsets[1] = m.Offsets[1]
				} else {
					acc = append(acc, OffsetsMatch{Offsets: m.Offsets, Match: false})
				}
			} else {
				acc = append(acc, OffsetsMatch{Offsets: m.Offsets, Match: false})
			}

			previousMatch = m.Match
		}
		splits = acc
	case ContiguousBehavior:
		previousMatch := false
		var acc []OffsetsMatch
		for _, m := range matches {
			if m.Match == previousMatch {
				if len(acc) > 0 {
					// update last item of acc
					acc[len(acc)-1].Offsets[1] = m.Offsets[1]
				} else {
					acc = append(acc, OffsetsMatch{Offsets: m.Offsets, Match: false})
				}
			} else {
				acc = append(acc, OffsetsMatch{Offsets: m.Offsets, Match: false})
			}

			previousMatch = m.Match
		}
		splits = acc

	case MergedWithNextBehavior:
		previousMatch := false
		var acc []OffsetsMatch
		// iterate reversively
		for i := len(matches) - 1; i >= 0; i-- {
			m := matches[i]
			if m.Match && !previousMatch {
				if len(acc) > 0 {
					// update last item of acc
					acc[len(acc)-1].Offsets[0] = m.Offsets[0]
				} else {
					acc = append(acc, OffsetsMatch{Offsets: m.Offsets, Match: false})
				}
			} else {
				acc = append(acc, OffsetsMatch{Offsets: m.Offsets, Match: false})
			}

			previousMatch = m.Match
		}

		// reverse back
		for i := len(acc) - 1; i >= 0; i-- {
			splits = append(splits, acc[i])
		}
	}

	// Then split according to the computed splits
	var slices []NormalizedString
	for _, split := range splits {
		if !split.Match {
			slice := n.Slice(NewRange(split.Offsets[0], split.Offsets[1], NormalizedTarget))
			if slice != nil {
				slices = append(slices, *slice)
			}
		}
	}

	// log.Printf("output: %+v\n", slices)

	return slices
}

// LStrip removes leading spaces
func (n *NormalizedString) LStrip() (retVal *NormalizedString) {
	return n.lrstrip(true, false)
}

// RStrip removes trailing spaces
func (n *NormalizedString) RStrip() (retVal *NormalizedString) {
	return n.lrstrip(false, true)
}

// Strip remove leading and trailing spaces
func (n *NormalizedString) Strip() (retVal *NormalizedString) {
	return n.lrstrip(true, true)
}

// lrstrip - Private func to help with exposed strip funcs
func (n *NormalizedString) lrstrip(left, right bool) (retVal *NormalizedString) {
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

// Len returns length (in bytes) of normalized string
func (n *NormalizedString) Len() int {
	return len(n.normalized)
}

// LenOriginal returns the length of Original string in bytes
func (n *NormalizedString) LenOriginal() int {
	return len(n.GetOriginal())
}

// IsEmpty returns whether the normalized string is empty
func (n *NormalizedString) IsEmpty() bool {
	return n.Len() == 0
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

// RangeOf returns a range of normalized string
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
	if start < 0 || start >= sLen || end > sLen || start >= end {
		return ""
	}

	slicedRunes := runes[start:end]
	return string(slicedRunes)
}

// expandAlignments returns an offsets slice covered by a slice of alignments.
// Return value is a slice of 2 elements (start, end).
func expandAlignments(alignments [][]int) (retVal []int) {
	if len(alignments) == 0 {
		return nil
	}
	start := alignments[0][0]
	end := alignments[len(alignments)-1][1]

	return []int{start, end}
}

// applySign adds or substracts a signed value on a origin value. Makes sure of avoiding
// any substraction overflow, flooring at 0.
func applySign(origin int, signed int) int {
	if result := origin + signed; result < 0 {
		return 0
	} else {
		return result
	}
}

func (n *NormalizedString) Replace(pattern Pattern, content string) (retVal *NormalizedString) {

	offset := 0
	matches := pattern.FindMatches(n.normalized)

	if len(matches) == 0 {
		// return nil
		return n
	}

	for _, m := range matches {
		if m.Match {
			start := m.Offsets[0]
			end := m.Offsets[1]
			r := []int{m.Offsets[0], m.Offsets[1]}
			r0 := applySign(r[0], offset)
			r1 := applySign(r[1], offset)

			newLen := 0
			var removedChars int // num of removed chars
			removedStr := n.normalized[r0:r1]
			removedChars = len([]rune(removedStr))

			var changeMap []ChangeMap
			for _, r := range []rune(content) {
				newLen += len([]byte(string(r)))
				changeMap = append(changeMap, ChangeMap{string(r), 1})
			}

			n = n.TransformRange(NewRange(r0, r1, NormalizedTarget), changeMap, removedChars)

			oldLen := end - start
			offset += newLen - oldLen
		}
	}

	return n
}

type byteIdxRune struct {
	byteIdx int
	runeIdx int
	char    rune
}

// BytesToChar converts a given range from bytes to `char`
func BytesToChar(s string, byteRange []int) (retVal []int) {
	var start, end int
	if reflect.DeepEqual(byteRange, []int{0, 0}) {
		start = 0
		end = 0
	} else {
		start = -1 // nil value
		end = -1   // nil value
	}

	var selectedChars []byteIdxRune
	var currRuneIdx int = 0
	for i, char := range []rune(s) {
		if i >= byteRange[0] && i <= byteRange[1] {
			selectedChars = append(selectedChars, byteIdxRune{
				byteIdx: i,
				runeIdx: currRuneIdx,
				char:    char,
			})
		}
		currRuneIdx++
	}

	for _, item := range selectedChars {
		if item.byteIdx == byteRange[0] {
			start = item.runeIdx
		}

		if item.byteIdx == byteRange[1] {
			end = item.runeIdx
		}

		if item.byteIdx+len([]byte(string(item.char))) == byteRange[1] {
			end = item.runeIdx + 1
		}
	}

	return []int{start, end}
}

// CharToBytes converts a given range from `char` to bytes
func CharToBytes(s string, charRange []int) (retVal []int) {
	var start, end int
	if reflect.DeepEqual(charRange, []int{0, 0}) {
		start = 0
		end = 0
	} else {
		start = -1
		end = -1
	}

	var chars []byteIdxRune
	var currRuneIdx int = 0
	for i, char := range []rune(s) {
		chars = append(chars, byteIdxRune{
			byteIdx: i,
			runeIdx: currRuneIdx,
			char:    char,
		})
		currRuneIdx++
	}

	if charRange[0] == charRange[1] {
		for i := range chars {
			if i == charRange[0] {
				start = chars[i+1].byteIdx
				end = chars[i+1].byteIdx
			}
		}
	} else {
		var selected []byteIdxRune
		for _, c := range chars {
			if c.byteIdx > charRange[0] && c.byteIdx <= charRange[1] {
				selected = append(selected, c)
			}
		}

		for _, c := range selected {
			if start == -1 {
				start = c.byteIdx
			}

			end = c.byteIdx + len([]byte(string(c.char)))
		}
	}

	return []int{start, end}
}
