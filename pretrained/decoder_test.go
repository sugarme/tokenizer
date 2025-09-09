package pretrained

import (
	"log"
	"testing"

	"github.com/gengzongjie/tokenizer"
)

/*
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
*/

func TestCreateDecoder(t *testing.T) {
	modelName := "hf-internal-testing/llama-tokenizer"
	config, err := loadConfig(modelName)
	if err != nil {
		panic(err)
	}

	log.Printf("config: %#v\n", config.Decoder)

	var d tokenizer.Decoder
	d, err = CreateDecoder(config.Decoder)
	if err != nil {
		panic(err)
	}

	log.Printf("decoder: %#v\n", d)

	// TODO. test decode a string here...
}

// Output:
// decoder: &decoder.Sequence{DecoderBase:(*decoder.DecoderBase)(0xc0002494f0), decoders:[]tokenizer.Decoder{(*normalizer.Replace)(0xc00007fda0), (*decoder.ByteFallback)(0xc000012468), (*decoder.Fuse)(0xc0000140c8), (*decoder.Strip)(0xc00007fdd0)}}
