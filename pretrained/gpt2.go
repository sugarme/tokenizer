package pretrained

import (
	"embed"
	"log"
	"os"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/decoder"
	"github.com/sugarme/tokenizer/model/bpe"
	"github.com/sugarme/tokenizer/pretokenizer"
	"github.com/sugarme/tokenizer/processor"
	"github.com/sugarme/tokenizer/util"
)

// GPT2 loads GPT2 (small) tokenizer from vocab and merges files.
//
// Params:
// - addPrefixSpace: set whether to add a leading space to the first word.
//   It allows to treat the leading word just as any other words.
// - trimOffsets: set Whether the post processing step should trim offsets
//   to avoid including whitespaces.
//
// Special tokens:
// - cls-token: "<s>"
// - sep token: "</s>"
// - pad token: "<pad>"
// - space token: "Ġ"
//
// Source:
// "https://cdn.huggingface.co/gpt2-merges.txt"
// "https://cdn.huggingface.co/gpt2-vocab.json"

//go:embed model/gpt2-merges.txt
//go:embed model/gpt2-vocab.json
var fs embed.FS

const (
	vocabFilename = "model/gpt2-vocab.json"
	mergeFilename = "model/gpt2-merges.txt"
)

func GPT2(addPrefixSpace bool, trimOffsets bool) *tokenizer.Tokenizer {

	currDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	util.CdToThis()
	defer util.CdBack(currDir)

	model, err := bpe.NewBPEFromFS(fs, vocabFilename, mergeFilename)
	if err != nil {
		log.Fatal(err)
	}

	tk := tokenizer.NewTokenizer(model)

	pretok := pretokenizer.NewByteLevel()
	pretok.SetAddPrefixSpace(addPrefixSpace)
	pretok.SetTrimOffsets(trimOffsets)
	tk.WithPreTokenizer(pretok)

	pprocessor := processor.NewByteLevelProcessing(pretok)
	tk.WithPostProcessor(pprocessor)

	bpeDecoder := decoder.NewBpeDecoder("Ġ")
	tk.WithDecoder(bpeDecoder)

	return tk
}
