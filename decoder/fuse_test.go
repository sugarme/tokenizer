package decoder

import (
	"reflect"
	"testing"
)

func TestFuse_DecodeChain(t *testing.T) {
	dec := NewFuse()

	tokens := []string{
		"Hey",
		" friend!",
	}
	got := dec.DecodeChain(tokens)
	want := []string{"Hey friend!"}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v got %v", want, got)
	}
}
