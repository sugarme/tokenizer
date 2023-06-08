package normalizer

import (
	"reflect"
	"testing"
)

func TestPrepend(t *testing.T) {
	original := "Hello"
	// normalized := "▁Hello"

	n := NewNormalizedFrom(original)
	prepend := NewPrepend("▁")

	out, err := prepend.Normalize(n)
	if err != nil {
		panic(err)
	}

	gotAlignments := out.Alignments()
	wantAlignments := [][]int{
		{0, 1},
		{0, 1},
		{0, 1},
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 4},
		{4, 5},
	}

	if !reflect.DeepEqual(wantAlignments, gotAlignments) {
		t.Errorf("want %v, got %v\n", wantAlignments, gotAlignments)
	}

	gotNormalized := out.GetNormalized()
	wantNormalized := "▁Hello"

	if !reflect.DeepEqual(wantNormalized, gotNormalized) {
		t.Errorf("want %v, got %v\n", wantNormalized, gotNormalized)
	}

	gotOriginal := out.AlignmentsOriginal()
	wantOriginal := [][]int{
		{0, 4},
		{4, 5},
		{5, 6},
		{6, 7},
		{7, 8},
	}

	if !reflect.DeepEqual(wantOriginal, gotOriginal) {
		t.Errorf("want %v, got %v\n", wantOriginal, gotOriginal)
	}

}
