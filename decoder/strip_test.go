package decoder

import (
	"reflect"
	"testing"
)

func TestStrip_DecodeChain(t *testing.T) {
	dec := NewStrip("H", 1, 0)

	tokens := []string{
		"Hey",
		"friend!",
		"HHH",
	}
	got1 := dec.DecodeChain(tokens)
	want1 := []string{
		"ey",
		"friend!",
		"HH",
	}

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want %v, got %v\n", want1, got1)
	}

	dec = NewStrip("y", 0, 1)
	tokens = []string{
		"Hey",
		"friend!",
	}

	got2 := dec.DecodeChain(tokens)
	want2 := []string{
		"He",
		"friend!",
	}

	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want %v, got %v\n", want2, got2)
	}
}
