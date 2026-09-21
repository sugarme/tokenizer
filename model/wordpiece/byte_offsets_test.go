package wordpiece_test

import (
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/model"
	"github.com/sugarme/tokenizer/model/wordpiece"
	"reflect"
	"testing"
)

func TestWordPieceUTF8ByteOffsets(t *testing.T) {
	vocab := model.Vocab{"[UNK]": 0, "a": 1, "å": 2, "##b": 3, "Trøndelag": 4}
	for _, tc := range []struct {
		name, text string
		limit      int
		want       [][]int
	}{
		{"whole multibyte word", "Trøndelag", 100, [][]int{{0, len("Trøndelag")}}},
		{"multibyte subwords", "åb", 100, [][]int{{0, 2}, {2, 3}}},
		{"unknown multibyte word", "😀", 100, [][]int{{0, 4}}},
		{"over character limit", "åbc", 2, [][]int{{0, 4}}},
		{"ASCII control", "ab", 100, [][]int{{0, 1}, {1, 2}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			wp := wordpiece.NewWordPieceBuilder().Vocab(&vocab).MaxInputCharsPerWord(tc.limit).Build()
			tokens, err := wp.Tokenize(tc.text)
			if err != nil {
				t.Fatal(err)
			}
			var offsets [][]int
			for _, token := range tokens {
				offsets = append(offsets, token.Offsets)
			}
			if !reflect.DeepEqual(offsets, tc.want) {
				t.Errorf("offsets = %v, want %v", offsets, tc.want)
			}
		})
	}
}

func TestWordPieceEncodingByteAndCharacterOffsets(t *testing.T) {
	vocab := model.Vocab{"[UNK]": 0, "å": 1, "##b": 2}
	wp := wordpiece.NewWordPieceBuilder().Vocab(&vocab).Build()
	tk := tokenizer.NewTokenizer(wp)
	bytes, err := tk.EncodeSingle("åb", false)
	if err != nil {
		t.Fatal(err)
	}
	if want := [][]int{{0, 2}, {2, 3}}; !reflect.DeepEqual(bytes.GetOffsets(), want) {
		t.Errorf("byte offsets = %v, want %v", bytes.GetOffsets(), want)
	}
	input := tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence("åb"))
	chars, err := tk.EncodeCharOffsets(input, false)
	if err != nil {
		t.Fatal(err)
	}
	if want := [][]int{{0, 1}, {1, 2}}; !reflect.DeepEqual(chars.GetOffsets(), want) {
		t.Errorf("character offsets = %v, want %v", chars.GetOffsets(), want)
	}
}

func TestByteToCharacterConverterEndBoundary(t *testing.T) {
	for _, tc := range []struct {
		text               string
		byteSpan, charSpan []int
	}{
		{"åb", []int{2, 3}, []int{1, 2}},
		{"😀", []int{0, 4}, []int{0, 1}},
		{"", []int{0, 0}, []int{0, 0}},
	} {
		converter := tokenizer.NewBytesToCharOffsetConverter(tc.text)
		got, err := converter.Convert(tc.byteSpan)
		if err != nil {
			t.Errorf("%q: %v", tc.text, err)
			continue
		}
		if !reflect.DeepEqual(got, tc.charSpan) {
			t.Errorf("%q: got %v, want %v", tc.text, got, tc.charSpan)
		}
		if _, err := converter.Convert([]int{0, len(tc.text) + 1}); err == nil {
			t.Errorf("%q: accepted an out-of-bounds end", tc.text)
		}
	}
}
