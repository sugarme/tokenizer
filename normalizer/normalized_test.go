package normalizer_test

import (
	// "fmt"
	"reflect"

	// "strings"
	"testing"
	// "unicode"

	// "golang.org/x/text/transform"
	// "golang.org/x/text/unicode/norm"

	"github.com/sugarme/tokenizer/normalizer"
	// "github.com/sugarme/tokenizer/util"
)

func TestNormalized_NewNormalizedFrom(t *testing.T) {
	gotN := normalizer.NewNormalizedFrom("élégant").NFD()

	want := []normalizer.Alignment{
		{0, 1},
		{0, 1},
		{1, 2},
		{2, 3},
		{2, 3},
		{3, 4},
		{4, 5},
		{5, 6},
		{6, 7},
	}
	got := gotN.Alignments()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

// Unchanged: Remove accents - Mark, nonspacing (Mn)
func TestNormalized_RemoveAccents(t *testing.T) {
	gotN := normalizer.NewNormalizedFrom("élégant").RemoveAccents()

	want := []normalizer.Alignment{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 4},
		{4, 5},
		{5, 6},
		{6, 7},
	}
	got := gotN.Alignments()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

// Removed Chars
func TestNormalized_Filter(t *testing.T) {
	gotN := normalizer.NewNormalizedFrom("élégant").Filter('n')

	want := []normalizer.Alignment{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 4},
		{4, 5},
		{6, 7},
	}
	got := gotN.Alignments()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

// Mixed addition and removal
func TestNormalized_Mixed(t *testing.T) {
	gotN := normalizer.NewNormalizedFrom("élégant").RemoveAccents().Filter('n')

	want := []normalizer.Alignment{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 4},
		{4, 5},
		{6, 7},
	}
	got := gotN.Alignments()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

// Range and Conversion
func TestNormalized_RangeConversion(t *testing.T) {
	gotN := normalizer.NewNormalizedFrom(`    __Hello__   `).Filter(' ').Lowercase()

	originalRange := normalizer.NewRange(6, 11, normalizer.OriginalTarget)
	got1 := gotN.Range(originalRange)
	want1 := "hello"
	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}

	// normalized string: '__hello__'
	normalizedRange := normalizer.NewRange(2, 7, normalizer.NormalizedTarget)
	got2 := gotN.Range(normalizedRange)
	want2 := "hello"
	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want: %v\n", want2)
		t.Errorf("Got: %v\n", got2)
	}

	helloN := gotN.ConvertOffset(normalizer.NewRange(6, 11, normalizer.OriginalTarget))
	got3 := helloN
	want3 := normalizer.NewRange(2, 7, normalizer.NormalizedTarget)
	if !reflect.DeepEqual(want3, got3) {
		t.Errorf("Want: %+v\n", want3)
		t.Errorf("Got: %+v\n", got3)
	}

	got4 := gotN.Range(helloN)
	want4 := "hello"
	if !reflect.DeepEqual(want4, got4) {
		t.Errorf("Want: '%v'\n", want4)
		t.Errorf("Got: '%v'\n", got4)
	}

	got5 := gotN.RangeOriginal(originalRange) // (6,11)
	want5 := "Hello"
	if !reflect.DeepEqual(want5, got5) {
		t.Errorf("Want: %v\n", want5)
		t.Errorf("Got: %v\n", got5)
	}
}

func TestNormalized_OriginalRange(t *testing.T) {
	n := normalizer.NewNormalizedFrom(`Hello_______ World!`).Filter('_').Lowercase()

	// normalized string: 'hello world!'
	normalizedRange := normalizer.NewRange(6, 11, normalizer.NormalizedTarget)
	worldN := n.Range(normalizedRange)
	wantWorldN := "world"
	if !reflect.DeepEqual(wantWorldN, worldN) {
		t.Errorf("Want normalized world: %v\n", wantWorldN)
		t.Errorf("Got normalized world: %v\n", worldN)
	}

	// original string: 'Hello_______ World!'
	originalRange := n.ConvertOffset(normalizedRange)
	// originalRange := normalizer.NewRange(13, 18, normalizer.OriginalTarget)
	worldO := n.RangeOriginal(originalRange)
	wantWorldO := "World"
	if !reflect.DeepEqual(wantWorldO, worldO) {
		t.Errorf("Want original world: %v\n", wantWorldO)
		t.Errorf("Got original world: %v\n", worldO)
	}
}

func TestNormalized_AddedAroundEdge(t *testing.T) {
	n := normalizer.NewNormalizedFrom(`Hello`)

	var changeMap []normalizer.ChangeMap = []normalizer.ChangeMap{
		{" ", 1},
		{"H", 0},
		{"e", 0},
		{"l", 0},
		{"l", 0},
		{"o", 0},
		{" ", 1},
	}

	n = n.Transform(changeMap, 0)

	want := " Hello "
	got := n.GetNormalized()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: '%v'\n", want)
		t.Errorf("Got: '%v'\n", got)
	}

	n1 := normalizer.NewNormalizedFrom(` Hello `)
	normalizedRange := normalizer.NewRange(1, len([]rune(n1.GetNormalized()))-1, normalizer.NormalizedTarget)
	gotO := n1.RangeOriginal(normalizedRange)
	wantO := "Hello"
	if !reflect.DeepEqual(wantO, gotO) {
		t.Errorf("Want: '%v'\n", wantO)
		t.Errorf("Got: '%v'\n", gotO)
	}
}

func TestNormalized_RemoveAtStart(t *testing.T) {

	n := normalizer.NewNormalizedFrom(`     Hello`).Filter(' ') // 5 white spaces

	nRange := normalizer.NewRange(1, len([]rune("Hello")), normalizer.NormalizedTarget)
	got := n.RangeOriginal(nRange)
	want := "ello"
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: '%v'\n", want)
		t.Errorf("Got: '%v'\n", got)
	}

}

func TestNormalized_ConvertOffset(t *testing.T) {
	// Test 1
	n1 := normalizer.NewNormalizedFrom(`    __Hello__   `).Filter(' ').Lowercase() // `__hello__`

	// Original -> Normalized
	oRange1a := normalizer.NewRange(6, 11, normalizer.OriginalTarget)
	nRange1a := n1.ConvertOffset(oRange1a)

	want1a := normalizer.NewRange(2, 7, normalizer.NormalizedTarget)
	got1a := nRange1a
	if !reflect.DeepEqual(want1a, got1a) {
		t.Errorf("Want: %v\n", want1a)
		t.Errorf("Got: %v\n", got1a)
	}

	// Normalized -> Original
	nRange1b := normalizer.NewRange(2, 7, normalizer.NormalizedTarget)
	oRange1b := n1.ConvertOffset(nRange1b)
	want1b := normalizer.NewRange(6, 11, normalizer.OriginalTarget)
	got1b := oRange1b
	if !reflect.DeepEqual(want1b, got1b) {
		t.Errorf("Want: %v\n", want1b)
		t.Errorf("Got: %v\n", got1b)
	}

	// Test 2
	n2 := normalizer.NewNormalizedFrom(`     Hello`).Filter(' ')

	oRange2a := normalizer.NewRange(6, 9, normalizer.OriginalTarget)
	nRange2a := n2.ConvertOffset(oRange2a)
	want2a := normalizer.NewRange(1, 4, normalizer.NormalizedTarget)
	got2a := nRange2a
	if !reflect.DeepEqual(want2a, got2a) {
		t.Errorf("Want: %v\n", want2a)
		t.Errorf("Got: %v\n", got2a)
	}

	nRange2b := normalizer.NewRange(1, 5, normalizer.NormalizedTarget) // `ello`
	oRange2b := n2.ConvertOffset(nRange2b)
	want2b := normalizer.NewRange(6, 10, normalizer.OriginalTarget)
	got2b := oRange2b
	if !reflect.DeepEqual(want2b, got2b) {
		t.Errorf("Want: %v\n", want2b)
		t.Errorf("Got: %v\n", got2b)
	}

	// Test 3
	n3 := normalizer.NewNormalizedFrom(`Hello_______ World!`).Filter('_').Lowercase()

	oRange3a := normalizer.NewRange(13, 18, normalizer.OriginalTarget) // `World`
	nRange3a := n3.ConvertOffset(oRange3a)
	want3a := normalizer.NewRange(6, 11, normalizer.NormalizedTarget)
	got3a := nRange3a
	if !reflect.DeepEqual(want3a, got3a) {
		t.Errorf("Want range: %v\n", want3a)
		t.Errorf("Got range: %v\n", got3a)
	}

	// normalized string: 'hello world!'
	nRange3b := normalizer.NewRange(6, 11, normalizer.NormalizedTarget)
	oRange3b := n3.ConvertOffset(nRange3b)
	want3b := normalizer.NewRange(13, 18, normalizer.OriginalTarget)
	got3b := oRange3b

	if !reflect.DeepEqual(want3b, got3b) {
		t.Errorf("Want range: %v\n", want3b)
		t.Errorf("Got range: %v\n", got3b)
	}
}

// TODO. more unit tests.
