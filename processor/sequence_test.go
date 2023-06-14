package processor

import (
	"reflect"
	"testing"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretokenizer"
)

func TestSequence(t *testing.T) {
	ids := []int{0, 0, 0, 0, 0}
	typeIds := []int{0, 0, 0, 0, 0}
	tokens := []string{
		"Ġ",
		"ĠĠĠĠHelloĠĠ",
		"ĠĠHello",
		"HelloĠĠ",
		"ĠĠĠĠ",
	}
	offsets := [][]int{
		{0, 1},
		{0, 11},
		{11, 18},
		{18, 25},
		{25, 29},
	}
	start := tokenizer.NewEncoding(ids, typeIds, tokens, offsets, nil, nil, nil)

	blPreTokenizer := pretokenizer.NewByteLevel()
	blPreTokenizer.TrimOffsets = true
	bl := NewByteLevelProcessing(blPreTokenizer)

	sequence := NewSequence([]tokenizer.PostProcessor{bl})

	got := bl.Process(start, nil, false)
	got2 := sequence.Process(start, nil, false)

	wantOffsets := [][]int{
		{0, 0},
		{4, 9},
		{13, 18},
		{18, 23},
		{29, 29},
	}

	sr := make(map[int]tokenizer.Range)
	sr[0] = tokenizer.NewRange(0, 5)
	sequenceRange := tokenizer.WithSequenceRangeEncodingOpt(sr)
	want := tokenizer.NewEncoding(ids, typeIds, tokens, wantOffsets, nil, nil, nil, sequenceRange)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %+v, got %+v\n", want, got)
	}
	if !reflect.DeepEqual(want, got2) {
		t.Errorf("want %+v, got %+v\n", want, got2)
	}

	pairIds := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	pairTypeIds := []int{0, 0, 0, 0, 0, 1, 1, 1, 1, 1}
	pairTokens := []string{
		"Ġ",
		"ĠĠĠĠHelloĠĠ",
		"ĠĠHello",
		"HelloĠĠ",
		"ĠĠĠĠ",
		"Ġ",
		"ĠĠĠĠHelloĠĠ",
		"ĠĠHello",
		"HelloĠĠ",
		"ĠĠĠĠ",
	}
	pairOffsets := [][]int{
		{0, 0},
		{4, 9},
		{13, 18},
		{18, 23},
		{29, 29},
		{0, 0},
		{4, 9},
		{13, 18},
		{18, 23},
		{29, 29},
	}

	sr[1] = tokenizer.NewRange(5, 10)
	srPair := tokenizer.WithSequenceRangeEncodingOpt(sr)

	wantPair := tokenizer.NewEncoding(pairIds, pairTypeIds, pairTokens, pairOffsets, nil, nil, nil, srPair)

	idsPair := []int{0, 0, 0, 0, 0}
	typeIdsPair := []int{0, 0, 0, 0, 0}
	tokensPair := []string{
		"Ġ",
		"ĠĠĠĠHelloĠĠ",
		"ĠĠHello",
		"HelloĠĠ",
		"ĠĠĠĠ",
	}
	offsetsPair := [][]int{
		{0, 1},
		{0, 11},
		{11, 18},
		{18, 25},
		{25, 29},
	}
	pairStart := tokenizer.NewEncoding(idsPair, typeIdsPair, tokensPair, offsetsPair, nil, nil, nil)
	gotPair := bl.Process(start, pairStart, false)
	gotPair2 := sequence.Process(start, pairStart, false)

	if !reflect.DeepEqual(wantPair, gotPair) {
		t.Errorf("want %#v\n, got %#v\n", wantPair, gotPair)
	}

	if !reflect.DeepEqual(wantPair, gotPair2) {
		t.Errorf("want %#v\n, got %#v\n", wantPair, gotPair2)
	}

}
