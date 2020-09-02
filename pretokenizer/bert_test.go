package pretokenizer_test

import (
	"reflect"
	"testing"

	tokenizer "github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/pretokenizer"
)

func TestBertPreTokenize(t *testing.T) {

	var bertPreTok pretokenizer.BertPreTokenizer

	input := tokenizer.NewPreTokenizedString("Hey friend!     How are you?!?")

	pretokenized, err := bertPreTok.PreTokenize(input)
	if err != nil {
		t.Error(err)
	}

	pretoks := pretokenized.GetSplits(normalizer.OriginalTarget)

	var got []struct {
		Value   string
		Offsets []int
	}

	for _, pretok := range pretoks {
		got = append(got, struct {
			Value   string
			Offsets []int
		}{pretok.Value, pretok.Offsets})
	}

	want := []struct {
		Value   string
		Offsets []int
	}{
		{Value: "Hey", Offsets: []int{0, 3}},
		{Value: "friend", Offsets: []int{4, 10}},
		{Value: "!", Offsets: []int{10, 11}},
		{Value: "How", Offsets: []int{16, 19}},
		{Value: "are", Offsets: []int{20, 23}},
		{Value: "you", Offsets: []int{24, 27}},
		{Value: "?", Offsets: []int{27, 28}},
		{Value: "!", Offsets: []int{28, 29}},
		{Value: "?", Offsets: []int{29, 30}},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want:\n%v\n Got:\n%v\n", want, got)
	}
}
