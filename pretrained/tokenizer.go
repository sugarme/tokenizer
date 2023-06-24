package pretrained

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sugarme/tokenizer"
)

// FromFile constructs a new Tokenizer from json data file (normally 'tokenizer.json')
func FromFile(file string) (*tokenizer.Tokenizer, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(f)

	var config *tokenizer.Config
	err = dec.Decode(&config)
	if err != nil {
		return nil, err
	}

	model, err := CreateModel(config)
	if err != nil {
		err := fmt.Errorf("Creating Model failed: %v", err)
		return nil, err
	}

	tk := tokenizer.NewTokenizer(model)

	// 2. Normalizer
	n, err := CreateNormalizer(config.Normalizer)
	if err != nil {
		err = fmt.Errorf("Creating Normalizer failed: %v", err)
		return nil, err
	}
	tk.WithNormalizer(n)

	// 3. PreTokenizer
	preTok, err := CreatePreTokenizer(config.PreTokenizer)
	if err != nil {
		err = fmt.Errorf("Creating PreTokenizer failed: %v", err)
		return nil, err
	}
	tk.WithPreTokenizer(preTok)

	// 4. PostProcessor
	postProcessor, err := CreatePostProcessor(config.PostProcessor)
	if err != nil {
		err = fmt.Errorf("Creating PostProcessor failed: %v", err)
		return nil, err
	}
	tk.WithPostProcessor(postProcessor)

	// 5. Decoder
	decoder, err := CreateDecoder(config.Decoder)
	if err != nil {
		err = fmt.Errorf("Creating Decoder failed: %v", err)
		return nil, err
	}
	tk.WithDecoder(decoder)

	// 6. AddedVocabulary
	specialAddedTokens, addedTokens := CreateAddedTokens(config.AddedTokens)
	if len(specialAddedTokens) > 0 {
		tk.AddSpecialTokens(specialAddedTokens)
	}
	if len(addedTokens) > 0 {
		tk.AddTokens(addedTokens)
	}

	// 7. TruncationParams
	truncParams, err := CreateTruncationParams(config.Truncation)
	if err != nil {
		err = fmt.Errorf("Creating TruncationParams failed: %v", err)
		return nil, err
	}
	tk.WithTruncation(truncParams)

	// 8. PaddingParams
	paddingParams, err := CreatePaddingParams(config.Padding)
	if err != nil {
		err = fmt.Errorf("Creating PaddingParams failed: %v", err)
		return nil, err
	}
	tk.WithPadding(paddingParams)

	return tk, nil
}
