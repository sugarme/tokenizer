package bpe

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	// "strconv"
	"log"
	"strings"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/model"
	"github.com/sugarme/tokenizer/util"
)

type Merges map[Pair]PairVal

type configFiles struct {
	vocab  string
	merges string
}

type Config struct {
	files                   *configFiles
	vocab                   *model.Vocab
	merges                  *Merges
	cacheCapacity           int
	dropout                 *float32
	unkToken                *string
	continuingSubwordPrefix *string
	endOfWordSuffix         *string
}

// BpeBuilder can be used to create a `BPE` model with
// a custom configuration
type BpeBuilder struct {
	config Config
}

func NewBpeBuilder() *BpeBuilder {
	var (
		vocab  model.Vocab = make(map[string]int)
		merges Merges      = make(map[Pair]PairVal)
	)
	return &BpeBuilder{
		config: Config{
			files:                   nil,
			vocab:                   &vocab,
			merges:                  &merges,
			cacheCapacity:           DefaultCacheCapacity,
			dropout:                 nil,
			unkToken:                nil,
			continuingSubwordPrefix: nil,
			endOfWordSuffix:         nil,
		},
	}
}

// Files sets input files for the model
func (bb *BpeBuilder) Files(vocab string, merges string) {
	bb.config.files = &configFiles{vocab, merges}
}

// VocabAndMerges sets vocab and merges
func (bb *BpeBuilder) VocabAndMerges(vocab model.Vocab, merges Merges) {
	bb.config.vocab = &vocab
	bb.config.merges = &merges
}

// CacheCapacity sets the cache capacity. Disable cache by setting it to 0
func (bb *BpeBuilder) CacheCapacity(capacity int) {
	bb.config.cacheCapacity = capacity
}

// Dropout set dropout for model
// Ref. https://arxiv.org/abs/1910.13267
func (bb *BpeBuilder) Dropout(dropout float32) {
	bb.config.dropout = &dropout
}

// UnkToken set the `UNK` token for the vocab
func (bb *BpeBuilder) UnkToken(unkTok string) {
	bb.config.unkToken = &unkTok
}

// ContinuingSubword set the `continuingSubwordPrefix` option.
func (bb *BpeBuilder) ContinuingSubwordPrefix(continuingSubwordPrefix string) {
	bb.config.continuingSubwordPrefix = &continuingSubwordPrefix
}

// EndOfWordSuffix set the `endOfWordSuffix` option.
func (bb *BpeBuilder) EndOfWordSuffix(endOfWordSuffix string) {
	bb.config.endOfWordSuffix = &endOfWordSuffix
}

// Build returns a `BPE` model that uses the BpeBuilder configuration
func (bb *BpeBuilder) Build() (*BPE, error) {
	var (
		err    error
		vocab  *model.Vocab
		merges *Merges
		vocabR model.VocabR = make(map[int]string)
		cache  *Cache
		bpe    BPE
	)

	vocab = bb.config.vocab
	merges = bb.config.merges

	// validate dropout
	if bb.config.dropout != nil {
		var p float32
		p = *bb.config.dropout
		if p <= 0.0 || p > 1.0 {
			err = errors.New("Error: Invalid dropout.")
			return nil, err
		}
	}

	// Read files if provided
	if bb.config.files != nil {
		vocab, merges, err = bpe.ReadFiles(bb.config.files.vocab, bb.config.files.merges)
		if err != nil {
			return nil, err
		}

		bb.config.vocab = vocab
		bb.config.merges = merges
	}

	for k, v := range *vocab {
		vocabR[v] = k
	}

	if bb.config.cacheCapacity != 0 {
		cache = NewCache(bb.config.cacheCapacity)
	} else {
		cache = nil
	}

	bpe = BPE{
		Vocab:                   vocab,
		VocabR:                  &vocabR,
		Merges:                  merges,
		Cache:                   cache,
		Dropout:                 bb.config.dropout,
		UnkToken:                bb.config.unkToken,
		ContinuingSubwordPrefix: bb.config.continuingSubwordPrefix,
		EndOfWordSuffix:         bb.config.endOfWordSuffix,
	}

	return &bpe, nil

}

// BPE is a struct for byte pair encoding model
// Ref. https://www.aclweb.org/anthology/P16-1162/
type BPE struct {
	// Vocab is the vocabulary assigns a number to each token.
	Vocab *model.Vocab

	// VocabR is Reversed vocabulary, to rebuild sentences.
	VocabR *model.VocabR

	// Merges contains the mapping between Pairs and their (rank, newId).
	Merges *Merges

	// Cache contains the cache for optimizing the encoding step.
	// It is a `map[string]Word`
	Cache *Cache

	// Dropout probability for merges.
	// 0 = no dropout is the default.
	// At 1.0, tokenization will perform no merges, so the result will just be characters.
	Dropout *float32

	// UnkToken is the unknown token to be used when we encounter an unknown char
	UnkToken *string

	// ContinuingSubwordPrefix is an optional prefix
	// to use on any subword that exist only behind another one
	ContinuingSubwordPrefix *string

	// EndOfWordSuffix is an optional suffix
	// to caracterize and end-of-word subword
	EndOfWordSuffix *string
}

func (b *BPE) builder() *BpeBuilder {
	return NewBpeBuilder()
}

// new create a BPE with default values
func (b *BPE) new() {
	var err error
	b, err = b.builder().Build()
	if err != nil {
		log.Fatal(err)
	}
}

// `Clone` can't be derive because it's not implemented for `Cache`.
// To keep things simple when we clone, the new BPE will start with a fresh cache.
func (b *BPE) clone() {
	newBpe := b
	newBpe.Cache.Fresh()
	b = newBpe
}

// newBPE create a default BPE from sratch using its pbeBuilder
func newBPE() (*BPE, error) {
	b := NewBpeBuilder()
	return b.Build()
}

func DefaultBPE() (*BPE, error) {
	return newBPE()
}

// NewBpeFromFiles create BPE model from vocab and merges files
func NewBpeFromFiles(vocab, merges string) (*BPE, error) {
	b := NewBpeBuilder()
	b.Files(vocab, merges)
	return b.Build()
}

// NewBPE creates new BPE model with given vocab and merges
func NewBPE(vocab model.Vocab, merges Merges) *BPE {
	b, err := newBPE()
	if err != nil {
		log.Fatal(err)
	}
	b.Vocab = &vocab

	var vocabR model.VocabR = make(map[int]string)
	for k, v := range vocab {
		vocabR[v] = k
	}

	b.VocabR = &vocabR

	b.Merges = &merges
	return b
}

// FromFile creates `BpeBuilder` from vocab and merges files.
func (b *BPE) FromFiles(vocab string, merges string) *BpeBuilder {
	builder := b.builder()
	builder.Files(vocab, merges)
	return builder
}

// ReadFiles read the given files to extract vocab and merges
func (b *BPE) ReadFiles(vocabF string, mergesF string) (*model.Vocab, *Merges, error) {
	var err error
	// read json file
	vocabBytes, err := ioutil.ReadFile(vocabF)
	if err != nil {
		return nil, nil, err
	}

	var (
		vocab  model.Vocab
		merges Merges = make(map[Pair]PairVal)
	)

	err = json.Unmarshal(vocabBytes, &vocab)
	if err != nil {
		return nil, nil, err
	}

	// Read merges file. Each line contains a Merges object(rank, )
	// Recall: Merges is map[Pair]PairVal (rank int, newId int)
	mFile, err := os.Open(mergesF)
	if err != nil {
		return nil, nil, err
	}
	defer mFile.Close()

	s := bufio.NewScanner(mFile)

	// `s.Scan()` advance scaning and return `false` if
	// end of file or hit any error. The error will be
	// access by s.Err. If error caused by EOF it's value is nil.
	var lineNum = 0
	for s.Scan() {
		line := s.Text()

		// Skip line with `#version`
		re := regexp.MustCompile(`#version`)
		if re.MatchString(line) {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			err = fmt.Errorf("Read merge file error: invalid data at line %d\n", lineNum)
			return nil, nil, err
		}

		a, ok := vocab[parts[0]]
		if !ok {
			// err = fmt.Errorf("Read merge file error: part a value for '%s' key not found.", parts[0])
			continue
			// return nil, nil, err
		}

		b, ok := vocab[parts[1]]
		if !ok {
			// err = fmt.Errorf("Read merge file error: part b value for '%s' key not found.", parts[1])
			continue
			// return nil, nil, err
		}

		pair := Pair{a, b}
		// newToken := fmt.Sprintf("%v%v", parts[0], parts[1])
		newToken := fmt.Sprintf("%v%v", parts[0], parts[1])
		newId, ok := vocab[newToken]
		if !ok {
			err = fmt.Errorf("Read merge file error: key value for token: \"%s\" not found.", newToken)
			return nil, nil, err
		}

		// newTokenInt, err := strconv.ParseInt(newToken, 10, 64)

		err = util.TraceError(err)
		if err != nil {
			return nil, nil, err
		}

		pairVal := PairVal{lineNum, newId}

		merges[pair] = pairVal

		lineNum += 1
	}

	if s.Err() != nil {
		return nil, nil, s.Err()
	}

	return &vocab, &merges, nil
}

// ClearCache reset the cache
func (b *BPE) ClearCache() {
	if b.Cache != nil {
		b.Cache.Clear()
	}
}

// GetVocab returns BPE vocab
// func (b *BPE) GetVocab() *model.Vocab {
func (b BPE) GetVocab() map[string]int {
	return *b.Vocab
}

// GetUnkToken returns `unk` token
func (b *BPE) GetUnkToken() *string {
	return b.UnkToken
}

// GetContinuingSubwordPrefix returns continuing subword prefix
func (b *BPE) GetContinuingSubwordPrfix() *string {
	return b.ContinuingSubwordPrefix
}

// MergeWord merges given word
func (b *BPE) MergeWord(w string) *Word {

	word := NewWord()
	var (
		prefix, suffix string
	)

	if b.ContinuingSubwordPrefix != nil {
		prefix = *b.ContinuingSubwordPrefix
	} else {
		prefix = ""
	}

	if b.EndOfWordSuffix != nil {
		suffix = *b.EndOfWordSuffix
	} else {
		suffix = ""
	}

	chars := []rune(w)
	currRuneIdx := 0
	for byteIdx, r := range w {
		var (
			s       string
			byteLen int
		)
		byteLen = len(string(r))

		// if first rune, add prefix
		if byteIdx == 0 {
			s = fmt.Sprintf("%v%v", prefix, string(r))
		} else if currRuneIdx == len(chars) { // last rune, add suffix
			s = fmt.Sprintf("%v%v", string(r), suffix)
		} else { // the rest
			s = string(r)
		}
		currRuneIdx++

		// If `s` exists in vocab, add its id, otherwise add id of `unk`
		vocab := *b.Vocab
		if id, ok := vocab[s]; ok { // found
			word.Add(id, byteLen)
		} else { // not found, add `unk`
			if b.UnkToken != nil {
				// get `unk` id
				unkId := (*b.Vocab)[*b.UnkToken]
				// add `unk`
				word.Add(unkId, byteLen)
			} else {
				fmt.Printf("cannot find '%s' in the vocab. \n", s)
				panic("Can't find `unk` token in the vocab. Have you added one when initiating the model?")
			}
		}
	}

	if b.Dropout != nil {
		word.MergeAll(*b.Merges, *b.Dropout)
	} else {
		word.MergeAll(*b.Merges)
	}

	return word
}

// WordToTokens slices word to tokens
func (b *BPE) WordToTokens(word Word) []tokenizer.Token {
	var tokens []tokenizer.Token
	chars := word.GetChars()
	offsets := word.GetOffsets()
	type zword struct { // zip id and offsets
		Id      int
		Offsets []int
	}
	var zWord []zword

	for i, char := range chars {
		zWord = append(zWord, zword{
			Id:      char,
			Offsets: offsets[i],
		})
	}

	for _, z := range zWord {
		tok := tokenizer.Token{
			Id:      z.Id,
			Value:   (*b.VocabR)[z.Id],
			Offsets: z.Offsets,
		}
		tokens = append(tokens, tok)
	}

	return tokens
}

// Implement Model interface for BPE
// Model has the following methods:
// 1. Tokenize(sequence string) ([]Token, error)
// 2. TokenToId(token string) (id int, ok bool)
// 3. IdToToken(id int) (token string, ok bool)
// 4. GetVocab() map[string]int
// 5. GetVocabSize() int
// 6. Save(path string, prefixOpt ...string) error

// Tokenize tokenizes sentences into tokens
// NOTE: sentence is []PreToken struct{Value string, Offsets Offsets}
func (b BPE) Tokenize(sequence string) (retVal []tokenizer.Token, err error) {
	if len(sequence) == 0 {

		return []tokenizer.Token{}, nil
	}

	if b.Dropout == nil {
		return b.TokenizeWithCache(sequence), nil
	}

	word := b.MergeWord(sequence)

	return b.WordToTokens(*word), nil
}

func (b BPE) TokenizeWithCache(sequence string) (retVal []tokenizer.Token) {

	if hit, ok := b.Cache.cmap[sequence]; ok {
		return b.WordToTokens(hit)
	} else {
		word := b.MergeWord(sequence)
		retVal = b.WordToTokens(*word)
		if b.Cache != nil {
			b.Cache.SetValues([]CacheItem{
				{sequence, *word},
			})
		}
		return retVal
	}
}

func (b BPE) TokenToId(token string) (id int, ok bool) {
	id, ok = (*b.Vocab)[token]
	return id, ok
}

func (b BPE) IdToToken(id int) (token string, ok bool) {
	token, ok = (*b.VocabR)[id]
	return token, ok
}

func (b BPE) GetVocabSize() int {
	return len(*b.Vocab)
}

func (b BPE) Save(dir string, nameOpt ...string) error {
	var vfile string
	var mfile string
	var err error
	if len(nameOpt) > 0 {
		vfile = fmt.Sprintf("%v/%v-vocab.json", dir, nameOpt[0])
		mfile = fmt.Sprintf("%v/%v-merges.txt", dir, nameOpt[0])
	} else {
		vfile = fmt.Sprintf("%v/vocab.json", dir)
		mfile = fmt.Sprintf("%v/merges.txt", dir)

	}
	// make filepath
	err = makeFilePath(vfile)
	if err != nil {
		return err
	}

	// Write vocab.json
	var vocabData []byte
	vocabData, err = json.Marshal(b.Vocab)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(vfile, vocabData, os.ModePerm)
	if err != nil {
		return err
	}

	// Write merges.txt
	// each line is a pair separated by a space
	var lines []string
	type pairRank struct {
		Pair Pair
		Rank int
	}
	var pairRanks []pairRank
	for pair, pairVal := range *b.Merges {
		pairRanks = append(pairRanks, pairRank{
			Pair: pair,
			Rank: pairVal.Rank,
		})
	}

	// sort pairRanks by `Rank` field in-place
	sort.Slice(pairRanks, func(i, j int) bool {
		return pairRanks[i].Rank < pairRanks[j].Rank
	})

	// Create lines of merges
	for _, p := range pairRanks {
		// line := fmt.Sprintf("%v %v", p.Pair.C1, p.Pair.C2)
		c1, _ := b.IdToToken(p.Pair.C1)
		c2, _ := b.IdToToken(p.Pair.C2)
		line := fmt.Sprintf("%v %v", c1, c2)
		lines = append(lines, line)
	}

	// write to file
	file, err := os.Create(mfile)
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

func deleteWord(a []Word, i int) ([]Word, error) {
	var err error
	if i < 0 || i > len(a) {
		err = errors.New("`i` index is out of bound.")
		return nil, err
	}

	return append(a[:i], a[i+1:]...), nil
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

func CreateMerges(vocab map[string]int, mergesData []string) (*Merges, error) {
	var (
		lineNum int    = 0
		merges  Merges = make(map[Pair]PairVal)
	)
	for _, line := range mergesData {
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			err := fmt.Errorf("Read merges error: invalid data at line %d\n", lineNum)
			return nil, err
		}

		a, ok := vocab[parts[0]]
		if !ok {
			// err = fmt.Errorf("Read merge file error: part a value for '%s' key not found.", parts[0])
			continue
			// return nil, nil, err
		}

		b, ok := vocab[parts[1]]
		if !ok {
			// err = fmt.Errorf("Read merge file error: part b value for '%s' key not found.", parts[1])
			continue
			// return nil, nil, err
		}

		pair := Pair{a, b}
		// newToken := fmt.Sprintf("%v%v", parts[0], parts[1])
		newToken := fmt.Sprintf("%v%v", parts[0], parts[1])
		newId, ok := vocab[newToken]
		if !ok {
			err := fmt.Errorf("Read merges error: key value for token: \"%s\" not found.", newToken)
			return nil, err
		}

		pairVal := PairVal{lineNum, newId}

		merges[pair] = pairVal

		lineNum += 1
	}

	return &merges, nil
}

// New create new BPE model.
func New(
	// vocab map[string]int,
	vocab model.Vocab,
	mergesData []string,
	dropout *float32,
	unkToken *string,
	continuingSubwordPrefix *string,
	endOfWordSuffix *string,
) (*BPE, error) {
	merges, err := CreateMerges(vocab, mergesData)
	if err != nil {
		return nil, err
	}

	// vc := interface{}(vocab).(model.Vocab)

	builder := &BpeBuilder{
		config: Config{
			files:                   nil,
			vocab:                   &vocab,
			merges:                  merges,
			cacheCapacity:           DefaultCacheCapacity,
			dropout:                 dropout,
			unkToken:                unkToken,
			continuingSubwordPrefix: continuingSubwordPrefix,
			endOfWordSuffix:         endOfWordSuffix,
		},
	}

	return builder.Build()
}
