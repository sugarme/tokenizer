package pretrained

import (
	"testing"

	"github.com/sugarme/tokenizer/pretokenizer"
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

func TestCreateByteLevelUseRegexOption(t *testing.T) {
	config := map[string]interface{}{
		"type":             "ByteLevel",
		"add_prefix_space": false,
		"trim_offsets":     true,
		"use_regex":        false,
	}

	pt, err := CreatePreTokenizer(config)
	if err != nil {
		t.Fatalf("CreatePreTokenizer error: %v", err)
	}

	bl, ok := pt.(*pretokenizer.ByteLevel)
	if !ok {
		t.Fatalf("expected *pretokenizer.ByteLevel, got %T", pt)
	}

	if bl.UseRegex {
		t.Fatalf("expected UseRegex=false from config")
	}
}
