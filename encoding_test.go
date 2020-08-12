package tokenizer_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/sugarme/tokenizer"
)

func TestTokenizer_MergeWith(t *testing.T) {
	a := tokenizer.Encoding{
		Ids:              []uint32{1},
		TypeIds:          []uint32{0},
		Tokens:           []string{fmt.Sprintf("%v", "Hello ")},
		Offsets:          []tokenizer.Offsets{{0, 6}},
		SpecialTokenMask: []uint32{0},
		AttentionMask:    []uint32{1},
		Overflowing:      make([]tokenizer.Encoding, 0),
		Words:            []uint32{0},
	}

	b := tokenizer.Encoding{
		Ids:     []uint32{2},
		TypeIds: []uint32{1},
		Tokens:  []string{fmt.Sprintf("%v", "World!")},
		Offsets: []tokenizer.Offsets{{
			Start: 0,
			End:   6},
		},
		SpecialTokenMask: []uint32{0},
		AttentionMask:    []uint32{1},
		Overflowing:      make([]tokenizer.Encoding, 0),
		Words:            []uint32{0},
	}

	got := a.MergeWith(b)
	want := tokenizer.Encoding{
		Ids:     []uint32{1, 2},
		TypeIds: []uint32{0, 1},
		Tokens:  []string{fmt.Sprintf("%v", "Hello "), fmt.Sprintf("%v", "World!")},
		Offsets: []tokenizer.Offsets{
			{Start: 0, End: 6},
			{Start: 6, End: 12},
		},
		SpecialTokenMask: []uint32{0, 0},
		AttentionMask:    []uint32{1, 1},
		Overflowing:      make([]tokenizer.Encoding, 0),
		Words:            []uint32{0, 1},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}

	/*
	 *   if !reflect.DeepEqual(want.Tokens, got.Tokens) {
	 *     t.Errorf("Want: %v\n", want.Tokens)
	 *     t.Errorf("Got: %v\n", got.Tokens)
	 *   }
	 *  */
}

func TestTokenizer_Truncate(t *testing.T) {
	a := tokenizer.Encoding{
		Ids:     []uint32{1, 2, 3},
		TypeIds: []uint32{0, 0, 0},
		Tokens: []string{
			fmt.Sprintf("%v", "Hello"),
			fmt.Sprintf("%v", "World"),
			fmt.Sprintf("%v", "!"),
		},
		Offsets:          []tokenizer.Offsets{{0, 5}, {6, 11}, {11, 12}},
		SpecialTokenMask: []uint32{0, 0, 0},
		AttentionMask:    []uint32{1, 1, 1},
		Overflowing:      make([]tokenizer.Encoding, 0),
		Words:            []uint32{0, 1, 2},
	}

	got, err := a.Truncate(2, 0)
	if err != nil {
		t.Error(err)
	}

	want := tokenizer.Encoding{
		Ids:     []uint32{1, 2},
		TypeIds: []uint32{0, 0},
		Tokens: []string{
			fmt.Sprintf("%v", "Hello"),
			fmt.Sprintf("%v", "World"),
		},
		Offsets:          []tokenizer.Offsets{{0, 5}, {6, 11}},
		SpecialTokenMask: []uint32{0, 0},
		AttentionMask:    []uint32{1, 1},
		Overflowing: []tokenizer.Encoding{
			{
				Ids:     []uint32{3},
				TypeIds: []uint32{0},
				Tokens: []string{
					fmt.Sprintf("%v", "!"),
				},
				Offsets:          []tokenizer.Offsets{{11, 12}},
				SpecialTokenMask: []uint32{0},
				AttentionMask:    []uint32{1},
				Overflowing:      make([]tokenizer.Encoding, 0),
				Words:            []uint32{2},
			},
		},
		Words: []uint32{0, 1},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}

}
