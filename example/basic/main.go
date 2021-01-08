package main

import (
	"flag"
)

var model string

func init() {
	flag.StringVar(&model, "model", "bpe", "select tokenizer model - 'bpe' or 'bert'")
}

func main() {
	flag.Parse()

	switch model {
	case "bert":
		runBERT()

	case "bpe":
		runBPE()
	case "word":
		runWordLevel()
	}
}
