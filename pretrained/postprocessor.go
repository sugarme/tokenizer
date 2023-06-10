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

func createTemplateProcessing(params *util.Params) tokenizer.PostProcessor {
	// TODO
	panic("NotImplementedError")
}

func createSequence(params *util.Params) tokenizer.PostProcessor {
	// TODO
	panic("NotImplementedError")
}
