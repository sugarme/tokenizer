package decoder

import (
	"reflect"
	// "strings"
	"testing"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretokenizer"
)

func TestSequence(t *testing.T) {
	decoders := []tokenizer.Decoder{
		DefaultCTC(),
		pretokenizer.DefaultMetaspace(),
	}

	dec := NewSequence(decoders)

	tokens := []string{"▁", "▁", "H", "H", "i", "i", "▁", "y", "o", "u"}

	// out := dec.DecodeChain(tokens)
	// got := strings.Join(out, "")
	got := dec.Decode(tokens)
	want := "Hi you"

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %q, got %q", want, got)
	}
}

/*
fn sequence_basic() {
        let decoders = vec![
            DecoderWrapper::CTC(CTC::default()),
            DecoderWrapper::Metaspace(Metaspace::default()),
        ];
        let decoder = Sequence::new(decoders);
        let tokens: Vec<String> = vec!["▁", "▁", "H", "H", "i", "i", "▁", "y", "o", "u"]
            .into_iter()
            .map(|s| s.to_string())
            .collect();
        let out_tokens = decoder.decode(tokens).unwrap();
        assert_eq!(out_tokens, "Hi you");
    }
*/
