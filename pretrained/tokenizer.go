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

	model, err := CreateModel(config.Model)
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

	// 2. PreTokenizer
	// 3. PostProcessor
	// 4. Decoder
	// 5. AddedVocabulary
	// 6. TruncationParams
	// 7. PaddingParams

	return tk, nil
}
