package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sugarme/tokenizer/pretokenizer"
	"github.com/sugarme/tokenizer/tokenizer"
	bpe "github.com/sugarme/tokenizer/tokenizer/model/bpe"
)

func main() {

	startTime := time.Now()

	files := []string{
		// "example/tokenizer/bpe/train/input/oscar.eo-50k.txt",
		// "example/tokenizer/bpe/train/input/adieu.txt",
		// "example/tokenizer/bpe/train/input/test.txt",
		// "example/tokenizer/bpe/train/input/test-eo.txt",

		"example/tokenizer/bpe/train/input/oscar.eo.txt",
		"example/tokenizer/bpe/train/input/epo_literature_2011_300K-sentences.txt",
		"example/tokenizer/bpe/train/input/epo_mixed_2012_1M-sentences.txt",
		"example/tokenizer/bpe/train/input/epo_newscrawl_2017_1M-sentences.txt",
		"example/tokenizer/bpe/train/input/epo_web_2011_100K-sentences.txt",
		"example/tokenizer/bpe/train/input/epo_web_2012_1M-sentences.txt",
		"example/tokenizer/bpe/train/input/epo_wikipedia_2007_300K-sentences.txt",
		"example/tokenizer/bpe/train/input/epo_wikipedia_2011_300K-sentences.txt",
		"example/tokenizer/bpe/train/input/epo_wikipedia_2012_300K-sentences.txt",
		"example/tokenizer/bpe/train/input/epo_wikipedia_2016_300K-sentences.txt",
	}

	model, err := bpe.NewBPE()
	if err != nil {
		log.Fatal(err)
	}

	trainer := bpe.NewBpeTrainer(2, 52000)

	tk := tokenizer.NewTokenizer(model)

	tk.AddSpecialTokens([]string{
		"<s>",
		"<pad>",
		"</s>",
		"<unk>",
		"<mask>",
	})

	bytelevel := pretokenizer.NewByteLevel()

	tk.WithPreTokenizer(bytelevel)

	err = tk.Train(trainer, files)
	if err != nil {
		log.Fatal(err)
	}

	trainedModel := tk.GetModel()

	trainedModel.Save("example/tokenizer/bpe/train/model", "es")

	trainedTime := time.Since(startTime).Seconds() / 60

	fmt.Printf("Training time (min): %f.2\n", trainedTime)

}
