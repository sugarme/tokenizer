package pretokenizer

import (
	"reflect"
	"testing"

	"github.com/season-studio/tokenizer"
	"github.com/season-studio/tokenizer/normalizer"
)

func TestWhitespace(t *testing.T) {
	pretok := DefaultWhitespace()

	tests := []struct {
		s   string
		res []tokenizer.PreToken
	}{
		{
			s: "Hey man!",
			res: []tokenizer.PreToken{
				{Value: "Hey", Offsets: []int{0, 3}, Tokens: nil},
				{Value: "man", Offsets: []int{4, 7}, Tokens: nil},
				{Value: "!", Offsets: []int{7, 8}, Tokens: nil},
			},
		},
		{
			s: "How are you doing?",
			res: []tokenizer.PreToken{
				{Value: "How", Offsets: []int{0, 3}, Tokens: nil},
				{Value: "are", Offsets: []int{4, 7}, Tokens: nil},
				{Value: "you", Offsets: []int{8, 11}, Tokens: nil},
				{Value: "doing", Offsets: []int{12, 17}, Tokens: nil},
				{Value: "?", Offsets: []int{17, 18}, Tokens: nil},
			},
		},
	}

	for _, data := range tests {
		pretokenized := tokenizer.NewPreTokenizedString(data.s)
		out, err := pretok.PreTokenize(pretokenized)
		if err != nil {
			t.Fail()
		}

		got := out.GetSplits(normalizer.OriginalTarget, tokenizer.Byte)
		want := data.res

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %#v\ngot %#v\n", want, got)
		}
	}
}

func TestWhitespaceSplit(t *testing.T) {
	pretok := NewWhitespaceSplit()
	tests := []struct {
		s   string
		res []tokenizer.PreToken
	}{
		{
			s: "Hey man!",
			res: []tokenizer.PreToken{
				{Value: "Hey", Offsets: []int{0, 3}, Tokens: nil},
				{Value: "man!", Offsets: []int{4, 8}, Tokens: nil},
			},
		},
		{
			s: "Hey, man, Good?",
			res: []tokenizer.PreToken{
				{Value: "Hey,", Offsets: []int{0, 4}, Tokens: nil},
				{Value: "man,", Offsets: []int{5, 9}, Tokens: nil},
				{Value: "Good?", Offsets: []int{10, 15}, Tokens: nil},
			},
		},
	}

	for _, data := range tests {
		pretokenized := tokenizer.NewPreTokenizedString(data.s)
		out, err := pretok.PreTokenize(pretokenized)
		if err != nil {
			t.Fail()
		}

		got := out.GetSplits(normalizer.OriginalTarget, tokenizer.Byte)
		want := data.res

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %#v\ngot %#v\n", want, got)
		}
	}
}
