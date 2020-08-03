package util_test

import (
	"reflect"
	"testing"

	slice "github.com/sugarme/tokenizer/util/slice"
)

func TestInsertStr(t *testing.T) {

	a := []string{"A", "B", "C", "D", "E"}
	want := []string{"A", "B", "x", "C", "D", "E"}

	var item string = "x"
	got, err := slice.InsertStr(a, item, 2)
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}

}
