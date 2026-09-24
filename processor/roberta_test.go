package processor

import (
	"reflect"
	"testing"

	"github.com/sugarme/tokenizer"
)

func TestRobertaProcessingWithoutOffsetTrimming(t *testing.T) {
	processing := NewRobertaProcessing(PostToken{Value: "</s>", Id: 2}, PostToken{Value: "<s>", Id: 0}, false, true)
	makeEncoding := func(id int, token string) *tokenizer.Encoding {
		encoding := tokenizer.NewEncoding([]int{id}, []int{0}, []string{token}, [][]int{{0, 2}}, []int{0}, []int{1}, nil)
		encoding.Overflowing = []tokenizer.Encoding{*tokenizer.NewEncoding([]int{id + 1}, []int{0}, []string{token}, [][]int{{2, 4}}, []int{0}, []int{1}, nil)}
		return encoding
	}

	for _, withPair := range []bool{false, true} {
		for _, special := range []bool{false, true} {
			t.Run(testName(withPair, special), func(t *testing.T) {
				first := makeEncoding(10, "Ġhi")
				var pair *tokenizer.Encoding
				if withPair {
					pair = makeEncoding(20, "Ġby")
				}
				got := processing.Process(first, pair, special)
				if got == nil {
					t.Fatal("Process returned nil encoding")
				}
				wantIDs := []int{10}
				if special {
					wantIDs = []int{0, 10, 2}
				}
				if withPair {
					if special {
						wantIDs = append(wantIDs, 2, 20, 2)
					} else {
						wantIDs = append(wantIDs, 20)
					}
				}
				if !reflect.DeepEqual(got.Ids, wantIDs) {
					t.Errorf("IDs = %v, want %v", got.Ids, wantIDs)
				}
				if len(got.Overflowing) == 0 {
					t.Fatal("lost overflow encoding")
				}
				if got.Offsets[boolOffset(special)][1] != 2 {
					t.Errorf("offsets were trimmed: %v", got.Offsets)
				}
			})
		}
	}
}

func testName(pair, special bool) string {
	if pair {
		if special {
			return "pair_with_special_tokens"
		}
		return "pair_without_special_tokens"
	}
	if special {
		return "single_with_special_tokens"
	}
	return "single_without_special_tokens"
}

func boolOffset(special bool) int {
	if special {
		return 1
	}
	return 0
}
