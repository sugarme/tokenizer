package wordpiece_test

import (
	"reflect"
	"testing"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/model/wordpiece"
)

func TestWordpieceBuilder(t *testing.T) {
	b := wordpiece.NewWordPieceBuilder()

	wp := b.Build()

	vocab := wp.GetVocab()
	got := len(vocab)
	want := 0

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}
}

func TestWordpieceFromFile(t *testing.T) {
	vocabFile, err := tokenizer.CachedPath("bert-base-uncased", "vocab.txt")
	if err != nil {
		t.Fail()
	}

	unkToken := "[UNK]"
	m, err := wordpiece.NewWordPieceFromFile(vocabFile, unkToken)
	if err != nil {
		t.Fail()
	}

	got := m.GetVocabSize()
	want := 30_522
	if want != got {
		t.Errorf("Want %v, got %v\n", want, got)
	}
}

func TestWordpieceTokenize(t *testing.T) {
	vocabFile, err := tokenizer.CachedPath("bert-base-uncased", "vocab.txt")
	if err != nil {
		t.Fail()
	}

	unkToken := "[UNK]"
	m, err := wordpiece.NewWordPieceFromFile(vocabFile, unkToken)
	if err != nil {
		t.Fail()
	}

	// NOTE. just a lower-case word (as testing purely the model without normalizer, pretokenizer, processor)
	seq := "gopher"

	got, err := m.Tokenize(seq)
	if err != nil {
		t.Fail()
	}

	want := []tokenizer.Token{
		{Id: 2175, Value: "go", Offsets: []int{0, 2}},
		{Id: 27921, Value: "##pher", Offsets: []int{2, 6}},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("\nwant %q,\ngot  %+v", want, got)
	}
}
