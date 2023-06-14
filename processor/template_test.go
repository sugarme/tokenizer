package processor

import (
	"reflect"
	"testing"
)

func TestPiece(t *testing.T) {
	seq0 := &SequencePiece{
		Id:     A,
		TypeId: 0,
	}

	seq0String := "$"
	got, err := NewPiece(seq0String)
	if err != nil {
		panic(err)
	}

	want := seq0

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v\n ", want, got)
	}

	seq1 := &SequencePiece{
		Id:     B,
		TypeId: 0,
	}

	seq1String := "$B"
	got, err = NewPiece(seq1String)
	if err != nil {
		panic(err)
	}

	want = seq1
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v\n ", want, got)
	}

	want = &SequencePiece{
		Id:     A,
		TypeId: 1,
	}

	got, err = NewPiece("$1")
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v\n ", want, got)
	}

	want = &SequencePiece{
		Id:     B,
		TypeId: 2,
	}

	got, err = NewPiece("$B:2")
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v\n ", want, got)
	}

	want = &SequencePiece{
		Id:     A,
		TypeId: 1,
	}

	got, err = NewPiece("$:1")
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v, got %#v\n ", want, got)
	}
}
