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
	gotN := normalizer.NewNormalizedFrom("élégant").NFD()

	want1 := [][]int{
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
	got1 := gotN.Alignments()

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}
}
