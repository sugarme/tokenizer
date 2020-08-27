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

func TestNormalized_NFDAddsNewChars(t *testing.T) {
	n := normalizer.NewNormalizedFrom("élégant").NFD()

	wantN := [][]int{
		{0, 2},
		{0, 2},
		{0, 2},
		{2, 3},
		{3, 5},
		{3, 5},
		{3, 5},
		{5, 6},
		{6, 7},
		{7, 8},
		{8, 9},
	}
	gotN := n.Alignments()

	wantO := [][]int{
		{0, 3},
		{0, 3},
		{3, 4},
		{4, 7},
		{4, 7},
		{7, 8},
		{8, 9},
		{9, 10},
		{10, 11},
	}
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
