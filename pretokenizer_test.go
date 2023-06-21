package tokenizer

import (
	"reflect"
	"testing"
)

func TestBytesToCharConverter(t *testing.T) {

	sequence := "Löwe 老虎 Léopard"
	converter := NewBytesToCharOffsetConverter(sequence)

	want := map[int]int{
		0:  0,
		1:  1,
		2:  1,
		3:  2,
		4:  3,
		5:  4,
		6:  5,
		7:  5,
		8:  5,
		9:  6,
		10: 6,
		11: 6,
		12: 7,
		13: 8,
		14: 9,
		15: 9,
		16: 10,
		17: 11,
		18: 12,
		19: 13,
		20: 14,
	}

	got := converter.b2c

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}
}
