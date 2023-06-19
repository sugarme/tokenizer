package tokenizer

import (
	"fmt"
	"log"
	"reflect"

	"github.com/sugarme/tokenizer/util"
)

type PaddingDirection int

const (
	Left PaddingDirection = iota
	Right
)

// Encoding represents the output of tokenizer
type Encoding struct {
	Ids              []int         // ID produced by the `tokenizer`
	TypeIds          []int         // Type of the ID
	Tokens           []string      // Tokens associated with each ID
	Offsets          [][]int       // Offsets of the token/ID from the NormalizedString
	SpecialTokenMask []int         // Mask identifying special tokens
	AttentionMask    []int         // Mask identifying padding tokens for the attention mechanism
	Overflowing      []Encoding    // A list of overflowing generated when being truncated
	Words            []int         // Optional - Indexes of the word associated with each token/ID. None value = -1
	SequenceRanges   map[int]Range // Range of tokens covered by each sequence. If empty -> only one sequence and covers the entire range.
}

type EncodingOpts struct {
	Words         []int
	SequenceRange map[int]Range
}

type EncodingOpt func(o *EncodingOpts)

func WithWordsEncodingOpt(v []int) EncodingOpt {
	return func(o *EncodingOpts) {
		o.Words = v
	}
}
func WithSequenceRangeEncodingOpt(v map[int]Range) EncodingOpt {
	return func(o *EncodingOpts) {
		o.SequenceRange = v
	}
}

func DefaultEncodingOpts() *EncodingOpts {
	return &EncodingOpts{
		Words:         nil,
		SequenceRange: make(map[int]Range),
	}
}

// NewEncoding initiate a new encoding from input data
func NewEncoding(ids []int, typeIds []int, tokens []string, offsets [][]int, specialTokenMask []int, attentionMask []int, overflowing []Encoding, opts ...EncodingOpt) *Encoding {
	o := DefaultEncodingOpts()
	for _, opt := range opts {
		opt(o)
	}

	return &Encoding{
		ids,
		typeIds,
		tokens,
		offsets,
		specialTokenMask,
		attentionMask,
		overflowing,
		o.Words,
		o.SequenceRange,
	}
}

func NewEncodingWithCapacity(l int) (retVal *Encoding) {
	return &Encoding{
		Ids:              make([]int, l),
		TypeIds:          make([]int, l),
		Tokens:           make([]string, l),
		Offsets:          make([][]int, l),
		SpecialTokenMask: make([]int, l),
		AttentionMask:    make([]int, l),
		Overflowing:      []Encoding{},
		Words:            make([]int, l),
		SequenceRanges:   make(map[int]Range),
	}
}

// Default creates an encoding with default values
func DefaultEncoding() *Encoding {
	return &Encoding{
		Ids:              []int{},
		TypeIds:          []int{},
		Tokens:           []string{},
		Offsets:          [][]int{},
		SpecialTokenMask: []int{},
		AttentionMask:    []int{},
		Overflowing:      []Encoding{},
		Words:            nil,
		SequenceRanges:   make(map[int]Range),
	}
}

// NewEncodingFromTokens initiate Encoding from input tokens
func NewEncodingFromTokens(tokens []Token, typeId int) (retVal *Encoding) {
	var (
		ids     []int
		offsets [][]int
		toks    []string
	)
	for _, t := range tokens {
		ids = append(ids, t.Id)
		offsets = append(offsets, t.Offsets)
		toks = append(toks, t.Value)
	}

	typeIds := make([]int, len(tokens))
	// words := make([]int, len(tokens))
	var words []int
	specialTokenMask := util.Repeat(0, len(tokens))
	attentionMask := util.Repeat(1, len(tokens))

	return &Encoding{
		Ids:              ids,
		TypeIds:          typeIds,
		Tokens:           toks,
		Offsets:          offsets,
		SpecialTokenMask: specialTokenMask,
		AttentionMask:    attentionMask,
		Overflowing:      []Encoding{},
		Words:            words,
		SequenceRanges:   make(map[int]Range),
	}
}

func (e *Encoding) Clone() *Encoding {
	out := new(Encoding)
	err := util.DeepCopy(e, out)
	if err != nil {
		panic(err)
	}

	return out
}

// IsEmpty returns whether Encoding is empty
func (e *Encoding) IsEmpty() (retVal bool) {
	return len(e.Ids) == 0
}

// Len returns number of encoding tokens
func (e *Encoding) Len() (retVal int) {
	return len(e.Ids)
}

// NSequences returns number of sequences combined in this encoding.
func (e *Encoding) NSequences() int {
	if len(e.SequenceRanges) == 0 {
		return 1
	}

	return len(e.SequenceRanges)
}

// GetToken returns tokens from encoding
func (e *Encoding) GetTokens() []string {
	return e.Tokens
}

// GetWords returns word indexes on normalized string
func (e *Encoding) GetWords() []int {
	return e.Words
}

// SetWord set word index value at given index of word in e.Words slice
func (e *Encoding) SetWord(index int, val int) {
	e.Words[index] = val
}

// SetSequenceIds set the given sequence id for the whole range of tokens contained in this Encoding
func (e *Encoding) SetSequenceIds(sequenceId int) {
	if e.Len() > 0 {
		e.SequenceRanges[sequenceId] = NewRange(0, e.Len())
	}
}

func (e *Encoding) GetSequenceIds() []int {
	sequences := make([]int, e.Len())
	for seqId := 0; seqId < e.NSequences(); seqId++ {
		r := e.SequenceRanges[seqId]
		seqLen := r.Len()
		var a []int
		for i := 0; i < seqLen; i++ {
			a = append(a, seqId)
		}
		// replace items in Range r with seqId
		start := r[0]
		end := r[len(r)-1]
		before := sequences[:start]
		after := sequences[end:]
		sequences = append(before, a...)
		sequences = append(sequences, after...)
	}

	return sequences
}

// GetIds returns Ids from encoding
func (e *Encoding) GetIds() []int {
	return e.Ids
}

// GetTypeIds returns type Ids from encoding
func (e *Encoding) GetTypeIds() []int {
	return e.TypeIds
}

func (e *Encoding) SetTypeIds(typeIds []int) {
	e.TypeIds = typeIds
}

// GetOffsets returns offsets from encoding
func (e *Encoding) GetOffsets() [][]int {
	return e.Offsets
}

// GetSpecialTokenMask returns specialTokenMask from encoding
func (e *Encoding) GetSpecialTokenMask() []int {
	return e.SpecialTokenMask
}

// GetAttentionMask returns attentionMask from encoding
func (e *Encoding) GetAttentionMask() []int {
	return e.AttentionMask
}

// GetOverflowing returns overflowing from encoding
func (e *Encoding) GetOverflowing() []Encoding {
	return e.Overflowing
}

// SetOverflowing set overflowing.
func (e *Encoding) SetOverflowing(overflowing []Encoding) {
	e.Overflowing = overflowing
}

// TakeOverflowing returns overflowing and reset it to empty at encoding
func (e *Encoding) TakeOverflowing() []Encoding {
	o := e.Overflowing
	e.Overflowing = []Encoding{}
	return o
}

// Word2Tokens gets the encoded tokens corresponding the word
// at the given index in the input sequence
// in the form `(startToken, endToken + 1)`
//
// NOTE. e.Words is optional, therefore, there's case of `none` result
// if `none` result, `ok` will be false.
func (e *Encoding) Word2Tokens(word int) (startTok, endTok int, ok bool) {

	var start, end int = -1, -1

	var inRangeWords []int
	for _, w := range e.Words {
		if w <= word {
			inRangeWords = append(inRangeWords, w)
		}
	}
	for i, w := range inRangeWords {
		if w == word && start < 0 {
			start = i
		}
	}

	end = len(inRangeWords)

	if start != -1 && end != -1 {
		return start, end, true
	} else {
		return startTok, endTok, false
	}
}

// Word2Chars get the offsets of the word at a given index in
// the input sequence
func (e *Encoding) Word2Chars(word int) (retVal []int, ok bool) {
	start, end, ok := e.Word2Tokens(word)
	if end == 0 {
		return retVal, false
	} else {
		oStart := e.Offsets[start][0]
		oEnd := e.Offsets[end-1][1]
		return []int{oStart, oEnd}, true // Should we check whether `ok`?
	}
}

// Token2Chars get the offsets of the token at the given index
func (e *Encoding) Token2Chars(tokenIdx int) (retVal []int, ok bool) {
	if tokenIdx < 0 || tokenIdx > len(e.Offsets) {
		return retVal, false
	} else {
		return e.Offsets[tokenIdx], true
	}
}

// Token2Word get the word index of corresponding token if existing
func (e *Encoding) Token2Word(tokenIdx int) (retVal int, ok bool) {
	// naive search. TODO. improve algorithm
	for i, w := range e.Words {
		if i == tokenIdx {
			return w, true
		}
	}
	return retVal, false
}

// Char2Token returns a token index that contains the given `char` index
func (e *Encoding) Char2Token(pos int) (retVal int, ok bool) {
	for i, o := range e.Offsets {
		if pos >= o[0] && pos < o[1] {
			return i, true
		}
	}

	return -1, false
}

// Char2Word get the word index that contain the given `char` index
func (e *Encoding) Char2Word(pos int) (retVal int, ok bool) {
	if idx, ok := e.Char2Token(pos); ok {
		return e.Token2Word(idx)
	}

	return -1, false
}

// Truncate truncates the current encoding
func (e *Encoding) Truncate(maxLen int, stride int) (retVal *Encoding, err error) {

	if stride >= maxLen || maxLen == 0 {
		return retVal, fmt.Errorf("Invalid input maxLen or stride (stride must be less than maxLen and maxLen must be greater than zero.)")
	}

	if maxLen >= len(e.Ids) {
		// do nothing
		return e, nil
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
	newWords := e.Words[0:maxLen]
	oWords := e.Words[maxLen:len(e.Words)]

	e.Ids = newIds
	e.TypeIds = newTypeIds
	e.Tokens = newTokens
	e.Offsets = newOffsets
	e.SpecialTokenMask = newSpeToks
	e.AttentionMask = newAttent
	e.Words = newWords

	// Separate the overflowing part into as many Encoding as needed
	partSize := maxLen - stride
	overflowing := make([]Encoding, 0)
	partId := 0
	prevEncoding := e

	// while loop
	for int(partSize)*partId < len(oIds) {
		o := Encoding{
			Ids:              reflect.ValueOf(getCurrentPart(prevEncoding.Ids, oIds, partSize, partId, stride)).Interface().([]int),
			TypeIds:          reflect.ValueOf(getCurrentPart(prevEncoding.TypeIds, oTypeIds, partSize, partId, stride)).Interface().([]int),
			Tokens:           reflect.ValueOf(getCurrentPart(prevEncoding.Tokens, oTokens, partSize, partId, stride)).Interface().([]string),
			Offsets:          reflect.ValueOf(getCurrentPart(prevEncoding.Offsets, oOffsets, partSize, partId, stride)).Interface().([][]int),
			SpecialTokenMask: reflect.ValueOf(getCurrentPart(prevEncoding.SpecialTokenMask, oSpeToks, partSize, partId, stride)).Interface().([]int),
			AttentionMask:    reflect.ValueOf(getCurrentPart(prevEncoding.AttentionMask, oAttent, partSize, partId, stride)).Interface().([]int),
			Words:            reflect.ValueOf(getCurrentPart(prevEncoding.Words, oWords, partSize, partId, stride)).Interface().([]int),
			Overflowing:      make([]Encoding, 0),
		}

		partId += 1
		overflowing = append(overflowing, o)
		prevEncoding = &overflowing[len(overflowing)-1]
	}

	e.Overflowing = overflowing

	return e, nil
}

// Merge merges all Encodings together
func (e *Encoding) Merge(encodings []Encoding, growingOffsets bool) (retVal *Encoding) {
	retVal = e
	for _, encoding := range encodings {
		retVal = retVal.MergeWith(&encoding, growingOffsets)
	}

	return retVal
}

// MergeWith merges the current encoding with other (pair) encoding
func (e *Encoding) MergeWith(pair *Encoding, growingOffsets bool) (retVal *Encoding) {
	// Merge overflowing
	var overflowings []Encoding
	var (
		en              Encoding   = *e
		pen             Encoding   = *pair
		enOverflowings  []Encoding = e.Overflowing
		penOverflowings []Encoding = pair.Overflowing
	)
	en.Overflowing = []Encoding{}
	pen.Overflowing = []Encoding{}

	// 1. All our overflowings with all other overflowings
	for _, o := range enOverflowings {
		nEncoding := o.Clone()
		// 1.1. The pair itself
		merge := nEncoding.MergeWith(pair.Clone(), growingOffsets)
		overflowings = append(overflowings, *merge)

		// 1.2. Its overflowings
		for _, otherO := range penOverflowings {
			nEncoding := o.Clone()
			merge := nEncoding.MergeWith(otherO.Clone(), growingOffsets)
			overflowings = append(overflowings, *merge)
		}
	}

	// 2. Ourself with all the other overflowings
	for _, otherO := range penOverflowings {
		nEncoding := e.Clone()
		merge := nEncoding.MergeWith(otherO.Clone(), growingOffsets)
		overflowings = append(overflowings, *merge)
	}

	e.Overflowing = overflowings

	// Merging others
	originalLen := e.Len()
	if len(pair.SequenceRanges) > 0 {
		for seqId, r := range pair.SequenceRanges {
			start := originalLen + r[0]
			end := originalLen + r[r.Len()-1] + 1
			newRange := NewRange(start, end)
			var oldRange Range
			util.DeepCopy(e.SequenceRanges[seqId], oldRange)
			e.SequenceRanges[seqId] = util.Merge(oldRange, newRange)
		}
	}

	e.Ids = util.Merge(e.Ids, pair.Ids)
	e.Tokens = util.Merge(e.Tokens, pair.Tokens)
	e.Words = util.Merge(e.Words, pair.Words)
	e.TypeIds = util.Merge(e.TypeIds, pair.TypeIds)
	e.SpecialTokenMask = util.Merge(e.SpecialTokenMask, pair.SpecialTokenMask)
	e.AttentionMask = util.Merge(e.AttentionMask, pair.AttentionMask)

	// Offsets
	var startingOffset int = 0
	offsets := e.Offsets
	if growingOffsets {
		if len(offsets) > 0 {
			last := offsets[len(offsets)-1]
			startingOffset = last[1]
		}
	}

	for _, o := range pair.Offsets {
		adjustedO := []int{
			o[0] + startingOffset,
			o[1] + startingOffset,
		}
		offsets = append(offsets, adjustedO)
	}
	e.Offsets = offsets

	return e
}

// mergeEncoding merges 2 encodings those have `Overflowing` field empty.
// Otherwise, it will be panic.
func mergeEncoding(en1, en2 Encoding, growingOffsets bool) Encoding {
	if len(en1.Overflowing) > 0 || len(en2.Overflowing) > 0 {
		log.Fatalf("Invalid input encodings. Input encodings must have 'Overflowing' field empty.\n")
	}

	var merge Encoding
	merge.Overflowing = []Encoding{}
	merge.Ids = util.Merge(en1.Ids, en2.Ids)
	merge.TypeIds = util.Merge(en1.TypeIds, en2.TypeIds)
	merge.Words = util.Merge(en1.Words, en2.Words)
	merge.Tokens = util.Merge(en1.Tokens, en2.Tokens)
	merge.SpecialTokenMask = util.Merge(en1.SpecialTokenMask, en2.SpecialTokenMask)
	merge.AttentionMask = util.Merge(en1.AttentionMask, en2.AttentionMask)

	// sequence range
	merge.SequenceRanges = make(map[int]Range)
	sequenceRanges := make(map[int]Range)
	originalLen := en1.Len()
	if len(en2.SequenceRanges) > 0 {
		for seqId, r := range en2.SequenceRanges {
			start := originalLen + r[0]
			end := originalLen + r[r.Len()-1] + 1
			newRange := NewRange(start, end)
			oldRange := en1.SequenceRanges[seqId]
			sequenceRanges[seqId] = append(oldRange, newRange...)
		}
	} else {
		sequenceRanges = en1.SequenceRanges
	}

	merge.SequenceRanges = sequenceRanges

	// Offsets
	offsets := en1.Offsets
	var startingOffset int = 0
	if growingOffsets {
		if len(offsets) > 0 {
			last := offsets[len(offsets)-1]
			startingOffset = last[1]
		}
	}

	for _, o := range en2.Offsets {
		adjustedO := []int{
			o[0] + startingOffset,
			o[1] + startingOffset,
		}
		offsets = append(offsets, adjustedO)
	}
	merge.Offsets = offsets

	return merge
}

// Pad pads current encoding with given length, values to either Left or Right direction
func (e *Encoding) Pad(targetLength, padId, padTypeId int, padToken string, direction PaddingDirection) *Encoding {
	// 1. Overflowing
	var overflowing []Encoding
	for _, o := range e.Overflowing {
		padded := o.pad(targetLength, padId, padTypeId, padToken, direction)
		overflowing = append(overflowing, *padded)
	}
	e.Overflowing = overflowing

	// 2. Check whether we should pad encoding itself
	// if wanted padding length is smaller, then do nothing
	if len(e.Ids) >= targetLength {
		return e
	}

	paddedEn := e.pad(targetLength, padId, padTypeId, padToken, direction)
	return paddedEn
}

func (e *Encoding) pad(targetLength, padId, padTypeId int, padToken string, direction PaddingDirection) *Encoding {
	padLength := targetLength - len(e.Ids)

	switch direction {
	case Left:
		newIds := make([]int, padLength)
		for i := 0; i < len(newIds); i++ {
			newIds[i] = padId
		}
		newIds = append(newIds, e.Ids...)
		e.Ids = newIds

		newTypeIds := make([]int, padLength)
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

		newSpecialTokenMask := make([]int, padLength)
		for i := 0; i < len(newSpecialTokenMask); i++ {
			newSpecialTokenMask[i] = 1
		}
		newSpecialTokenMask = append(newSpecialTokenMask, e.SpecialTokenMask...)
		e.SpecialTokenMask = newSpecialTokenMask

		newAttentionMask := make([]int, padLength)
		for i := 0; i < len(newAttentionMask); i++ {
			newAttentionMask[i] = 0
		}
		newAttentionMask = append(newAttentionMask, e.AttentionMask...)
		e.AttentionMask = newAttentionMask

		newOffsets := make([][]int, padLength)
		for i := 0; i < len(newIds); i++ {
			newOffsets[i] = []int{0, 0}
		}
		newOffsets = append(newOffsets, e.Offsets...)
		e.Offsets = newOffsets

		newWords := make([]int, padLength)
		for i := 0; i < len(newWords); i++ {
			newWords[i] = -1
		}
		newWords = append(newWords, e.Words...)
		e.Words = newWords

	case Right:
		for i := 0; i < padLength; i++ {
			e.Ids = append(e.Ids, padId)
			e.TypeIds = append(e.TypeIds, padTypeId)
			e.Tokens = append(e.Tokens, padToken)
			e.SpecialTokenMask = append(e.SpecialTokenMask, 1)
			e.AttentionMask = append(e.AttentionMask, 0)
			e.Offsets = append(e.Offsets, []int{0, 0})
			e.Words = append(e.Words, -1)
		}
	}

	return e
}

func getCurrentPart(previous, current interface{}, size, idx, stride int) interface{} {

	switch current.(type) {
	case []int:
		var curr, prev []int
		if int((idx+1)*size) > reflect.ValueOf(current).Len() {
			curr = current.([]int)[(idx * size):]
		} else {
			curr = current.([]int)[(idx * size) : (idx+1)*size]
		}
		prev = previous.([]int)[len(previous.([]int))-stride:]
		return append(prev, curr...)
	case []string:
		var curr, prev []string
		if (idx+1)*size > reflect.ValueOf(current).Len() {
			curr = current.([]string)[(idx * size):]
		} else {
			curr = current.([]string)[(idx * size) : (idx+1)*size]
		}
		prev = previous.([]string)[len(previous.([]string))-stride:]
		return append(prev, curr...)
	case [][]int:
		var curr, prev [][]int
		if (idx+1)*size > reflect.ValueOf(current).Len() {
			curr = current.([][]int)[(idx * size):]
		} else {
			curr = current.([][]int)[(idx * size) : (idx+1)*size]
		}
		prev = previous.([][]int)[len(previous.([][]int))-stride:]
		return append(prev, curr...)
	default:
		log.Fatalf("getCurrentPart method call: invalid type\n")
	}

	return nil
}

// Token2Sequence returns the index of the sequence containing the given token.
func (e *Encoding) Token2Sequence(token int) (int, bool) {
	if token > e.Len() {
		return -1, false
	} else if len(e.SequenceRanges) == 0 {
		return 0, true
	} else {
		for seqId, r := range e.SequenceRanges {
			if r.Contains(token) {
				return seqId, true
			}
		}

		return -1, false
	}
}

// SequenceRange returns the range to target to retrieve something (word id, offsets, ...)
// related to the given sequence id.
func (e *Encoding) SequenceRange(sequencId int) (Range, error) {
	r, ok := e.SequenceRanges[sequencId]
	if !ok {
		err := fmt.Errorf("input 'sequence_id' is out of range.\n")
		return nil, err
	}

	return r[0:e.Len()], nil
}
