package normalizer_test

import (
	"fmt"
	"reflect"

	// "strings"
	"testing"
	"unicode"

	// "golang.org/x/text/transform"
	// "golang.org/x/text/unicode/norm"

	"github.com/sugarme/tokenizer/normalizer"
	// "github.com/sugarme/tokenizer/util"
)

func TestNormalized_NFDAddsNewChars(t *testing.T) {
	n := normalizer.NewNormalizedFrom("élégant").NFD()

	wantN := [][]int{{0, 2}, {0, 2}, {0, 2}, {2, 3}, {3, 5}, {3, 5}, {3, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}}
	gotN := n.Alignments()

	wantO := [][]int{{0, 3}, {0, 3}, {3, 4}, {4, 7}, {4, 7}, {7, 8}, {8, 9}, {9, 10}, {10, 11}}
	gotO := n.AlignmentsOriginal()

	if !reflect.DeepEqual(wantN, gotN) {
		t.Errorf("Want normalized: %v\n", wantN)
		t.Errorf("Got normalized: %v\n", gotN)
	}

	if !reflect.DeepEqual(wantO, gotO) {
		t.Errorf("Want original: %v\n", wantO)
		t.Errorf("Got original: %v\n", gotO)
	}
}

func TestNormalized_RemoveCharsAddedByNFD(t *testing.T) {
	n := normalizer.NewNormalizedFrom("élégant").NFD()
	/*
	 *   n = n.Filter(func(r rune) bool {
	 *     return unicode.Is(unicode.Mn, r)
	 *   })
	 *  */
	n = n.RemoveAccents()
	wantN := [][]int{{0, 2}, {2, 3}, {3, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}}
	gotN := n.Alignments()
	wantO := [][]int{{0, 1}, {0, 1}, {1, 2}, {2, 3}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}}
	gotO := n.AlignmentsOriginal()

	if !reflect.DeepEqual(wantN, gotN) {
		t.Errorf("Want normalized: %v\n", wantN)
		t.Errorf("Got normalized: %v\n", gotN)
	}

	if !reflect.DeepEqual(wantO, gotO) {
		t.Errorf("Want original: %v\n", wantO)
		t.Errorf("Got original: %v\n", gotO)
	}
}

func TestNormalized_RemoveChars(t *testing.T) {
	n := normalizer.NewNormalizedFrom("élégant")

	n = n.Filter(func(r rune) bool {
		return r != 'n'
	})

	wantN := [][]int{{0, 2}, {0, 2}, {2, 3}, {3, 5}, {3, 5}, {5, 6}, {6, 7}, {8, 9}}
	gotN := n.Alignments()
	wantO := [][]int{{0, 2}, {0, 2}, {2, 3}, {3, 5}, {3, 5}, {5, 6}, {6, 7}, {7, 7}, {7, 8}}
	gotO := n.AlignmentsOriginal()

	if !reflect.DeepEqual(wantN, gotN) {
		t.Errorf("Want normalized: %v\n", wantN)
		t.Errorf("Got normalized: %v\n", gotN)
	}

	if !reflect.DeepEqual(wantO, gotO) {
		t.Errorf("Want original: %v\n", wantO)
		t.Errorf("Got original: %v\n", gotO)
	}
}

func TestNormalized_MixedAdditionRemoval(t *testing.T) {
	n := normalizer.NewNormalizedFrom("élégant").NFD()

	n = n.Filter(func(r rune) bool {
		return r != 'n' && !unicode.Is(unicode.Mn, r) // Mark non-spacing
	})

	wantN := [][]int{{0, 2}, {2, 3}, {3, 5}, {5, 6}, {6, 7}, {8, 9}}
	gotN := n.Alignments()
	wantO := [][]int{{0, 1}, {0, 1}, {1, 2}, {2, 3}, {2, 3}, {3, 4}, {4, 5}, {5, 5}, {5, 6}}
	gotO := n.AlignmentsOriginal()

	if !reflect.DeepEqual(wantN, gotN) {
		t.Errorf("Want normalized: %v\n", wantN)
		t.Errorf("Got normalized: %v\n", gotN)
	}

	if !reflect.DeepEqual(wantO, gotO) {
		t.Errorf("Want original: %v\n", wantO)
		t.Errorf("Got original: %v\n", gotO)
	}
}

func TestNormalized_RangeConversion(t *testing.T) {
	n := normalizer.NewNormalizedFrom("    __Hello__   ")
	n = n.Filter(func(r rune) bool {
		return r != ' '
	}).Lowercase()

	fmt.Printf("n: %+v\n", n)

	helloRange := n.ConvertOffset(normalizer.NewRange(6, 11, normalizer.OriginalTarget))
	gotRange := helloRange.Values()
	wantRange := []int{2, 7}

	gotN := n.Range(helloRange)
	wantN := "hello"
	gotO := n.RangeOriginal(helloRange)
	wantO := "Hello"

	if !reflect.DeepEqual(wantRange, gotRange) {
		t.Errorf("Want range: %v\n", wantRange)
		t.Errorf("Got range: %v\n", gotRange)
	}

	if !reflect.DeepEqual(wantN, gotN) {
		t.Errorf("Want normalized: %v\n", wantN)
		t.Errorf("Got normalized: %v\n", gotN)
	}

	if !reflect.DeepEqual(wantO, gotO) {
		t.Errorf("Want original: %v\n", wantO)
		t.Errorf("Got original: %v\n", gotO)
	}

	// Make sure we get None only in specific cases
	fmt.Printf("Len original: %v\n", n.LenOriginal())
	testRange(t, n, []int{0, 0}, []int{0, 0}, normalizer.OriginalTarget)
	testRange(t, n, []int{3, 3}, []int{3, 3}, normalizer.OriginalTarget)
	testRange(t, n, []int{15, n.LenOriginal()}, []int{9, 9}, normalizer.OriginalTarget)
	testRange(t, n, []int{16, n.LenOriginal() + 1}, []int{16, 16}, normalizer.OriginalTarget)
	testRange(t, n, []int{17, n.LenOriginal() + 1}, nil, normalizer.OriginalTarget)
	testRange(t, n, []int{0, 0}, []int{0, 0}, normalizer.NormalizedTarget)
	testRange(t, n, []int{3, 3}, []int{3, 3}, normalizer.NormalizedTarget)
	testRange(t, n, []int{9, n.Len() + 1}, []int{9, 9}, normalizer.NormalizedTarget)
	testRange(t, n, []int{10, n.Len() + 1}, nil, normalizer.NormalizedTarget)
}

func testRange(t *testing.T, n *normalizer.NormalizedString, input, wantR []int, indexOn normalizer.IndexOn) {
	gotR := n.ConvertOffset(normalizer.NewRange(input[0], input[1], indexOn)).Values()
	if !reflect.DeepEqual(wantR, gotR) {
		t.Errorf("Want range: %v\n", wantR)
		t.Errorf("Got range: %v\n", gotR)
	}
}

func TestNormalized_AddedAroundEdge(t *testing.T) {
	n := normalizer.NewNormalizedFrom("Hello")

	changeMap := []normalizer.ChangeMap{
		{" ", 1},
		{"H", 0},
		{"e", 0},
		{"l", 0},
		{"l", 0},
		{"o", 0},
		{" ", 1},
	}

	n.Transform(changeMap, 0)

	want := " Hello "
	got := n.GetNormalized()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

func TestNormalized_AddedCharactersAlignment(t *testing.T) {
	n := normalizer.NewNormalizedFrom("野口 No")

	chars := []rune(n.GetNormalized())
	var changeMap []normalizer.ChangeMap
	for _, r := range chars {
		if r > 0x4E00 {
			changeMap = append(changeMap, []normalizer.ChangeMap{
				{" ", 0},
				{string(r), 1},
				{" ", 1},
			}...)
		} else {
			changeMap = append(changeMap, normalizer.ChangeMap{string(r), 0})
		}
	}

	n.Transform(changeMap, 0)

	original := "野口 No"
	normalized := " 野  口  No"
	alignments := [][]int{{0, 3}, {0, 3}, {0, 3}, {0, 3}, {0, 3}, {3, 6}, {3, 6}, {3, 6}, {3, 6}, {3, 6}, {6, 7}, {7, 8}, {8, 9}}
	alignmentsOriginal := [][]int{{0, 5}, {0, 5}, {0, 5}, {5, 10}, {5, 10}, {5, 10}, {10, 11}, {11, 12}, {12, 13}}
	originalShift := 0

	want := normalizer.NewNormalizedString(original, normalized, alignments, alignmentsOriginal, originalShift)
	got := n
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

func TestNormalized_RemoveAtBeginning(t *testing.T) {
	n := normalizer.NewNormalizedFrom("     Hello")
	n.Filter(func(r rune) bool {
		return r != ' '
	})

	got1 := n.RangeOriginal(normalizer.NewRange(1, len("Hello"), normalizer.NormalizedTarget))
	want1 := "ello"

	got2 := n.RangeOriginal(normalizer.NewRange(0, len(n.GetNormalized()), normalizer.NormalizedTarget))
	want2 := "Hello"

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}

	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want: %v\n", want2)
		t.Errorf("Got: %v\n", got2)
	}
}

func TestNormalized_RemoveAtEnd(t *testing.T) {
	n := normalizer.NewNormalizedFrom("Hello    ")
	n.Filter(func(r rune) bool {
		return r != ' '
	})

	got1 := n.RangeOriginal(normalizer.NewRange(0, 4, normalizer.NormalizedTarget))
	want1 := "Hell"

	got2 := n.RangeOriginal(normalizer.NewRange(0, len(n.GetNormalized()), normalizer.NormalizedTarget))
	want2 := "Hello"

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}

	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want: %v\n", want2)
		t.Errorf("Got: %v\n", got2)
	}
}

func TestNormalized_AroundBothEdges(t *testing.T) {
	n := normalizer.NewNormalizedFrom("  Hello  ")
	n.Filter(func(r rune) bool {
		return r != ' '
	})

	got0 := n.GetNormalized()
	want0 := "Hello"

	got1 := n.RangeOriginal(normalizer.NewRange(0, len("Hello"), normalizer.NormalizedTarget))
	want1 := "Hello"

	got2 := n.RangeOriginal(normalizer.NewRange(1, len("Hell"), normalizer.NormalizedTarget))
	want2 := "ell"

	if !reflect.DeepEqual(want0, got0) {
		t.Errorf("Want: %v\n", want0)
		t.Errorf("Got: %v\n", got0)
	}

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}

	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want: %v\n", want2)
		t.Errorf("Got: %v\n", got2)
	}
}
