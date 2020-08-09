package normalizer_test

import (
	"fmt"
	"reflect"

	// "strings"
	"testing"
	// "unicode"

	// "golang.org/x/text/transform"
	// "golang.org/x/text/unicode/norm"

	"github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/util"
)

func TestNormalized_NewNormalizedFrom(t *testing.T) {
	gotN := normalizer.NewNormalizedFrom("élégant")
	gotN.NFD()

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
	got := gotN.Get().Alignments

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

// Unchanged: Remove accents - Mark, nonspacing (Mn)
func TestNormalized_RemoveAccents(t *testing.T) {
	gotN := normalizer.NewNormalizedFrom("élégant")
	gotN.RemoveAccents()

	want := []normalizer.Alignment{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 4},
		{4, 5},
		{5, 6},
		{6, 7},
	}
	got := gotN.Get().Alignments

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

// Removed Chars
func TestNormalized_Filter(t *testing.T) {
	gotN := normalizer.NewNormalizedFrom("élégant")

	gotN.Filter('n')

	want := []normalizer.Alignment{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 4},
		{4, 5},
		{6, 7},
	}
	got := gotN.Get().Alignments

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

// Mixed addition and removal
func TestNormalized_Mixed(t *testing.T) {
	gotN := normalizer.NewNormalizedFrom("élégant")

	gotN.RemoveAccents()
	gotN.Filter('n')

	want := []normalizer.Alignment{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 4},
		{4, 5},
		{6, 7},
	}
	got := gotN.Get().Alignments

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

// Range and Conversion
func TestNormalized_RangeConversion(t *testing.T) {
	gotN := normalizer.NewNormalizedFrom(`    __Hello__   `)

	gotN.Filter(' ')
	gotN.Lowercase()

	originalRange := normalizer.NewRange(6, 11, normalizer.OriginalTarget)
	got1 := gotN.Range(originalRange)
	want1 := "Hello"
	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}

	normalizedRange := normalizer.NewRange(2, 7, normalizer.NormalizedTarget)
	got2 := gotN.Range(normalizedRange)
	want2 := "hello"
	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want: %v\n", want2)
		t.Errorf("Got: %v\n", got2)
	}

	// helloN will have `indexOn` = NormalizedTarget
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

	got5 := gotN.RangeOriginal(util.MakeRange(6, 11))
	want5 := "Hello"
	if !reflect.DeepEqual(want5, got5) {
		t.Errorf("Want: %v\n", want5)
		t.Errorf("Got: %v\n", got5)
	}
}

func TestNormalized_OriginalRange(t *testing.T) {
	n := normalizer.NewNormalizedFrom(`Hello_______ World!`)
	n.Filter('_')
	n.Lowercase()
	fmt.Printf("Original: '%v'\n", n.GetOriginal())
	fmt.Printf("Normalized: '%v'\n", n.GetNormalized())
	fmt.Printf("All alignments: %+v\n", n.Get().Alignments)

	normalizedRange := normalizer.NewRange(6, 11, normalizer.NormalizedTarget)
	worldN := n.Range(normalizedRange)
	rangeOriginal := n.ConvertOffset(normalizedRange)
	fmt.Printf("rangeOriginal: %+v\n", rangeOriginal)
	worldO := n.Range(rangeOriginal)

	wantWorldN := "world"
	wantWorldO := "World"

	if !reflect.DeepEqual(wantWorldN, worldN) {
		t.Errorf("Want normalized world: %v\n", wantWorldN)
		t.Errorf("Got normalized world: %v\n", worldN)
	}

	if !reflect.DeepEqual(wantWorldO, worldO) {
		t.Errorf("Want original world: %v\n", wantWorldO)
		t.Errorf("Got original world: %v\n", worldO)
	}
}

// TODO. more unit tests.
