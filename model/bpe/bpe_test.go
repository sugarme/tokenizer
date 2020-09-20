package bpe_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	// "reflect"
	// "strings"
	"testing"

	"github.com/sugarme/tokenizer"
	bpe "github.com/sugarme/tokenizer/model/bpe"
	"github.com/sugarme/tokenizer/util"
)

func TestBPE_FromFiles(t *testing.T) {
	// Ensure `BadMerges` error is returned when there is an invalid line in the
	// merges.txt file.

	// 1. Set up vocab file
	// 1.1. Create temp vocab file
	// Ref. https://yourbasic.org/golang/temporary-file-directory/
	vf, err := ioutil.TempFile("/tmp", "vocab")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer os.Remove(vf.Name())

	// 1.2. Write some values as bytes to it
	var vocab map[string]int = make(map[string]int)
	vocab["a"] = 0
	vocab["b"] = 1
	vocab["c"] = 2
	vocab["ab"] = 3

	vocabBytes, err := json.Marshal(vocab)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	_, err = vf.Write(vocabBytes)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// 2. Setup a merge file with bad line
	// 2.1. Create temp merge file
	mf, err := ioutil.TempFile("/tmp", "merge")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer os.Remove(mf.Name())

	// 2.2. Write a bad line to it
	// First line: `#version: 0.2` is ok
	// Second line: `a b` is ok
	// Third line `c` is invalid
	badLine := []byte("#version: 0.2\na b\nc")
	_, err = mf.Write(badLine)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	_, err = bpe.NewBpeFromFiles(vf.Name(), mf.Name())

	got := util.TraceError(err)
	want := "Read merge file error: invalid data at line 1 \n"

	if util.ErrorContains(got, want) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

// Test tokenization. With dropout set to 0 tokenization is deterministic,
// so we know exactly what the result should be.
//
// To test this, we'll build a simple model to tokenize the word 'unrelated'.
func TestBPE_TokenizeWithAndWithoutDropout(t *testing.T) {
	var vocab map[string]int = make(map[string]int)
	vocab["u"] = 0
	vocab["n"] = 1
	vocab["r"] = 2
	vocab["e"] = 3
	vocab["l"] = 4
	vocab["a"] = 5
	vocab["t"] = 6
	vocab["d"] = 7
	vocab["re"] = 8
	vocab["at"] = 9
	vocab["ed"] = 10
	vocab["un"] = 11
	vocab["ated"] = 12
	vocab["rel"] = 13
	vocab["related"] = 14
	vocab["unrelated"] = 15

	var merges bpe.Merges = make(map[bpe.Pair]bpe.PairVal)
	merges[bpe.Pair{C1: vocab["r"], C2: vocab["e"]}] = bpe.PairVal{Rank: 1, NewId: vocab["re"]}
	merges[bpe.Pair{C1: vocab["a"], C2: vocab["t"]}] = bpe.PairVal{Rank: 2, NewId: vocab["at"]}
	merges[bpe.Pair{C1: vocab["e"], C2: vocab["d"]}] = bpe.PairVal{Rank: 3, NewId: vocab["ed"]}
	merges[bpe.Pair{C1: vocab["u"], C2: vocab["n"]}] = bpe.PairVal{Rank: 4, NewId: vocab["un"]}
	merges[bpe.Pair{C1: vocab["at"], C2: vocab["ed"]}] = bpe.PairVal{Rank: 5, NewId: vocab["ated"]}
	merges[bpe.Pair{C1: vocab["re"], C2: vocab["l"]}] = bpe.PairVal{Rank: 6, NewId: vocab["rel"]}
	merges[bpe.Pair{C1: vocab["rel"], C2: vocab["ated"]}] = bpe.PairVal{Rank: 7, NewId: vocab["related"]}
	merges[bpe.Pair{C1: vocab["un"], C2: vocab["related"]}] = bpe.PairVal{Rank: 8, NewId: vocab["unrelated"]}

	fmt.Printf("merges: %+v\n", merges)
	fmt.Printf("vocab: %+v\n", vocab)

	model := bpe.NewBPE(vocab, merges)

	// With no dropout:
	got, err := model.Tokenize("unrelated")
	if err != nil {
		t.Error(err)
	}
	want := []tokenizer.Token{
		{Id: 15, Value: "unrelated", Offsets: []int{0, 9}},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %v\n", want)
		t.Errorf("got: %v\n", got)
	}

	// Now set dropout to 1.0. Result should be no merges performed.
	dropout := float32(1.0)
	model.Dropout = &dropout
	got, err = model.Tokenize("unrelated")
	if err != nil {
		t.Error(err)
	}

	want = []tokenizer.Token{
		{Id: 0, Value: "u", Offsets: []int{0, 1}},
		{Id: 1, Value: "n", Offsets: []int{1, 2}},
		{Id: 2, Value: "r", Offsets: []int{2, 3}},
		{Id: 3, Value: "e", Offsets: []int{3, 4}},
		{Id: 4, Value: "l", Offsets: []int{4, 5}},
		{Id: 5, Value: "a", Offsets: []int{5, 6}},
		{Id: 6, Value: "t", Offsets: []int{6, 7}},
		{Id: 3, Value: "e", Offsets: []int{7, 8}},
		{Id: 7, Value: "d", Offsets: []int{8, 9}},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %v\n", want)
		t.Errorf("got: %v\n", got)
	}

	// Now try with dropout between 0 and 1.
	dropout = float32(0.5)
	tokens, err := model.Tokenize("unrelated")
	if err != nil {
		t.Error(err)
	}

	if len(tokens) == 0 || len(tokens) > 9 {
		t.Errorf("want: %v\n", "len(tokens) not empty, and len(tokens) <=9")
		t.Errorf("got: %v\n", got)
	}

}
