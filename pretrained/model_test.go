package pretrained

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/util"
)

func TestCreateBPE(t *testing.T) {
	// BPE
	file, err := tokenizer.CachedPath("hf-internal-testing/llama-tokenizer", "tokenizer.json")
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

	modelParams := util.NewParams(config.Model)
	m, err := createBPE(modelParams)
	if err != nil {
		panic(err)
	}

	got := m.GetVocabSize()
	want := 32_000

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}
}

func TestCreateWordPiece(t *testing.T) {
	// file, err := tokenizer.CachedPath("hf-internal-testing/tiny-random-bert", "tokenizer.json")
	file, err := tokenizer.CachedPath("bert-base-uncased", "tokenizer.json")
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

	modelParams := util.NewParams(config.Model)
	m, err := createWordPiece(modelParams)
	if err != nil {
		panic(err)
	}

	got := m.GetVocabSize()
	want := 30_522

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}
}

func TestCreateModel(t *testing.T) {
	file, err := tokenizer.CachedPath("bert-base-uncased", "tokenizer.json")
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

	m, err := CreateModel(config.Model)
	if err != nil {
		panic(err)
	}

	got := m.GetVocabSize()
	want := 30_522

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}
}
