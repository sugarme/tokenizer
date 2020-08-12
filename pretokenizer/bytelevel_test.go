package pretokenizer_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/pretokenizer"
	"github.com/sugarme/tokenizer/tokenizer"
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

		normalized, res := bytelevel.PreTokenize(&n)

		nwant := "ĠHelloĠmyĠfriend,ĠhowĠisĠyourĠdayĠgoing?"
		ngot := normalized.GetNormalized()

		pwant := []tokenizer.PreToken{
			{Value: "ĠHello", Offsets: tokenizer.Offsets{Start: 0, End: 6}},
			{Value: "Ġmy", Offsets: tokenizer.Offsets{Start: 6, End: 9}},
			{Value: "Ġfriend", Offsets: tokenizer.Offsets{Start: 9, End: 16}},
			{Value: ",", Offsets: tokenizer.Offsets{Start: 16, End: 17}},
			{Value: "Ġhow", Offsets: tokenizer.Offsets{Start: 17, End: 21}},
			{Value: "Ġis", Offsets: tokenizer.Offsets{Start: 21, End: 24}},
			{Value: "Ġyour", Offsets: tokenizer.Offsets{Start: 24, End: 29}},
			{Value: "Ġday", Offsets: tokenizer.Offsets{Start: 29, End: 33}},
			{Value: "Ġgoing", Offsets: tokenizer.Offsets{Start: 33, End: 39}},
			{Value: "?", Offsets: tokenizer.Offsets{Start: 39, End: 40}},
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
		var normalized normalizer.NormalizedString
		normalized = normalizer.NewNormalizedFrom(l)

		_, preTokenized := bytelevel.PreTokenize(&normalized)

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

	var normalized normalizer.NormalizedString
	normalized = normalizer.NewNormalizedFrom("Hello there\nHello there")

	_, preTokenized := bytelevel.PreTokenize(&normalized)

	var separatedTokens []string
	for _, preTok := range *preTokenized {
		chars := strings.Split(preTok.Value, "")
		separatedTokens = append(separatedTokens, chars...)
	}

	want := []tokenizer.PreToken{
		{Value: "Hello", Offsets: tokenizer.Offsets{Start: 0, End: 5}},
		{Value: "Ġthere", Offsets: tokenizer.Offsets{Start: 5, End: 11}},
		{Value: "Ċ", Offsets: tokenizer.Offsets{Start: 11, End: 12}},
		{Value: "Hello", Offsets: tokenizer.Offsets{Start: 12, End: 17}},
		{Value: "Ġthere", Offsets: tokenizer.Offsets{Start: 17, End: 23}},
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

	var normalized normalizer.NormalizedString
	normalized = normalizer.NewNormalizedFrom("Hello there       dear")

	_, preTokenized := bytelevel.PreTokenize(&normalized)

	var separatedTokens []string
	for _, preTok := range *preTokenized {
		chars := strings.Split(preTok.Value, "")
		separatedTokens = append(separatedTokens, chars...)
	}

	want := []tokenizer.PreToken{
		{Value: "Hello", Offsets: tokenizer.Offsets{Start: 0, End: 5}},
		{Value: "Ġthere", Offsets: tokenizer.Offsets{Start: 5, End: 11}},
		{Value: "ĠĠĠĠĠĠ", Offsets: tokenizer.Offsets{Start: 11, End: 17}},
		{Value: "Ġdear", Offsets: tokenizer.Offsets{Start: 17, End: 22}},
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

	var normalized normalizer.NormalizedString
	normalized = normalizer.NewNormalizedFrom("i⭢j")

	_, preTokenized := bytelevel.PreTokenize(&normalized)

	var separatedTokens []string
	for _, preTok := range *preTokenized {
		chars := strings.Split(preTok.Value, "")
		separatedTokens = append(separatedTokens, chars...)
	}

	want := []tokenizer.PreToken{
		{Value: "i", Offsets: tokenizer.Offsets{Start: 0, End: 1}},
		{Value: "ŸŃ¢", Offsets: tokenizer.Offsets{Start: 1, End: 4}},
		{Value: "j", Offsets: tokenizer.Offsets{Start: 4, End: 5}},
	}
	got := *preTokenized

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

func TestProcessorTrimsOffsets(t *testing.T) {

	normalized := normalizer.NewNormalizedFrom("")

	start := tokenizer.NewEncoding(
		normalized,
		[]uint32{}, []uint32{}, []string{
			"ĠĠĠĠHelloĠĠ",
			"ĠĠHello",
			"HelloĠĠ",
			"ĠĠĠĠ",
		},
		[]tokenizer.Offsets{
			{Start: 0, End: 11},
			{Start: 11, End: 18},
			{Start: 18, End: 25},
			{Start: 25, End: 29},
		},
		[]uint32{}, []uint32{},
		[]tokenizer.Encoding{},
	)

	want := tokenizer.NewEncoding(
		normalized,
		[]uint32{}, []uint32{}, []string{
			"ĠĠĠĠHelloĠĠ",
			"ĠĠHello",
			"HelloĠĠ",
			"ĠĠĠĠ",
		},
		[]tokenizer.Offsets{
			{Start: 4, End: 9},
			{Start: 13, End: 18},
			{Start: 18, End: 23},
			{Start: 29, End: 29},
		},
		[]uint32{}, []uint32{},
		[]tokenizer.Encoding{},
	)

	bytelevel := pretokenizer.NewByteLevel()
	bytelevel.SetTrimOffsets(true)

	got := bytelevel.Process(start, false)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}

	pairWant := want
	pairWant.MergeWith(want)

	pairGot := bytelevel.Process(start, false, start)

	if !reflect.DeepEqual(pairWant, pairGot) {
		t.Errorf("Want: %v\n", pairWant)
		t.Errorf("Got: %v\n", pairGot)
	}
}
