package decoder

import (
	// "fmt"
	"reflect"
	"testing"
)

func TestByteFallbackDecodeChain(t *testing.T) {
	dec := NewByteFallback()

	tokens := []string{
		"<0x61>",
	}
	got1 := dec.DecodeChain(tokens)
	want1 := []string{"a"}

	if !reflect.DeepEqual(got1, want1) {
		t.Errorf("Want %v - Got %v\n", want1, got1)
	}

	tokens = []string{
		"<0xE5>",
	}
	got2 := dec.DecodeChain(tokens)
	want2 := []string{"�"}
	if !reflect.DeepEqual(got2, want2) {
		t.Errorf("Want %x - Got %x\n", want2, got2)
	}

	tokens = []string{
		"<0xE5>",
		"<0x8f>",
	}

	got3 := dec.DecodeChain(tokens)
	want3 := []string{"�", "�"}
	if !reflect.DeepEqual(got3, want3) {
		t.Errorf("Want %x - Got %x\n", want3, got3)
	}

	tokens = []string{
		"<0xE5>",
		"<0x8f>",
		"<0xab>",
	}

	got4 := dec.DecodeChain(tokens)
	want4 := []string{"叫"}
	if !reflect.DeepEqual(got4, want4) {
		t.Errorf("Want %x - Got %x\n", want4, got4)
	}

	tokens = []string{
		"<0xE5>",
		"<0x8f>",
		"<0xab>",
		"a",
	}

	got5 := dec.DecodeChain(tokens)
	want5 := []string{"叫", "a"}
	if !reflect.DeepEqual(got5, want5) {
		t.Errorf("Want %x - Got %x\n", want5, got5)
	}

	tokens = []string{
		"<0xE5>",
		"<0x8f>",
		"a",
	}

	got6 := dec.DecodeChain(tokens)
	want6 := []string{"�", "�", "a"}
	if !reflect.DeepEqual(got6, want6) {
		t.Errorf("Want %x - Got %x\n", want6, got6)
	}
}
