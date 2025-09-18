package wordpiece

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/model"
	"github.com/sugarme/tokenizer/model/bpe"
	"github.com/sugarme/tokenizer/util"
)

type config struct {
	files                   string
	vocab                   *model.Vocab
	unkToken                string
	continuingSubwordPrefix string
	maxInputCharsPerWord    int
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
func (wpb WordPieceBuilder) MaxInputCharsPerWord(maxInputCharsPerWord int) (retVal WordPieceBuilder) {
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

	// update `unk_token`
	// if _, ok := vocab[wpb.config.unkToken]; !ok {
	// vocab[wpb.config.unkToken] = len(vocab)
	// }

	var vocabR model.VocabR = make(map[int]string)
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
	maxInputCharsPerWord  int
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
		idx   int = 0
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

// NewWordPieceFromFile initializes a WordPiece model from a mapping file
func NewWordPieceFromFile(vocabFile string, unkToken string, maxInputCharsPerWordOpt ...int) (retVal WordPiece, err error) {

	filePath, err := filepath.Abs(vocabFile)
	if err != nil {
		return retVal, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return retVal, err
	}
	defer file.Close()

	var (
		vocab model.Vocab = make(map[string]int)
		line  string
		idx   int = 0
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		vocab[line] = idx
		idx += 1
	}

	if err := scanner.Err(); err != nil {
		return retVal, err
	}

	wp := NewWordPiece()
	builder := wp.Builder().Vocab(&vocab).UnkToken(unkToken)
	if len(maxInputCharsPerWordOpt) > 0 {
		maxInputCharsPerWord := maxInputCharsPerWordOpt[0]
		builder = builder.MaxInputCharsPerWord(maxInputCharsPerWord)
	}

	return builder.Build(), nil
}

// WordPieceBuilderFromBPE create a WordPieceBuilder from BPE model
func NewWordPieceFromBPE(bpe bpe.BPE) (retVal WordPiece) {
	wpb := NewWordPieceBuilder()
	var vocab model.Vocab = bpe.GetVocab()
	wpb.config.vocab = &vocab

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

// Implement Model interface for WordPiece:
// ========================================

func (wp WordPiece) GetVocab() (retVal map[string]int) {
	return *wp.vocab
}

func (wp WordPiece) GetVocabSize() (retVal int) {
	return len(*wp.vocab)
}

func (wp WordPiece) Tokenize(sequence string) (retVal []tokenizer.Token, err error) {

	// fmt.Printf("input sequence: %v\n", sequence)

	var outputTokens []tokenizer.Token

	chars := []rune(sequence)
	charLen := len(chars)
	if charLen > wp.maxInputCharsPerWord {
		id, ok := (*wp.vocab)[wp.unkToken]
		if !ok {
			err := fmt.Errorf("WordPiece error: Missing [UNK] token. Unknown token value %q not found in the vocab\n", wp.unkToken)
			return retVal, err
		}
		token := tokenizer.Token{
			Value:   wp.unkToken,
			Id:      id,
			Offsets: []int{0, charLen},
		}
		outputTokens = append(outputTokens, token)

		return outputTokens, nil
	}

	var (
		isBad     bool = false
		start     int  = 0
		subTokens []tokenizer.Token
	)

	for start < charLen {
		end := charLen
		var currStr *tokenizer.Token

		for start < end {
			substr := string(chars[start:end])
			if start > 0 {
				substr = fmt.Sprintf("%v%v", wp.continueSubwordPrefix, substr)
			}

			if id, ok := (*wp.vocab)[substr]; ok {
				currStr = &tokenizer.Token{
					Id:      id,
					Value:   substr,
					Offsets: []int{start, end},
				}

				break
			}
			end -= 1
		}
		if currStr == nil {
			isBad = true
			break
		}

		subTokens = append(subTokens, *currStr)
		start = end
	}

	if isBad {
		id, ok := (*wp.vocab)[wp.unkToken]
		if !ok {
			err := fmt.Errorf("WordPiece error: Missing [UNK] token. Unknown token value %q not found in the vocab\n", wp.unkToken)
			return retVal, err
		}
		token := tokenizer.Token{
			Value:   wp.unkToken,
			Id:      id,
			Offsets: []int{0, charLen},
		}

		outputTokens = append(outputTokens, token)
	} else {
		outputTokens = append(outputTokens, subTokens...)
	}

	return outputTokens, nil
}

func (wp WordPiece) TokenToId(token string) (retVal int, ok bool) {
	retVal, ok = (*wp.vocab)[token]
	return
}

func (wp WordPiece) IdToToken(id int) (retVal string, ok bool) {
	retVal, ok = (*wp.vocabR)[id]
	return retVal, ok
}

func (wp WordPiece) Save(dir string, nameOpt ...string) (err error) {
	var vfile string
	if len(nameOpt) > 0 {
		vfile = fmt.Sprintf("%v/%v-vocab.txt", dir, nameOpt[0])
	} else {
		vfile = fmt.Sprintf("%v/vocab.txt", dir)
	}

	// make filepath
	err = makeFilePath(vfile)
	if err != nil {
		return err
	}

	// Write vocab.txt
	// each line is a pair separated by a space
	var lines []string
	vocab := *wp.vocab

	// sort vocab by map's value (int)
	type kv struct {
		Key   string
		Value int
	}
	var sVocab []kv
	for k, v := range vocab {
		sVocab = append(sVocab, kv{k, v})
	}

	// sort sVocab by `Rank` field in-place
	sort.Slice(sVocab, func(i, j int) bool {
		return sVocab[i].Value < sVocab[j].Value
	})

	// Create vocab lines
	for _, item := range sVocab {
		line := fmt.Sprintf("%v", item.Key)
		lines = append(lines, line)
	}

	// write to file
	file, err := os.Create(vfile)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()

}

// makeFilePath creates a filePath. If dir not existing, create it
func makeFilePath(filename string) error {
	var err error
	dirName := filepath.Dir(filename)
	if _, err = os.Stat(dirName); err != nil {
		return err
	}
	return os.MkdirAll(dirName, os.ModePerm)
}

// New creates WordPiece model from input data.
func New(
	vocab model.Vocab,
	opts *util.Params,
) (*WordPiece, error) {
	// Default values:
	unkToken := "[UNK]"
	continuingSubwordPrefix := "##"
	maxInputCharsPerWord := 100
	if opts.Has("unk_token") {
		unkToken = opts.Get("unk_token").(string)
	}
	if opts.Has("continuingSubwordPrefix") {
		continuingSubwordPrefix = opts.Get("continuing_subword_prefix").(string)
	}
	if opts.Has("max_input_chars_per_word") {
		maxInputCharsPerWord = opts.Get("max_input_chars_per_word").(int)
	}

	builder := WordPieceBuilder{
		config: config{
			files:                   "",
			vocab:                   &vocab,
			unkToken:                unkToken,
			continuingSubwordPrefix: continuingSubwordPrefix,
			maxInputCharsPerWord:    maxInputCharsPerWord,
		},
	}

	m := builder.Build()

	return &m, nil
}
