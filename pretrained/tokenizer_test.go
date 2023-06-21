package pretrained

import (
	"log"
	"testing"

	"github.com/sugarme/tokenizer"
)

func TestFromFile(t *testing.T) {
	configFile, err := tokenizer.CachedPath("hf-internal-testing/llama-tokenizer", "tokenizer.json")
	if err != nil {
		panic(err)
	}

	tk, err := FromFile(configFile)
	if err != nil {
		panic(err)
	}

	log.Printf("tokenizer: %#v\n", tk)

	log.Printf("vocab size (excl. added tokens): %v\n", tk.GetVocabSize(false))
	log.Printf("vocab size (incl. added tokens): %v\n", tk.GetVocabSize(true))
	log.Printf("special tokens: %q\n", tk.GetSpecialTokens())
}

/*
{
  "version": "1.0",
  "truncation": null,
  "padding": null,
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
  "normalizer": {
    "type": "Sequence",
    "normalizers": [
      {
        "type": "Prepend",
        "prepend": "▁"
      },
      {
        "type": "Replace",
        "pattern": {
          "String": " "
        },
        "content": "▁"
      }
    ]
  },
  "pre_tokenizer": null,
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
        "SpecialToken": {
          "id": "<s>",
          "type_id": 1
        }
      },
      {
        "Sequence": {
          "id": "B",
          "type_id": 1
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
  "decoder": {
    "type": "Sequence",
    "decoders": [
      {
        "type": "Replace",
        "pattern": {
          "String": "▁"
        },
        "content": " "
      },
      {
        "type": "ByteFallback"
      },
      {
        "type": "Fuse"
      },
      {
        "type": "Strip",
        "content": " ",
        "start": 1,
        "stop": 0
      }
    ]
  },
  "model": {
    "type": "BPE",
    "dropout": null,
    "unk_token": "<unk>",
    "continuing_subword_prefix": null,
    "end_of_word_suffix": null,
    "fuse_unk": true,
    "byte_fallback": true,
    "vocab": {
      "<unk>": 0,
      "<s>": 1,
      "</s>": 2,
      ...
    }
  }
}

*/
