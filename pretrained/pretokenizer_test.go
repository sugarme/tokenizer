package pretrained

import (
	"testing"
)

func TestCreatePreTokenizer(t *testing.T) {
	modelName := "facebook/bart-base"
	config, err := loadConfig(modelName)
	if err != nil {
		panic(err)
	}

	_, err = CreatePreTokenizer(config.PreTokenizer)
	if err != nil {
		panic(err)
	}
}

func TestNullPreTokenizer(t *testing.T) {
	modelName := "hf-internal-testing/llama-tokenizer"
	config, err := loadConfig(modelName)
	if err != nil {
		panic(err)
	}

	_, err = CreatePreTokenizer(config.PreTokenizer)
	if err != nil {
		panic(err)
	}
}
