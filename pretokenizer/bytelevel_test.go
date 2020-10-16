package pretokenizer_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/pretokenizer"
)

type charidx struct {
	s string
	o []int
}

func TestBytesChar(t *testing.T) {

	want := "!"

	bc := pretokenizer.GenerateBytesChar()
	got := bc[33]

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}

	// Test generating bytesChar map
	bcGot := pretokenizer.BytesChar[33]
	// bcWant := make(map[uint8]string)
	bcWant := "!"
	if !reflect.DeepEqual(bcWant, bcGot) {
		t.Errorf("Want: %v\n", bcWant)
		t.Errorf("Got: %v\n", bcGot)
	}

	// Testing generate charBytes map
	cbGot := pretokenizer.CharBytes["!"]
	var cbWant uint8 = 33
	if !reflect.DeepEqual(cbWant, cbGot) {
		t.Errorf("Want: %v\n", cbWant)
		t.Errorf("Got: %v\n", cbGot)
	}

}

func TestDecoding(t *testing.T) {

	bytelevel := pretokenizer.NewByteLevel()
	bytelevel.SetAddPrefixSpace(false)

	want := "Hello my friend, how is your day going?"

	toks := []string{
		"Hello",
		"Ġmy",
		"Ġfriend",
		",",
		"Ġhow",
		"Ġis",
		"Ġyour",
		"Ġday",
		"Ġgoing",
		"?",
	}

	got := bytelevel.Decode(toks)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

func TestAddPrefixSpace(t *testing.T) {

	bytelevel := pretokenizer.NewByteLevel()
	bytelevel.SetAddPrefixSpace(true)

	lines := []string{
		" Hello my friend, how is your day going?",
		"Hello my friend, how is your day going?",
	}

	for _, l := range lines {
		pretokenized := tokenizer.NewPreTokenizedString(l)
		pretok, err := bytelevel.PreTokenize(pretokenized)
		if err != nil {
			t.Error(err)
		}

		pretokens := pretok.GetSplits(normalizer.NormalizedTarget)

		var want, got []charidx

		for _, pretoken := range pretokens {
			got = append(got, charidx{pretoken.Value, pretoken.Offsets})
		}

		want = []charidx{
			{"ĠHello", []int{0, 7}},
			{"Ġmy", []int{7, 11}},
			{"Ġfriend", []int{11, 19}},
			{",", []int{19, 20}},
			{"Ġhow", []int{20, 25}},
			{"Ġis", []int{25, 29}},
			{"Ġyour", []int{29, 35}},
			{"Ġday", []int{35, 40}},
			{"Ġgoing", []int{40, 47}},
			{"?", []int{47, 48}},
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("Want: %v\n", want)
			t.Errorf("Got: %v\n", got)
		}
	}

}

func TestDecodeWorksOnSeparatedTokens(t *testing.T) {

	bytelevel := pretokenizer.NewByteLevel()
	bytelevel.SetAddPrefixSpace(false)

	lines := []string{
		"A Nuskhuri abbreviation of იესუ ქრისტე ( iesu kriste ) \" Jesus Christ \"",
		"An equal number have descenders , like p or q in English : გ , დ , ე , ვ , კ , ლ , ჟ , ტ , უ , ფ , ღ , ყ , ც",
	}

	for _, l := range lines {

		pretokenized := tokenizer.NewPreTokenizedString(l)
		pretok, err := bytelevel.PreTokenize(pretokenized)
		if err != nil {
			t.Error(err)
		}

		var separatedTokens []string
		for _, preTok := range pretok.GetSplits(normalizer.OriginalTarget) {
			chars := strings.Split(preTok.Value, "")
			separatedTokens = append(separatedTokens, chars...)
		}

		want := l
		got := bytelevel.Decode(separatedTokens)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("Want: %v\n", want)
			t.Errorf("Got: %v\n", got)
		}
	}
}

func TestHandlingOfNewLines(t *testing.T) {

	bytelevel := pretokenizer.NewByteLevel()
	bytelevel.SetAddPrefixSpace(false)

	pretokenized := tokenizer.NewPreTokenizedString("Hello there\nHello there")

	pretok, err := bytelevel.PreTokenize(pretokenized)
	if err != nil {
		t.Error(err)
	}

	var got []charidx

	for _, preTok := range pretok.GetSplits(normalizer.OriginalTarget) {
		got = append(got, charidx{s: preTok.Value, o: preTok.Offsets})
	}

	want := []charidx{
		{"Hello", []int{0, 5}},
		{"Ġthere", []int{5, 11}},
		{"Ċ", []int{11, 12}},
		{"Hello", []int{12, 17}},
		{"Ġthere", []int{17, 23}},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

func TestHandlingOfMultipleSpaces(t *testing.T) {

	bytelevel := pretokenizer.NewByteLevel()
	bytelevel.SetAddPrefixSpace(false)

	pretokenized := tokenizer.NewPreTokenizedString("Hello there       dear")

	pretok, err := bytelevel.PreTokenize(pretokenized)
	if err != nil {
		t.Error(err)
	}

	var got []charidx

	for _, preTok := range pretok.GetSplits(normalizer.OriginalTarget) {
		got = append(got, charidx{s: preTok.Value, o: preTok.Offsets})
	}

	want := []charidx{
		{"Hello", []int{0, 5}},
		{"Ġthere", []int{5, 11}},
		{"ĠĠĠĠĠĠ", []int{11, 17}},
		{"Ġdear", []int{17, 22}},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

func TestOffsetsWhenCharSplitUp(t *testing.T) {

	bytelevel := pretokenizer.NewByteLevel()
	bytelevel.SetAddPrefixSpace(false)

	input := "i⭢j"
	pretokenized := tokenizer.NewPreTokenizedString(input)

	pretok, err := bytelevel.PreTokenize(pretokenized)
	if err != nil {
		t.Error(err)
	}

	var got1 []charidx

	for _, preTok := range pretok.GetSplits(normalizer.OriginalTarget) {
		got1 = append(got1, charidx{s: preTok.Value, o: preTok.Offsets})
	}

	want1 := []charidx{
		{"i", []int{0, 1}},
		{"âŃ¢", []int{1, 4}},
		{"j", []int{4, 5}},
	}

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}

	var got2 []charidx

	for _, preTok := range pretok.GetSplits(normalizer.NormalizedTarget) {
		got2 = append(got2, charidx{s: preTok.Value, o: preTok.Offsets})
	}

	want2 := []charidx{
		{"i", []int{0, 1}},
		{"âŃ¢", []int{1, 7}},
		{"j", []int{7, 8}},
	}

	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want: %v\n", want2)
		t.Errorf("Got: %v\n", got2)
	}

	var got3 []string

	for _, preTok := range pretok.GetSplits(normalizer.OriginalTarget) {
		o := preTok.Offsets
		got3 = append(got3, input[o[0]:o[1]])
	}

	want3 := []string{"i", "⭢", "j"}

	if !reflect.DeepEqual(want3, got3) {
		t.Errorf("Want: %v\n", want3)
		t.Errorf("Got: %v\n", got3)
	}
}

func TestProcessorTrimsOffsets(t *testing.T) {
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

	wantOffsets := [][]int{
		{0, 0},
		{4, 9},
		{13, 18},
		{18, 23},
		{29, 29},
	}

	start := tokenizer.NewEncoding(nil, nil, tokens, offsets, nil, nil, nil)
	want := tokenizer.NewEncoding(nil, nil, tokens, wantOffsets, nil, nil, nil)

	bytelevel := pretokenizer.NewByteLevel()
	bytelevel.SetTrimOffsets(true)

	got := bytelevel.Process(start, nil, false)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}

	pairWant := want
	pair := tokenizer.NewEncoding(nil, nil, tokens, wantOffsets, nil, nil, nil)
	pairWant.MergeWith(pair, false)

	start1 := tokenizer.NewEncoding(nil, nil, tokens, offsets, nil, nil, nil)
	startClone := tokenizer.NewEncoding(nil, nil, tokens, offsets, nil, nil, nil)
	pairGot := bytelevel.Process(startClone, start1, false)

	if !reflect.DeepEqual(pairWant, pairGot) {
		t.Errorf("Want: %+v\n", pairWant)
		t.Errorf("Got: %+v\n", pairGot)
	}
}
