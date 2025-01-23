package pretrained

import (
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/util"
)

func CreatePaddingParams(config map[string]interface{}) (*tokenizer.PaddingParams, error) {
	if config == nil {
		return nil, nil
	}

	params := util.NewParams(config)
	var strategy *tokenizer.PaddingStrategy

	strategyName := params.Get("strategy").(string)

	switch strategyName {
	case "BatchLongest":
		opt := tokenizer.WithBatchLongest()
		strategy = tokenizer.NewPaddingStrategy(opt)
	case "Fixed":
		strategySize := int(params.Get("size").(float64))
		opt := tokenizer.WithFixed(strategySize)
		strategy = tokenizer.NewPaddingStrategy(opt)
	}

	directionName := params.Get("direction").(string)
	var direction tokenizer.PaddingDirection
	switch directionName {
	case "left", "Left":
		direction = tokenizer.Left
	case "right", "Right":
		direction = tokenizer.Right
	}

	id := int(params.Get("pad_id").(float64))
	typeId := int(params.Get("pad_type_id").(float64))
	token := params.Get("pad_token").(string)

	return &tokenizer.PaddingParams{
		Strategy:  *strategy,
		Direction: direction,
		PadId:     id,
		PadTypeId: typeId,
		PadToken:  token,
	}, nil
}
