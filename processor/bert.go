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

// Process post-processes input encoding(s) by adding special tokens if specifying.
func (bp *BertProcessing) Process(encoding, pairEncoding *tokenizer.Encoding, addSpecialTokens bool) (retVal *tokenizer.Encoding) {
	if !addSpecialTokens {
		return tokenizer.DefaultProcess(encoding, pairEncoding, addSpecialTokens)
	}

	// add special token "[CLS]" and "[SEP]" to encoding itself
	newEncoding := bp.addSpecialToken(encoding)

	// add special token "[CLS]" and "[SEP]" to overflowing encodings if it has
	var overflowing []tokenizer.Encoding
	for _, en := range encoding.Overflowing {
		newEn := bp.addSpecialToken(&en)
		overflowing = append(overflowing, *newEn)
	}
	newEncoding.Overflowing = overflowing

	if pairEncoding != nil {

		// Add special token "[SEP]" to pair encoding itself
		newPairEncoding := bp.pairAddSpecialToken(pairEncoding)

		// Add special token "[SEP]" at the end of its overflowing if it has
		var pairOverflowing []tokenizer.Encoding
		for _, en := range pairEncoding.Overflowing {
			newEnc := bp.pairAddSpecialToken(&en)
			pairOverflowing = append(pairOverflowing, *newEnc)
		}
		newPairEncoding.Overflowing = pairOverflowing

		// Merge newPairEncoding with newEncoding
		newEncoding.MergeWith(newPairEncoding, false)
	}

	return newEncoding
}

// addSpecialToken adds special token "[CLS]" and "[SEP]" to input encoding. It ignores
// `Overflowing` field.
func (bp *BertProcessing) addSpecialToken(encoding *tokenizer.Encoding) *tokenizer.Encoding {
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

	var words []int
	words = append(words, -1)
	words = append(words, encoding.GetWords()...)
	words = append(words, -1)

	var offsets [][]int
	offsets = append(offsets, []int{0, 0})
	offsets = append(offsets, encoding.GetOffsets()...)
	offsets = append(offsets, []int{0, 0})

	var specialTokens []int
	specialTokens = append(specialTokens, 1)
	for i := 0; i < len(encoding.Ids); i++ {
		specialTokens = append(specialTokens, 0)
	}
	specialTokens = append(specialTokens, 1)

	// As all tokens are non-padded tokens, just assign 1
	nonPaddedTokens := len(encoding.Ids) + 2
	var attentionMask []int
	for i := 0; i < nonPaddedTokens; i++ {
		attentionMask = append(attentionMask, 1)
	}

	wordsOpt := tokenizer.WithWordsEncodingOpt(words)
	return tokenizer.NewEncoding(ids, typeIds, tokens, offsets, specialTokens, attentionMask, []tokenizer.Encoding{}, wordsOpt)
}

// pairAddSpecialToken adds special token "[SEP]" to input encoding. It ignores
// `Overflowing` field.
func (bp *BertProcessing) pairAddSpecialToken(pairEncoding *tokenizer.Encoding) *tokenizer.Encoding {
	var pairIds []int
	pairIds = append(pairIds, pairEncoding.Ids...)
	pairIds = append(pairIds, bp.sep.Id)

	var pairTypeIds []int
	pairTypeIds = append(pairTypeIds, pairEncoding.GetTypeIds()...)
	pairTypeIds = append(pairTypeIds, 1)

	var pairTokens []string
	pairTokens = append(pairTokens, pairEncoding.GetTokens()...)
	pairTokens = append(pairTokens, bp.sep.Value)

	var pairWords []int
	pairWords = append(pairWords, pairEncoding.GetWords()...)
	pairWords = append(pairWords, -1)

	var pairOffsets [][]int
	pairOffsets = append(pairOffsets, pairEncoding.GetOffsets()...)
	pairOffsets = append(pairOffsets, []int{0, 0})

	var pairSpecialTokens []int
	for i := 0; i < len(pairEncoding.Ids); i++ {
		pairSpecialTokens = append(pairSpecialTokens, 0)
	}
	pairSpecialTokens = append(pairSpecialTokens, 1)

	var pairAttentionMask []int
	for i := 0; i < len(pairEncoding.Ids); i++ {
		pairAttentionMask = append(pairAttentionMask, 1)
	}
	pairAttentionMask = append(pairAttentionMask, 1)

	pairWordsOpt := tokenizer.WithWordsEncodingOpt(pairWords)

	return tokenizer.NewEncoding(pairIds, pairTypeIds, pairTokens, pairOffsets, pairSpecialTokens, pairAttentionMask, []tokenizer.Encoding{}, pairWordsOpt)
}
