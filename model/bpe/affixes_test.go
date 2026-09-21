package bpe_test

import (
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/model"
	"github.com/sugarme/tokenizer/model/bpe"
	"reflect"
	"testing"
)

func TestBPEWordAffixPlacement(t *testing.T) {
	for _, tc := range []struct {
		name, text, prefix, suffix string
		vocab                      model.Vocab
		want                       []tokenizer.Token
	}{
		{"single character suffix", "a", "", "</w>", model.Vocab{"[UNK]": 0, "a</w>": 1}, []tokenizer.Token{{Id: 1, Value: "a</w>", Offsets: []int{0, 1}}}},
		{"last character suffix", "ab", "", "</w>", model.Vocab{"[UNK]": 0, "a": 1, "b</w>": 2}, []tokenizer.Token{{Id: 1, Value: "a", Offsets: []int{0, 1}}, {Id: 2, Value: "b</w>", Offsets: []int{1, 2}}}},
		{"Unicode suffix", "aé", "", "</w>", model.Vocab{"[UNK]": 0, "a": 1, "é</w>": 2}, []tokenizer.Token{{Id: 1, Value: "a", Offsets: []int{0, 1}}, {Id: 2, Value: "é</w>", Offsets: []int{1, 3}}}},
		{"continuation prefix", "ab", "##", "", model.Vocab{"[UNK]": 0, "a": 1, "##b": 2}, []tokenizer.Token{{Id: 1, Value: "a", Offsets: []int{0, 1}}, {Id: 2, Value: "##b", Offsets: []int{1, 2}}}},
		{"combined affixes", "ab", "##", "</w>", model.Vocab{"[UNK]": 0, "a": 1, "##b</w>": 2}, []tokenizer.Token{{Id: 1, Value: "a", Offsets: []int{0, 1}}, {Id: 2, Value: "##b</w>", Offsets: []int{1, 2}}}},
		{"single character with both options", "é", "##", "</w>", model.Vocab{"[UNK]": 0, "é</w>": 1}, []tokenizer.Token{{Id: 1, Value: "é</w>", Offsets: []int{0, 2}}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := bpe.NewBPE(tc.vocab, bpe.Merges{})
			unk := "[UNK]"
			m.UnkToken = &unk
			m.ContinuingSubwordPrefix = &tc.prefix
			m.EndOfWordSuffix = &tc.suffix
			for i := 0; i < 2; i++ {
				got, err := m.Tokenize(tc.text)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("tokens = %v, want %v", got, tc.want)
				}
			}
		})
	}
}

func TestBPESuffixParticipatesInMerges(t *testing.T) {
	vocab := model.Vocab{"a": 0, "b</w>": 1, "ab</w>": 2}
	merges := bpe.Merges{bpe.Pair{C1: 0, C2: 1}: bpe.PairVal{Rank: 0, NewId: 2}}
	m := bpe.NewBPE(vocab, merges)
	suffix := "</w>"
	m.EndOfWordSuffix = &suffix
	got, err := m.Tokenize("ab")
	if err != nil {
		t.Fatal(err)
	}
	want := []tokenizer.Token{{Id: 2, Value: "ab</w>", Offsets: []int{0, 2}}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tokens = %v, want %v", got, want)
	}
}
