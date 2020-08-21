package tokenizer

// wrapper for subpart of NormalizedString

import (
	"fmt"
	"log"
	"reflect"

	"github.com/sugarme/tokenizer/normalizer"
)

// SubString contains the underlying `NormalizedString` as well as
// its offsets in the original string. These offsets are in the
// `original` referential.
type SubString struct {
	// Normalized is the underlying `NormalizedString`. Each SubString
	// is represented by a `NormalizedString`. In the end, there might
	// be many SubStrings representing various parts of the original
	// input string.
	Normalized *normalizer.NormalizedString

	// OriginalOffsets is the Offsets of `NormalizedString` in the `original`
	// input string. These is useful to find sub-part of the input string
	// represented by `NormalizedString`
	OriginalOffsets Offsets
}

// NewSubString creates a new SubString with input of a NormalizedString and its
// offsets on `original` input string
func NewSubString(normalized *normalizer.NormalizedString, originalOffsets Offsets) (retVal SubString) {
	return SubString{normalized, originalOffsets}
}

// PreTokenizedString contains SubStrings. It helps to keep track of offsets
// during the whole normalization and pre-tokenization steps.
type PreTokenizedString struct {
	parts   []SubString
	nextIdx int
}

// SplitFn takes a `NormalizedString` and returns an iterator over the
// produced `NormalizedString`.
//
// NOTE. SplitFn is free of modifying these `NormalizedString` as long as:
// The produced `NormalizedString`, if combined back together, must have
// the same `original` string as the original one given to `SplitFn`. This
// means that for the offsets tracking to work as expected, `SplitFn` must
// produce "splits" of the ORIGINAL string.
type SplitFn func(int, *normalizer.NormalizedString) []*normalizer.NormalizedString

// Split splits the `PreTokenizedString` by providing a `SplitFn` which is in
// charge of splitting each substring (`NormalizedString`) into multiple parts.
func (pt *PreTokenizedString) Split(splitFn SplitFn) (err error) {

	var newParts []SubString
	for i, sub := range pt.parts {
		originalLen := sub.Normalized.LenOriginal()
		originalOffsets := sub.OriginalOffsets

		newLen := 0
		for _, normalized := range splitFn(i, sub.Normalized) {
			len := normalized.LenOriginal()
			start := originalOffsets.Start + newLen
			end := originalOffsets.Start + newLen + len
			newS := NewSubString(normalized, Offsets{start, end})

			newParts = append(newParts, newS)
			newLen += len
		}

		if originalLen != newLen {
			return fmt.Errorf("Split pre-tokenized string must represent the entire original string.\nOriginal length %v - new length %v\n", originalLen, newLen)
		}
	}

	pt.parts = newParts

	return nil
}

// Next implement iterator interface for `PreTokenizedString`
func (pt *PreTokenizedString) Next() (retVal SubString, ok bool) {
	if pt.nextIdx == len(pt.parts) {
		return retVal, false
	}

	retVal = pt.parts[pt.nextIdx]
	pt.nextIdx += 1

	return retVal, true
}

// Returns a list of normalized string and the associated offsets,
// either in original or normalized referential
func (pt *PreTokenizedString) GetNormalized(offsetType normalizer.IndexOn) (retVal []PreToken) {
	var (
		offset  int = 0
		preToks []PreToken
	)

	for _, sub := range pt.parts {
		var offsets Offsets
		switch offsetType {
		case normalizer.OriginalTarget:
			offsets = Offsets{
				Start: sub.OriginalOffsets.Start,
				End:   sub.OriginalOffsets.Start + sub.Normalized.LenOriginal(),
			}
		case normalizer.NormalizedTarget:
			length := sub.Normalized.Len()
			offset += length
			offsets = Offsets{Start: offset - length, End: offset}
		}
		preToks = append(preToks, PreToken{Value: sub.Normalized.GetNormalized(), Offsets: offsets})
	}

	return preToks
}

// IntoMerged merges back to a NormalizedString
func (pt *PreTokenizedString) IntoMerged() (retVal *normalizer.NormalizedString) {

	var end int = 0
	if len(pt.parts) > 0 {
		end = pt.parts[len(pt.parts)-1].OriginalOffsets.End
	}

	offsets := Offsets{0, end}
	var normalized *normalizer.NormalizedString
	for i, sub := range pt.parts {
		if i == 0 {
			normalized = sub.Normalized
		} else {
			normalized.MergeWith(sub.Normalized)
		}
	}

	if !reflect.DeepEqual(offsets, Offsets{0, normalized.LenOriginal()}) {
		log.Fatalf("Merging Error: original string length and merged string length are mismatche.\n")
	}

	return normalized
}

// NewNormalizedStringFromNS creates a PreTokenizedString from input
// NormalizedString
func NewPreTokenizedStringFromNS(n *normalizer.NormalizedString) (retVal PreTokenizedString) {
	originalOffsets := Offsets{0, n.LenOriginal()}

	return PreTokenizedString{
		parts: []SubString{
			{
				Normalized:      n,
				OriginalOffsets: originalOffsets,
			},
		},
		nextIdx: 0,
	}
}

// NewPreTokenizedString create a new PreTokenizedString from input string
func NewPreTokenizedString(s string) (retVal PreTokenizedString) {
	n := normalizer.NewNormalizedFrom(s)
	return NewPreTokenizedStringFromNS(n)
}
