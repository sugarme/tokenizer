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
}
