package pretrained

import (
	"strings"
	"testing"

	"github.com/sugarme/tokenizer"
)

const benchText = `The quick brown fox jumps over the lazy dog. In 2023, researchers at the Université de Montréal ` +
	`published a paper titled "Efficient Tokenization for Large-Scale Language Models", which reported a 3.5x speed-up. ` +
	`Don't forget: tokenizers split text into sub-word units; e.g. "unbelievably" -> "un", "##believ", "##ably". ` +
	`Prices rose from $1,299.99 to $1,499.00 (an increase of ~15%) over the past twelve months! `

func benchTokenizer(b *testing.B, model string) *tokenizer.Tokenizer {
	b.Helper()
	f, err := tokenizer.CachedPath(model, "tokenizer.json")
	if err != nil {
		b.Skipf("cannot get %s: %v", model, err)
	}
	tk, err := FromFile(f)
	if err != nil {
		b.Fatal(err)
	}
	return tk
}

func BenchmarkEncode(b *testing.B) {
	long := strings.Repeat(benchText, 20)
	// shortcut: "bench-xlmr" is xlm-roberta-base's tokenizer.json with `"type": "Unigram"` added to "model",
	// put by hand in the cache dir (it skips if absent). Use "xlm-roberta-base" once FromFile detects Unigram without "type".
	for _, model := range []string{"bert-base-uncased", "gpt2", "bench-xlmr"} {
		tk := benchTokenizer(b, model)
		for _, in := range []struct {
			name string
			text string
		}{{"short", benchText}, {"long", long}} {
			b.Run(model+"/"+in.name, func(b *testing.B) {
				b.SetBytes(int64(len(in.text)))
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					if _, err := tk.EncodeSingle(in.text, true); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}
