package pretokenizer

import (
	"reflect"
	"testing"

	"github.com/season-studio/tokenizer"
	"github.com/season-studio/tokenizer/normalizer"
)

func TestPunctuation(t *testing.T) {
	pretok := DefaultPunctuation()

	pretokenized := tokenizer.NewPreTokenizedString("Hey friend!     How are you?!?")

	out, err := pretok.PreTokenize(pretokenized)
	if err != nil {
		panic(err)
	}

	got := out.GetSplits(normalizer.OriginalTarget, tokenizer.Byte)

	want := []tokenizer.PreToken{
		{Value: "Hey friend", Offsets: []int{0, 10}, Tokens: nil},
		{Value: "!", Offsets: []int{10, 11}, Tokens: nil},
		{Value: "     How are you", Offsets: []int{11, 27}, Tokens: nil},
		{Value: "?", Offsets: []int{27, 28}, Tokens: nil},
		{Value: "!", Offsets: []int{28, 29}, Tokens: nil},
		{Value: "?", Offsets: []int{29, 30}, Tokens: nil},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %#v\ngot %#v\n", want, got)
	}
}
