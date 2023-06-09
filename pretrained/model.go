package pretrained

import (
	"fmt"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/util"
)

// This file provides functions to create tokenizer.Model from input data.

func CreateModel(config map[string]interface{}) (tokenizer.Model, error) {
	if config == nil {
		return nil, nil
	}

	params := util.NewParams(config)

	typ := params.Get("type").(string)

	switch typ {
	case "BPE":
		return createBPE(params)
	case "WordPiece":
		return createWordPiece(params)
	case "WordLevel":
		return createWordLevel(params)
	case "Unigram":
		return createUnigram(params)
	case "SentencePiece":
		return createSentencePiece(params)

	default:
		err := fmt.Errorf("Could not construct tokenizer.Model from input data: %#v\n", config)
		return nil, err
	}
}

func createBPE(params *util.Params) (tokenizer.Model, error) {
	panic("NotImplementedError")
}

func createWordPiece(params *util.Params) (tokenizer.Model, error) {
	panic("NotImplementedError")
}

func createWordLevel(params *util.Params) (tokenizer.Model, error) {
	panic("NotImplementedError")
}

func createUnigram(params *util.Params) (tokenizer.Model, error) {
	panic("NotImplementedError")
}

func createSentencePiece(params *util.Params) (tokenizer.Model, error) {
	panic("NotImplementedError")
}
