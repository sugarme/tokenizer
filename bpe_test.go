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
	// "github.com/sugarme/tokenizer/processor"
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
	tk.WithPreTokenizer(pretok)

	// TODO: adde bytelevel (post)processor

	return tk
}

func checkOffsets(t *testing.T, input string, output *tokenizer.Encoding, offset int, want string) {
	offsets := output.GetOffsets()[offset]

	got := input[offsets[0]:offsets[1]]

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %v\n", want)
		t.Errorf("Got: %v\n", got)
	}
}

func TestByteLevelBasic(t *testing.T) {
	tk := getByteLevel(true, false)

	input := "Hello there, how are you?"

	inputSeq := tokenizer.NewInputSequence(input)
	output, err := tk.Encode(tokenizer.NewSingleEncodeInput(inputSeq), false)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("output: %+v\n", output)
	checkOffsets(t, input, output, 0, "Hello")
}
