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

func NewRobertaProcessing(sep, cls PostToken) *RobertaProcessing {
	r := DefaultRobertaProcessing()
	r.sep = sep
	r.cls = cls

	return r
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

	var ids []int
	ids = append(ids, rp.cls.Id)
	ids = append(ids, newEncoding.GetIds()...)
	ids = append(ids, rp.sep.Id)

	var typeIds []int
	typeIds = append(typeIds, 0)
	typeIds = append(typeIds, newEncoding.GetTypeIds()...)
	typeIds = append(typeIds, 0)

	var tokens []string
	tokens = append(tokens, rp.cls.Value)
	tokens = append(tokens, newEncoding.GetTokens()...)
	tokens = append(tokens, rp.sep.Value)

	var words []int
	words = append(words, -1)
	words = append(words, newEncoding.GetWords()...)
	words = append(words, -1)

	var offsets [][]int
	offsets = append(offsets, []int{0, 0})
	offsets = append(offsets, newEncoding.GetOffsets()...)
	offsets = append(offsets, []int{0, 0})

	var specialTokens []int
	specialTokens = append(specialTokens, 1)
	specialTokens = append(specialTokens, 0)
	specialTokens = append(specialTokens, len(newEncoding.GetIds()))
	specialTokens = append(specialTokens, 1)

	var attentionMask []int = []int{1, len(ids)}

	finalEncoding := tokenizer.NewEncoding(ids, typeIds, tokens, offsets, specialTokens, attentionMask, newEncoding.TakeOverflowing(), words)

	if pairEncoding != nil {
		var pairIds []int
		pairIds = append(pairIds, newPairEncoding.GetTypeIds()...)
		pairIds = append(pairIds, rp.sep.Id)

		var pairTypeIds []int
		pairTypeIds = append(pairTypeIds, newPairEncoding.GetTypeIds()...)
		pairTypeIds = append(pairTypeIds, 1)

		var pairTokens []string
		pairTokens = append(pairTokens, newPairEncoding.GetTokens()...)
		pairTokens = append(pairTokens, rp.sep.Value)

		var pairWords []int
		pairWords = append(pairWords, newPairEncoding.GetWords()...)
		pairWords = append(pairWords, -1)

		var pairOffsets [][]int
		pairOffsets = append(pairOffsets, newPairEncoding.GetOffsets()...)
		pairOffsets = append(pairOffsets, []int{0, 0})

		var pairSpecialTokens []int
		pairSpecialTokens = append(pairSpecialTokens, 0)
		pairSpecialTokens = append(pairSpecialTokens, len(newPairEncoding.GetTypeIds()))
		pairSpecialTokens = append(pairSpecialTokens, 1)

		var pairAttentionMask []int
		pairAttentionMask = append(pairAttentionMask, 1)
		pairAttentionMask = append(pairAttentionMask, len(pairIds))

		finalPairEncoding := tokenizer.NewEncoding(pairIds, pairTypeIds, pairTokens, pairOffsets, pairSpecialTokens, pairAttentionMask, newPairEncoding.TakeOverflowing(), pairWords)

		// Merge with pair
		finalEncoding.MergeWith(finalPairEncoding, false)

	}

	return finalEncoding
}

// TODO: implement Serialize interface for RobertaProcessing
