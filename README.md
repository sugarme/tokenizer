# Tokenizer [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/sugarme/tokenizer?tab=doc)[![Travis CI](https://api.travis-ci.org/sugarme/tokenizer.svg?branch=master)](https://travis-ci.org/sugarme/tokenizer)[![Go Report Card](https://goreportcard.com/badge/github.com/sugarme/tokenizer)](https://goreportcard.com/report/github.com/sugarme/tokenizer) 

## Overview

`tokenizer` is pure Go package to facilitate applying Natural Language Processing (NLP) models train/test and inference in Go. 

It is heavily inspired by and based on the popular [HuggingFace Tokenizers](https://github.com/huggingface/tokenizers). 

`tokenizer` is part of an ambitious goal (together with [**transformer**](https://github.com/sugarme/transformer) and [**gotch**](https://github.com/sugarme/gotch)) to bring more AI/deep-learning tools to Gophers so that they can stick to the language they love and build faster software in production. 

## Features

`tokenizer` is built in modules located in sub-packages. 
1. Normalizer
2. Pretokenizer
3. Tokenizer
4. Post-processing

It implements various tokenizer models: 
- [x] Word level model
- [x] Wordpiece model
- [x] Byte Pair Encoding (BPE)

It can be used for both **training** new models from scratch or **fine-tuning** existing models. See [examples](./example) detail.

## Basic example

```go
import (
	"fmt"
	"log"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/model/wordpiece"
	"github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/pretokenizer"
	"github.com/sugarme/tokenizer/processor"
)

func getBert() *tokenizer.Tokenizer {

	vocabFile := "./bert-base-uncased-vocab.txt"
	model, err := wordpiece.NewWordPieceFromFile(vocabFile, "[UNK]")
	if err != nil {
		log.Fatal(err)
	}

	tk := tokenizer.NewTokenizer(model)
	fmt.Printf("Vocab size: %v\n", tk.GetVocabSize(false))

	bertNormalizer := normalizer.NewBertNormalizer(true, true, true, true)
	tk.WithNormalizer(bertNormalizer)

	bertPreTokenizer := pretokenizer.NewBertPreTokenizer()
	tk.WithPreTokenizer(bertPreTokenizer)

	bertPreTokenizer := pretokenizer.NewBertPreTokenizer()
	tk.WithPreTokenizer(bertPreTokenizer)

	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("[MASK]", true)})

	return tk
}

func main() {

	tk := getBert()
	sentence := `Yesterday I saw a [MASK] far away`
	en, err := tk.EncodeSingle(sentence)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tokens: %v\n", en.GetTokens())
	fmt.Printf("offsets: %v\n", en.GetOffsets())

	// Output
	// tokens: [yesterday i saw a [MASK] far away]
	// offsets: [{0 9} {10 11} {12 15} {16 17} {18 24} {25 28} {29 33}]
}
```

## Getting Started

- See [pkg.go.dev](https://pkg.go.dev/github.com/sugarme/tokenizer?tab=doc) for detail APIs 


## License

`tokenizer` is Apache 2.0 licensed.


## Acknowledgement

- This project has been inspired and used many concepts from [HuggingFace Tokenizers](https://github.com/huggingface/tokenizers).


