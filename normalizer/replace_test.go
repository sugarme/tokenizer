package normalizer

import (
	"reflect"
	"testing"
)

func TestReplace_Normalize(t *testing.T) {
	// 1. RegexReplace
	original := "This     is   a         test"
	n := NewNormalizedFrom(original)

	r := NewReplace(Regex, `\s+`, " ")

	out, err := r.Normalize(n)
	if err != nil {
		panic(err)
	}

	got := out.GetNormalized()
	want := "This is a test"

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}

	// 2. StringReplace
	original = "This is a ''test''"
	n = NewNormalizedFrom(original)

	r = NewReplace(String, "''", "\"")
	out, err = r.Normalize(n)

	got = out.GetNormalized()
	want = "This is a \"test\""

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}
}

func TestReplace_Decode(t *testing.T) {
	tokens := []string{
		"hello",
		"_hello",
	}

	r := NewReplace(String, "_", " ")

	got := r.DecodeChain(tokens)
	want := []string{
		"hello",
		" hello",
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}
}
