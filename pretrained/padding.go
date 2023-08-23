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

	var strategyName string
	var strategySize int
	for k, v := range params.Get("strategy").(map[string]interface{}) {
		strategyName = k
		strategySize = int(v.(float64))
	}

	switch strategyName {
	case "BatchLongest":
		opt := tokenizer.WithBatchLongest()
		strategy = tokenizer.NewPaddingStrategy(opt)
	case "Fixed":
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
