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

	want1 := "Hello"
	got1 := n.RangeOriginal(normalizer.NewRange(0, n.Len(), normalizer.NormalizedTarget))
	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: '%v'\n", want1)
		t.Errorf("Got: '%v'\n", got1)
	}
}

func TestNormalized_RemoveAtEnd(t *testing.T) {

	n := normalizer.NewNormalizedFrom(`Hello    `).Filter(' ')

	nRange := normalizer.NewRange(0, 4, normalizer.NormalizedTarget)
	got := n.RangeOriginal(nRange)
	want := "Hell"
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: '%v'\n", want)
		t.Errorf("Got: '%v'\n", got)
	}
	want1 := "Hello"
	got1 := n.RangeOriginal(normalizer.NewRange(0, n.Len(), normalizer.NormalizedTarget))
	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: '%v'\n", want1)
		t.Errorf("Got: '%v'\n", got1)
	}
}

func TestNormalized_RemoveAround(t *testing.T) {

	n := normalizer.NewNormalizedFrom(`  Hello  `).Filter(' ')

	want := "Hello"
	got := n.GetNormalized()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: '%v'\n", want)
		t.Errorf("Got: '%v'\n", got)
	}

	want1 := "Hello"
	got1 := n.RangeOriginal(normalizer.NewRange(0, len([]rune("Hello")), normalizer.NormalizedTarget))
	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: '%v'\n", want1)
		t.Errorf("Got: '%v'\n", got1)
	}

	want2 := "ell"
	got2 := n.RangeOriginal(normalizer.NewRange(1, len([]rune("Hell")), normalizer.NormalizedTarget))
	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want: '%v'\n", want2)
		t.Errorf("Got: '%v'\n", got2)
	}
}

func TestNormalized_Lstrip(t *testing.T) {

	n := normalizer.NewNormalizedFrom(`  This is an example  `).LStrip()

	want := "This is an example  "
	got := n.GetNormalized()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: '%v'\n", want)
		t.Errorf("Got: '%v'\n", got)
	}

	want1 := "This is an example  "
	got1 := n.RangeOriginal(normalizer.NewRange(0, n.Len(), normalizer.NormalizedTarget))
	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: '%v'\n", want1)
		t.Errorf("Got: '%v'\n", got1)
	}
}

func TestNormalized_Rstrip(t *testing.T) {

	n := normalizer.NewNormalizedFrom(`  This is an example  `).RStrip()

	want := "  This is an example"
	got := n.GetNormalized()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: '%v'\n", want)
		t.Errorf("Got: '%v'\n", got)
	}

	want1 := "  This is an example"
	got1 := n.RangeOriginal(normalizer.NewRange(0, n.Len(), normalizer.NormalizedTarget))
	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: '%v'\n", want1)
		t.Errorf("Got: '%v'\n", got1)
	}
}

func TestNormalized_Strip(t *testing.T) {

	n := normalizer.NewNormalizedFrom(`  This is an example  `).Strip()

	want := "This is an example"
	got := n.GetNormalized()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: '%v'\n", want)
		t.Errorf("Got: '%v'\n", got)
	}
}

func TestNormalized_Prepend(t *testing.T) {

	n := normalizer.NewNormalizedFrom(`there`)
	n = n.Prepend("Hey ")

	want := "Hey there"
	got := n.GetNormalized()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: '%v'\n", want)
		t.Errorf("Got: '%v'\n", got)
	}

	want1 := []normalizer.Alignment{
		{0, 0},
		{0, 0},
		{0, 0},
		{0, 0},
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 4},
		{4, 5},
	}
	got1 := n.Alignments()
	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: '%v'\n", want1)
		t.Errorf("Got: '%v'\n", got1)
	}
}

func TestNormalized_Append(t *testing.T) {

	n := normalizer.NewNormalizedFrom(`Hey`)
	n = n.Append(" there")

	want := "Hey there"
	got := n.GetNormalized()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: '%v'\n", want)
		t.Errorf("Got: '%v'\n", got)
	}

	want1 := []normalizer.Alignment{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 3},
		{3, 3},
		{3, 3},
		{3, 3},
		{3, 3},
		{3, 3},
	}
	got1 := n.Alignments()
	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: '%v'\n", want1)
		t.Errorf("Got: '%v'\n", got1)
	}

	want2 := normalizer.NewRange(3, 3, normalizer.OriginalTarget)
	r := normalizer.NewRange(3, len([]rune(" there")), normalizer.NormalizedTarget)
	got2 := n.ConvertOffset(r)
	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want: '%v'\n", want2)
		t.Errorf("Got: '%v'\n", got2)
	}
}

func TestNormalized_GetRange(t *testing.T) {

	s := "Hello my name is John 👋"
	runes := []rune(s)

	want := string(runes[:])
	got := normalizer.RangeOf(s, util.MakeRange(0, 100))
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: '%v'\n", want)
		t.Errorf("Got: '%v'\n", got)
	}

	want1 := "John 👋"
	got1 := normalizer.RangeOf(s, util.MakeRange(17, 100))
	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: '%v'\n", want1)
		t.Errorf("Got: '%v'\n", got1)
	}
}

func TestNormalized_Merge(t *testing.T) {

	want := " A sentence that will be merged"
	merged := normalizer.NewNormalizedFrom("A sentence")
	s2 := normalizer.NewNormalizedFrom(" that will")
	s3 := normalizer.NewNormalizedFrom(" be merged")

	n := merged.Prepend(" ").MergeWith(s2).MergeWith(s3)

	got := n.GetNormalized()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: '%v'\n", want)
		t.Errorf("Got: '%v'\n", got)
	}
}

func TestNormalized_Slice(t *testing.T) {

	s := "𝔾𝕠𝕠𝕕 𝕞𝕠𝕣𝕟𝕚𝕟𝕘"
	n := normalizer.NewNormalizedFrom(s).NFKC()

	got := n.Slice(normalizer.NewRange(0, 4, normalizer.OriginalTarget))
	fmt.Printf("got: %+v\n", got)
	wantO := "𝔾𝕠𝕠𝕕"
	wantN := "Good"
	wantA := []normalizer.Alignment{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 4},
	}

	gotO := got.GetOriginal()
	gotN := got.GetNormalized()
	gotA := got.Alignments()

	if !reflect.DeepEqual(wantO, gotO) {
		t.Errorf("Want: '%v'\n", wantO)
		t.Errorf("Got: '%v'\n", gotO)
	}
	if !reflect.DeepEqual(wantN, gotN) {
		t.Errorf("Want: '%v'\n", wantN)
		t.Errorf("Got: '%v'\n", gotN)
	}
	if !reflect.DeepEqual(wantA, gotA) {
		t.Errorf("Want: '%v'\n", wantA)
		t.Errorf("Got: '%v'\n", gotA)
	}

	// Make sure the sliced NormalizedString is still aligned as expected
	s1 := normalizer.NewNormalizedFrom("   Good Morning!   ").Strip()

	// If we keep the whole slice
	slice1 := s1.Slice(normalizer.NewRange(0, 100, normalizer.OriginalTarget))
	want1 := "Good"
	got1 := slice1.RangeOriginal(normalizer.NewRange(0, 4, normalizer.NormalizedTarget))
	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: '%v'\n", want1)
		t.Errorf("Got: '%v'\n", got1)
	}
	slice2 := s1.Slice(normalizer.NewRange(0, 100, normalizer.NormalizedTarget))
	want2 := "Good"
	got2 := slice2.RangeOriginal(normalizer.NewRange(0, 4, normalizer.NormalizedTarget))
	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want: '%v'\n", want2)
		t.Errorf("Got: '%v'\n", got2)
	}

	// If we keep after the modified piece
	slice3 := s1.Slice(normalizer.NewRange(4, 15, normalizer.OriginalTarget))
	want3 := "ood"
	got3 := slice3.RangeOriginal(normalizer.NewRange(0, 3, normalizer.NormalizedTarget))
	if !reflect.DeepEqual(want3, got3) {
		t.Errorf("Want: '%v'\n", want3)
		t.Errorf("Got: '%v'\n", got3)
	}

	// If we keep only the modified piece
	slice4 := s1.Slice(normalizer.NewRange(3, 16, normalizer.OriginalTarget))
	want4 := "Good"
	got4 := slice4.RangeOriginal(normalizer.NewRange(0, 4, normalizer.NormalizedTarget))
	if !reflect.DeepEqual(want4, got4) {
		t.Errorf("Want: '%v'\n", want4)
		t.Errorf("Got: '%v'\n", got4)
	}
}

// TODO. more unit tests.
