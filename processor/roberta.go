package processor

import (
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretokenizer"
)

// RobertaProcessing is a post post processor for Roberta model
type RobertaProcessing struct {
	sep            PostToken
	cls            PostToken
	trimOffsets    bool
	addPrefixSpace bool
}

// DefaultRobertaProcessing creates a RobertaProcessing with default values
func DefaultRobertaProcessing() *RobertaProcessing {

	return &RobertaProcessing{
		sep:            PostToken{Value: "</s>", Id: 2},
		cls:            PostToken{Value: "<s>", Id: 0},
		trimOffsets:    true,
		addPrefixSpace: true,
	}
}

func NewRobertaProcessing(sep, cls PostToken, trimOffsets bool, addPrefixSpace bool) *RobertaProcessing {
	return &RobertaProcessing{
		sep:            sep,
		cls:            cls,
		trimOffsets:    trimOffsets,
		addPrefixSpace: addPrefixSpace,
	}
}

// TrimOffsets set whether the processor will trim offsets
func (rp *RobertaProcessing) TrimOffsets(trimOffsets bool) {
	rp.trimOffsets = trimOffsets
}

// AddPrefixSpace set whether the processor will add a prefix space
func (rp *RobertaProcessing) AddPrefixSpace(addPrefixSpace bool) {
	rp.addPrefixSpace = addPrefixSpace
}

// Implement PostProcessor interface for RobertaProcessing:
// ========================================================

func (rp *RobertaProcessing) AddedTokens(isPair bool) int {
	if isPair {
		return 4
	} else {
		return 2
	}
}

// Process post-processes input encoding(s) by adding special tokens if instructed to do so.
//
// Specifically, if addSpecialToken=true, it will add special tokens patterns
// - Single encoding: <s> Sequence </s>
// - Pair encoding: <s> SequenceA </s> </s> SequenceB </s>
func (rp *RobertaProcessing) Process(encoding, pairEncoding *tokenizer.Encoding, addSpecialTokens bool) *tokenizer.Encoding {

	var (
		newEncoding             *tokenizer.Encoding
		newOverflowEncodings    []tokenizer.Encoding
		newPairEncoding         *tokenizer.Encoding
		newOverflowPairEncoding []tokenizer.Encoding
	)
	if rp.trimOffsets {
		newEncoding = pretokenizer.ProcessOffsets(encoding, rp.addPrefixSpace)

		overflowEncodings := newEncoding.GetOverflowing()
		for _, e := range overflowEncodings {
			newEn := pretokenizer.ProcessOffsets(&e, rp.addPrefixSpace)
			newOverflowEncodings = append(newOverflowEncodings, *newEn)
		}
		newEncoding.Overflowing = newOverflowEncodings

		if pairEncoding != nil {
			newPairEncoding = pretokenizer.ProcessOffsets(pairEncoding, rp.addPrefixSpace)
			for _, en := range newPairEncoding.Overflowing {
				newEn := pretokenizer.ProcessOffsets(&en, rp.addPrefixSpace)
				newOverflowPairEncoding = append(newOverflowPairEncoding, *newEn)
			}
			newPairEncoding.Overflowing = newOverflowPairEncoding
		}

	}

	if !addSpecialTokens {
		return tokenizer.DefaultProcess(newEncoding, newPairEncoding, addSpecialTokens)
	}

	// add special token to itself
	finalEncoding := rp.addSpecialToken(newEncoding)

	// add special token to its overflowing
	var overflowing []tokenizer.Encoding
	for _, en := range newEncoding.Overflowing {
		newEn := rp.addSpecialToken(&en)
		overflowing = append(overflowing, *newEn)
	}
	finalEncoding.Overflowing = overflowing

	if pairEncoding != nil {

		// add special tokens for pair itself
		finalPairEncoding := rp.pairAddSpecialToken(newPairEncoding)
		// add special tokens for pair's overflowing
		var pairOverflowing []tokenizer.Encoding
		for _, en := range newPairEncoding.Overflowing {
			newEn := rp.pairAddSpecialToken(&en)
			pairOverflowing = append(pairOverflowing, *newEn)
		}
		finalPairEncoding.Overflowing = pairOverflowing

		// Merge with pair
		finalEncoding.MergeWith(finalPairEncoding, false)

	}

	return finalEncoding
}

// addSpecialToken adds special tokens to input encoding. It ignores the `Overflowing` field
// of input encoding.
//
// Specifically, it adds: <s> Sequence </s>
func (rp *RobertaProcessing) addSpecialToken(encoding *tokenizer.Encoding) *tokenizer.Encoding {
	var ids []int
	ids = append(ids, rp.cls.Id)
	ids = append(ids, encoding.Ids...)
	ids = append(ids, rp.sep.Id)

	var typeIds []int
	typeIds = append(typeIds, 0)
	typeIds = append(typeIds, encoding.TypeIds...)
	typeIds = append(typeIds, 0)

	var tokens []string
	tokens = append(tokens, rp.cls.Value)
	tokens = append(tokens, encoding.Tokens...)
	tokens = append(tokens, rp.sep.Value)

	var words []int
	words = append(words, -1)
	words = append(words, encoding.Words...)
	words = append(words, -1)

	var offsets [][]int
	offsets = append(offsets, []int{0, 0})
	offsets = append(offsets, encoding.Offsets...)
	offsets = append(offsets, []int{0, 0})

	var specialTokens []int
	specialTokens = append(specialTokens, 1)
	for i := 0; i < len(encoding.SpecialTokenMask); i++ {
		specialTokens = append(specialTokens, 0)
	}
	specialTokens = append(specialTokens, 1)

	var attentionMask []int
	attentionMask = append(attentionMask, 1)
	for i := 0; i < len(encoding.AttentionMask); i++ {
		attentionMask = append(attentionMask, 1)
	}
	attentionMask = append(attentionMask, 1)

	wordsOpt := tokenizer.WithWordsEncodingOpt(words)
	return tokenizer.NewEncoding(ids, typeIds, tokens, offsets, specialTokens, attentionMask, []tokenizer.Encoding{}, wordsOpt)
}

// addSpecialToken adds special tokens to input pair encoding. It ignores the `Overflowing` field
// of input pair encoding.
//
// Specifically, it adds </s> PairEncoding </s>
func (rp *RobertaProcessing) pairAddSpecialToken(pair *tokenizer.Encoding) *tokenizer.Encoding {
	var pairIds []int
	pairIds = append(pairIds, rp.sep.Id)
	pairIds = append(pairIds, pair.Ids...)
	pairIds = append(pairIds, rp.sep.Id)

	var pairTypeIds []int
	pairTypeIds = append(pairTypeIds, 1)
	pairTypeIds = append(pairTypeIds, pair.TypeIds...)
	pairTypeIds = append(pairTypeIds, 1)

	var pairTokens []string
	pairTokens = append(pairTokens, rp.sep.Value)
	pairTokens = append(pairTokens, pair.Tokens...)
	pairTokens = append(pairTokens, rp.sep.Value)

	var pairWords []int
	pairWords = append(pairWords, -1)
	pairWords = append(pairWords, pair.Words...)
	pairWords = append(pairWords, -1)

	var pairOffsets [][]int
	pairOffsets = append(pairOffsets, []int{0, 0})
	pairOffsets = append(pairOffsets, pair.Offsets...)
	pairOffsets = append(pairOffsets, []int{0, 0})

	var pairSpecialTokens []int
	pairSpecialTokens = append(pairSpecialTokens, 1)
	for i := 0; i < len(pair.SpecialTokenMask); i++ {
		pairSpecialTokens = append(pairSpecialTokens, 0)
	}
	pairSpecialTokens = append(pairSpecialTokens, 1)

	var pairAttentionMask []int
	pairAttentionMask = append(pairAttentionMask, 1)
	for i := 0; i < len(pair.AttentionMask); i++ {
		pairAttentionMask = append(pairAttentionMask, 1)
	}
	pairAttentionMask = append(pairAttentionMask, 1)

	pairWordsOpt := tokenizer.WithWordsEncodingOpt(pairWords)
	return tokenizer.NewEncoding(pairIds, pairTypeIds, pairTokens, pairOffsets, pairSpecialTokens, pairAttentionMask, []tokenizer.Encoding{}, pairWordsOpt)
}

// TODO: implement Serialize interface for RobertaProcessing
