package pretrained

// This file provide some helpers for unit testing only.

import (
	"encoding/json"
	"github.com/sugarme/tokenizer"
	"os"
)

func loadConfig(modelName string) (*tokenizer.Config, error) {
	file, err := tokenizer.CachedPath(modelName, "tokenizer.json")
	if err != nil {
		panic(err)
	}

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	dec := json.NewDecoder(f)

	var config *tokenizer.Config
	err = dec.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config, nil
}
