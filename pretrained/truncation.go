package pretrained

import (
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/util"
)

func CreateTruncationParams(config map[string]interface{}) (*tokenizer.TruncationParams, error) {
	if config == nil {
		return nil, nil
	}

	params := util.NewParams(config)
	// direction := params.Get("direction").(string)
	maxLen := int(params.Get("max_length").(float64))
	stride := int(params.Get("stride").(float64))
	strategyName := params.Get("strategy").(string)

	var strategy tokenizer.TruncationStrategy
	switch strategyName {
	case "LongestFirst":
		strategy = tokenizer.LongestFirst
	case "OnlyFirst":
		strategy = tokenizer.OnlyFirst
	case "OnlySecond":
		strategy = tokenizer.OnlySecond
	}

	return &tokenizer.TruncationParams{
		MaxLength: maxLen,
		Strategy:  strategy,
		Stride:    stride,
	}, nil
}
