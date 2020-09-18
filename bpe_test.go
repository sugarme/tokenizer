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

	// TODO: adde bytelevel (post)processor
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

	fmt.Printf("output========: %+v\n", output1)
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

	fmt.Printf("output========: %+v\n", output2)

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

	fmt.Printf("output: %+v\n", output)

	checkOffsets(t, input, output, 1, "⭢")
	checkOffsets(t, input, output, 2, "⭢")
	checkOffsets(t, input, output, 3, "⭢")
}
