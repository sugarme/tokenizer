package pretrained

// This file provides functions to create tokenizer.PreTokenizer
// 1. BertPreTokenizer
// 2. ByteLevel
// 3. Delimiter
// 4. Metaspace
// 5. Whitespace
// 6. Sequence
// 7. Split
// 8. Punctuation
// 9. WhitespaceSplit
// 10. Digits
// 11. UnicodeScripts

import (
	"fmt"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/pretokenizer"
	"github.com/sugarme/tokenizer/util"
)

func CreatePreTokenizer(config map[string]interface{}) (tokenizer.PreTokenizer, error) {
	if config == nil {
		return nil, nil
	}

	params := util.NewParams(config)
	typ := params.Get("type").(string)

	switch typ {
	case "BertPreTokenizer":
		return pretokenizer.NewBertPreTokenizer(), nil
	case "ByteLevel":
		return createByteLevelPreTokenizer(params)
	case "Delimiter":
		return createDelimiterPreTokenizer(params)
	case "Metaspace":
		return createMetaspacePreTokenizer(params)
	case "Whitespace":
		return createWhitespacePreTokenizer(params)
	case "Sequence":
		return createSequencePreTokenizer(params)
	case "WhitespaceSplit":
		return createWhitespaceSplitPreTokenizer(params)
	case "Punctuation":
		return createPunctuationPreTokenizer(params)
	case "Digits":
		return createDigitsPreTokenizer(params)
	case "UnicodeScripts":
		return createUnicodeScriptsPreTokenizer(params)
	case "Split":
		return createSplitPreTokenizer(params)

	default:
		err := fmt.Errorf("Could not create PreTokenizer from input data: %#v\n", config)
		return nil, err
	}
}

func createByteLevelPreTokenizer(params *util.Params) (tokenizer.PreTokenizer, error) {
	if params == nil {
		return nil, nil
	}

	addPrefixSpace := params.Get("add_prefix_space", false).(bool)
	trimOffsets := params.Get("trim_offsets", false).(bool)

	return &pretokenizer.ByteLevel{
		AddPrefixSpace: addPrefixSpace,
		TrimOffsets:    trimOffsets,
	}, nil
}

func createDelimiterPreTokenizer(params *util.Params) (tokenizer.PreTokenizer, error) {
	// TODO. verify key `delimiter`
	delimiter := []rune(params.Get("delimiter").(string))[0]
	return pretokenizer.NewCharDelimiterSplit(delimiter), nil
}

func createMetaspacePreTokenizer(params *util.Params) (tokenizer.PreTokenizer, error) {
	if params == nil {
		return nil, nil
	}

	replacement := params.Get("replacement", "").(string)
	addPrefixSpace := params.Get("add_prefix_space", false).(bool)

	return pretokenizer.NewMetaspace(replacement, addPrefixSpace), nil
}

func createWhitespacePreTokenizer(params *util.Params) (tokenizer.PreTokenizer, error) {
	return pretokenizer.NewWhitespace(), nil
}

func createWhitespaceSplitPreTokenizer(params *util.Params) (tokenizer.PreTokenizer, error) {
	return pretokenizer.NewWhitespaceSplit(), nil
}

/*
	{
	       "type": "Punctuation",
	       "behavior": "Contiguous"
	     },
*/
func createPunctuationPreTokenizer(params *util.Params) (tokenizer.PreTokenizer, error) {
	behaviorVal := params.Get("behavior").(string)

	var b normalizer.SplitDelimiterBehavior
	switch behaviorVal {
	case "Removed":
		b = normalizer.RemovedBehavior
	case "Isolated":
		b = normalizer.IsolatedBehavior
	case "MergedWithNext":
		b = normalizer.MergedWithNextBehavior
	case "MergedWithPrevious":
		b = normalizer.MergedWithPreviousBehavior
	case "Contiguous":
		b = normalizer.ContiguousBehavior

	default:
		err := fmt.Errorf("Unsupported behavior: %#v\n", behaviorVal)
		return nil, err
	}

	return pretokenizer.NewPunctuation(b), nil
}

/*
	{
	        "type": "Digits",
	        "individual_digits": false
	      },
*/
func createDigitsPreTokenizer(params *util.Params) (tokenizer.PreTokenizer, error) {
	individualDigits := params.Get("individual_digits").(bool)

	return pretokenizer.NewDigits(individualDigits), nil
}

func createUnicodeScriptsPreTokenizer(params *util.Params) (tokenizer.PreTokenizer, error) {
	return pretokenizer.NewUnicodeScript(), nil
}

/*
	{
	    "type": "Split",
	    "pattern": {
	      "Regex": "[0-9][0-9][0-9]"
	    },
	    "behavior": "Isolated",
	    "invert": false
	  }
*/
func createSplitPreTokenizer(params *util.Params) (tokenizer.PreTokenizer, error) {
	if params == nil {
		return nil, nil
	}

	patternMap := params.Get("pattern").(map[string]interface{})
	behaviorVal := params.Get("behavior").(string)
	invert := params.Get("invert").(bool)

	var pattern normalizer.Pattern
	if v, ok := patternMap["Regex"]; ok {
		pattern = normalizer.NewRegexpPattern(v.(string))
	} else if v, ok := patternMap["String"]; ok {
		pattern = normalizer.NewRegexpPattern(v.(string))
	} else {
		err := fmt.Errorf("Unsupported pattern: %#v\n", patternMap)
		return nil, err
	}

	var b normalizer.SplitDelimiterBehavior
	switch behaviorVal {
	case "Removed":
		b = normalizer.RemovedBehavior
	case "Isolated":
		b = normalizer.IsolatedBehavior
	case "MergedWithNext":
		b = normalizer.MergedWithNextBehavior
	case "MergedWithPrevious":
		b = normalizer.MergedWithPreviousBehavior
	case "Contiguous":
		b = normalizer.ContiguousBehavior

	default:
		err := fmt.Errorf("Unsupported behavior: %#v\n", behaviorVal)
		return nil, err
	}

	return pretokenizer.NewSplit(pattern, b, invert), nil
}

func createSequencePreTokenizer(params *util.Params) (tokenizer.PreTokenizer, error) {
	var pretoks []tokenizer.PreTokenizer

	data := params.Get("pretokenizers").([]interface{})
	for _, d := range data {
		pretok, err := CreatePreTokenizer(d.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		pretoks = append(pretoks, pretok)
	}

	out := pretokenizer.NewSequence(pretoks)

	return out, nil
}
