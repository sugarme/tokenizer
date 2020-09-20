package wordpiece

import (
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/model/bpe"
)

// WordPieceTrainerBuilder can be used to create a `WordPieceTrainer` with a custom
// configuration.
type WordPieceTrainerBuilder struct {
	bpeTrainerBuilder bpe.BpeTrainerBuilder
}

// NewWordPieceTrainerBuilder create a new WordPieceTrainerBuilder
func NewWordPieceTrainerBuilder() (retVal WordPieceTrainerBuilder) {
	bpeTrainerBuilder := *bpe.NewBPETrainerBuilder()
	bpeTrainerBuilder.ContinuingSubwordPrefix("##")

	return WordPieceTrainerBuilder{
		bpeTrainerBuilder: bpeTrainerBuilder,
	}
}

// MinFrequency set the frequency threshold for the trainer
func (wptb WordPieceTrainerBuilder) MinFrequency(frequency int) (retVal WordPieceTrainerBuilder) {
	wptb.bpeTrainerBuilder.MinFrequency(frequency)

	return wptb
}

// VocabSize set the vocabulary size
func (wptb WordPieceTrainerBuilder) VocabSize(size int) (retVal WordPieceTrainerBuilder) {

	wptb.bpeTrainerBuilder.VocabSize(size)
	return wptb
}

// ShowProgress set whether to show progress
func (wptb WordPieceTrainerBuilder) ShowProgress(show bool) (retVal WordPieceTrainerBuilder) {

	wptb.bpeTrainerBuilder.ShowProgress(show)
	return wptb
}

// SpecialTokens set the special tokens
func (wptb WordPieceTrainerBuilder) SpecialTokens(tokens []tokenizer.AddedToken) (retVal WordPieceTrainerBuilder) {

	wptb.bpeTrainerBuilder.SpecialTokens(tokens)
	return wptb
}

// LimitAlphabet set whether to limit the alphabet
func (wptb WordPieceTrainerBuilder) LimitAlphabet(limit int) (retVal WordPieceTrainerBuilder) {

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

// WordPieceTrainer is a trainer for WordPiece model
type WordPieceTrainer struct {
	bpeTrainer bpe.BpeTrainer
}

// Builder creates WordPieceTrainerBuilder
func (wpt WordPieceTrainer) Builder() (retVal WordPieceTrainerBuilder) {
	return NewWordPieceTrainerBuilder()
}

// Implement Trainer interface for WordPieceTrainer:
// =================================================

func (wpt WordPieceTrainer) Train(wordCounts map[string]int) (retVal tokenizer.Model) {

	bpeModel, _ := wpt.bpeTrainer.Train(wordCounts)

	return NewWordPieceFromBPE(bpeModel.(bpe.BPE))
}

func (wpt WordPieceTrainer) ProcessTokens(words map[string]int, tokens []string) {
	wpt.bpeTrainer.ProcessTokens(words, tokens)
}

func (wpt WordPieceTrainer) WithProgressBar() (retVal bool) {
	return wpt.bpeTrainer.WithProgressBar()
}
