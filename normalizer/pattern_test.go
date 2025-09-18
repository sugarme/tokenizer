package normalizer_test

import (
	"reflect"
	"testing"

	"github.com/sugarme/tokenizer/normalizer"
)

func doTest(t *testing.T, p normalizer.Pattern, inside string, want []normalizer.OffsetsMatch) {

	var got []normalizer.OffsetsMatch

	typeStr := reflect.TypeOf(p).String()

	switch typeStr {
	case "*normalizer.RunePattern":
		got = p.(*normalizer.RunePattern).FindMatches(inside)
	case "*normalizer.StringPattern":
		got = p.(*normalizer.StringPattern).FindMatches(inside)
	case "*normalizer.FnPattern":
		got = p.(*normalizer.FnPattern).FindMatches(inside)
	case "*normalizer.RegexpPattern":
		got = p.(*normalizer.RegexpPattern).FindMatches(inside)

	default:
		panic("Invalid type\n")
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

func TestString(t *testing.T) {
	p := normalizer.NewStringPattern("a")
	inside := "aba"
	want := []normalizer.OffsetsMatch{
		{Offsets: []int{0, 1}, Match: true},
		{Offsets: []int{1, 2}, Match: false},
		{Offsets: []int{2, 3}, Match: true},
	}
	doTest(t, p, inside, want)

	p = normalizer.NewStringPattern("a")
	inside = "bbbba"
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 4}, Match: false},
		{Offsets: []int{4, 5}, Match: true},
	}
	doTest(t, p, inside, want)

	p = normalizer.NewStringPattern("a")
	inside = "aabbb"
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 1}, Match: true},
		{Offsets: []int{1, 2}, Match: true},
		{Offsets: []int{2, 5}, Match: false},
	}
	doTest(t, p, inside, want)

	p = normalizer.NewStringPattern("ab")
	inside = "aabbb"
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 1}, Match: false},
		{Offsets: []int{1, 3}, Match: true},
		{Offsets: []int{3, 5}, Match: false},
	}
	doTest(t, p, inside, want)

	p = normalizer.NewStringPattern("ab")
	inside = "aabbab"
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 1}, Match: false},
		{Offsets: []int{1, 3}, Match: true},
		{Offsets: []int{3, 4}, Match: false},
		{Offsets: []int{4, 6}, Match: true},
	}
	doTest(t, p, inside, want)

	p = normalizer.NewStringPattern("")
	inside = ""
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 0}, Match: false},
	}
	doTest(t, p, inside, want)

	p = normalizer.NewStringPattern("")
	inside = "aaa"
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 3}, Match: false},
	}
	doTest(t, p, inside, want)

	p = normalizer.NewStringPattern("b")
	inside = "aaa"
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 3}, Match: false},
	}
	doTest(t, p, inside, want)

}

func TestFunctions(t *testing.T) {
	var fn normalizer.PatternFn = func(r rune) bool {
		return r == 'b'
	}

	p := normalizer.NewFnPattern(fn)
	inside := "aba"
	want := []normalizer.OffsetsMatch{
		{Offsets: []int{0, 1}, Match: false},
		{Offsets: []int{1, 2}, Match: true},
		{Offsets: []int{2, 3}, Match: false},
	}
	doTest(t, p, inside, want)

	inside = "aaaab"
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 4}, Match: false},
		{Offsets: []int{4, 5}, Match: true},
	}
	doTest(t, p, inside, want)

	inside = "bbaaa"
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 1}, Match: true},
		{Offsets: []int{1, 2}, Match: true},
		{Offsets: []int{2, 5}, Match: false},
	}
	doTest(t, p, inside, want)

	inside = ""
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 0}, Match: false},
	}
	doTest(t, p, inside, want)

	inside = "aaa"
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 3}, Match: false},
	}
	doTest(t, p, inside, want)
}

func TestRegexp(t *testing.T) {

	// isWhitespace
	p := normalizer.NewRegexpPattern(`\s+`)
	inside := "a   b"
	want := []normalizer.OffsetsMatch{
		{Offsets: []int{0, 1}, Match: false},
		{Offsets: []int{1, 4}, Match: true},
		{Offsets: []int{4, 5}, Match: false},
	}
	doTest(t, p, inside, want)

	inside = "   a   b   "
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 3}, Match: true},
		{Offsets: []int{3, 4}, Match: false},
		{Offsets: []int{4, 7}, Match: true},
		{Offsets: []int{7, 8}, Match: false},
		{Offsets: []int{8, 11}, Match: true},
	}
	doTest(t, p, inside, want)

	inside = ""
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 0}, Match: false},
	}
	doTest(t, p, inside, want)

	inside = "𝔾𝕠𝕠𝕕 𝕞𝕠𝕣𝕟𝕚𝕟𝕘"
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 16}, Match: false},
		{Offsets: []int{16, 17}, Match: true},
		{Offsets: []int{17, 45}, Match: false},
	}
	doTest(t, p, inside, want)

	inside = "aaa"
	want = []normalizer.OffsetsMatch{
		{Offsets: []int{0, 3}, Match: false},
	}
	doTest(t, p, inside, want)
}
