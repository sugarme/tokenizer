package tokenizer

import (
	"encoding/json"
	"fmt"
	"os"
)

func ExampleConfig() {
	tokFile, err := CachedPath("hf-internal-testing/llama-tokenizer", "tokenizer.json")
	if err != nil {
		panic(err)
	}

	f, err := os.Open(tokFile)
	if err != nil {
		panic(err)
	}

	dec := json.NewDecoder(f)

	var config *Config

	err = dec.Decode(&config)
	if err != nil {
		panic(err)
	}

	modelType := config.Model.Type
	fmt.Println(modelType)

	// Output:
	// BPE
}
