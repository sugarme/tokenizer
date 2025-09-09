package main

import (
	"fmt"
	"log"

	"github.com/gengzongjie/tokenizer"
	"github.com/gengzongjie/tokenizer/model/bpe"
	"github.com/gengzongjie/tokenizer/pretokenizer"
	"github.com/gengzongjie/tokenizer/processor"
	"github.com/gengzongjie/tokenizer/util"
)

func runBPE() {
	tk := getByteLevel(true, false)

	input := "Hello there, how are you?"

	inputSeq := tokenizer.NewInputSequence(input)
	output, err := tk.Encode(tokenizer.NewSingleEncodeInput(inputSeq), false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("encoding: %v\n", output)
}

func getByteLevelBPE() (retVal *tokenizer.Tokenizer) {

	util.CdToThis()
	vocabFile := "../../data/gpt2-vocab.json"
	mergeFile := "../../data/gpt2-merges.txt"

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
