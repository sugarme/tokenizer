package processor

import "github.com/sugarme/tokenizer"

type Sequence struct {
	processors []tokenizer.PostProcessor
}

var _ tokenizer.PostProcessor = new(Sequence)

func NewSequence(processors []tokenizer.PostProcessor) *Sequence {
	return &Sequence{processors}
}

// Implement tokenizer.PostProcessor for Sequence

func (seq *Sequence) AddedTokens(isPair bool) (retVal int) {
	var count int
	for _, p := range seq.processors {
		count += p.AddedTokens(isPair)
	}

	return count
}

func (seq *Sequence) Process(encoding, pairEncoding *tokenizer.Encoding, addSpecialTokens bool) (retVal *tokenizer.Encoding) {
	// return blp.pretok.Process(encoding, pairEncoding, addSpecialTokens)
	var encodings *tokenizer.Encoding = encoding
	for _, p := range seq.processors {
		encodings = p.Process(encodings, pairEncoding, addSpecialTokens)
	}

	return encodings
}
