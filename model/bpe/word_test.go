package bpe_test

import (
	"reflect"
	"testing"

	bpe "github.com/sugarme/tokenizer/model/bpe"
)

func TestMerge_Merge(t *testing.T) {
	// Let's say we have the word 'hello' and a word-to-id vocab that looks
	// like this: {'h': 0, 'e': 1, 'l': 2, 'o': 3}.
	word := bpe.NewWord()
	word.Add(0, 1) // 'h'
	word.Add(1, 1) // 'e'
	word.Add(2, 1) // 'l'
	word.Add(2, 1) // 'l'
	word.Add(3, 1) // 'o'

	// We're going to perform a merge on the pair ('l', 'l') ~= (2, 2). Let's
	// say that 'll' has the ID of 4 in the updated word-to-id vocab.
	changesGot, err := word.Merge(2, 2, 4)
	if err != nil {
		t.Errorf("Err: %v", err)
	}

	// So the word should now look like this:
	charsGot := word.GetChars()
	charsWant := []int{
		0, // 'h'
		1, // 'e'
		4, // 'll'
		3, // 'o'
	}

	// The return value `changes` will be used to update the pair counts during
	// training. This merge affects the counts for the pairs
	// ('e', 'l') ~= (1, 2),
	// ('e', 'll') ~= (1, 4),
	// ('l', 'o') ~= (2, 3), and
	// ('ll', 'o') ~= (4, 3).
	// So the changes should reflect that:
	changesWant := []bpe.WChange{
		{1, 2, -1}, // count for ('e', 'l') should be decreased by 1.
		{1, 4, 1},  // count for ('e', 'll') should be increased by 1.
		{2, 3, -1}, // count for ('l', 'o') should be decreased by 1.
		{4, 3, 1},  // count for ('ll', 'o') should be increased by 1.
	}

	if !reflect.DeepEqual(charsWant, charsGot) {
		t.Errorf("Want: %v\n", charsWant)
		t.Errorf("Got: %v\n", charsGot)
	}

	if !reflect.DeepEqual(changesWant, changesGot) {
		t.Errorf("Want: %v\n", changesWant)
		t.Errorf("Got: %v\n", changesGot)
	}
}
