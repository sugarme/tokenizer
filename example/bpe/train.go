package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/model/bpe"
	"github.com/sugarme/tokenizer/pretokenizer"
)

func runTrain() {

	startTime := time.Now()

	files := []string{
		// "example/tokenizer/bpe/train/input/oscar.eo-50k.txt",
		// "example/tokenizer/bpe/train/input/adieu.txt",
		// "example/tokenizer/bpe/train/input/test.txt",
		// "example/tokenizer/bpe/train/input/test-eo.txt",

		"input/oscar.eo.txt",
		"input/epo_literature_2011_300K-sentences.txt",
		"input/epo_mixed_2012_1M-sentences.txt",
		"input/epo_newscrawl_2017_1M-sentences.txt",
		"input/epo_web_2011_100K-sentences.txt",
		"input/epo_web_2012_1M-sentences.txt",
		"input/epo_wikipedia_2007_300K-sentences.txt",
		"input/epo_wikipedia_2011_300K-sentences.txt",
		"input/epo_wikipedia_2012_300K-sentences.txt",
		"input/epo_wikipedia_2016_300K-sentences.txt",
	}

	var vocab map[string]int = make(map[string]int)
	vocab["<s>"] = 0
	vocab["<pad>"] = 1
	vocab["</s>"] = 2
	vocab["<unk>"] = 3
	vocab["<mask>"] = 4

	var merges bpe.Merges = make(map[bpe.Pair]bpe.PairVal)

	model := bpe.NewBPE(vocab, merges)

	// model, err := bpe.DefaultBPE()
	// if err != nil {
	// log.Fatal(err)
	// }

	unkToken := "<unk>"
	model.UnkToken = &unkToken

	trainer := bpe.NewBpeTrainer(2, 52000)

	tk := tokenizer.NewTokenizer(model)

	// specialToks := []tokenizer.AddedToken{
	// tokenizer.NewAddedToken("<s>", true),
	// tokenizer.NewAddedToken("<pad>", true),
	// tokenizer.NewAddedToken("</s>", true),
	// tokenizer.NewAddedToken("<unk>", true),
	// tokenizer.NewAddedToken("<mask>", true),
	// }
	//
	// tk.AddSpecialTokens(specialToks)

	bytelevel := pretokenizer.NewByteLevel()

	tk.WithPreTokenizer(bytelevel)

	err := tk.Train(trainer, files)
	if err != nil {
		log.Fatal(err)
	}

	trainedModel := tk.GetModel()

	trainedModel.Save("example/tokenizer/bpe/train/model", "es")

	trainedTime := time.Since(startTime).Seconds() / 60

	fmt.Printf("Training time (min): %f.2\n", trainedTime)

}
