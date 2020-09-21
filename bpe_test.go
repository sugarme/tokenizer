package tokenizer_test

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/model/bpe"

	// "github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/pretokenizer"
	"github.com/sugarme/tokenizer/processor"
	"github.com/sugarme/tokenizer/util"
)

func getByteLevelBPE() (retVal *tokenizer.Tokenizer) {

	util.CdToThis()
	vocabFile := "data/gpt2-vocab.json"
	mergeFile := "data/gpt2-merges.txt"

	model, err := bpe.NewBpeFromFiles(vocabFile, mergeFile)
	if err != nil {
		log.Fatal(err)
	}

	tk := tokenizer.NewTokenizer(model)
	fmt.Printf("Vocab size: %v\n", tk.GetVocabSize(false))

	return tk
}

func getByteLevel(addPrefixSpace bool, trimOffsets bool) *tokenizer.Tokenizer {

	tk := getByteLevelBPE()

	pretok := pretokenizer.NewByteLevel()
	pretok.SetAddPrefixSpace(addPrefixSpace)
	pretok.SetTrimOffsets(trimOffsets)
	tk.WithPreTokenizer(pretok)

	pprocessor := processor.NewByteLevelProcessing(pretok)
	tk.WithPostProcessor(pprocessor)

	return tk
}

func checkOffsets(t *testing.T, input string, output *tokenizer.Encoding, offset int, want string) {
	offsets := output.GetOffsets()[offset]

	got := input[offsets[0]:offsets[1]]

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: '%v'\n", want)
		t.Errorf("Got: '%v'\n", got)
	}
}

func TestByteLevelBasic(t *testing.T) {
	tk1 := getByteLevel(true, false)

	input := "Hello there, how are you?"

	inputSeq := tokenizer.NewInputSequence(input)
	output1, err := tk1.Encode(tokenizer.NewSingleEncodeInput(inputSeq), false)
	if err != nil {
		t.Error(err)
	}

	// fmt.Printf("output========: %+v\n", output1)
	checkOffsets(t, input, output1, 0, "Hello")
	checkOffsets(t, input, output1, 1, " there")
	checkOffsets(t, input, output1, 2, ",")
	checkOffsets(t, input, output1, 3, " how")
	checkOffsets(t, input, output1, 4, " are")
	checkOffsets(t, input, output1, 5, " you")
	checkOffsets(t, input, output1, 6, "?")

	// And when trimming offsets:
	tk2 := getByteLevel(true, true)
	inputSeq = tokenizer.NewInputSequence(input)
	output2, err := tk2.Encode(tokenizer.NewSingleEncodeInput(inputSeq), false)
	if err != nil {
		t.Error(err)
	}

	// fmt.Printf("output========: %+v\n", output2)

	checkOffsets(t, input, output2, 0, "Hello")
	checkOffsets(t, input, output2, 1, "there")
	checkOffsets(t, input, output2, 2, ",")
	checkOffsets(t, input, output2, 3, "how")
	checkOffsets(t, input, output2, 4, "are")
	checkOffsets(t, input, output2, 5, "you")
	checkOffsets(t, input, output2, 6, "?")

}

func TestByteLevelUnicode(t *testing.T) {
	tk := getByteLevel(true, false)

	input := "i⭢j"

	inputSeq := tokenizer.NewInputSequence(input)
	output, err := tk.Encode(tokenizer.NewSingleEncodeInput(inputSeq), false)
	if err != nil {
		t.Error(err)
	}

	// fmt.Printf("output: %+v\n", output)

	checkOffsets(t, input, output, 1, "⭢")
	checkOffsets(t, input, output, 2, "⭢")
	checkOffsets(t, input, output, 3, "⭢")
}

func TestByteLevelDoubleSequence(t *testing.T) {

	input1 := "My name is Anthony"
	inputSeq1 := tokenizer.NewInputSequence(input1)
	input2 := "What is my name?"
	inputSeq2 := tokenizer.NewInputSequence(input2)

	// Without trimming offsets
	tk := getByteLevel(true, false)
	output, err := tk.Encode(tokenizer.NewDualEncodeInput(inputSeq1, inputSeq2), false)
	if err != nil {
		t.Error(err)
	}

	// fmt.Printf("output: %+v\n", output)

	got1 := output.GetOffsets()
	want1 := [][]int{
		{0, 2},
		{2, 7},
		{7, 10},
		{10, 18},
		{0, 4},
		{4, 7},
		{7, 10},
		{10, 15},
		{15, 16},
	}

	if !reflect.DeepEqual(got1, want1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}

	// word Ids
	got2 := output.GetWords()
	want2 := []int{0, 1, 2, 3, 0, 1, 2, 3, 4}
	if !reflect.DeepEqual(got2, want2) {
		t.Errorf("Want: %v\n", want2)
		t.Errorf("Got: %v\n", got2)
	}

	// type Ids
	got3 := output.GetTypeIds()
	want3 := []int{0, 0, 0, 0, 1, 1, 1, 1, 1}
	if !reflect.DeepEqual(got3, want3) {
		t.Errorf("Want: %v\n", want3)
		t.Errorf("Got: %v\n", got3)
	}

	// When trimming offsets
	tk2 := getByteLevel(true, true)
	output2, err := tk2.Encode(tokenizer.NewDualEncodeInput(inputSeq1, inputSeq2), false)
	if err != nil {
		t.Error(err)
	}
	got4 := output2.GetOffsets()
	want4 := [][]int{
		{0, 2},
		{3, 7},
		{8, 10},
		{11, 18},
		{0, 4},
		{5, 7},
		{8, 10},
		{11, 15},
		{15, 16},
	}
	if !reflect.DeepEqual(got4, want4) {
		t.Errorf("Want: %v\n", want4)
		t.Errorf("Got: %v\n", got4)
	}

}

func TestByteLevelPreTokenizedSequence(t *testing.T) {
	input := []string{"My", "name", "is", "Anthonino"}
	inputSeq := tokenizer.NewInputSequence(input)

	// Without trimming offsets
	tk := getByteLevel(true, false)
	output, err := tk.Encode(tokenizer.NewSingleEncodeInput(inputSeq), false)
	if err != nil {
		t.Error(err)
	}

	// fmt.Printf("output: %+v\n", output)

	got1 := output.GetOffsets()
	want1 := [][]int{{0, 2}, {0, 4}, {0, 2}, {0, 4}, {4, 6}, {6, 9}}
	if !reflect.DeepEqual(got1, want1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}

	got2 := output.GetWords()
	want2 := []int{0, 1, 2, 3, 3, 3}
	if !reflect.DeepEqual(got2, want2) {
		t.Errorf("Want: %v\n", want2)
		t.Errorf("Got: %v\n", got2)
	}

	got3 := output.GetTokens()
	want3 := []string{"ĠMy", "Ġname", "Ġis", "ĠAnth", "on", "ino"}
	if !reflect.DeepEqual(got3, want3) {
		t.Errorf("Want: %v\n", want3)
		t.Errorf("Got: %v\n", got3)
	}
}

func TestByteLevelPreTokenizedSequenceWithTrimming(t *testing.T) {
	input := []string{"My", "name", "is", "Anthonino"}
	inputSeq := tokenizer.NewInputSequence(input)

	// When trimming offsets
	tk := getByteLevel(true, true)
	output, err := tk.Encode(tokenizer.NewSingleEncodeInput(inputSeq), false)
	if err != nil {
		t.Error(err)
	}

	// fmt.Printf("output: %+v\n", output)

	got1 := output.GetOffsets()
	want1 := [][]int{{0, 2}, {1, 4}, {1, 2}, {1, 4}, {4, 6}, {6, 9}}
	if !reflect.DeepEqual(got1, want1) {
		t.Errorf("Want: %v\n", want1)
		t.Errorf("Got: %v\n", got1)
	}

	got2 := output.GetWords()
	want2 := []int{0, 1, 2, 3, 3, 3}
	if !reflect.DeepEqual(got2, want2) {
		t.Errorf("Want: %v\n", want2)
		t.Errorf("Got: %v\n", got2)
	}

	got3 := output.GetTokens()
	want3 := []string{"ĠMy", "Ġname", "Ġis", "ĠAnth", "on", "ino"}
	if !reflect.DeepEqual(got3, want3) {
		t.Errorf("Want: %v\n", want3)
		t.Errorf("Got: %v\n", got3)
	}
}
