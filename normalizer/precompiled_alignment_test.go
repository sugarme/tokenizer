package normalizer

import (
	"github.com/sugarme/tokenizer/spm"
	"reflect"
	"testing"
)

func TestPrecompiledExpansionPreservesAlignment(t *testing.T) {
	compiled, err := spm.NewPrecompiledFrom(spm.NmtNfkc())
	if err != nil {
		t.Fatal(err)
	}
	normalizer := &Precompiled{compiled}
	for _, tc := range []struct {
		input, want string
		align       [][]int
	}{
		{"a™b", "aTMb", [][]int{{0, 1}, {1, 4}, {1, 4}, {4, 5}}},
		{"™", "TM", [][]int{{0, 3}, {0, 3}}},
		{"plain", "plain", [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}}},
		{"Ａb", "Ab", [][]int{{0, 3}, {3, 4}}},
		{"™ﬁ", "TMfi", [][]int{{0, 3}, {0, 3}, {3, 6}, {3, 6}}},
		{"aﬁb", "afib", [][]int{{0, 1}, {1, 4}, {1, 4}, {4, 5}}},
	} {
		t.Run(tc.input, func(t *testing.T) {
			defer func() {
				if p := recover(); p != nil {
					t.Fatalf("Normalize panicked: %v", p)
				}
			}()
			got, err := normalizer.Normalize(NewNormalizedFrom(tc.input))
			if err != nil {
				t.Fatal(err)
			}
			if got.GetNormalized() != tc.want {
				t.Errorf("text = %q, want %q", got.GetNormalized(), tc.want)
			}
			if !reflect.DeepEqual(got.Alignments(), tc.align) {
				t.Errorf("alignments = %v, want %v", got.Alignments(), tc.align)
			}
		})
	}
}
