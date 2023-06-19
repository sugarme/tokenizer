package pretrained

// This file provides functions to create tokenizer.PostProcessor
// 1. RobertaProcessing
// 2. BertProcessing
// 3. ByteLevel
// 4. TemplateProcessing
// 5. Sequence

import (
	"fmt"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretokenizer"
	"github.com/sugarme/tokenizer/processor"
	"github.com/sugarme/tokenizer/util"
)

func CreatePostProcessor(config map[string]interface{}) (tokenizer.PostProcessor, error) {
	if config == nil {
		return nil, nil
	}

	params := util.NewParams(config)

	typ := params.Get("type").(string)

	switch typ {
	case "RobertaProcessing":
		return createRobertaProcessing(params), nil
	case "BertProcessing":
		return createBertProcessing(params), nil
	case "ByteLevel":
		return createByteLevel(params)
	case "TemplateProcessing":
		return createTemplateProcessing(params), nil
	case "Sequence":
		return createSequence(params), nil

	default:
		err := fmt.Errorf("Could not create tokenizer.PostProcessor from input data %#v\n", config)
		return nil, err
	}
}

// RobertaProcessing json data e.g.:
// "sep":["</s>",2],
// "cls":["<s>",0],
// "trim_offsets":true,
// "add_prefix_space":false
func createRobertaProcessing(params *util.Params) tokenizer.PostProcessor {
	sep := getPostToken(params, "sep")
	cls := getPostToken(params, "cls")
	trimOffsets := params.Get("trim_offsets").(bool)
	addPrefixSpace := params.Get("add_prefix_space").(bool)

	return processor.NewRobertaProcessing(sep, cls, trimOffsets, addPrefixSpace)
}

func getPostToken(params *util.Params, name string) processor.PostToken {
	sepData := params.Get(name).([]interface{})[0].(map[string]float64)
	var tok processor.PostToken
	for k, v := range sepData {
		tok = processor.PostToken{
			Value: k,
			Id:    int(v),
		}
		break
	}

	return tok
}

func createBertProcessing(params *util.Params) tokenizer.PostProcessor {
	sep := getPostToken(params, "sep")
	cls := getPostToken(params, "cls")

	return processor.NewBertProcessing(sep, cls)
}

func createByteLevel(params *util.Params) (tokenizer.PostProcessor, error) {
	pretok, err := createByteLevelPreTokenizer(params)
	if err != nil {
		return nil, err
	}
	return processor.NewByteLevelProcessing(pretok.(*pretokenizer.ByteLevel)), nil
}

// e.g. `TheBloke/guanaco-7B-HF`
// https://huggingface.co/TheBloke/guanaco-7B-HF/raw/main/tokenizer.json
/*
 "post_processor": {
    "type": "TemplateProcessing",
    "single": [
      {
        "SpecialToken": {
          "id": "<s>",
          "type_id": 0
        }
      },
      {
        "Sequence": {
          "id": "A",
          "type_id": 0
        }
      }
    ],
    "pair": [
      {
        "SpecialToken": {
          "id": "<s>",
          "type_id": 0
        }
      },
      {
        "Sequence": {
          "id": "A",
          "type_id": 0
        }
      },
      {
        "Sequence": {
          "id": "B",
          "type_id": 0
        }
      }
    ],
    "special_tokens": {
      "<s>": {
        "id": "<s>",
        "ids": [
          1
        ],
        "tokens": [
          "<s>"
        ]
      }
    }
  },
*/
func createTemplateProcessing(params *util.Params) tokenizer.PostProcessor {
	var (
		single        processor.Template
		pair          processor.Template
		specialTokens *processor.Tokens
	)

	// Single
	if params.Has("single") {
		singleData := params.Get("single").([]interface{})
		for _, p := range singleData {
			ps := util.NewParams(p.(map[string]interface{}))
			if ps.Has("Sequence") {
				item := ps.Get("Sequence").(map[string]interface{})
				id := item["id"].(string) // always "A"
				typeId := int(item["type_id"].(float64))
				piece := processor.NewSequencePiece(id, typeId)
				single = append(single, piece)
			}

			if ps.Has("SpecialToken") {
				item := ps.Get("SpecialToken").(map[string]interface{})
				id := item["id"].(string)
				typeId := int(item["type_id"].(float64))
				piece := processor.NewSpecialTokenPiece(id, typeId)
				single = append(single, piece)
			}
		}
	}

	// Pair
	if params.Has("pair") {
		pairData := params.Get("pair").([]interface{})
		for _, p := range pairData {
			ps := util.NewParams(p.(map[string]interface{}))
			if ps.Has("Sequence") {
				item := ps.Get("Sequence").(map[string]interface{})
				id := item["id"].(string)
				typeId := int(item["type_id"].(float64))
				piece := processor.NewSequencePiece(id, typeId)
				pair = append(pair, piece)
			}

			if ps.Has("SpecialToken") {
				item := ps.Get("SpecialToken").(map[string]interface{})
				id := item["id"].(string)
				typeId := int(item["type_id"].(float64))
				piece := processor.NewSpecialTokenPiece(id, typeId)
				pair = append(pair, piece)
			}
		}
	}

	// SpecialTokens
	if params.Has("special_tokens") {
		data := params.Get("special_tokens").(map[string]interface{})
		var toks []processor.SpecialToken
		for _, v := range data {
			d := v.(map[string]interface{})
			id := d["id"].(string)
			vals := d["ids"].([]float64)
			var ids []int
			for _, val := range vals {
				ids = append(ids, int(val))
			}
			tokens := d["tokens"].([]string)
			tok := processor.NewSpecialToken(id, ids, tokens)
			toks = append(toks, *tok)
		}

		specialTokens = processor.NewTokensFrom(toks)
	}

	return processor.NewTemplateProcessing(single, pair, specialTokens)
}

func createSequence(params *util.Params) tokenizer.PostProcessor {
	// TODO
	panic("NotImplementedError")
}
