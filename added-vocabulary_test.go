package tokenizer_test

import (
	// "fmt"
	"reflect"
	"testing"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
)

type ModelMock struct {
	vocab  map[string]int
	vocabR map[int]string
}

func newModelMock(toks []string, ids []int) ModelMock {
	vocab := make(map[string]int, 0)
	vocabR := make(map[int]string, 0)
	for i := 0; i < len(toks); i++ {
		vocab[toks[i]] = ids[i]
		vocabR[ids[i]] = toks[i]

	}
	return ModelMock{vocab, vocabR}
}

// implement Model interface for ModelMock

func (mm ModelMock) Tokenize(sequence string) (retVal []tokenizer.Token, err error) {
	return // not implement
}

func (mm ModelMock) IdToToken(id int) (retVal string, ok bool) {
	retVal, ok = mm.vocabR[id]
	return
}

func (mm ModelMock) TokenToId(tok string) (retVal int, ok bool) {
	retVal, ok = mm.vocab[tok]
	return
}

func (mm ModelMock) GetVocab() (retVal map[string]int) {
	return mm.vocab
}

func (mm ModelMock) GetVocabSize() (retVal int) {
	return len(mm.vocab)
}

func (mm ModelMock) Save(dir string, prefixOpt ...string) (err error) {
	return // not implement
}

func TestCanAddTokens(t *testing.T) {
	model := newModelMock([]string{"test", "tost"}, []int{0, 1})
	vocab := tokenizer.NewAddedVocabulary()

	addedToks := []tokenizer.AddedToken{
		tokenizer.NewAddedToken("added_token_1", false),
		tokenizer.NewAddedToken("added_token_2", false),
	}

	got := vocab.AddTokens(addedToks, model, nil)
	want := 2

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want %v\n", want)
		t.Errorf("Got %v\n", got)
	}
}

func TestAddSpecialTokens(t *testing.T) {
	model := newModelMock([]string{"test", "tost"}, []int{0, 1})
	vocab := tokenizer.NewAddedVocabulary()

	specialTok := []tokenizer.AddedToken{
		tokenizer.NewAddedToken("added_token_1", true),
	}

	got1 := vocab.AddSpecialTokens(specialTok, model, nil)
	want1 := 1

	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("Want %v\n", want1)
		t.Errorf("Got %v\n", got1)
	}

	// Does not add multiple time the same token
	otherSpecialToks := []tokenizer.AddedToken{
		tokenizer.NewAddedToken("added_token_2", true),
		tokenizer.NewAddedToken("added_token_2", true),
	}

	vocab.AddSpecialTokens(otherSpecialToks, model, nil)
	got2 := vocab.Len()
	want2 := 2

	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("Want %v\n", want2)
		t.Errorf("Got %v\n", got2)
	}

	// Can add tokens already covered by the model
	got3 := vocab.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("test", true)}, model, nil)
	want3 := 0
	if !reflect.DeepEqual(want3, got3) {
		t.Errorf("Want %v\n", want3)
		t.Errorf("Got %v\n", got3)
	}

	got4 := vocab.Len() // Did not add new token as it exists in the original model
	want4 := 2
	if !reflect.DeepEqual(want4, got4) {
		t.Errorf("Want %v\n", want4)
		t.Errorf("Got %v\n", got4)
	}

	got5 := vocab.IsSpecialToken("test")
	want5 := true
	if !reflect.DeepEqual(want5, got5) {
		t.Errorf("Want %v\n", want5)
		t.Errorf("Got %v\n", got5)
	}
}

func TestCanExtractAddedTokens(t *testing.T) {
	// Able to extract both normal and special tokens
	model := newModelMock([]string{}, []int{})
	vocab := tokenizer.NewAddedVocabulary()

	addedToks := []tokenizer.AddedToken{
		tokenizer.NewAddedToken("my", false),
		tokenizer.NewAddedToken("name", false),
	}

	addedSpecialToks := []tokenizer.AddedToken{
		tokenizer.NewAddedToken("[CLS]", true),
		tokenizer.NewAddedToken("[SEP]", true),
	}

	vocab.AddTokens(addedToks, model, nil)
	vocab.AddSpecialTokens(addedSpecialToks, model, nil)

	result := vocab.ExtractAndNormalize("[CLS] My name is Anthony [SEP]", nil)

	type tokenid struct {
		tokens string
		ids    []int
	}

	var got []tokenid
	pretoks := result.GetSplits(normalizer.OriginalTarget)
	for _, pretok := range pretoks {
		var tokIds []int
		for _, tok := range pretok.Tokens {
			tokIds = append(tokIds, tok.Id)
		}
		got = append(got, tokenid{pretok.Value, tokIds})
	}

	want := []tokenid{
		{"[CLS]", []int{2}},
		{" My ", nil},
		{"name", []int{1}},
		{" is Anthony ", nil},
		{"[SEP]", []int{3}},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want %+v\n", want)
		t.Errorf("Got %+v\n", got)
	}
}

func TestOptionUseCases(t *testing.T) {
	// Is able to extract both normal and special tokens, with various options (lstrip, rstrip,
	// single_word, normalized)
	model := newModelMock([]string{}, []int{})
	vocab := tokenizer.NewAddedVocabulary()
	n := normalizer.Lowercase()

	addedToks := []tokenizer.AddedToken{
		tokenizer.NewAddedToken("my", false).SetLStrip(true).SetRStrip(true),
		tokenizer.NewAddedToken("name", false),
		tokenizer.NewAddedToken("ony", false).SetSingleWord(true),
	}

	addedSpecialToks := []tokenizer.AddedToken{
		tokenizer.NewAddedToken("[CLS]", true),
		tokenizer.NewAddedToken("[SEP]", true),
	}

	vocab.AddTokens(addedToks, model, nil)
	vocab.AddSpecialTokens(addedSpecialToks, model, nil)

	result := vocab.ExtractAndNormalize("[CLS] My name is Anthony [SEP]", n)

	type tokenid struct {
		token string
		ids   []int
	}

	var got []tokenid
	pretoks := result.GetSplits(normalizer.OriginalTarget)
	for _, pretok := range pretoks {
		var tokIds []int
		for _, tok := range pretok.Tokens {
			tokIds = append(tokIds, tok.Id)
		}
		got = append(got, tokenid{pretok.Value, tokIds})
	}

	want := []tokenid{
		{"[CLS]", []int{3}},
		// This one includes both spaces because of the lstrip & rstrip
		// And it matches because normalized == true
		{" my ", []int{0}},
		{"name", []int{1}},
		// `ony` is not extracted here thanks to single_word
		{" is anthony ", nil},
		{"[SEP]", []int{4}},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want %+v\n", want)
		t.Errorf("Got %+v\n", got)
	}
}
