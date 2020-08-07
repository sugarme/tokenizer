package wordpiece

import (
	"github.com/sugarme/tokenizer/model/bpe"
	"github.com/sugarme/tokenizer/tokenizer"
)

// A `WordPieceTrainerBuilder` can be used to create a `WordPieceTrainer` with a custom
// configuration.
type WordPieceTrainerBuilder struct {
	bpeTrainerBuilder bpe.BpeTrainerBuilder
}

func NewWordPieceTrainerBuilder() (retVal WordPieceTrainerBuilder) {

	bpeTrainerBuilder := *bpe.NewBPETrainerBuilder()
	bpeTrainerBuilder.ContinuingSubwordPrefix("##")

	return WordPieceTrainerBuilder{
		bpeTrainerBuilder: bpeTrainerBuilder,
	}
}

func (wptb WordPieceTrainerBuilder) MinFrequency(frequency uint32) (retVal WordPieceTrainerBuilder) {
	wptb.bpeTrainerBuilder.MinFrequency(frequency)

	return wptb
}

// VocabSize Set the vocabulary size
func (wptb WordPieceTrainerBuilder) VocabSize(size uint) (retVal WordPieceTrainerBuilder) {

	wptb.bpeTrainerBuilder.VocabSize(size)
	return wptb
}

// ShowProgress set whether to show progress
func (wptb WordPieceTrainerBuilder) ShowProgress(show bool) (retVal WordPieceTrainerBuilder) {

	wptb.bpeTrainerBuilder.ShowProgress(show)
	return wptb
}

// SpecialTokens Set the special tokens
func (wptb WordPieceTrainerBuilder) SpecialTokens(tokens []string) (retVal WordPieceTrainerBuilder) {

	wptb.bpeTrainerBuilder.SpecialTokens(tokens)
	return wptb
}

// LimitAlphabet Set whether to limit the alphabet
func (wptb WordPieceTrainerBuilder) LimitAlphabet(limit uint) (retVal WordPieceTrainerBuilder) {

	wptb.bpeTrainerBuilder.LimitAlphabet(limit)
	return wptb
}

// InitialAlphabet set the initial alphabet
func (wptb WordPieceTrainerBuilder) InitialAlphabet(alphabet bpe.CharSet) (retVal WordPieceTrainerBuilder) {

	wptb.bpeTrainerBuilder.InitialAlphabet(alphabet)
	return wptb
}

// ContinuingSubwordPrefix set the continuing_subword_prefix
func (wptb WordPieceTrainerBuilder) ContinuingSubwordPrefix(prefix string) (retVal WordPieceTrainerBuilder) {

	wptb.bpeTrainerBuilder.ContinuingSubwordPrefix(prefix)
	return wptb
}

// EndOfWordSuffix set the end_of_word_suffix
func (wptb WordPieceTrainerBuilder) EndOfWordSuffix(suffix string) (retVal WordPieceTrainerBuilder) {

	wptb.bpeTrainerBuilder.EndOfWordSuffix(suffix)
	return wptb
}

// Build constructs the final BpeTrainer
func (wptb WordPieceTrainerBuilder) Build() (retVal WordPieceTrainer) {

	bpeTrainer := *wptb.bpeTrainerBuilder.Build()
	return WordPieceTrainer{bpeTrainer: bpeTrainer}

}

type WordPieceTrainer struct {
	bpeTrainer bpe.BpeTrainer
}

func (wpt WordPieceTrainer) Builder() (retVal WordPieceTrainerBuilder) {

	return NewWordPieceTrainerBuilder()
}

// Implement Trainer interface for WordPieceTrainer:
// =================================================

func (wpt WordPieceTrainer) Train(wordCounts map[string]uint32) (retVal tokenizer.Model) {

	bpeModel, _ := wpt.bpeTrainer.Train(wordCounts)

	return NewWordPieceFromBPE(bpeModel.(bpe.BPE))
}

func (wpt WordPieceTrainer) ProcessTokens(words map[string]uint32, tokens []string) {
	wpt.bpeTrainer.ProcessTokens(words, tokens)
}

func (wpt WordPieceTrainer) WithProgressBar() (retVal bool) {
	return wpt.bpeTrainer.WithProgressBar()
}
