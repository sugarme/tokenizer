package pretrained

import (
	"log"
	"os"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/decoder"
	"github.com/sugarme/tokenizer/model/bpe"
	"github.com/sugarme/tokenizer/pretokenizer"
	"github.com/sugarme/tokenizer/processor"
	"github.com/sugarme/tokenizer/util"
)

// RobertaBase loads pretrained RoBERTa tokenizer.
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
// - vocab: "https://cdn.huggingface.co/roberta-base-vocab.json",
// - merges: "https://cdn.huggingface.co/roberta-base-merges.txt",
func RobertaBase(addPrefixSpace, trimOffsets bool) *tokenizer.Tokenizer {
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	util.CdToThis()
	defer util.CdBack(currDir)

	vocabFile := "model/roberta-base-vocab.json"
	mergesFile := "model/roberta-base-merges.txt"
	model, err := bpe.NewBpeFromFiles(vocabFile, mergesFile)
	if err != nil {
		log.Fatal(err)
	}

	tk := tokenizer.NewTokenizer(model)

	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("<s>", true)})
	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("</s>", true)})

	blPreTokenizer := pretokenizer.NewByteLevel()
	blPreTokenizer.SetAddPrefixSpace(addPrefixSpace)
	blPreTokenizer.SetTrimOffsets(trimOffsets)
	tk.WithPreTokenizer(blPreTokenizer)

	postProcess := processor.DefaultRobertaProcessing()
	// postProcess.TrimOffsets(trimOffsets)
	tk.WithPostProcessor(postProcess)

	bpeDecoder := decoder.NewBpeDecoder("Ġ")
	tk.WithDecoder(bpeDecoder)

	return tk
}

// RobertaBaseSquad2 loads pretrained RoBERTa fine-tuned SQuAD Question Answering tokenizer.
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
// - vocab: "https://cdn.huggingface.co/deepset/roberta-base-squad2/vocab.json",
// - merges: "https://cdn.huggingface.co/deepset/roberta-base-squad2/merges.txt",
func RobertaBaseSquad2(addPrefixSpace, trimOffsets bool) *tokenizer.Tokenizer {
	util.CdToThis()
	vocabFile := "model/roberta-base-squad2-vocab.json"
	mergesFile := "model/roberta-base-squad2-merges.txt"
	model, err := bpe.NewBpeFromFiles(vocabFile, mergesFile)
	if err != nil {
		log.Fatal(err)
	}

	tk := tokenizer.NewTokenizer(model)

	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("<s>", true)})
	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("</s>", true)})

	blPreTokenizer := pretokenizer.NewByteLevel()
	blPreTokenizer.SetAddPrefixSpace(addPrefixSpace)
	blPreTokenizer.SetTrimOffsets(trimOffsets)
	tk.WithPreTokenizer(blPreTokenizer)

	postProcess := processor.DefaultRobertaProcessing()
	tk.WithPostProcessor(postProcess)

	bpeDecoder := decoder.NewBpeDecoder("Ġ")
	tk.WithDecoder(bpeDecoder)

	return tk
}
