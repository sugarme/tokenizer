package pretokenizer

import (
	"reflect"
	"testing"

	"github.com/season-studio/tokenizer"
	"github.com/season-studio/tokenizer/normalizer"
)

func TestGetScript(t *testing.T) {

	tests := []struct {
		script string
		name   rune
	}{
		{"Han", '京'},
		{"Han", '太'},
		{"Hiragana", 'い'},
		{"Katakana", 'グ'},
		{"Common", 'ー'},
		{"Latin", 'a'},
		{"Latin", 'A'},
		{"Common", '0'},
		{"Common", '$'},
		{"Common", '@'},
		{"Common", '-'},
		{"Common", ' '},
		{"Common", '�'},
	}

	for _, d := range tests {
		got := GetScript(d.name)
		want := d.script

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %q, got %q\n", want, got)
		}
	}
}

func TestUnicodeScript(t *testing.T) {
	pretok := DefaultUnicodeScript()
	pretokenized := tokenizer.NewPreTokenizedString("どこで生れ。Yes")

	out, err := pretok.PreTokenize(pretokenized)
	if err != nil {
		t.Fail()
	}

	got := out.GetSplits(normalizer.NormalizedTarget, tokenizer.Byte)

	want := []tokenizer.PreToken{
		{Value: "どこで生れ", Offsets: []int{0, 15}, Tokens: nil},
		{Value: "。", Offsets: []int{15, 18}, Tokens: nil},
		{Value: "Yes", Offsets: []int{18, 21}, Tokens: nil},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %#v\ngot %#v\n", want, got)
	}

	got = out.GetSplits(normalizer.OriginalTarget, tokenizer.Byte)

	want = []tokenizer.PreToken{
		{Value: "どこで生れ", Offsets: []int{0, 15}, Tokens: nil},
		{Value: "。", Offsets: []int{15, 18}, Tokens: nil},
		{Value: "Yes", Offsets: []int{18, 21}, Tokens: nil},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %#v\ngot %#v\n", want, got)
	}

}

func TestUnicodeScriptSpaces(t *testing.T) {
	pretok := DefaultUnicodeScript()
	pretokenized := tokenizer.NewPreTokenizedString("Apples are りんご 林檎")

	out, err := pretok.PreTokenize(pretokenized)
	if err != nil {
		t.Fail()
	}

	got := out.GetSplits(normalizer.NormalizedTarget, tokenizer.Byte)

	want := []tokenizer.PreToken{
		{Value: "Apples are ", Offsets: []int{0, 11}, Tokens: nil},
		{Value: "りんご 林檎", Offsets: []int{11, 27}, Tokens: nil},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %#v\ngot %#v\n", want, got)
	}

	got = out.GetSplits(normalizer.OriginalTarget, tokenizer.Byte)

	want = []tokenizer.PreToken{
		{Value: "Apples are ", Offsets: []int{0, 11}, Tokens: nil},
		{Value: "りんご 林檎", Offsets: []int{11, 27}, Tokens: nil},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %#v\ngot %#v\n", want, got)
	}

}

func TestFixedScript(t *testing.T) {
	tests := []struct {
		script string
		name   rune
	}{
		{"Han", '京'},
		{"Han", '太'},
		{"Han", 'い'},
		{"Han", 'グ'},
		{"Han", 'ー'},
		{"Latin", 'a'},
		{"Latin", 'A'},
		{"Common", '0'},
		{"Common", '$'},
		{"Common", '@'},
		{"Common", '-'},
		{"Any", ' '},
	}

	for _, d := range tests {
		got := FixedScript(d.name)
		want := d.script

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %q, got %q\n", want, got)
		}
	}
}
