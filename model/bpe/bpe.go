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
	"strings"

	"github.com/sugarme/tokenizer/model"
	"github.com/sugarme/tokenizer/tokenizer"
	"github.com/sugarme/tokenizer/util"
)

type Merges map[Pair]PairVal

type ConfigFiles struct {
	Vocab  string
	Merges string
}

type Config struct {
	Files                   *ConfigFiles
	Vocab                   *model.Vocab
	Merges                  *Merges
	CacheCapacity           uint
	Dropout                 *float32
	UnkToken                *string
	ContinuingSubwordPrefix *string
	EndOfWordSuffix         *string
}

// BpeBuilder can be used to create a `BPE` model with
// a custom configuration
type BpeBuilder struct {
	Config Config
}

func NewBpeBuilder() *BpeBuilder {
	var (
		vocab  *model.Vocab = new(model.Vocab)
		merges *Merges      = new(Merges)
	)
	return &BpeBuilder{
		Config: Config{
			Files:                   nil,
			Vocab:                   vocab,
			Merges:                  merges,
			CacheCapacity:           DefaultCacheCapacity,
			Dropout:                 nil,
			UnkToken:                nil,
			ContinuingSubwordPrefix: nil,
			EndOfWordSuffix:         nil,
		},
	}
}

// Files sets input files for the model
func (bb *BpeBuilder) Files(vocab string, merges string) {
	bb.Config.Files = &ConfigFiles{vocab, merges}
}

// VocabAndMerges sets vocab and merges
func (bb *BpeBuilder) VocabAndMerges(vocab model.Vocab, merges Merges) {
	bb.Config.Vocab = &vocab
	bb.Config.Merges = &merges
}

// CacheCapacity sets the cache capacity. Disable cache by setting it to 0
func (bb *BpeBuilder) CacheCapacity(capacity uint) {
	bb.Config.CacheCapacity = capacity
}

// Dropout set dropout for model
// Ref. https://arxiv.org/abs/1910.13267
func (bb *BpeBuilder) Dropout(dropout float32) {
	bb.Config.Dropout = &dropout
}

// UnkToken set the `UNK` token for the vocab
func (bb *BpeBuilder) UnkToken(unkTok string) {
	bb.Config.UnkToken = &unkTok
}

// ContinuingSubword set the `continuingSubwordPrefix` option.
func (bb *BpeBuilder) ContinuingSubwordPrefix(continuingSubwordPrefix string) {
	bb.Config.ContinuingSubwordPrefix = &continuingSubwordPrefix
}

// EndOfWordSuffix set the `endOfWordSuffix` option.
func (bb *BpeBuilder) EndOfWordSuffix(endOfWordSuffix string) {
	bb.Config.EndOfWordSuffix = &endOfWordSuffix
}

// Build returns a `BPE` model that uses the BpeBuilder configuration
func (bb *BpeBuilder) Build() (*BPE, error) {
	var (
		err    error
		vocab  *model.Vocab
		merges *Merges
		vocabR model.VocabR = make(map[uint32]string)
		cache  *Cache
		bpe    BPE
	)

	vocab = bb.Config.Vocab
	merges = bb.Config.Merges

	// validate dropout
	if bb.Config.Dropout != nil {
		var p float32
		p = *bb.Config.Dropout
		if p <= 0.0 || p > 1.0 {
			err = errors.New("Error: Invalid dropout.")
			return nil, err
		}
	}

	// Read files if provided
	if bb.Config.Files != nil {
		vocab, merges, err = bpe.ReadFiles(bb.Config.Files.Vocab, bb.Config.Files.Merges)
		if err != nil {
			return nil, err
		}

		bb.Config.Vocab = vocab
		bb.Config.Merges = merges
	}

	for k, v := range *vocab {
		vocabR[v] = k
	}

	if bb.Config.CacheCapacity != 0 {
		cache = NewCache(bb.Config.CacheCapacity)
	} else {
		cache = nil
	}

	bpe = BPE{
		Vocab:                   vocab,
		VocabR:                  &vocabR,
		Merges:                  merges,
		Cache:                   cache,
		Dropout:                 bb.Config.Dropout,
		UnkToken:                bb.Config.UnkToken,
		ContinuingSubwordPrefix: bb.Config.ContinuingSubwordPrefix,
		EndOfWordSuffix:         bb.Config.EndOfWordSuffix,
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
	// TODO: handling error
	b, _ = b.builder().Build()
}

// `Clone` can't be derive because it's not implemented for `Cache`.
// To keep things simple when we clone, the new BPE will start with a fresh cache.
func (b *BPE) clone() {
	newBpe := b
	newBpe.Cache.Fresh()
	b = newBpe
}

// NewBPE create a default BPE from sratch using its pbeBuilder
func NewBPE() (*BPE, error) {
	b := NewBpeBuilder()
	return b.Build()
}

// NewBpeFromFiles create BPE model from vocab and merges files
func NewBpeFromFiles(vocab, merges string) (*BPE, error) {
	b := NewBpeBuilder()
	b.Files(vocab, merges)
	return b.Build()
}

// New creates new BPE model with given vocab and merges
func (b *BPE) New(vocab model.Vocab, merges Merges) {
	b.new()
	b.Vocab = &vocab
	b.Merges = &merges
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
	// Recall: Merges is map[Pair]PairVal (rank uint32, newId uint32)
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

		// pairVal := PairVal{uint32(newTokenInt), newId}
		pairVal := PairVal{uint32(lineNum), newId}

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
func (b *BPE) GetVocab() *model.Vocab {
	return b.Vocab
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

	chars := strings.Split(w, "")
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

	for i, c := range chars {
		var s string = c
		// Add `continuingSubwordPrefix` if relevant
		if i > 0 && i < len(chars) {
			s = fmt.Sprintf("%v%v", prefix, s)
		} else if i == len(chars) { // last `char`
			s = fmt.Sprintf("%v%v", s, suffix)
		}

		// Look its id up

		vocab := *b.Vocab
		// if id, ok := (*b.Vocab)[s]; ok { // found
		if id, ok := vocab[s]; ok { // found
			word.Add(id)
		} else { // not found, add `unk`
			if b.UnkToken != nil {
				// get `unk` id
				unkId := (*b.Vocab)[*b.UnkToken]
				// add `unk`
				word.Add(unkId)
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
func (b *BPE) WordToTokens(word Word, initialOffsets tokenizer.Offsets) ([]tokenizer.Token, error) {

	var tokens []tokenizer.Token
	chars := word.GetChars()
	offsets := word.GetOffsets()
	type zword struct {
		Id      uint32
		Offsets tokenizer.Offsets
	}
	var zWord []zword

	// err := util.Zip(chars, offsets, &zWord)
	// if err != nil {
	// return nil, err
	// }

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

	return tokens, nil
}

// Implement Model interface for BPE
// Model has the following methods:
// 1. Tokenize(tokens []PreToken) ([]Token, error)
// 2. TokenToId(token string) uint32
// 3. IdToToken(id uint32) string
// 4. GetVocabSize() uint
// 5. Save(path string, name string) error

// Tokenize tokenizes sentences into tokens
// NOTE: sentence is []PreToken struct{Value string, Offsets Offsets}
func (b *BPE) Tokenize(sentence []tokenizer.PreToken) ([]tokenizer.Token, error) {

	if len(sentence) == 0 {
		return []tokenizer.Token{}, nil
	}

	var encoded []tokenizer.Token
	var cachedWords []Word
	var keys []string

	// if using dropout, we don't use the cache
	if b.Dropout != nil {
		cachedWords = nil
	} else {
		// for _, k := range sentence {
		// keys = append(keys, k.Value)
		// }
		// cachedWords = b.Cache.GetValues(keys)
	}

	var shouldUpdateCache bool = false

	for i, preTok := range sentence {
		var (
			tokens []tokenizer.Token
			err    error
		)

		// not using cache as we're using dropout
		if cachedWords == nil {
			word := b.MergeWord(preTok.Value)
			tokens, err = b.WordToTokens(*word, preTok.Offsets)
			if err != nil {
				return nil, err
			}
		} else {
			if i > len(cachedWords) {
				// no cache hit, let's recompute merges
				word := b.MergeWord(preTok.Value)
				tokens, err = b.WordToTokens(*word, preTok.Offsets)
				// Add to cache
				cachedWords[i] = *word
				shouldUpdateCache = true
			} else {
				word := cachedWords[i]
				tokens, err = b.WordToTokens(word, preTok.Offsets)
				if err != nil {
					return nil, err
				}
				// Remove this entry so we don't needlessly try to update it
				// in the cache below.
				cachedWords, err = deleteWord(cachedWords, i)
				if err != nil {
					return nil, err
				}

			}

		}

		encoded = append(encoded, tokens...)

		// Updae the cache if we need
		if cachedWords != nil {
			if shouldUpdateCache {
				var cachedItems []CacheItem

				err = util.Zip(keys, cachedWords, cachedItems)
				if err != nil {
					return nil, err
				}

				b.Cache.SetValues(cachedItems)
			}
		}

	}

	return encoded, nil
}

func (b *BPE) TokenToId(token string) uint32 {
	return (*b.Vocab)[token]
}

func (b *BPE) IdToToken(id uint32) string {
	return (*b.VocabR)[id]
}

func (b *BPE) GetVocabSize() uint {
	vocabLen := len(*b.Vocab)
	return uint(vocabLen)
}

func (b *BPE) Save(dir string, nameOpt ...string) error {
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
		Rank uint32
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
		c1 := b.IdToToken(p.Pair.C1)
		c2 := b.IdToToken(p.Pair.C2)
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
