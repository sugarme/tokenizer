package processor

import (
	"github.com/sugarme/tokenizer"
)

type PostToken struct {
	Value string
	Id    int
}

type BertProcessing struct {
	sep PostToken
	cls PostToken
}

func NewBertProcessing(sep, cls PostToken) (retVal *BertProcessing) {
	return &BertProcessing{
		sep: sep,
		cls: cls,
	}
}

// Implement PostProcessor interface for BertProcessing:
// =====================================================

func (bp *BertProcessing) AddedTokens(isPair bool) (retVal int) {
	if isPair {
		return 3
	} else {
		return 2
	}
}

func (bp *BertProcessing) Process(encoding, pairEncoding *tokenizer.Encoding, addSpecialTokens bool) (retVal *tokenizer.Encoding) {

	if !addSpecialTokens {
		return tokenizer.DefaultProcess(encoding, pairEncoding, addSpecialTokens)
	}

	var ids []int
	ids = append(ids, bp.cls.Id)
	ids = append(ids, encoding.GetIds()...)
	ids = append(ids, bp.sep.Id)

	var typeIds []int
	typeIds = append(typeIds, 0)
	typeIds = append(typeIds, encoding.GetTypeIds()...)
	typeIds = append(typeIds, 0)

	var tokens []string
	tokens = append(tokens, bp.cls.Value)
	tokens = append(tokens, encoding.GetTokens()...)
	tokens = append(tokens, bp.sep.Value)

	var offsets []tokenizer.Offsets
	offsets = append(offsets, tokenizer.Offsets{Start: 0, End: 0})
	offsets = append(offsets, encoding.GetOffsets()...)
	offsets = append(offsets, tokenizer.Offsets{Start: 0, End: 0})

	var specialTokens []int
	specialTokens = append(specialTokens, 1)
	specialTokens = append(specialTokens, 0)
	specialTokens = append(specialTokens, len(encoding.GetIds()))
	specialTokens = append(specialTokens, 1)

	var attentionMask []int = []int{1, len(ids)}

	newEncoding := tokenizer.NewEncoding(ids, typeIds, tokens, offsets, specialTokens, attentionMask, encoding.TakeOverflowing())

	if pairEncoding != nil {
		var pairIds []int
		pairIds = append(pairIds, pairEncoding.GetTypeIds()...)
		pairIds = append(pairIds, bp.sep.Id)

		var pairTypeIds []int
		pairTypeIds = append(pairTypeIds, pairEncoding.GetTypeIds()...)
		pairTypeIds = append(pairTypeIds, 1)

		var pairTokens []string
		pairTokens = append(pairTokens, pairEncoding.GetTokens()...)
		pairTokens = append(pairTokens, bp.sep.Value)

		var pairOffsets []tokenizer.Offsets
		pairOffsets = append(pairOffsets, pairEncoding.GetOffsets()...)
		pairOffsets = append(pairOffsets, tokenizer.Offsets{Start: 0, End: 0})

		var pairSpecialTokens []int
		pairSpecialTokens = append(pairSpecialTokens, 0)
		pairSpecialTokens = append(pairSpecialTokens, len(pairEncoding.GetTypeIds()))
		pairSpecialTokens = append(pairSpecialTokens, 1)

		var pairAttentionMask []int
		pairAttentionMask = append(pairAttentionMask, 1)
		pairAttentionMask = append(pairAttentionMask, len(pairIds))

		newPairEncoding := tokenizer.NewEncoding(pairIds, pairTypeIds, pairTokens, pairOffsets, pairSpecialTokens, pairAttentionMask, pairEncoding.TakeOverflowing())

		// Merge newPairEncoding with newEncoding
		newEncoding.MergeWith(newPairEncoding, false)

	}

	return newEncoding
}
