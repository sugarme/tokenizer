package tokenizer_test

import (
	"testing"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
)

func TestBatchKeepsOrder(t *testing.T) {
	tk := pretrained.BertBaseUncased()
	sentences := []string{"hello world", "good morning", "the quick brown fox", "see you later", "thank you"}

	var inputs []tokenizer.EncodeInput
	for i := 0; i < 500; i++ {
		inputs = append(inputs, tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(sentences[i%len(sentences)])))
	}
	encodings, err := tk.EncodeBatch(inputs, false)
	if err != nil {
		t.Fatal(err)
	}

	ids := make([][]int, len(encodings))
	for i := range encodings {
		ids[i] = encodings[i].GetIds()
	}
	got := tk.DecodeBatch(ids, true)
	if len(got) != len(inputs) {
		t.Fatalf("DecodeBatch returned %d results, want %d", len(got), len(inputs))
	}
	for i, s := range got {
		if want := sentences[i%len(sentences)]; s != want {
			t.Fatalf("result %d = %q, want %q", i, s, want)
		}
	}
}
