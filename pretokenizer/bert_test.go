package pretokenizer_test

import (
	"reflect"
	"testing"

	tokenizer "github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretokenizer"
)

func TestBertPreTokenize(t *testing.T) {

	var preTok pretokenizer.BertPreTokenizer

	input := tokenizer.NewPreTokenizedString("Hey friend!     How are you?!?")

	pretokenized, err := preTok.PreTokenize(input)
	if err != nil {
		t.Error(err)
	}

	var got []tokenizer.PreToken

	for {
		sub, ok := pretokenized.Next()
		if !ok {
			break
		}

		pretok := tokenizer.PreToken{
			Value:   (*sub.Normalized).GetNormalized(),
			Offsets: sub.OriginalOffsets,
		}
		got = append(got, pretok)
	}

	want := &[]tokenizer.PreToken{
		{Value: "Hey", Offsets: tokenizer.Offsets{Start: 0, End: 3}},
		{Value: "friend", Offsets: tokenizer.Offsets{Start: 4, End: 10}},
		{Value: "!", Offsets: tokenizer.Offsets{Start: 10, End: 11}},
		{Value: "How", Offsets: tokenizer.Offsets{Start: 16, End: 19}},
		{Value: "are", Offsets: tokenizer.Offsets{Start: 20, End: 23}},
		{Value: "you", Offsets: tokenizer.Offsets{Start: 24, End: 27}},
		{Value: "?", Offsets: tokenizer.Offsets{Start: 27, End: 28}},
		{Value: "!", Offsets: tokenizer.Offsets{Start: 28, End: 29}},
		{Value: "?", Offsets: tokenizer.Offsets{Start: 29, End: 30}},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want:\n%v\n Got:\n%v\n", want, got)
	}
}
