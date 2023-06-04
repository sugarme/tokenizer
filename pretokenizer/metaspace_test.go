package pretokenizer

import (
	"reflect"
	"testing"
)

func TestMetaspace_Decode(t *testing.T) {
	dec := DefaultMetaspace()
	// dec := NewMetaspace("_", true)

	tokens := []string{
		"_Hey",
		"_friend!",
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

/*
 */
