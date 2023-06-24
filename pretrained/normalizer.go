package pretrained

// This file provides functions to create Normalizer from json data
// 1. BertNormalizer
// 2. StripNormalizer
// 3. StripAccents
// 4. NFC
// 5. NFD
// 6. NFKC
// 7. NFKD
// 8. Sequence
// 9. Lowercase
// 10. Nmt (TODO)
// 11. Precompiled (TODO)
// 12. Replace
// 13. Prepend

import (
	"fmt"

	"github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/util"
)

// CreateNormalizer creates Normalizer from config data.
func CreateNormalizer(config map[string]interface{}) (normalizer.Normalizer, error) {
	// No Normalizer at all
	if config == nil {
		return nil, nil
	}

	params := util.NewParams(config)

	typ := params.Get("type").(string)

	switch typ {
	case "BertNormalizer":
		return createBertNormalizer(params)
	case "StripNormalizer":
		return createStripNormalizer(params)
	case "StripAccents":
		return createStripAccents(params)

		// unicode normalizers
	case "NFC":
		return normalizer.NewNFC(), nil
	case "NFD":
		return normalizer.NewNFD(), nil
	case "NFKC":
		return normalizer.NewNFKC(), nil
	case "NFKD":
		return normalizer.NewNFKD(), nil

	case "Sequence":
		return createSequenceNormalizer(params)

	case "Lowercase":
		return normalizer.Lowercase(), nil

	case "Nmt":
		return createNmtNormalizer(params)

	case "Precompiled":
		return createPrecompiledNormalizer(params)

	case "Replace":
		return createReplaceNormalizer(params)

	case "Prepend":
		return createPrependNormalizer(params)

	default:
		msg := fmt.Errorf("Could not create Normalizer from config: %#v", config)
		return nil, msg
	}
}

// BertNormalizer json data:
// -------------------------
// "type":"BertNormalizer"
// "clean_text":true
// "handle_chinese_chars":true
// "strip_accents":null
// "lowercase":true
func createBertNormalizer(params *util.Params) (normalizer.Normalizer, error) {
	if params == nil {
		return nil, nil
	}

	cleanText := params.Get("clean_text", false).(bool)
	handleChineseChars := params.Get("handle_chinese_chars", false).(bool)
	stripAccents := params.Get("strip_accents", false).(bool)
	lowercase := params.Get("lowercase", false).(bool)

	return normalizer.NewBertNormalizer(cleanText, lowercase, handleChineseChars, stripAccents), nil
}

func createReplaceNormalizer(params *util.Params) (normalizer.Normalizer, error) {
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

func createPrependNormalizer(params *util.Params) (normalizer.Normalizer, error) {
	var prepend string
	if params.Has("prepend") {
		prepend = params.Get("prepend").(string)
	}

	return normalizer.NewPrepend(prepend), nil
}

func createStripNormalizer(params *util.Params) (normalizer.Normalizer, error) {
	stripLeft := params.Get("strip_left", false).(bool)
	stripRight := params.Get("strip_right", false).(bool)

	return normalizer.NewStrip(stripLeft, stripRight), nil
}

func createStripAccents(params *util.Params) (normalizer.Normalizer, error) {
	return normalizer.NewStripAccents(), nil
}

func createPrecompiledNormalizer(params *util.Params) (normalizer.Normalizer, error) {
	// TODO
	panic("NotImplementedError")
}

func createNmtNormalizer(params *util.Params) (normalizer.Normalizer, error) {
	// TODO
	panic("NotImplementedError")
}

func createSequenceNormalizer(params *util.Params) (normalizer.Normalizer, error) {
	var data []interface{}
	if params.Has("normalizers") {
		data = params.Get("normalizers").([]interface{})
	}

	var norms []normalizer.Normalizer
	for _, d := range data {
		n, err := CreateNormalizer(d.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		norms = append(norms, n)
	}

	seq := normalizer.NewSequence(norms)

	return seq, nil
}
