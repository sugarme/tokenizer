package tokenizer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/season-studio/tokenizer/util"
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

	modelConfig := util.NewParams(config.Model)

	modelType := modelConfig.Get("type", "").(string)
	fmt.Println(modelType)

	// Output:
	// BPE
}
