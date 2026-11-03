package decoder

import (
	"strings"
	"testing"

	"github.com/sugarme/tokenizer"
)

// Decode on a Default* decoder used to panic with a nil DecoderBase.
func TestDefaultDecodersDecode(t *testing.T) {
	tokens := []string{"he", "##llo", "</w>", "<pad>", "|"}
	for name, dec := range map[string]tokenizer.Decoder{
		"wordpiece": DefaultWordpieceDecoder(),
		"bpe":       DefaultBpeDecoder(),
		"ctc":       DefaultCTC(),
	} {
		want := strings.Join(dec.DecodeChain(tokens), "")
		if got := dec.Decode(tokens); got != want {
			t.Errorf("%s: Decode = %q, want %q", name, got, want)
		}
	}
}
