package pretrained

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/model/bpe"
	"github.com/sugarme/tokenizer/model/wordlevel"
	"github.com/sugarme/tokenizer/model/wordpiece"
	"github.com/sugarme/tokenizer/normalizer"
)

// FromFile instantiates a new Tokenizer from the given file
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

	var (
		dropout                 float32
		unkToken                string
		continuingSubwordPrefix string
		endOfWordSuffix         string
		maxInputCharsPerWord    int
	)

	if config.Model.Dropout != nil {
		dropout = config.Model.Dropout.(float32)
	}
	if config.Model.UnkToken != "" {
		unkToken = config.Model.UnkToken
	}

	if config.Model.ContinuingSubwordPrefix != nil {
		continuingSubwordPrefix = config.Model.ContinuingSubwordPrefix.(string)
	}

	if config.Model.EndOfWordSuffix != nil {
		endOfWordSuffix = config.Model.EndOfWordSuffix.(string)
	}

	if config.Model.MaxInputCharsPerWord != nil {
		maxInputCharsPerWord = config.Model.MaxInputCharsPerWord.(int)
	}

	var model tokenizer.Model
	switch config.Model.Type {
	case "BPE":
		model, err = bpe.New(config.Model.Vocab, config.Model.Merges, &dropout, &unkToken, &continuingSubwordPrefix, &endOfWordSuffix)
		if err != nil {
			return nil, err
		}

	case "WordPiece":
		model, err = wordpiece.New(config.Model.Vocab, &unkToken, &continuingSubwordPrefix, &maxInputCharsPerWord)
		if err != nil {
			return nil, err
		}

	case "WordLevel":
		model, err = wordlevel.New(config.Model.Vocab, &unkToken)
		if err != nil {
			return nil, err
		}

	default:
		err := fmt.Errorf("Unsupported model type: %q\n", config.Model.Type)
		return nil, err
	}

	tk := tokenizer.NewTokenizer(model)

	// TODO. continue with config
	// 1. normalizer.Normalizer
	var n normalizer.Normalizer
	if config.Normalizer.Type != "" {
		var norms []normalizer.Normalizer
		// TODO. build normalizers from config
		switch config.Normalizer.Type {
		case "Sequence":
			n = normalizer.NewSequence(norms)
		}
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
