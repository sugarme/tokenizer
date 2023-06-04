package decoder

import (
	"strings"
)

// Fuse constructs Fuse decoder
// It's simply fuses all tokens into one big string.
type Fuse struct {
	*DecoderBase
}

func NewFuse() *Fuse {
	base := new(DecoderBase)

	return &Fuse{base}
}

func (f *Fuse) DecodeChain(tokens []string) []string {
	str := strings.Join(tokens, "")

	return []string{str}
}
