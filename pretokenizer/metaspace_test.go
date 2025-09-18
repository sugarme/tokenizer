package pretokenizer

import (
	"reflect"
	"testing"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
)

func TestMetaspace_Decode(t *testing.T) {
	dec := DefaultMetaspace()
	// dec := NewMetaspace("_", true)

	tokens := []string{
		"▁Hey",
		"▁friend!",
	}
	got := dec.DecodeChain(tokens)
	want := []string{
		"Hey",
		" friend!",
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %q got %q\n", want, got)
	}
}

func TestMetaspace_PreTokenize(t *testing.T) {
	pt := NewMetaspace("▁", true)
	pretokenized := tokenizer.NewPreTokenizedString("Hey   friend!")

	out, err := pt.PreTokenize(pretokenized)
	if err != nil {
		panic(err)
	}

	pretoks := out.GetSplits(normalizer.NormalizedTarget, tokenizer.Byte)

	var gotValues []string
	var gotOffsets [][]int
	for _, pretok := range pretoks {
		gotValues = append(gotValues, pretok.Value)
		gotOffsets = append(gotOffsets, pretok.Offsets)
	}

	wantValues := []string{
		"▁Hey",
		"▁",
		"▁",
		"▁friend!",
	}

	wantOffsets := [][]int{
		{0, 6},
		{6, 9},
		{9, 12},
		{12, 22},
	}

	if !reflect.DeepEqual(wantValues, gotValues) {
		t.Errorf("Want %q got %q\n", wantValues, gotValues)
	}

	if !reflect.DeepEqual(wantOffsets, gotOffsets) {
		t.Errorf("Want %v got %v\n", wantOffsets, gotOffsets)
	}

	// ----------------
	gotValues = nil
	gotOffsets = nil
	pretoks = out.GetSplits(normalizer.OriginalTarget, tokenizer.Byte)
	for _, pretok := range pretoks {
		gotValues = append(gotValues, pretok.Value)
		gotOffsets = append(gotOffsets, pretok.Offsets)
	}

	wantValues = []string{
		"▁Hey",
		"▁",
		"▁",
		"▁friend!",
	}

	wantOffsets = [][]int{
		{0, 3},
		{3, 4},
		{4, 5},
		{5, 13},
	}

	if !reflect.DeepEqual(wantValues, gotValues) {
		t.Errorf("Want %q got %q\n", wantValues, gotValues)
	}

	if !reflect.DeepEqual(wantOffsets, gotOffsets) {
		t.Errorf("Want %v got %v\n", wantOffsets, gotOffsets)
	}
}
