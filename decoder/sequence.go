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
	base := new(DecoderBase)

	seq := &Sequence{
		DecoderBase: base,
		decoders:    decoders,
	}

	seq.DecoderBase.Decoder = interface{}(seq).(tokenizer.Decoder)

	return seq
}

// Decode implements `tokenizer.Decoder` interface.
func (d *Sequence) DecodeChain(tokens []string) []string {
	var input []string
	for _, token := range tokens {
		input = append(input, token)
	}
	for _, dec := range d.decoders {
		tmp := dec.DecodeChain(input)
		input = []string{}
		input = tmp
	}

	return input
}
