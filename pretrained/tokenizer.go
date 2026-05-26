package pretrained

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/sugarme/tokenizer"
)

// FromFile constructs a new Tokenizer from json data file (normally 'tokenizer.json')
func FromFile(file string) (*tokenizer.Tokenizer, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tk, err := FromReader(f)
	if err != nil {
		err := fmt.Errorf("FromReader: %w", err)
		return nil, err
	}
	return tk, nil
}

// FromReader constructs a new Tokenizer from json data reader.
func FromReader(r io.Reader) (*tokenizer.Tokenizer, error) {
	var config *tokenizer.Config
	if err := json.NewDecoder(r).Decode(&config); err != nil {
		return nil, err
	}

	model, err := CreateModel(config)
	if err != nil {
		err = fmt.Errorf("CreateModel: %w", err)
		return nil, err
	}

	tk := tokenizer.NewTokenizer(model)

	// 2. Normalizer
	n, err := CreateNormalizer(config.Normalizer)
	if err != nil {
		err = fmt.Errorf("CreateNormalizer: %v", err)
		return nil, err
	}
	tk.WithNormalizer(n)

	// 3. PreTokenizer
	preTok, err := CreatePreTokenizer(config.PreTokenizer)
	if err != nil {
		err = fmt.Errorf("CreatePreTokenizer: %v", err)
		return nil, err
	}
	tk.WithPreTokenizer(preTok)

	// 4. PostProcessor
	postProcessor, err := CreatePostProcessor(config.PostProcessor)
	if err != nil {
		err = fmt.Errorf("CreatePostProcessor: %v", err)
		return nil, err
	}
	tk.WithPostProcessor(postProcessor)

	// 5. Decoder
	decoder, err := CreateDecoder(config.Decoder)
	if err != nil {
		err = fmt.Errorf("CreateDecoder: %v", err)
		return nil, err
	}
	tk.WithDecoder(decoder)

	// 6. AddedVocabulary — use ID-preserving path so that compacted
	//    tokenizers keep the exact added-token IDs from tokenizer.json.
	addedTokensWithIds := CreateAddedTokensWithIds(config.AddedTokens)
	if len(addedTokensWithIds) > 0 {
		tk.AddTokensWithIds(addedTokensWithIds)
	}

	// 7. TruncationParams
	truncParams, err := CreateTruncationParams(config.Truncation)
	if err != nil {
		err = fmt.Errorf("CreatingTruncationParams: %v", err)
		return nil, err
	}
	tk.WithTruncation(truncParams)

	// 8. PaddingParams
	paddingParams, err := CreatePaddingParams(config.Padding)
	if err != nil {
		err = fmt.Errorf("CreatePaddingParams: %v", err)
		return nil, err
	}
	tk.WithPadding(paddingParams)

	return tk, nil
}
