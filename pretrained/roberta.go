package pretrained

import (
	"log"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/decoder"
	"github.com/sugarme/tokenizer/model/bpe"
	"github.com/sugarme/tokenizer/pretokenizer"
	"github.com/sugarme/tokenizer/processor"
	"github.com/sugarme/tokenizer/util"
)

// RobertaBase loads pretrained RoBERTa tokenizer.
//
// It uses special tokens:
// - classifier: "<s>"
// - seperator: "</s>"
// - padding: "<pad>"
// - suffix: "Ġ"
// Source:
// - vocab: "https://cdn.huggingface.co/deepset/roberta-base-squad2/vocab.json",
// - merges: "https://cdn.huggingface.co/roberta-base-merges.txt",
func RobertaBase() *tokenizer.Tokenizer {
	util.CdToThis()
	vocabFile := "../../data/roberta-base-vocab.json"
	mergesFile := "../../data/roberta-base-merges.txt"
	model, err := bpe.NewBpeFromFiles(vocabFile, mergesFile)
	if err != nil {
		log.Fatal(err)
	}

	tk := tokenizer.NewTokenizer(model)

	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("<s>", true)})
	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("</s>", true)})

	blPreTokenizer := pretokenizer.NewByteLevel()
	tk.WithPreTokenizer(blPreTokenizer)

	postProcess := processor.DefaultRobertaProcessing()
	tk.WithPostProcessor(postProcess)

	bpeDecoder := decoder.NewBpeDecoder("Ġ")
	tk.WithDecoder(bpeDecoder)

	return tk
}

// RobertaBaseSquad2 loads pretrained RoBERTa fine-tuned SQuAD Question Answering tokenizer.
//
// It uses special tokens:
// - classifier: "<s>"
// - seperator: "</s>"
// - padding: "<pad>"
// - suffix: "Ġ"
// Source:
// - vocab: "https://cdn.huggingface.co/deepset/roberta-base-squad2/vocab.json",
// - merges: "https://cdn.huggingface.co/deepset/roberta-base-squad2/merges.txt",
func RobertaBaseSquad2() *tokenizer.Tokenizer {
	util.CdToThis()
	vocabFile := "../../data/roberta-qa-vocab.json"
	mergesFile := "../../data/roberta-qa-merges.txt"
	model, err := bpe.NewBpeFromFiles(vocabFile, mergesFile)
	if err != nil {
		log.Fatal(err)
	}

	tk := tokenizer.NewTokenizer(model)

	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("<s>", true)})
	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("</s>", true)})

	blPreTokenizer := pretokenizer.NewByteLevel()
	tk.WithPreTokenizer(blPreTokenizer)

	postProcess := processor.DefaultRobertaProcessing()
	tk.WithPostProcessor(postProcess)

	bpeDecoder := decoder.NewBpeDecoder("Ġ")
	tk.WithDecoder(bpeDecoder)

	return tk
}
