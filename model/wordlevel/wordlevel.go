package wordlevel

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/sugarme/tokenizer"
)

type config struct {
	vocab    map[string]int
	unkToken string
}

// WordLevelBuilder is a builder for WordLevel model
type WordLevelBuilder struct {
	config config
}

// defaultWordLevelBuilder creates a WordLevelBuilder with default values
func defaultWordLevelBuilder() *WordLevelBuilder {
	return &WordLevelBuilder{
		config: config{
			vocab:    make(map[string]int),
			unkToken: "<unk>",
		},
	}
}

// NewWordLevelBuilder creates a WordLevelBuilder with default values
func NewWordLevelBuilder() *WordLevelBuilder {
	return defaultWordLevelBuilder()
}

// Vocab set the vocab (token -> id) mapping
func (wlb *WordLevelBuilder) Vocab(vocab map[string]int) {
	wlb.config.vocab = vocab
}

// UnkToken set `UNK` token for the vocab
func (wlb *WordLevelBuilder) UnkToken(unkToken string) {
	wlb.config.unkToken = unkToken
}

// Build builds a WordLevel using configuration
func (wlb *WordLevelBuilder) Build() *WordLevel {
	var vocabR map[int]string = make(map[int]string)
	for k, v := range wlb.config.vocab {
		vocabR[v] = k
	}

	return &WordLevel{
		vocab:    wlb.config.vocab,
		vocabR:   vocabR,
		unkToken: wlb.config.unkToken,
	}
}

// WordLevel is a model for building WordLevel tokenizer
type WordLevel struct {
	vocab    map[string]int
	vocabR   map[int]string
	unkToken string
}

// NewWordLevelFromFile initializes a WordLevel from file
func NewWorldLevelFromFile(vocabFile string, unkToken string) (*WordLevel, error) {
	filePath, err := filepath.Abs(vocabFile)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var (
		vocab map[string]int = make(map[string]int)
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
		return nil, err
	}

	wlb := NewWordLevelBuilder()
	wlb.config.vocab = vocab
	wlb.config.unkToken = unkToken

	return wlb.Build(), nil
}

// NewWordLevel initiates a new WordLevel
func NewWordLevel() *WordLevel {
	return &WordLevel{
		vocab:    make(map[string]int),
		vocabR:   make(map[int]string),
		unkToken: "<unk>",
	}
}

// Implement Model interface for WordLevel
// =======================================

// Tokenize transforms given input to token
func (wl WordLevel) Tokenize(token string) []tokenizer.Token {

	var output []tokenizer.Token
	var (
		id        int
		ok, unkOk bool
	)

	id, ok = wl.vocab[token]
	if !ok {
		id, unkOk = wl.vocab[wl.unkToken]
		if !unkOk {
			log.Fatalf("Missing 'unk' token in vocab.\n")
		}
	}

	output = append(output, tokenizer.Token{
		Id:      id,
		Value:   token,
		Offsets: []int{0, len(token)},
	})

	return output
}

// TokenToId returns id of a given token if existing
func (wl WordLevel) TokenToId(token string) (int, bool) {
	id, ok := wl.vocab[token]
	return id, ok
}

// IdToToken gets token of given id if existing
func (wl WordLevel) IdToToken(id int) (string, bool) {
	tok, ok := wl.vocabR[id]
	return tok, ok
}

// Vocab returns a map of model vocab
func (wl WordLevel) Vocab() map[string]int {
	return wl.vocab
}

// Save saves vocab to a file
func (wl WordLevel) Save(dir string, nameOpt ...string) (err error) {
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
	vocab := wl.vocab

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
