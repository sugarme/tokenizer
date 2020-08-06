package wordpiece

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sugarme/tokenizer/model"
	"github.com/sugarme/tokenizer/model/bpe"
	"github.com/sugarme/tokenizer/tokenizer"
)

var MissingUnkToken = fmt.Errorf("WordPiece error: Missing [UNK] token from the vocabulary\n")

type config struct {
	files                   string
	vocab                   *model.Vocab
	unkToken                string
	continuingSubwordPrefix string
	maxInputCharsPerWord    uint
}

// WordPieceBuilder can be used to create a WordPiece model with a custom
// configuration
type WordPieceBuilder struct {
	config config
}

func NewWordPieceBuilder() (retVal WordPieceBuilder) {

	return WordPieceBuilder{
		config: config{
			files:                   "",
			vocab:                   new(model.Vocab),
			unkToken:                "[UNK]",
			continuingSubwordPrefix: "##",
			maxInputCharsPerWord:    100,
		},
	}
}

// Files sets the input files
func (wpb WordPieceBuilder) Files(vocab string) (retVal WordPieceBuilder) {
	wpb.config.files = vocab

	return wpb
}

// Vocab set the vocab (token -> ID) mapping.
func (wpb WordPieceBuilder) Vocab(vocab *model.Vocab) (retVal WordPieceBuilder) {
	wpb.config.vocab = vocab

	return wpb
}

// UnkToken set the `UNK` token for the vocab.
func (wpb WordPieceBuilder) UnkToken(unkToken string) (retVal WordPieceBuilder) {
	wpb.config.unkToken = unkToken

	return wpb
}

// ContinueSubwordPrefix set the prefix for continuing subwords.
func (wpb WordPieceBuilder) ContinuingSubwordPrefix(continueSubwordPrefix string) (retVal WordPieceBuilder) {
	wpb.config.continuingSubwordPrefix = continueSubwordPrefix

	return wpb
}

// Set the maximum number of input characters per word.
func (wpb WordPieceBuilder) MaxInputCharsPerWord(maxInputCharsPerWord uint) (retVal WordPieceBuilder) {
	wpb.config.maxInputCharsPerWord = maxInputCharsPerWord

	return wpb
}

// Build contructs a `WordPiece` model that uses the `WordPieceBuilder`'s configuration.
func (wpb WordPieceBuilder) Build() (retVal WordPiece) {

	var wp = NewWordPiece()

	files := wpb.config.files

	if files != "" {
		vocab := wp.ReadFiles(files)
		wpb.config.vocab = &vocab
	}

	vocab := *wpb.config.vocab
	var vocabR model.VocabR
	for k, v := range vocab {
		vocabR[v] = k
	}

	return WordPiece{
		vocab:                 &vocab,
		vocabR:                &vocabR,
		unkToken:              wpb.config.unkToken,
		continueSubwordPrefix: wpb.config.continuingSubwordPrefix,
		maxInputCharsPerWord:  wpb.config.maxInputCharsPerWord,
	}
}

// WordPiece model:
// ================

// WordPiece is a WordPiece model
// Ref.https://static.googleusercontent.com/media/research.google.com/en//pubs/archive/37842.pdf
type WordPiece struct {
	vocab                 *model.Vocab
	vocabR                *model.VocabR
	unkToken              string
	continueSubwordPrefix string
	maxInputCharsPerWord  uint
}

// NewWordPiece initiates a new WordPiece with default values.
func NewWordPiece() (retVal WordPiece) {
	return WordPiece{
		vocab:                 new(model.Vocab),
		vocabR:                new(model.VocabR),
		unkToken:              "[UNK]",
		continueSubwordPrefix: "##",
		maxInputCharsPerWord:  100,
	}
}

// Builder gets a WordPieceBuilder
func (wp WordPiece) Builder() (retVal WordPieceBuilder) {
	return NewWordPieceBuilder()
}

// ReadFiles reads the given file to extract the vocab
func (wp WordPiece) ReadFiles(filename string) (retVal model.Vocab) {
	filePath, err := filepath.Abs(filename)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var (
		vocab model.Vocab
		line  string
		idx   uint32 = 0
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		vocab[line] = idx
		idx += 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return vocab
}

// NewWordPieceBuilderFromFile initializes a WordPieceBuilder from a vocab mapping file
func NewWordPieceBuilderFromFile(filename string) (retVal WordPieceBuilder) {
	wp := NewWordPiece()
	wpb := wp.Builder()

	return wpb.Files(filename)
}

// WordPieceBuilderFromBPE create a WordPieceBuilder from BPE model
func NewWordPieceFromBPE(bpe bpe.BPE) (retVal WordPiece) {
	wpb := NewWordPieceBuilder()
	wpb.config.vocab = bpe.GetVocab()

	wp := wpb.Build()
	unk := bpe.GetUnkToken()
	if unk != nil {
		wp.unkToken = *unk
	}

	continueSubwordPrefix := bpe.GetContinuingSubwordPrfix()
	if continueSubwordPrefix != nil {
		wp.continueSubwordPrefix = *continueSubwordPrefix
	}

	return wp
}

// Implement Model for WordPiece:
// ==============================
func (wp WordPiece) GetVocab() (retVal *model.Vocab) {
	return wp.vocab
}

func (wp WordPiece) VocabSize() (retVal int) {
	return len(*wp.vocab)
}

func (wp WordPiece) Tokenize(sentence tokenizer.PreToken) (retVal []tokenizer.Token) {

	// TODO. continue

	return
}
