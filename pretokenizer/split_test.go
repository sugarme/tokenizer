package pretokenizer

import (
	"reflect"
	"testing"

	"github.com/gengzongjie/tokenizer"
	"github.com/gengzongjie/tokenizer/normalizer"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		behavior normalizer.SplitDelimiterBehavior
		s        string
		res      []tokenizer.PreToken
	}{
		{
			behavior: normalizer.RemovedBehavior,
			s:        "How are you doing?",
			res: []tokenizer.PreToken{
				{Value: "How", Offsets: []int{0, 3}, Tokens: nil},
				{Value: "are", Offsets: []int{4, 7}, Tokens: nil},
				{Value: "you", Offsets: []int{8, 11}, Tokens: nil},
				{Value: "doing", Offsets: []int{12, 17}, Tokens: nil},
				{Value: "?", Offsets: []int{17, 18}, Tokens: nil},
			},
		},
		{
			behavior: normalizer.IsolatedBehavior,
			s:        "How are you doing?",
			res: []tokenizer.PreToken{
				{Value: "How", Offsets: []int{0, 3}, Tokens: nil},
				{Value: " ", Offsets: []int{3, 4}, Tokens: nil},
				{Value: "are", Offsets: []int{4, 7}, Tokens: nil},
				{Value: " ", Offsets: []int{7, 8}, Tokens: nil},
				{Value: "you", Offsets: []int{8, 11}, Tokens: nil},
				{Value: " ", Offsets: []int{11, 12}, Tokens: nil},
				{Value: "doing", Offsets: []int{12, 17}, Tokens: nil},
				{Value: "?", Offsets: []int{17, 18}, Tokens: nil},
			},
		},
		{
			behavior: normalizer.MergedWithPreviousBehavior,
			s:        "How are you doing?",
			res: []tokenizer.PreToken{
				{Value: "How ", Offsets: []int{0, 4}, Tokens: nil},
				{Value: "are ", Offsets: []int{4, 8}, Tokens: nil},
				{Value: "you ", Offsets: []int{8, 12}, Tokens: nil},
				{Value: "doing", Offsets: []int{12, 17}, Tokens: nil},
				{Value: "?", Offsets: []int{17, 18}, Tokens: nil},
			},
		},
		{
			behavior: normalizer.MergedWithNextBehavior,
			s:        "How are you doing?",
			res: []tokenizer.PreToken{
				{Value: "How", Offsets: []int{0, 3}, Tokens: nil},
				{Value: " are", Offsets: []int{3, 7}, Tokens: nil},
				{Value: " you", Offsets: []int{7, 11}, Tokens: nil},
				{Value: " doing", Offsets: []int{11, 17}, Tokens: nil},
				{Value: "?", Offsets: []int{17, 18}, Tokens: nil},
			},
		},
		{
			behavior: normalizer.ContiguousBehavior,
			s:        "How are you doing?",
			res: []tokenizer.PreToken{
				{Value: "How", Offsets: []int{0, 3}, Tokens: nil},
				{Value: " ", Offsets: []int{3, 4}, Tokens: nil},
				{Value: "are", Offsets: []int{4, 7}, Tokens: nil},
				{Value: " ", Offsets: []int{7, 8}, Tokens: nil},
				{Value: "you", Offsets: []int{8, 11}, Tokens: nil},
				{Value: " ", Offsets: []int{11, 12}, Tokens: nil},
				{Value: "doing?", Offsets: []int{12, 18}, Tokens: nil},
			},
		},
	}

	// Whitespace regex
	pattern := normalizer.NewRegexpPattern(`\w+|[^\w\s]+`)

	for _, d := range tests {
		pretokenized := tokenizer.NewPreTokenizedString(d.s)
		pretok := NewSplit(pattern, d.behavior, true)

		out, err := pretok.PreTokenize(pretokenized)
		if err != nil {
			t.Fail()
		}

		got := out.GetSplits(normalizer.OriginalTarget, tokenizer.Byte)
		want := d.res

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %#v\ngot %#v\n", want, got)
		}
	}
}

func TestSplitRegexString(t *testing.T) {
	pretokenized := tokenizer.NewPreTokenizedString("Hey, man!")

	// pre-tokenizer splits on " " - one from Regex, one from string
	rePattern := normalizer.NewRegexpPattern(`\s+`)
	rePretok := NewSplit(rePattern, normalizer.RemovedBehavior, false)
	strPattern := normalizer.NewStringPattern(" ")
	strPretok := NewSplit(strPattern, normalizer.RemovedBehavior, false)

	got1, err := rePretok.PreTokenize(pretokenized)
	if err != nil {
		t.Fail()
	}
	got2, err := strPretok.PreTokenize(pretokenized)
	if err != nil {
		t.Fail()
	}

	if !reflect.DeepEqual(got1, got2) {
		t.Errorf("Expected got1 and got1 are equal. But \ngot1: %#v\ngot2: %#v\n", got1, got2)
	}
}

func TestSplitInvert(t *testing.T) {
	pretokenized := tokenizer.NewPreTokenizedString("Hello Hello Hello")

	// one pre-tokenizer splits on " " - one splits inverted on "Hello"
	pattern1 := normalizer.NewStringPattern(" ")
	pretok1 := NewSplit(pattern1, normalizer.RemovedBehavior, false)
	pattern2 := normalizer.NewStringPattern("Hello")
	pretok2 := NewSplit(pattern2, normalizer.RemovedBehavior, true)

	got1, err := pretok1.PreTokenize(pretokenized)
	if err != nil {
		t.Fail()
	}
	got2, err := pretok2.PreTokenize(pretokenized)
	if err != nil {
		t.Fail()
	}

	if !reflect.DeepEqual(got1, got2) {
		t.Errorf("Expected got1 and got1 are equal. But \ngot1: %#v\ngot2: %#v\n", got1, got2)
	}
}
