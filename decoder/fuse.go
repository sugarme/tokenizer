package decoder

import (
	"strings"

	"github.com/sugarme/tokenizer"
)

// Fuse constructs Fuse decoder
// It's simply fuses all tokens into one big string.
type Fuse struct {
	*DecoderBase
}

func NewFuse() *Fuse {
	base := new(DecoderBase)

	d := &Fuse{base}

	d.DecoderBase.Decoder = interface{}(d).(tokenizer.Decoder)

	return d
}

func (f *Fuse) DecodeChain(tokens []string) []string {
	str := strings.Join(tokens, "")

	return []string{str}
}
