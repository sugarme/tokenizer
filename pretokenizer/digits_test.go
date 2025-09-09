package pretokenizer

import (
	"reflect"
	"testing"

	"github.com/gengzongjie/tokenizer"
	"github.com/gengzongjie/tokenizer/normalizer"
)

func TestNumbers(t *testing.T) {
	pretok := NewDigits(false)
	pretokenized := tokenizer.NewPreTokenizedString("Hey 123 friend!")

	out, err := pretok.PreTokenize(pretokenized)
	if err != nil {
		t.Fail()
	}

	got := out.GetSplits(normalizer.NormalizedTarget, tokenizer.Byte)
	want := []tokenizer.PreToken{
		{Value: "Hey ", Offsets: []int{0, 4}, Tokens: nil},
		{Value: "123", Offsets: []int{4, 7}, Tokens: nil},
		{Value: " friend!", Offsets: []int{7, 15}, Tokens: nil},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %#v\ngot %#v\n", want, got)
	}

	got = out.GetSplits(normalizer.OriginalTarget, tokenizer.Byte)
	want = []tokenizer.PreToken{
		{Value: "Hey ", Offsets: []int{0, 4}, Tokens: nil},
		{Value: "123", Offsets: []int{4, 7}, Tokens: nil},
		{Value: " friend!", Offsets: []int{7, 15}, Tokens: nil},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %#v\ngot %#v\n", want, got)
	}
}

func TestIndividualDigits(t *testing.T) {
	pretok := NewDigits(true)
	pretokenized := tokenizer.NewPreTokenizedString("Hey 123 friend!")

	out, err := pretok.PreTokenize(pretokenized)
	if err != nil {
		t.Fail()
	}

	got := out.GetSplits(normalizer.NormalizedTarget, tokenizer.Byte)
	want := []tokenizer.PreToken{
		{Value: "Hey ", Offsets: []int{0, 4}, Tokens: nil},
		{Value: "1", Offsets: []int{4, 5}, Tokens: nil},
		{Value: "2", Offsets: []int{5, 6}, Tokens: nil},
		{Value: "3", Offsets: []int{6, 7}, Tokens: nil},
		{Value: " friend!", Offsets: []int{7, 15}, Tokens: nil},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %#v\ngot %#v\n", want, got)
	}

	got = out.GetSplits(normalizer.OriginalTarget, tokenizer.Byte)
	want = []tokenizer.PreToken{
		{Value: "Hey ", Offsets: []int{0, 4}, Tokens: nil},
		{Value: "1", Offsets: []int{4, 5}, Tokens: nil},
		{Value: "2", Offsets: []int{5, 6}, Tokens: nil},
		{Value: "3", Offsets: []int{6, 7}, Tokens: nil},
		{Value: " friend!", Offsets: []int{7, 15}, Tokens: nil},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %#v\ngot %#v\n", want, got)
	}
}
