package pretokenizer_test

import (
	"reflect"
	"testing"

	tokenizer "github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/pretokenizer"
)

func TestBertPreTokenize(t *testing.T) {

	var preTok pretokenizer.BertPreTokenizer

	input := tokenizer.NewPreTokenizedString("Hey friend!     How are you?!?")

	pretokenized, err := preTok.PreTokenize(input)
	if err != nil {
		t.Error(err)
	}

	got := pretokenized.GetNormalized(normalizer.OriginalTarget)

	want := []tokenizer.PreToken{
		{Value: "Hey", Offsets: tokenizer.Offsets{Start: 0, End: 3}},
		{Value: "", Offsets: tokenizer.Offsets{Start: 3, End: 4}},
		{Value: "friend", Offsets: tokenizer.Offsets{Start: 4, End: 10}},
		{Value: "!", Offsets: tokenizer.Offsets{Start: 10, End: 11}},
		{Value: "", Offsets: tokenizer.Offsets{Start: 11, End: 12}},
		{Value: "", Offsets: tokenizer.Offsets{Start: 12, End: 13}},
		{Value: "", Offsets: tokenizer.Offsets{Start: 13, End: 14}},
		{Value: "", Offsets: tokenizer.Offsets{Start: 14, End: 15}},
		{Value: "", Offsets: tokenizer.Offsets{Start: 15, End: 16}},
		{Value: "How", Offsets: tokenizer.Offsets{Start: 16, End: 19}},
		{Value: "", Offsets: tokenizer.Offsets{Start: 19, End: 20}},
		{Value: "are", Offsets: tokenizer.Offsets{Start: 20, End: 23}},
		{Value: "", Offsets: tokenizer.Offsets{Start: 23, End: 24}},
		{Value: "you", Offsets: tokenizer.Offsets{Start: 24, End: 27}},
		{Value: "?", Offsets: tokenizer.Offsets{Start: 27, End: 28}},
		{Value: "!", Offsets: tokenizer.Offsets{Start: 28, End: 29}},
		{Value: "?", Offsets: tokenizer.Offsets{Start: 29, End: 30}},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want:\n%v\n Got:\n%v\n", want, got)
	}
}
