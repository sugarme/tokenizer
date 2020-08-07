package wordpiece_test

import (
	"reflect"
	"testing"

	"github.com/sugarme/tokenizer/model/wordpiece"
)

func TestErrorDisplay(t *testing.T) {

	want := "WordPiece error: Missing [UNK] token from the vocabulary\n"
	got := wordpiece.MissingUnkToken.Error()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want %v - got %v\n", want, got)
	}
}
