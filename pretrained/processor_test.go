package pretrained

import (
	"log"
	"testing"

	"github.com/sugarme/tokenizer"
)

// e.g. `hf-internal-testing/llama-tokenizer`
// https://huggingface.co/hf-internal-testing/llama-tokenizer/raw/main/tokenizer.json
/*
 "post_processor": {
    "type": "TemplateProcessing",
    "single": [
      {
        "SpecialToken": {
          "id": "<s>",
          "type_id": 0
        }
      },
      {
        "Sequence": {
          "id": "A",
          "type_id": 0
        }
      }
    ],
    "pair": [
      {
        "SpecialToken": {
          "id": "<s>",
          "type_id": 0
        }
      },
      {
        "Sequence": {
          "id": "A",
          "type_id": 0
        }
      },
      {
        "Sequence": {
          "id": "B",
          "type_id": 0
        }
      }
    ],
    "special_tokens": {
      "<s>": {
        "id": "<s>",
        "ids": [
          1
        ],
        "tokens": [
          "<s>"
        ]
      }
    }
  },
*/

func TestCreatePostProcessor(t *testing.T) {
	modelName := "hf-internal-testing/llama-tokenizer"
	config, err := loadConfig(modelName)
	if err != nil {
		panic(err)
	}

	log.Printf("config: %#v\n", config.PostProcessor)

	var p tokenizer.PostProcessor
	p, err = CreatePostProcessor(config.PostProcessor)
	if err != nil {
		panic(err)
	}

	log.Printf("processor: %#v\n", p)

	got := p.AddedTokens(true)
	want := 2

	if got != want {
		t.Fatalf("Expected %v, got %v\n", want, got)
	}

	got = p.AddedTokens(false)
	want = 1

	if got != want {
		t.Fatalf("Expected %v, got %v\n", want, got)
	}

}
