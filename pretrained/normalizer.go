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
// 10. Nmt
// 11. Precompiled
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

	default:
		msg := fmt.Errorf("Could not create Normalizer from config: %#v", config)
		return nil, msg
	}
}

func createBertNormalizer(params *util.Params) (normalizer.Normalizer, error) {
	// TODO
	panic("NotImplementedError")
}

func createReplaceNormalizer(params *util.Params) (normalizer.Normalizer, error) {
	// TODO
	panic("NotImplementedError")
}

func createPrependNormalizer(params *util.Params) (normalizer.Normalizer, error) {
	// TODO
	panic("NotImplementedError")
}

func createStripNormalizer(params *util.Params) (normalizer.Normalizer, error) {
	// TODO
	panic("NotImplementedError")
}

func createStripAccents(params *util.Params) (normalizer.Normalizer, error) {
	// TODO
	panic("NotImplementedError")
}

func createMetaspaceNormalizer(params *util.Params) (normalizer.Normalizer, error) {
	// TODO
	panic("NotImplementedError")
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
	// TODO
	panic("NotImplementedError")
}
