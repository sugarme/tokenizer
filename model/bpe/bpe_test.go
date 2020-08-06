package bpe_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	// "reflect"
	// "strings"
	"testing"

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
	var vocab map[string]uint32 = make(map[string]uint32)
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
