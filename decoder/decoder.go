package decoder

import (
	"strings"
)

type DecoderBase struct{}

func (d *DecoderBase) Decode(tokens []string) string {
	return strings.Join(tokens, "")
}
