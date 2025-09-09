package processor

import (
	"github.com/gengzongjie/tokenizer"
	"github.com/gengzongjie/tokenizer/pretokenizer"
)

type ByteLevelProcessing struct {
	pretok *pretokenizer.ByteLevel
}

func NewByteLevelProcessing(pretok *pretokenizer.ByteLevel) (retVal *ByteLevelProcessing) {
	return &ByteLevelProcessing{
		pretok: pretok,
	}
}

// Implement PostProcessor interface for ByteLevelProcessing:
// =====================================================

func (blp *ByteLevelProcessing) AddedTokens(isPair bool) (retVal int) {
	return blp.pretok.AddedToken(isPair)
}

func (blp *ByteLevelProcessing) Process(encoding, pairEncoding *tokenizer.Encoding, addSpecialTokens bool) (retVal *tokenizer.Encoding) {
	return blp.pretok.Process(encoding, pairEncoding, addSpecialTokens)
}
