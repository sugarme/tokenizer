package pretrained

import (
	"fmt"
	"log"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/model"
	"github.com/sugarme/tokenizer/model/bpe"
	"github.com/sugarme/tokenizer/model/wordlevel"
	"github.com/sugarme/tokenizer/model/wordpiece"
	"github.com/sugarme/tokenizer/util"
)

// This file provides functions to create tokenizer.Model from input data.

func CreateModel(config map[string]interface{}) (tokenizer.Model, error) {
	if config == nil {
		return nil, nil
	}

	params := util.NewParams(config)

	var typ string
	if params.Has("type") {
		typ = params.Get("type").(string)
	} else {
		log.Printf("INFO: there is no field 'type' in model json data, a default 'WordPiece' model will be trying to create...\n")
		typ = "WordPiece" // Default to `WordPiece` model as in BERT "tokenizer.json", there's not field "type"
	}

	switch typ {
	case "BPE":
		return createBPE(params)
	case "WordPiece":
		return createWordPiece(params)
	case "WordLevel":
		return createWordLevel(params)
	case "Unigram":
		return createUnigram(params)

	default:
		err := fmt.Errorf("Could not construct tokenizer.Model from input data: %#v\n", config)
		return nil, err
	}
}

// BPE json format:
// ----------------
// "type": "BPE",
// "dropout": null,
// "unk_token": null,
// "continuing_subword_prefix": null,
// "end_of_word_suffix": null,
// "fuse_unk": false,
// "byte_fallback": false,
// "vocab": {}
// "merges": []

func createBPE(params *util.Params) (tokenizer.Model, error) {
	var dropout *float32
	if params.Has("dropout") {
		val := float32(params.Get("dropout").(float64))
		dropout = &val
	}

	var unkToken *string
	if params.Has("unk_token") {
		v := params.Get("unk_token").(string)
		unkToken = &v
	}
	var continuingSubwordPrefix *string
	if params.Has("continuing_subword_prefix") {
		v := params.Get("continuing_subword_prefix").(string)
		continuingSubwordPrefix = &v
	}

	var endOfWordSuffix *string
	if params.Has("end_of_word_suffix") {
		v := params.Get("end_of_word_suffix").(string)
		endOfWordSuffix = &v
	}
	// fuseUnk := params.Get("use_unk").(bool)
	// byteFallback := params.Get("byte_fallback").(bool)

	vocab := castVocab(params.Get("vocab").(map[string]interface{}))
	merges := castMerge(params.Get("merges").([]interface{}))

	return bpe.New(vocab, merges, dropout, unkToken, continuingSubwordPrefix, endOfWordSuffix)
}

// WordPiece json format:
// ----------------------
// "unk_token": "[UNK]"
// "continuing_subword_prefix":"##"
// "max_input_chars_per_word":100
// "vocab": {}
// "decoder":{"type":"WordPiece","prefix":"##","cleanup":true},

func createWordPiece(params *util.Params) (tokenizer.Model, error) {
	var unkToken *string
	if params.Has("unk_token") {
		v := params.Get("unk_token").(string)
		unkToken = &v
	}
	var continuingSubwordPrefix *string
	if params.Has("continuing_subword_prefix") {
		v := params.Get("continuing_subword_prefix").(string)
		continuingSubwordPrefix = &v
	}

	var maxInputCharsPerWord *int
	if params.Has("max_input_chars_per_word") {
		v := int(params.Get("max_input_chars_per_word").(float64))
		maxInputCharsPerWord = &v
	}

	vocab := castVocab(params.Get("vocab").(map[string]interface{}))

	return wordpiece.New(vocab, unkToken, continuingSubwordPrefix, maxInputCharsPerWord)
}

func createWordLevel(params *util.Params) (tokenizer.Model, error) {
	var unkToken *string
	if params.Has("unk_token") {
		v := params.Get("unk_token").(string)
		unkToken = &v
	}

	vocab := castVocab(params.Get("vocab").(map[string]interface{}))

	return wordlevel.New(vocab, unkToken)
}

func createUnigram(params *util.Params) (tokenizer.Model, error) {
	panic("NotImplementedError")
}

func castVocab(input map[string]interface{}) model.Vocab {
	out := make(map[string]int)
	for k, v := range input {
		out[k] = int(v.(float64))
	}

	return out
}

func castMerge(input []interface{}) []string {
	out := make([]string, len(input))
	for i, v := range input {
		out[i] = v.(string)
	}

	return out
}
