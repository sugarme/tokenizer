package tokenizer

import (
	"errors"
	"reflect"

	"github.com/sugarme/tokenizer/normalizer"
)

type PaddingDirection int

const (
	Left PaddingDirection = iota
	Right
)

// Encoding represents the output of tokenizer
type Encoding struct {
	Normalized       normalizer.Normalized
	Ids              []uint32
	TypeIds          []uint32
	Tokens           []string
	Offsets          []Offsets
	SpecialTokenMask []uint32
	AttentionMask    []uint32
	Overflowing      []Encoding
}

// NewEncoding initiate a new encoding from input data
func NewEncoding(normalized normalizer.Normalized, ids []uint32, typeIds []uint32, tokens []string, offsets []Offsets, specialTokenMask []uint32, attentionMask []uint32, overflowing []Encoding) Encoding {
	return Encoding{
		normalized,
		ids,
		typeIds,
		tokens,
		offsets,
		specialTokenMask,
		attentionMask,
		overflowing,
	}
}

// Default creates an encoding with default values
func DefaultEncoding() Encoding {
	return Encoding{
		Normalized:       *normalizer.NewNormalizedFrom(""),
		Ids:              []uint32{0},
		TypeIds:          []uint32{0},
		Tokens:           []string{},
		Offsets:          []Offsets{},
		SpecialTokenMask: []uint32{},
		AttentionMask:    []uint32{},
		Overflowing:      []Encoding{},
	}
}

// GetNormalized returns normalized string from encoding
func (e *Encoding) GetNormalized() normalizer.Normalized {
	return e.Normalized
}

// GetTokenreturns tokens from encoding
func (e *Encoding) GetTokens() []string {
	return e.Tokens
}

// GetIds returns Ids from encoding
func (e *Encoding) GetIds() []uint32 {
	return e.Ids
}

// GetTypeIds returns type Ids from encoding
func (e *Encoding) GetTypeIds() []uint32 {
	return e.TypeIds
}

// GetOffsets returns offsets from encoding
func (e *Encoding) GetOffsets() []Offsets {
	return e.Offsets
}

// GetSpecialTokenMask returns specialTokenMask from encoding
func (e *Encoding) GetSpecialTokenMask() []uint32 {
	return e.SpecialTokenMask
}

// GetAttentionMask returns attentionMask from encoding
func (e *Encoding) GetAttentionMask() []uint32 {
	return e.AttentionMask
}

// GetOverflowing returns overflowing from encoding
func (e *Encoding) GetOverflowing() []Encoding {
	return e.Overflowing
}

// TakeOverflowing returns overflowing and reset it to empty at encoding
func (e *Encoding) TakeOverflowing() []Encoding {
	o := e.Overflowing
	e.Overflowing = []Encoding{}
	return o
}

// Truncate truncates the current encoding
func (e *Encoding) Truncate(maxLen uint, stride uint) error {

	if stride >= maxLen || maxLen == 0 {
		return errors.New("Invalid input maxLen or stride (stride must be less than maxLen and maxLen must be greater than zero.)")
	}

	if maxLen >= uint(len(e.Ids)) {
		// do nothing
		return nil
	}

	// Truncating at maxLen (exclusive) to keep.
	// The rest (overflowing) from maxLen (inclusive)
	newIds := e.Ids[0:maxLen]
	oIds := e.Ids[maxLen:len(e.Ids)] // overflowing
	newTypeIds := e.TypeIds[0:maxLen]
	oTypeIds := e.TypeIds[maxLen:len(e.TypeIds)]
	newTokens := e.Tokens[0:maxLen]
	oTokens := e.Tokens[maxLen:len(e.Tokens)]
	newOffsets := e.Offsets[0:maxLen]
	oOffsets := e.Offsets[maxLen:len(e.Offsets)]
	newSpeToks := e.SpecialTokenMask[0:maxLen]
	oSpeToks := e.SpecialTokenMask[maxLen:len(e.SpecialTokenMask)]
	newAttent := e.AttentionMask[0:maxLen]
	oAttent := e.AttentionMask[maxLen:len(e.AttentionMask)]

	e.Ids = newIds
	e.TypeIds = newTypeIds
	e.Tokens = newTokens
	e.Offsets = newOffsets
	e.SpecialTokenMask = newSpeToks
	e.AttentionMask = newAttent

	// Separate the overflowing part into as many Encoding as needed
	partSize := maxLen - stride
	overflowing := make([]Encoding, 0)
	partId := 0
	prevEncoding := e

	// while loop
	for int(partSize)*partId < len(oIds) {
		o := Encoding{
			Normalized: e.Normalized,
			// Which way is better? using reflect or just type assertion
			// Ids:        (getCurrentPart(prevEncoding.Ids, oIds, partSize, uint(partId), stride)).([]uint32),
			Ids:              reflect.ValueOf(getCurrentPart(prevEncoding.Ids, oIds, partSize, uint(partId), stride)).Interface().([]uint32),
			TypeIds:          reflect.ValueOf(getCurrentPart(prevEncoding.TypeIds, oTypeIds, partSize, uint(partId), stride)).Interface().([]uint32),
			Tokens:           reflect.ValueOf(getCurrentPart(prevEncoding.Tokens, oTokens, partSize, uint(partId), stride)).Interface().([]string),
			Offsets:          reflect.ValueOf(getCurrentPart(prevEncoding.Offsets, oOffsets, partSize, uint(partId), stride)).Interface().([]Offsets),
			SpecialTokenMask: reflect.ValueOf(getCurrentPart(prevEncoding.SpecialTokenMask, oSpeToks, partSize, uint(partId), stride)).Interface().([]uint32),
			AttentionMask:    reflect.ValueOf(getCurrentPart(prevEncoding.AttentionMask, oAttent, partSize, uint(partId), stride)).Interface().([]uint32),
			Overflowing:      make([]Encoding, 0),
		}

		partId += 1
		overflowing = append(overflowing, o)
		prevEncoding = &overflowing[len(overflowing)-1]
	}

	e.Overflowing = overflowing

	return nil

}

// MergeWith merges the current encoding with other (pair) encoding
func (e *Encoding) MergeWith(pair Encoding) {
	// Merge overflowing
	overflowings := make([]Encoding, 0)
	// 1. All current overflowing with all other overflowing
	for _, o := range e.Overflowing {
		currO := o
		// 1.1. The pair itself
		currO.MergeWith(pair) // recursively call
		overflowings = append(overflowings, currO)
		currO = o // reset

		// 1.2. The pair's overflowing
		for _, otherO := range pair.Overflowing {
			currO.MergeWith(otherO)
			overflowings = append(overflowings, currO)
			currO = o // reset
		}
	}

	// 2. Current encoding with all other overflowing
	for _, otherO := range pair.Overflowing {
		newE := e
		newE.MergeWith(otherO)
		overflowings = append(overflowings, *newE)
	}

	// 3. Current encoding and other encoding
	e.Normalized.MergeWith(pair.Normalized.Get())
	e.Ids = append(e.Ids, pair.Ids...)
	e.TypeIds = append(e.TypeIds, pair.TypeIds...)
	e.Tokens = append(e.Tokens, pair.Tokens...)
	// Offsets
	var startingOffset int = 0
	for _, o := range e.Offsets {
		if o.End > startingOffset {
			startingOffset = o.End
		}
	}
	for _, o := range pair.Offsets {
		adjustedO := Offsets{
			Start: o.Start + startingOffset,
			End:   o.End + startingOffset,
		}
		e.Offsets = append(e.Offsets, adjustedO)
	}
	e.SpecialTokenMask = append(e.SpecialTokenMask, pair.SpecialTokenMask...)
	e.AttentionMask = append(e.AttentionMask, pair.AttentionMask...)
	e.Overflowing = overflowings

}

// Pad pads current encoding with given length, values to either Left or Right direction
func (e *Encoding) Pad(targetLength uint, padId uint32, padTypeId uint32, padToken string, direction PaddingDirection) {
	// 1. Recursively call for overflowing part
	for _, o := range e.Overflowing {
		o.Pad(targetLength, padId, padTypeId, padToken, direction)
	}

	// 2. Check whether we should pad encoding itself
	// if wanted padding length is smaller, then do nothing
	if len(e.Ids) >= int(targetLength) {
		return
	}

	padLength := int(targetLength) - len(e.Ids)

	switch direction {
	case Left:
		newIds := make([]uint32, padLength)
		for i := 0; i < len(newIds); i++ {
			newIds[i] = padId
		}
		newIds = append(newIds, e.Ids...)
		e.Ids = newIds

		newTypeIds := make([]uint32, padLength)
		for i := 0; i < len(newTypeIds); i++ {
			newTypeIds[i] = padTypeId
		}
		newTypeIds = append(newTypeIds, e.Ids...)
		e.TypeIds = newTypeIds

		newTokens := make([]string, padLength)
		for i := 0; i < len(newTokens); i++ {
			newTokens[i] = padToken
		}
		newTokens = append(newTokens, e.Tokens...)
		e.Tokens = newTokens

		newSpecialTokenMask := make([]uint32, padLength)
		for i := 0; i < len(newSpecialTokenMask); i++ {
			newSpecialTokenMask[i] = 1
		}
		newSpecialTokenMask = append(newSpecialTokenMask, e.SpecialTokenMask...)
		e.SpecialTokenMask = newSpecialTokenMask

		newAttentionMask := make([]uint32, padLength)
		for i := 0; i < len(newAttentionMask); i++ {
			newAttentionMask[i] = 0
		}
		newAttentionMask = append(newAttentionMask, e.AttentionMask...)
		e.AttentionMask = newAttentionMask

		newOffsets := make([]Offsets, padLength)
		for i := 0; i < len(newIds); i++ {
			newOffsets[i] = Offsets{0, 0}
		}
		newOffsets = append(newOffsets, e.Offsets...)
		e.Offsets = newOffsets

	case Right:
		for i := 0; i < padLength; i++ {
			e.Ids = append(e.Ids, padId)
			e.TypeIds = append(e.TypeIds, padTypeId)
			e.Tokens = append(e.Tokens, padToken)
			e.SpecialTokenMask = append(e.SpecialTokenMask, 1)
			e.AttentionMask = append(e.AttentionMask, 0)
			e.Offsets = append(e.Offsets, Offsets{0, 0})
		}

	}
}

func getCurrentPart(previous, current interface{}, size, idx, stride uint) interface{} {

	switch current.(type) {
	case []uint32:
		var curr, prev []uint32
		if int((idx+1)*size) > reflect.ValueOf(current).Len() {
			curr = current.([]uint32)[(idx * size):]
		} else {
			curr = current.([]uint32)[(idx * size) : (idx+1)*size]
		}
		prev = previous.([]uint32)[len(previous.([]uint32))-int(stride):]
		// concat
		return append(prev, curr...)
	case []string:
		var curr, prev []string
		if int((idx+1)*size) > reflect.ValueOf(current).Len() {
			curr = current.([]string)[(idx * size):]
		} else {
			curr = current.([]string)[(idx * size) : (idx+1)*size]
		}
		prev = previous.([]string)[len(previous.([]string))-int(stride):]
		// concat
		return append(prev, curr...)
	case []Offsets:
		var curr, prev []Offsets
		if int((idx+1)*size) > reflect.ValueOf(current).Len() {
			curr = current.([]Offsets)[(idx * size):]
		} else {
			curr = current.([]Offsets)[(idx * size) : (idx+1)*size]
		}
		prev = previous.([]Offsets)[len(previous.([]Offsets))-int(stride):]
		// concat
		return append(prev, curr...)

	}

	return nil

}

/*
 * func InterfaceSlice(slice interface{}) []interface{} {
 *   s := reflect.ValueOf(slice)
 *   if s.Kind() != reflect.Slice {
 *     panic("InterfaceSlice() given a non-slice type")
 *   }
 *
 *   ret := make([]interface{}, s.Len())
 *
 *   for i := 0; i < s.Len(); i++ {
 *     ret[i] = s.Index(i).Interface()
 *   }
 *
 *   return ret
 * } */
