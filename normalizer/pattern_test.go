package normalizer_test

import (
	"reflect"
	"testing"

	"github.com/sugarme/tokenizer/normalizer"
)

func doTest(t *testing.T, p normalizer.Pattern, inside string, want []normalizer.OffsetsMatch) {

	var got []normalizer.OffsetsMatch

	switch reflect.TypeOf(p).Name() {
	case "RunePattern":
		got = p.(normalizer.RunePattern).FindMatches(inside)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %+v\n", want)
		t.Errorf("Got: %+v\n", got)
	}
}

func TestChar(t *testing.T) {
	p := normalizer.NewRunePattern('a')
	inside := "aba"
	want := []normalizer.OffsetsMatch{
		{Offsets: []int{0, 1}, Match: true},
		{Offsets: []int{1, 2}, Match: false},
		{Offsets: []int{2, 3}, Match: true},
	}
	doTest(t, p, inside, want)

	p = normalizer.NewRunePattern('a')
	inside = "bbbba"
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 4}, Match: false},
		{Offsets: []int{4, 5}, Match: true},
	}
	doTest(t, p, inside, want)

	p = normalizer.NewRunePattern('a')
	inside = "aabbb"
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 1}, Match: true},
		{Offsets: []int{1, 2}, Match: true},
		{Offsets: []int{2, 5}, Match: false},
	}
	doTest(t, p, inside, want)

	p = normalizer.NewRunePattern('a')
	inside = ""
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 0}, Match: false},
	}
	doTest(t, p, inside, want)

	p = normalizer.NewRunePattern('b')
	inside = "aaa"
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 3}, Match: false},
	}
	doTest(t, p, inside, want)
}
