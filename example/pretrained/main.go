package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/season-studio/tokenizer"
	"github.com/season-studio/tokenizer/pretrained"
)

var (
	modelName string
)

func init() {
	flag.StringVar(&modelName, "model", "bert-base-uncased", "model name as at Huggingface model hub e.g. 'tiiuae/falcon-7b'. Default='bert-base-uncased'")
}

func main() {
	flag.Parse()

	// any model with file `tokenizer.json` available. Eg. `tiiuae/falcon-7b`, `TheBloke/guanaco-7B-HF`, `mosaicml/mpt-7b-instruct`
	// configFile, err := tokenizer.CachedPath("bert-base-uncased", "tokenizer.json")
	configFile, err := tokenizer.CachedPath(modelName, "tokenizer.json")
	if err != nil {
		panic(err)
	}

	tk, err := pretrained.FromFile(configFile)
	if err != nil {
		panic(err)
	}

	sentence := `The Gophers craft code using [MASK] language.`
	en, err := tk.EncodeSingle(sentence, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-10s: %q\n", "Tokens", en.Tokens)
	fmt.Printf("%-10s: %v\n", "Ids", en.Ids)
	fmt.Printf("%-10s: %v\n", "Offsets", en.Offsets)
}
