package pretrained

import (
	"reflect"

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

	// Handle different formats of strategy field
	strategyValue := params.Get("strategy")
	switch v := strategyValue.(type) {
	case map[string]interface{}:
		// Handle object format like {"Fixed": 128}
		for k, val := range v {
			strategyName = k
			// Handle different numeric types
			switch numVal := val.(type) {
			case float64:
				strategySize = int(numVal)
			case int:
				strategySize = numVal
			default:
				// Try to convert using reflection
				valValue := reflect.ValueOf(val)
				if valValue.Kind() == reflect.Float64 || valValue.Kind() == reflect.Int {
					strategySize = int(valValue.Int())
				} else {
					strategySize = 512 // Default fallback
				}
			}
		}
	case string:
		// Handle string format like "Fixed"
		strategyName = v
		// Look for size field
		if sizeVal, ok := config["size"]; ok {
			switch s := sizeVal.(type) {
			case float64:
				strategySize = int(s)
			case int:
				strategySize = s
			default:
				strategySize = 512 // Default fallback
			}
		} else {
			strategySize = 512 // Default size if not specified
		}
	default:
		// Handle unexpected type with a reasonable default
		strategyName = "Fixed"
		strategySize = 512
	}

	// Create strategy based on name
	switch strategyName {
	case "BatchLongest":
		opt := tokenizer.WithBatchLongest()
		strategy = tokenizer.NewPaddingStrategy(opt)
	case "Fixed":
		opt := tokenizer.WithFixed(strategySize)
		strategy = tokenizer.NewPaddingStrategy(opt)
	default:
		// Default to Fixed with size
		opt := tokenizer.WithFixed(strategySize)
		strategy = tokenizer.NewPaddingStrategy(opt)
	}

	// Get direction with fallback to Right
	var direction tokenizer.PaddingDirection
	directionVal := params.Get("direction")
	if dirVal, ok := directionVal.(string); ok {
		switch dirVal {
		case "left", "Left":
			direction = tokenizer.Left
		case "right", "Right", "":
			direction = tokenizer.Right
		default:
			direction = tokenizer.Right
		}
	} else {
		direction = tokenizer.Right // Default
	}

	// Get other parameters with fallbacks
	var id, typeId int
	var token string

	if idVal := params.Get("pad_id"); idVal != nil {
		if fVal, ok := idVal.(float64); ok {
			id = int(fVal)
		}
	}

	if typeIdVal := params.Get("pad_type_id"); typeIdVal != nil {
		if fVal, ok := typeIdVal.(float64); ok {
			typeId = int(fVal)
		}
	}

	if tokenVal := params.Get("pad_token"); tokenVal != nil {
		if sVal, ok := tokenVal.(string); ok {
			token = sVal
		} else {
			token = "[PAD]" // Default pad token
		}
	} else {
		token = "[PAD]" // Default pad token
	}

	return &tokenizer.PaddingParams{
		Strategy:  *strategy,
		Direction: direction,
		PadId:     id,
		PadTypeId: typeId,
		PadToken:  token,
	}, nil
}
