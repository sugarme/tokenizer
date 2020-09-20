package tokenizer_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/sugarme/tokenizer"
)

func TestTokenizer_MergeWith(t *testing.T) {

	a := tokenizer.Encoding{
		Ids:              []int{1},
		TypeIds:          []int{0},
		Tokens:           []string{fmt.Sprintf("%v", "Hello ")},
		Offsets:          [][]int{{0, 6}},
		SpecialTokenMask: []int{0},
		AttentionMask:    []int{1},
		Overflowing:      make([]tokenizer.Encoding, 0),
		Words:            []int{0},
	}

	b := tokenizer.Encoding{
		Ids:              []int{2},
		TypeIds:          []int{1},
		Tokens:           []string{fmt.Sprintf("%v", "World!")},
		Offsets:          [][]int{{0, 6}},
		SpecialTokenMask: []int{0},
		AttentionMask:    []int{1},
		Overflowing:      make([]tokenizer.Encoding, 0),
		Words:            []int{0},
	}

	got := a.MergeWith(&b, true)

	want := &tokenizer.Encoding{
		Ids:              []int{1, 2},
		TypeIds:          []int{0, 1},
		Tokens:           []string{fmt.Sprintf("%v", "Hello "), fmt.Sprintf("%v", "World!")},
		Offsets:          [][]int{{0, 6}, {6, 12}},
		SpecialTokenMask: []int{0, 0},
		AttentionMask:    []int{1, 1},
		Overflowing:      make([]tokenizer.Encoding, 0),
		Words:            []int{0, 0},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

func TestTokenizer_Truncate(t *testing.T) {
	a := tokenizer.Encoding{
		Ids:     []int{1, 2, 3},
		TypeIds: []int{0, 0, 0},
		Tokens: []string{
			fmt.Sprintf("%v", "Hello"),
			fmt.Sprintf("%v", "World"),
			fmt.Sprintf("%v", "!"),
		},
		Offsets:          [][]int{{0, 5}, {6, 11}, {11, 12}},
		SpecialTokenMask: []int{0, 0, 0},
		AttentionMask:    []int{1, 1, 1},
		Overflowing:      make([]tokenizer.Encoding, 0),
		Words:            []int{0, 1, 2},
	}

	got, err := a.Truncate(2, 0)
	if err != nil {
		t.Error(err)
	}

	want := &tokenizer.Encoding{
		Ids:     []int{1, 2},
		TypeIds: []int{0, 0},
		Tokens: []string{
			fmt.Sprintf("%v", "Hello"),
			fmt.Sprintf("%v", "World"),
		},
		Offsets:          [][]int{{0, 5}, {6, 11}},
		SpecialTokenMask: []int{0, 0},
		AttentionMask:    []int{1, 1},
		Overflowing: []tokenizer.Encoding{
			{
				Ids:     []int{3},
				TypeIds: []int{0},
				Tokens: []string{
					fmt.Sprintf("%v", "!"),
				},
				Offsets:          [][]int{{11, 12}},
				SpecialTokenMask: []int{0},
				AttentionMask:    []int{1},
				Overflowing:      make([]tokenizer.Encoding, 0),
				Words:            []int{2},
			},
		},
		Words: []int{0, 1},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

func TestTokenizer_Mapping(t *testing.T) {
	encoding := tokenizer.DefaultEncoding()
	encoding.Tokens = []string{"He", "llo", "won", "der", "ful", "friend", "!"}
	encoding.Offsets = [][]int{
		{0, 2},
		{2, 5},
		{7, 10},
		{10, 13},
		{13, 16},
		{17, 23},
		{23, 24},
	}
	encoding.Words = []int{0, 0, 1, 1, 1, 2, 3}

	var (
		start, end int
		ok         bool
	)

	if start, end, ok = encoding.Word2Tokens(0); !ok {
		start, end = -1, -1
	}
	testMapping(t, []int{start, end}, []int{0, 2})

	if start, end, ok = encoding.Word2Tokens(1); !ok {
		start, end = -1, -1
	}
	testMapping(t, []int{start, end}, []int{2, 5})

	if start, end, ok = encoding.Word2Tokens(2); !ok {
		start, end = -1, -1
	}
	testMapping(t, []int{start, end}, []int{5, 6})

	if start, end, ok = encoding.Word2Tokens(3); !ok {
		start, end = -1, -1
	}
	testMapping(t, []int{start, end}, []int{6, 7})

	var chars []int
	if chars, ok = encoding.Word2Chars(0); !ok {
		chars = []int{-1, -1}
	}
	testMapping(t, []int{chars[0], chars[1]}, []int{0, 5})

	if chars, ok = encoding.Word2Chars(1); !ok {
		chars = []int{-1, -1}
	}
	testMapping(t, []int{chars[0], chars[1]}, []int{7, 16})

	if chars, ok = encoding.Token2Chars(0); !ok {
		chars = []int{-1, -1}
	}
	testMapping(t, []int{chars[0], chars[1]}, []int{0, 2})

	var word int
	if word, ok = encoding.Token2Word(1); !ok {
		word = -1
	}
	testMapping(t, word, 0)

	if word, ok = encoding.Token2Word(2); !ok {
		word = -1
	}
	testMapping(t, word, 1)

	if word, ok = encoding.Token2Word(7); !ok {
		word = -1
	}
	testMapping(t, word, -1)

	var token int
	if token, ok = encoding.Char2Token(3); !ok {
		token = -1
	}
	testMapping(t, token, 1)

	if token, ok = encoding.Char2Token(8); !ok {
		token = -1
	}
	testMapping(t, token, 2)

	if token, ok = encoding.Char2Token(16); !ok {
		token = -1
	}
	testMapping(t, token, -1)

	if token, ok = encoding.Char2Token(23); !ok {
		token = -1
	}
	testMapping(t, token, 6)

	if word, ok = encoding.Char2Word(3); !ok {
		word = -1
	}
	testMapping(t, word, 0)

	if word, ok = encoding.Char2Word(8); !ok {
		word = -1
	}
	testMapping(t, word, 1)

	if word, ok = encoding.Char2Word(16); !ok {
		word = -1
	}
	testMapping(t, word, -1)

	if word, ok = encoding.Char2Word(23); !ok {
		word = -1
	}
	testMapping(t, word, 3)
}

func testMapping(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}
