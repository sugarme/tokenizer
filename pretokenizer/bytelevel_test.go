package pretokenizer_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/pretokenizer"
)

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
		var normalized *normalizer.NormalizedString
		n := normalizer.NewNormalizedFrom(l)

		normalized, res := bytelevel.PreTokenize(n)

		nwant := "ĠHelloĠmyĠfriend,ĠhowĠisĠyourĠdayĠgoing?"
		ngot := normalized.GetNormalized()

		pwant := []tokenizer.PreToken{
			{Value: "ĠHello", Offsets: []int{0, 6}},
			{Value: "Ġmy", Offsets: []int{6, 9}},
			{Value: "Ġfriend", Offsets: []int{9, 16}},
			{Value: ",", Offsets: []int{16, 17}},
			{Value: "Ġhow", Offsets: []int{17, 21}},
			{Value: "Ġis", Offsets: []int{21, 24}},
			{Value: "Ġyour", Offsets: []int{24, 29}},
			{Value: "Ġday", Offsets: []int{29, 33}},
			{Value: "Ġgoing", Offsets: []int{33, 39}},
			{Value: "?", Offsets: []int{39, 40}},
		}

		pgot := *res

		if !reflect.DeepEqual(nwant, ngot) {
			t.Errorf("nWant: %v\n", nwant)
			t.Errorf("nGot: %v\n", ngot)
		}

		if !reflect.DeepEqual(pwant, pgot) {
			t.Errorf("pWant: %v\n", pwant)
			t.Errorf("pGot: %v\n", pgot)
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
		var normalized *normalizer.NormalizedString
		normalized = normalizer.NewNormalizedFrom(l)

		_, preTokenized := bytelevel.PreTokenize(normalized)

		var separatedTokens []string
		for _, preTok := range *preTokenized {
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

	var normalized *normalizer.NormalizedString
	normalized = normalizer.NewNormalizedFrom("Hello there\nHello there")

	_, preTokenized := bytelevel.PreTokenize(normalized)

	var separatedTokens []string
	for _, preTok := range *preTokenized {
		chars := strings.Split(preTok.Value, "")
		separatedTokens = append(separatedTokens, chars...)
	}

	want := []tokenizer.PreToken{
		{Value: "Hello", Offsets: []int{0, 5}},
		{Value: "Ġthere", Offsets: []int{5, 11}},
		{Value: "Ċ", Offsets: []int{11, 12}},
		{Value: "Hello", Offsets: []int{12, 17}},
		{Value: "Ġthere", Offsets: []int{17, 23}},
	}
	got := *preTokenized

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

func TestHandlingOfMultipleSpaces(t *testing.T) {

	bytelevel := pretokenizer.NewByteLevel()
	bytelevel.SetAddPrefixSpace(false)

	var normalized *normalizer.NormalizedString
	normalized = normalizer.NewNormalizedFrom("Hello there       dear")

	_, preTokenized := bytelevel.PreTokenize(normalized)

	var separatedTokens []string
	for _, preTok := range *preTokenized {
		chars := strings.Split(preTok.Value, "")
		separatedTokens = append(separatedTokens, chars...)
	}

	want := []tokenizer.PreToken{
		{Value: "Hello", Offsets: []int{0, 5}},
		{Value: "Ġthere", Offsets: []int{5, 11}},
		{Value: "ĠĠĠĠĠĠ", Offsets: []int{11, 17}},
		{Value: "Ġdear", Offsets: []int{17, 22}},
	}
	got := *preTokenized

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

func TestOffsetsWhenCharSplitUp(t *testing.T) {

	bytelevel := pretokenizer.NewByteLevel()
	bytelevel.SetAddPrefixSpace(false)

	var normalized *normalizer.NormalizedString
	normalized = normalizer.NewNormalizedFrom("i⭢j")

	_, preTokenized := bytelevel.PreTokenize(normalized)

	var separatedTokens []string
	for _, preTok := range *preTokenized {
		chars := strings.Split(preTok.Value, "")
		separatedTokens = append(separatedTokens, chars...)
	}

	want := []tokenizer.PreToken{
		{Value: "i", Offsets: []int{0, 1}},
		{Value: "ŸŃ¢", Offsets: []int{1, 4}},
		{Value: "j", Offsets: []int{4, 5}},
	}
	got := *preTokenized

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

func TestProcessorTrimsOffsets(t *testing.T) {

	start := tokenizer.NewEncoding(
		[]int{}, []int{}, []string{
			"ĠĠĠĠHelloĠĠ",
			"ĠĠHello",
			"HelloĠĠ",
			"ĠĠĠĠ",
		},
		[][]int{
			{0, 11},
			{11, 18},
			{18, 25},
			{25, 29},
		},
		[]int{}, []int{},
		[]tokenizer.Encoding{},
	)

	want := tokenizer.NewEncoding(
		[]int{}, []int{}, []string{
			"ĠĠĠĠHelloĠĠ",
			"ĠĠHello",
			"HelloĠĠ",
			"ĠĠĠĠ",
		},
		[][]int{
			{4, 9},
			{13, 18},
			{18, 23},
			{29, 29},
		},
		[]int{}, []int{},
		[]tokenizer.Encoding{},
	)

	bytelevel := pretokenizer.NewByteLevel()
	bytelevel.SetTrimOffsets(true)

	got := bytelevel.Process(start, nil, false)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}

	pairWant := want
	pairWant.MergeWith(want, true)

	pairGot := bytelevel.Process(start, start, false)

	if !reflect.DeepEqual(pairWant, pairGot) {
		t.Errorf("Want: %v\n", pairWant)
		t.Errorf("Got: %v\n", pairGot)
	}
}
