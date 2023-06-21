package decoder

import (
	"reflect"
	"strings"
	"testing"
)

func TestCTC1(t *testing.T) {
	dec := DefaultCTC()

	tokens := strings.Split("<pad> <pad> h e e l l <pad> l o o o <pad>", " ")

	got := dec.DecodeChain(tokens)
	want := []string{"h", "e", "l", "l", "o"}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %q got %q\n", want, got)
	}
}

func TestCTC2(t *testing.T) {
	dec := DefaultCTC()

	tokens := strings.Split("<pad> <pad> h e e l l <pad> l o o o <pad> <pad> | <pad> w o o o r <pad> <pad> l l d <pad> <pad> <pad> <pad>", " ")

	got := dec.DecodeChain(tokens)
	want := []string{"h", "e", "l", "l", "o", " ", "w", "o", "r", "l", "d"}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %q got %q\n", want, got)
	}
}

func TestCTC3(t *testing.T) {
	dec := DefaultCTC()

	str := "<pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> A | | <pad> M <pad> <pad> <pad> <pad> A <pad> <pad> N <pad> <pad> <pad> | | | <pad> <pad> <pad> <pad> S <pad> <pad> <pad> A I <pad> D D | | T T <pad> O <pad> | | T H E E | | | <pad> U U <pad> N N <pad> I <pad> <pad> V <pad> <pad> <pad> E R R <pad> <pad> <pad> S E E | | <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> S S <pad> <pad> <pad> <pad> I <pad> R R <pad> <pad> | | | <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> I <pad> <pad> <pad> | <pad> <pad> <pad> E X <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> I <pad> S <pad> <pad> T <pad> <pad> | | <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad> <pad>"

	tokens := strings.Split(str, " ")

	got := dec.DecodeChain(tokens)
	want := []string{
		"A", " ", "M", "A", "N", " ", "S", "A", "I", "D", " ", "T", "O", " ", "T", "H",
		"E", " ", "U", "N", "I", "V", "E", "R", "S", "E", " ", "S", "I", "R", " ", "I",
		" ", "E", "X", "I", "S", "T", " ",
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %q got %q\n", want, got)
	}
}
