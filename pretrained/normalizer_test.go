package pretrained

import (
	"testing"
)

func TestCreateSequenceNormalizer(t *testing.T) {
	modelName := "hf-internal-testing/llama-tokenizer"
	config, err := loadConfig(modelName)
	if err != nil {
		panic(err)
	}

	_, err = CreateNormalizer(config.Normalizer)
	if err != nil {
		panic(err)
	}
}

func TestCreateBertNormalizer(t *testing.T) {
	modelName := "bert-base-uncased"
	config, err := loadConfig(modelName)
	if err != nil {
		panic(err)
	}

	_, err = CreateNormalizer(config.Normalizer)
	if err != nil {
		panic(err)
	}
}

func TestCreateNFCNormalizer(t *testing.T) {
	modelName := "mosaicml/mpt-7b"
	config, err := loadConfig(modelName)
	if err != nil {
		panic(err)
	}

	_, err = CreateNormalizer(config.Normalizer)
	if err != nil {
		panic(err)
	}
}

func TestNullNormalizer(t *testing.T) {
	modelName := "tiiuae/falcon-7b"
	config, err := loadConfig(modelName)
	if err != nil {
		panic(err)
	}

	_, err = CreateNormalizer(config.Normalizer)
	if err != nil {
		panic(err)
	}
}

func TestCreatePrecompiledNormalizer(t *testing.T) {
	modelName := "MoritzLaurer/deberta-v3-base-zeroshot-v1"
	config, err := loadConfig(modelName)
	if err != nil {
		panic(err)
	}

	_, err = CreateNormalizer(config.Normalizer)
	if err != nil {
		panic(err)
	}
}
