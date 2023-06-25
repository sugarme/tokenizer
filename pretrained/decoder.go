package pretrained

// This file provides functions to create Decoder from json data
// 1. BPE(BPEDecoder),
// 2. ByteLevel(ByteLevel),
// 3. WordPiece(WordPiece),
// 4. Metaspace(Metaspace),
// 5. CTC(CTC),
// 6. Sequence(Sequence),
// 7. Replace(Replace),
// 8. Fuse(Fuse),
// 9. Strip(Strip),
// 10. ByteFallback(ByteFallback),

import (
	"fmt"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/decoder"
	"github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/pretokenizer"
	"github.com/sugarme/tokenizer/util"
)

func CreateDecoder(config map[string]interface{}) (tokenizer.Decoder, error) {
	if config == nil {
		return nil, nil
	}

	params := util.NewParams(config)
	typ := params.Get("type").(string)

	switch typ {
	case "BPE":
		return createBPEDecoder(params)
	case "ByteLevel":
		return createByteLevelDecoder(params)
	case "WordPiece":
		return createWordPieceDecoder(params)
	case "Metaspace":
		return createMetaspaceDecoder(params)
	case "CTC":
		return createCTCDecoder(params)
	case "Sequence":
		return createSequenceDecoder(params)
	case "Replace":
		return createReplaceDecoder(params)
	case "Fuse":
		return createFuseDecoder(params)
	case "Strip":
		return createStripDecoder(params)
	case "ByteFallback":
		return createByteFallbackDecoder(params)
	default:
		err := fmt.Errorf("Could not create tokenizer.Decoder from input data %#v\n", config)
		return nil, err
	}
}

func createBPEDecoder(params *util.Params) (*decoder.BpeDecoder, error) {
	suffix := params.Get("suffix").(string)

	return decoder.NewBpeDecoder(suffix), nil
}

func createByteLevelDecoder(params *util.Params) (*pretokenizer.ByteLevel, error) {
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

func createMetaspaceDecoder(params *util.Params) (*pretokenizer.Metaspace, error) {
	if params == nil {
		return nil, nil
	}

	replacement := params.Get("replacement", "").(string)
	addPrefixSpace := params.Get("add_prefix_space", false).(bool)

	return pretokenizer.NewMetaspace(replacement, addPrefixSpace), nil
}

func createCTCDecoder(params *util.Params) (*decoder.CTC, error) {
	if params == nil {
		return nil, nil
	}

	padToken := params.Get("pad_token").(string)
	wordDelimiter := params.Get("word_delimiter").(string)
	cleanup := params.Get("cleanup").(bool)

	return decoder.NewCTC(padToken, wordDelimiter, cleanup), nil
}

// e.g. `Bert` model
// "decoder":{"type":"WordPiece","prefix":"##","cleanup":true}
func createWordPieceDecoder(params *util.Params) (*decoder.WordPieceDecoder, error) {
	prefix := params.Get("prefix").(string)
	cleanup := params.Get("cleanup").(bool)
	return decoder.NewWordPieceDecoder(prefix, cleanup), nil
}

// create a Sequence Decoder e.g. in llama model
/*
"decoder": {
    "type": "Sequence",
    "decoders": [
      {
        "type": "Replace",
        "pattern": {
          "String": "▁"
        },
        "content": " "
      },
      {
        "type": "ByteFallback"
      },
      {
        "type": "Fuse"
      },
      {
        "type": "Strip",
        "content": " ",
        "start": 1,
        "stop": 0
      }
    ]
  },
*/
func createSequenceDecoder(params *util.Params) (*decoder.Sequence, error) {
	data := params.Get("decoders").([]interface{})
	var decs []tokenizer.Decoder
	for _, v := range data {
		d := v.(map[string]interface{})
		dec, err := CreateDecoder(d)
		if err != nil {
			return nil, err
		}

		decs = append(decs, dec)
	}

	seqDec := decoder.NewSequence(decs)

	return seqDec, nil
}

func createReplaceDecoder(params *util.Params) (*normalizer.Replace, error) {
	if params == nil {
		return nil, nil
	}

	var pattern string
	var patternType normalizer.ReplacePattern
	patternParams := params.Get("pattern").(map[string]interface{})
	pparams := util.NewParams(patternParams)
	switch {
	case pparams.Has("String"):
		pattern = pparams.Get("String").(string)
		patternType = normalizer.String

	case params.Has("Regex"):
		pattern = pparams.Get("Regex").(string)
		patternType = normalizer.String
	}

	content := params.Get("content").(string)

	return normalizer.NewReplace(patternType, pattern, content), nil
}

func createFuseDecoder(params *util.Params) (*decoder.Fuse, error) {
	return decoder.NewFuse(), nil
}

func createByteFallbackDecoder(params *util.Params) (*decoder.ByteFallback, error) {
	return decoder.NewByteFallback(), nil
}

func createStripDecoder(params *util.Params) (*decoder.Strip, error) {

	content := params.Get("content").(string)
	start := int(params.Get("start").(float64))
	stop := int(params.Get("stop").(float64))

	return decoder.NewStrip(content, start, stop), nil
}
