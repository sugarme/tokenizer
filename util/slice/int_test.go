package util_test

import (
	"reflect"
	"testing"

	slice "github.com/sugarme/tokenizer/util/slice"
)

func TestInsertInt(t *testing.T) {

	a := []int{3, 5, 7, 9}
	want := []int{3, 5, 0, 7, 9}

	var item int = 0
	got, err := slice.InsertInt(a, item, 2)
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}

}
