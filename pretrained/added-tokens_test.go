package pretrained

import (
	"log"
	"testing"
)

// llama model
/*
 "added_tokens": [
    {
      "id": 0,
      "content": "<unk>",
      "single_word": false,
      "lstrip": false,
      "rstrip": false,
      "normalized": true,
      "special": true
    },
    {
      "id": 1,
      "content": "<s>",
      "single_word": false,
      "lstrip": false,
      "rstrip": false,
      "normalized": true,
      "special": true
    },
    {
      "id": 2,
      "content": "</s>",
      "single_word": false,
      "lstrip": false,
      "rstrip": false,
      "normalized": true,
      "special": true
    }
  ],
*/

func TestCreateAddedTokens(t *testing.T) {
	modelName := "hf-internal-testing/llama-tokenizer"
	config, err := loadConfig(modelName)
	if err != nil {
		panic(err)
	}

	log.Printf("config: %#v\n", config.AddedTokens)

	specialToks, toks := CreateAddedTokens(config.AddedTokens)

	log.Printf("specialAddedTokens: %#v\n", specialToks)
	log.Printf("addedTokens: %#v\n", toks)
}

// Output:
