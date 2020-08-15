package tokenizer_test

import (
	"reflect"
	"testing"

	"github.com/sugarme/tokenizer"
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
