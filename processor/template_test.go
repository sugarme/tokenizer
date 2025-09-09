package processor

import (
	"reflect"
	"testing"

	"github.com/gengzongjie/tokenizer"
)

func TestPiece(t *testing.T) {
	seq0 := &SequencePiece{
		Id:     A,
		TypeId: 0,
	}

	seq0String := "$"
	got, err := NewPiece(seq0String)
	if err != nil {
		panic(err)
	}

	want := seq0

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v\n ", want, got)
	}

	seq1 := &SequencePiece{
		Id:     B,
		TypeId: 0,
	}

	seq1String := "$B"
	got, err = NewPiece(seq1String)
	if err != nil {
		panic(err)
	}

	want = seq1
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v\n ", want, got)
	}

	want = &SequencePiece{
		Id:     A,
		TypeId: 1,
	}

	got, err = NewPiece("$1")
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v\n ", want, got)
	}

	want = &SequencePiece{
		Id:     B,
		TypeId: 2,
	}

	got, err = NewPiece("$B:2")
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v\n ", want, got)
	}

	want = &SequencePiece{
		Id:     A,
		TypeId: 1,
	}

	got, err = NewPiece("$:1")
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v\n ", want, got)
	}
}

func getBertTemplate() *TemplateProcessing {
	builder := DefaultTemplateProcessing().Builder()
	builder.NewSingle([]string{"[CLS]", "$0", "[SEP]"})
	builder.NewPair("[CLS]:0 $A:0 [SEP]:0 $B:1 [SEP]:1")
	builder.NewSpecialTokens([]tokenizer.Token{
		{Id: 1, Value: "[CLS]", Offsets: nil},
		{Id: 0, Value: "[SEP]", Offsets: nil},
	})

	return builder.Build()
}

func TestTemplateProcessing(t *testing.T) {
	processor := getBertTemplate()

	got := processor.AddedTokens(false)
	want := 2

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}

	got = processor.AddedTokens(true)
	want = 3

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}

	encoding := tokenizer.NewEncodingFromTokens([]tokenizer.Token{
		{Id: 12, Value: "Hello", Offsets: []int{0, 5}},
		{Id: 14, Value: "there", Offsets: []int{6, 11}},
	}, 0)

	singleEncoding := processor.Process(encoding, nil, true)

	ids := []int{1, 12, 14, 0}
	tokens := []string{
		"[CLS]",
		"Hello",
		"there",
		"[SEP]",
	}

	typeIds := []int{0, 0, 0, 0}

	offsets := [][]int{
		{0, 0},
		{0, 5},
		{6, 11},
		{0, 0},
	}
	specialTokenMask := []int{1, 0, 0, 1}
	attentionMask := []int{1, 1, 1, 1}

	sr := make(map[int]tokenizer.Range)
	sr[0] = tokenizer.NewRange(1, 3)
	sequenceRange := tokenizer.WithSequenceRangeEncodingOpt(sr)

	want1 := tokenizer.NewEncoding(ids, typeIds, tokens, offsets, specialTokenMask, attentionMask, nil, sequenceRange)
	got1 := singleEncoding

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("want %+v, \ngot %+v\n", want1, got1)
	}

	seq, ok := singleEncoding.Token2Sequence(2)
	if !ok || seq != 0 {
		t.Errorf("want %v, got %v\n", 0, seq)
	}

	seq, ok = singleEncoding.Token2Sequence(3)
	if ok {
		t.Errorf("want %v, got %v\n", 0, seq)
	}

	pair := tokenizer.NewEncodingFromTokens([]tokenizer.Token{
		{Id: 15, Value: "pair", Offsets: []int{0, 4}},
	}, 0)

	pairEncoding := processor.Process(encoding, pair, true)

	ids = []int{1, 12, 14, 0, 15, 0}
	typeIds = []int{0, 0, 0, 0, 1, 1}
	tokens = []string{
		"[CLS]",
		"Hello",
		"there",
		"[SEP]",
		"pair",
		"[SEP]",
	}

	offsets = [][]int{
		{0, 0},
		{0, 5},
		{6, 11},
		{0, 0},
		{0, 4},
		{0, 0},
	}
	specialTokenMask = []int{1, 0, 0, 1, 0, 1}
	attentionMask = []int{1, 1, 1, 1, 1, 1}

	sr = make(map[int]tokenizer.Range)
	sr[0] = tokenizer.NewRange(1, 3)
	sr[1] = tokenizer.NewRange(4, 5)
	sequenceRange = tokenizer.WithSequenceRangeEncodingOpt(sr)

	want2 := tokenizer.NewEncoding(ids, typeIds, tokens, offsets, specialTokenMask, attentionMask, nil, sequenceRange)
	got2 := pairEncoding

	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("want %+v, \ngot %+v\n", want2, got2)
	}

	seq, ok = pairEncoding.Token2Sequence(2)
	if !ok || seq != 0 {
		t.Errorf("want %v, got %v\n", 0, seq)
	}

	seq, ok = pairEncoding.Token2Sequence(3)
	if ok {
		t.Errorf("want %v (nil), got %v\n", -1, seq)
	}

	seq, ok = pairEncoding.Token2Sequence(4)
	if !ok || seq != 1 {
		t.Errorf("want %v, got %v\n", 1, seq)
	}

	seq, ok = pairEncoding.Token2Sequence(5)
	if ok {
		t.Errorf("want %v (nil), got %v\n", -1, seq)
	}
}

func TestTemplateProcessingOverflowing(t *testing.T) {
	processor := getBertTemplate()

	got := processor.AddedTokens(false)
	want := 2

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}

	got = processor.AddedTokens(true)
	want = 3

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}

	encoding := tokenizer.NewEncodingFromTokens([]tokenizer.Token{
		{Id: 12, Value: "Hello", Offsets: []int{0, 5}},
		{Id: 14, Value: "there", Offsets: []int{6, 11}},
	}, 0)

	overflowing := tokenizer.NewEncodingFromTokens([]tokenizer.Token{
		{Id: 13, Value: "you", Offsets: []int{12, 15}},
	}, 0)

	encoding.SetOverflowing([]tokenizer.Encoding{*overflowing})

	singleEncoding := processor.Process(encoding, nil, true)

	ids := []int{1, 12, 14, 0}
	tokens := []string{
		"[CLS]",
		"Hello",
		"there",
		"[SEP]",
	}

	typeIds := []int{0, 0, 0, 0}

	offsets := [][]int{
		{0, 0},
		{0, 5},
		{6, 11},
		{0, 0},
	}
	specialTokenMask := []int{1, 0, 0, 1}
	attentionMask := []int{1, 1, 1, 1}

	sr := make(map[int]tokenizer.Range)
	sr[0] = tokenizer.NewRange(1, 3)
	sequenceRange := tokenizer.WithSequenceRangeEncodingOpt(sr)

	wantOverflowing := tokenizer.NewEncoding(
		[]int{1, 13, 0},
		[]int{0, 0, 0},
		[]string{"[CLS]", "you", "[SEP]"},
		[][]int{{0, 0}, {12, 15}, {0, 0}},
		[]int{1, 0, 1},
		[]int{1, 1, 1},
		nil,
		tokenizer.WithSequenceRangeEncodingOpt(map[int]tokenizer.Range{0: tokenizer.NewRange(1, 2)}),
	)

	want1 := tokenizer.NewEncoding(ids, typeIds, tokens, offsets, specialTokenMask, attentionMask, []tokenizer.Encoding{*wantOverflowing}, sequenceRange)
	got1 := singleEncoding

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("want %#v, \ngot %#v\n", want1, got1)
	}

	got2, ok := singleEncoding.Token2Sequence(2)
	want2 := 0
	if !ok || got2 != want2 {
		t.Errorf("want %v, got %v\n", want2, got2)
	}

	got2, ok = singleEncoding.Token2Sequence(3)
	want2 = 0
	if ok || got2 == want2 {
		t.Errorf("want %v, got %v\n", want2, got2)
	}

	pair := tokenizer.NewEncodingFromTokens([]tokenizer.Token{
		{Id: 15, Value: "pair", Offsets: []int{0, 4}},
		{Id: 16, Value: "with", Offsets: []int{5, 9}},
	}, 0)

	pairOverflowing := tokenizer.NewEncodingFromTokens([]tokenizer.Token{
		{Id: 17, Value: "info", Offsets: []int{10, 14}},
	}, 0)

	pair.SetOverflowing([]tokenizer.Encoding{*pairOverflowing})

	gotPairEncoding := processor.Process(encoding, pair, true)

	ids = []int{1, 12, 14, 0, 15, 16, 0}
	tokens = []string{
		"[CLS]",
		"Hello",
		"there",
		"[SEP]",
		"pair",
		"with",
		"[SEP]",
	}

	typeIds = []int{0, 0, 0, 0, 1, 1, 1}
	offsets = [][]int{
		{0, 0},
		{0, 5},
		{6, 11},
		{0, 0},
		{0, 4},
		{5, 9},
		{0, 0},
	}

	specialTokenMask = []int{1, 0, 0, 1, 0, 0, 1}
	attentionMask = []int{1, 1, 1, 1, 1, 1, 1}

	sr = make(map[int]tokenizer.Range)
	sr[0] = tokenizer.NewRange(1, 3)
	sr[1] = tokenizer.NewRange(4, 6)
	sequenceRange = tokenizer.WithSequenceRangeEncodingOpt(sr)

	wantPairOverflowing := []tokenizer.Encoding{
		*tokenizer.NewEncoding(
			[]int{1, 13, 0, 15, 16, 0},
			[]int{0, 0, 0, 1, 1, 1},
			[]string{"[CLS]", "you", "[SEP]", "pair", "with", "[SEP]"},
			[][]int{{0, 0}, {12, 15}, {0, 0}, {0, 4}, {5, 9}, {0, 0}},
			[]int{1, 0, 1, 0, 0, 1},
			[]int{1, 1, 1, 1, 1, 1},
			[]tokenizer.Encoding{
				*tokenizer.NewEncoding(
					[]int{1, 13, 0, 17, 0},
					[]int{0, 0, 0, 0, 1},
					[]string{"[CLS]", "you", "[SEP]", "info", "[SEP]"},
					[][]int{{0, 0}, {12, 15}, {0, 0}, {10, 14}, {0, 0}},
					[]int{1, 0, 1, 0, 1},
					[]int{1, 1, 1, 1, 1},
					nil,
					tokenizer.WithSequenceRangeEncodingOpt(map[int]tokenizer.Range{
						0: tokenizer.NewRange(1, 2),
						1: tokenizer.NewRange(3, 4),
					}),
				),
			},
			tokenizer.WithSequenceRangeEncodingOpt(map[int]tokenizer.Range{
				0: tokenizer.NewRange(1, 2),
				1: tokenizer.NewRange(3, 5),
			})),

		*tokenizer.NewEncoding(
			[]int{1, 13, 0, 17, 0},
			[]int{0, 0, 0, 0, 1},
			[]string{"[CLS]", "you", "[SEP]", "info", "[SEP]"},
			[][]int{{0, 0}, {12, 15}, {0, 0}, {10, 14}, {0, 0}},
			[]int{1, 0, 1, 0, 1},
			[]int{1, 1, 1, 1, 1},
			nil,
			tokenizer.WithSequenceRangeEncodingOpt(map[int]tokenizer.Range{
				0: tokenizer.NewRange(1, 2),
				1: tokenizer.NewRange(3, 4),
			})),

		*tokenizer.NewEncoding(
			[]int{1, 12, 14, 0, 17, 0},
			[]int{0, 0, 0, 0, 0, 1},
			[]string{"[CLS]", "Hello", "there", "[SEP]", "info", "[SEP]"},
			[][]int{{0, 0}, {0, 5}, {6, 11}, {0, 0}, {10, 14}, {0, 0}},
			[]int{1, 0, 0, 1, 0, 1},
			[]int{1, 1, 1, 1, 1, 1},
			[]tokenizer.Encoding{
				*tokenizer.NewEncoding(
					[]int{1, 13, 0, 17, 0},
					[]int{0, 0, 0, 0, 1},
					[]string{"[CLS]", "you", "[SEP]", "info", "[SEP]"},
					[][]int{{0, 0}, {12, 15}, {0, 0}, {10, 14}, {0, 0}},
					[]int{1, 0, 1, 0, 1},
					[]int{1, 1, 1, 1, 1},
					nil,
					tokenizer.WithSequenceRangeEncodingOpt(map[int]tokenizer.Range{
						0: tokenizer.NewRange(1, 2),
						1: tokenizer.NewRange(3, 4),
					}),
				),
			},
			tokenizer.WithSequenceRangeEncodingOpt(map[int]tokenizer.Range{
				0: tokenizer.NewRange(1, 3),
				1: tokenizer.NewRange(4, 5),
			})),
	}

	wantPairEncoding := tokenizer.NewEncoding(ids, typeIds, tokens, offsets, specialTokenMask, attentionMask, wantPairOverflowing, sequenceRange)

	if !reflect.DeepEqual(wantPairEncoding, gotPairEncoding) {
		t.Errorf("\nwant %#v, \ngot %#v", wantPairEncoding, gotPairEncoding)
	}
}
