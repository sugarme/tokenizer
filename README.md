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

This tokenizer package is compatible to load pretrained models from Huggingface. Some of them can be loaded using `pretrained` subpackage.

```go
import (
	"fmt"
	"log"

	"github.com/sugarme/tokenizer/pretrained"
)

func main() {
    // Download and cache pretrained tokenizer. In this case `bert-base-uncased` from Huggingface
    // can be any model with `tokenizer.json` available. E.g. `tiiuae/falcon-7b`
	configFile, err := tokenizer.CachedPath("bert-base-uncased", "tokenizer.json")
	if err != nil {
		panic(err)
	}

	tk, err := pretrained.FromFile(configFile)
	if err != nil {
		panic(err)
	}

	sentence := `The Gophers craft code using [MASK] language.`
	en, err := tk.EncodeSingle(sentence)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tokens: %q\n", en.Tokens)
	fmt.Printf("offsets: %v\n", en.Offsets)

	// Output
	// tokens: ["the" "go" "##pher" "##s" "craft" "code" "using" "[MASK]" "language" "."]
	// offsets: [[0 3] [4 6] [6 10] [10 11] [12 17] [18 22] [23 28] [29 35] [36 44] [44 45]]
}
```

All models can be loaded from files manually. [pkg.go.dev](https://pkg.go.dev/github.com/sugarme/tokenizer?tab=doc) for detail APIs.


## Getting Started

- See [pkg.go.dev](https://pkg.go.dev/github.com/sugarme/tokenizer?tab=doc) for detail APIs 


## License

`tokenizer` is Apache 2.0 licensed.


## Acknowledgement

- This project has been inspired and used many concepts from [HuggingFace Tokenizers](https://github.com/huggingface/tokenizers).


