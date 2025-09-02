package pretrained

import (
	"reflect"
	"testing"

	"github.com/season-studio/tokenizer/util"
)

func TestCreateBPE(t *testing.T) {
	modelName := "hf-internal-testing/llama-tokenizer"
	config, err := loadConfig(modelName)
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
	modelName := "bert-base-uncased"
	config, err := loadConfig(modelName)
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
	modelName := "bert-base-uncased"
	config, err := loadConfig(modelName)
	if err != nil {
		panic(err)
	}

	m, err := CreateModel(config)
	if err != nil {
		panic(err)
	}

	got := m.GetVocabSize()
	want := 30_522

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}
}
