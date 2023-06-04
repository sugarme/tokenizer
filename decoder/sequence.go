package decoder

import (
	"github.com/sugarme/tokenizer"
)

type Sequence struct {
	*DecoderBase

	decoders []tokenizer.Decoder
}

var _ tokenizer.Decoder = new(Sequence)

func NewSequence(decoders []tokenizer.Decoder) *Sequence {
	return &Sequence{
		decoders: decoders,
	}
}

// Decode implements `tokenizer.Decoder` interface.
func (d *Sequence) DecodeChain(tokens []string) []string {
	var toks = tokens
	for _, dec := range d.decoders {
		toks = dec.DecodeChain(toks)
	}

	return toks
}
