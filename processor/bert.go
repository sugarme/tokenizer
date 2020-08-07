package processor

import (
	"github.com/sugarme/tokenizer/tokenizer"
)

type PostToken struct {
	Value string
	Id    uint32
}

type BertProcessing struct {
	sep PostToken
	cls PostToken
}

func NewBertProcessing(sep, cls PostToken) (retVal BertProcessing) {
	return BertProcessing{
		sep: sep,
		cls: cls,
	}
}

// Implement PostProcessor interface for BertProcessing:
// =====================================================

func (bp BertProcessing) AddedTokens(isPair bool) (retVal int) {
	if isPair {
		return 3
	} else {
		return 2
	}
}

func (bp BertProcessing) Process(encoding tokenizer.Encoding, pairEncodingOpt ...tokenizer.Encoding) (retVal tokenizer.Encoding) {

	var ids []uint32
	ids = append(ids, bp.cls.Id)
	ids = append(ids, encoding.GetIds()...)
	ids = append(ids, bp.sep.Id)

	var typeIds []uint32
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

	var specialTokens []uint32
	specialTokens = append(specialTokens, 1)
	specialTokens = append(specialTokens, 0)
	specialTokens = append(specialTokens, uint32(len(encoding.GetIds())))
	specialTokens = append(specialTokens, 1)

	var attentionMask []uint32 = []uint32{1, uint32(len(ids))}

	newEncoding := tokenizer.NewEncoding(encoding.GetNormalized(), ids, typeIds, tokens, offsets, specialTokens, attentionMask, encoding.TakeOverflowing())

	var pairEncoding tokenizer.Encoding
	if len(pairEncodingOpt) > 0 {
		pairEncoding = pairEncodingOpt[0]

		var pairIds []uint32
		pairIds = append(pairIds, pairEncoding.GetTypeIds()...)
		pairIds = append(pairIds, bp.sep.Id)

		var pairTypeIds []uint32
		pairTypeIds = append(pairTypeIds, pairEncoding.GetTypeIds()...)
		pairTypeIds = append(pairTypeIds, 1)

		var pairTokens []string
		pairTokens = append(pairTokens, pairEncoding.GetTokens()...)
		pairTokens = append(pairTokens, bp.sep.Value)

		var pairOffsets []tokenizer.Offsets
		pairOffsets = append(pairOffsets, pairEncoding.GetOffsets()...)
		pairOffsets = append(pairOffsets, tokenizer.Offsets{Start: 0, End: 0})

		var pairSpecialTokens []uint32
		pairSpecialTokens = append(pairSpecialTokens, 0)
		pairSpecialTokens = append(pairSpecialTokens, uint32(len(pairEncoding.GetTypeIds())))
		pairSpecialTokens = append(pairSpecialTokens, 1)

		var pairAttentionMask []uint32
		pairAttentionMask = append(pairAttentionMask, 1)
		pairAttentionMask = append(pairAttentionMask, uint32(len(pairIds)))

		newPairEncoding := tokenizer.NewEncoding(pairEncoding.GetNormalized(), pairIds, pairTypeIds, pairTokens, pairOffsets, pairSpecialTokens, pairAttentionMask, pairEncoding.TakeOverflowing())

		// Merge newPairEncoding with newEncoding
		newEncoding.MergeWith(newPairEncoding)

	}

	return newEncoding
}
