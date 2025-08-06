package unigram

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/util"
)

// TokenScore represents a token and its score in the Unigram model
type TokenScore struct {
	Token string
	Score float64
}

// Config holds the configuration for the Unigram model
type Config struct {
	vocab         []TokenScore
	unkID         *int
	bytesFallback bool
	fuseUnk       bool
	// Cache for tokenization
	cache map[string][]string
}

// Unigram implements the Unigram language model for tokenization
type Unigram struct {
	vocab         []TokenScore
	tokenToIDs    map[string]int
	unkID         *int
	bytesFallback bool
	fuseUnk       bool
	// Cache for tokenization
	cache map[string][]string
}

// UnigramBuilder can be used to create a Unigram model with a custom configuration
type UnigramBuilder struct {
	config Config
}

// NewUnigramBuilder creates a new UnigramBuilder with default configuration
func NewUnigramBuilder() *UnigramBuilder {
	return &UnigramBuilder{
		config: Config{
			vocab:         []TokenScore{},
			unkID:         nil,
			bytesFallback: false,
			fuseUnk:       true, // Default to true to match Rust implementation
			cache:         make(map[string][]string),
		},
	}
}

// Vocab sets the vocabulary for the Unigram model
func (ub *UnigramBuilder) Vocab(vocab []TokenScore) *UnigramBuilder {
	ub.config.vocab = vocab
	return ub
}

// UnkID sets the unknown token ID for the Unigram model
func (ub *UnigramBuilder) UnkID(unkID int) *UnigramBuilder {
	ub.config.unkID = &unkID
	return ub
}

// BytesFallback sets whether to use byte fallback for unknown tokens
func (ub *UnigramBuilder) BytesFallback(bytesFallback bool) *UnigramBuilder {
	ub.config.bytesFallback = bytesFallback
	return ub
}

// FuseUnk sets whether to fuse unknown tokens together
func (ub *UnigramBuilder) FuseUnk(fuseUnk bool) *UnigramBuilder {
	ub.config.fuseUnk = fuseUnk
	return ub
}

// Build creates a new Unigram model with the configured parameters
func (ub *UnigramBuilder) Build() (*Unigram, error) {
	// Create token to ID mapping
	tokenToIDs := make(map[string]int, len(ub.config.vocab))
	for i, ts := range ub.config.vocab {
		tokenToIDs[ts.Token] = i
	}

	// Validate unkID if provided
	if ub.config.unkID != nil {
		if *ub.config.unkID >= len(ub.config.vocab) {
			return nil, fmt.Errorf("unkID %d is out of vocabulary range (size: %d)", *ub.config.unkID, len(ub.config.vocab))
		}
	}

	return &Unigram{
		vocab:         ub.config.vocab,
		tokenToIDs:    tokenToIDs,
		unkID:         ub.config.unkID,
		bytesFallback: ub.config.bytesFallback,
		fuseUnk:       ub.config.fuseUnk,
		cache:         make(map[string][]string),
	}, nil
}

// New creates a new Unigram model with the given vocabulary and options
func New(vocab []TokenScore, opts *util.Params) (*Unigram, error) {
	builder := NewUnigramBuilder().Vocab(vocab)

	if opts != nil {
		if opts.Has("unk_id") {
			// Handle different types for unk_id
			unkIDValue := opts.Get("unk_id")
			var unkID int

			switch v := unkIDValue.(type) {
			case float64:
				unkID = int(v)
			case int:
				unkID = v
			default:
				return nil, fmt.Errorf("unk_id must be an integer, got %T", unkIDValue)
			}

			builder.UnkID(unkID)
		}

		if opts.Has("byte_fallback") {
			bytesFallback := opts.Get("byte_fallback").(bool)
			builder.BytesFallback(bytesFallback)
		}

		if opts.Has("fuse_unk") {
			fuseUnk := opts.Get("fuse_unk").(bool)
			builder.FuseUnk(fuseUnk)
		}
	}

	return builder.Build()
}

// GetVocab returns the vocabulary mapping (token -> ID)
func (u *Unigram) GetVocab() map[string]int {
	return u.tokenToIDs
}

// GetVocabSize returns the size of the vocabulary
func (u *Unigram) GetVocabSize() int {
	return len(u.vocab)
}

// TokenToId returns the ID for the given token
func (u *Unigram) TokenToId(token string) (int, bool) {
	id, ok := u.tokenToIDs[token]
	return id, ok
}

// getMinScore returns the minimum score in the vocabulary
// This is used for calculating the unknown token penalty
func (u *Unigram) getMinScore() float64 {
	if len(u.vocab) == 0 {
		return 0.0
	}

	minScore := u.vocab[0].Score
	for _, ts := range u.vocab {
		if ts.Score < minScore {
			minScore = ts.Score
		}
	}

	return minScore
}

// IdToToken returns the token for the given ID
func (u *Unigram) IdToToken(id int) (string, bool) {
	if id < 0 || id >= len(u.vocab) {
		return "", false
	}
	return u.vocab[id].Token, true
}

// Save saves the Unigram model to the given directory
func (u *Unigram) Save(dir string, prefixOpt ...string) error {
	var prefix string
	if len(prefixOpt) > 0 {
		prefix = prefixOpt[0]
	} else {
		prefix = "tokenizer"
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Save vocab.json
	vocabPath := filepath.Join(dir, prefix+"-vocab.json")
	vocabFile, err := os.Create(vocabPath)
	if err != nil {
		return err
	}
	defer vocabFile.Close()

	// Create a serializable representation of the vocabulary
	type VocabEntry struct {
		Token string  `json:"token"`
		Score float64 `json:"score"`
	}
	vocabEntries := make([]VocabEntry, len(u.vocab))
	for i, ts := range u.vocab {
		vocabEntries[i] = VocabEntry{
			Token: ts.Token,
			Score: ts.Score,
		}
	}

	// Create the model config
	modelConfig := map[string]interface{}{
		"type":          "Unigram",
		"vocab":         vocabEntries,
		"unk_id":        u.unkID,
		"byte_fallback": u.bytesFallback,
	}

	encoder := json.NewEncoder(vocabFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(modelConfig); err != nil {
		return err
	}

	return nil
}

// Tokenize tokenizes the given sequence into multiple tokens
func (u *Unigram) Tokenize(sequence string) ([]tokenizer.Token, error) {
	// Check cache first
	if tokens, ok := u.cache[sequence]; ok {
		return u.tokensToTokenizer(tokens, sequence), nil
	}

	// If byte fallback is enabled, always use it
	if u.bytesFallback {
		tokens := u.tokenizeWithByteFallback(sequence)
		u.cache[sequence] = tokens
		return u.tokensToTokenizer(tokens, sequence), nil
	}

	// Tokenize using the Viterbi algorithm
	tokens, err := u.tokenizeWithViterbi(sequence)
	if err != nil {
		return nil, err
	}

	// Cache the result
	u.cache[sequence] = tokens

	return u.tokensToTokenizer(tokens, sequence), nil
}

// tokensToTokenizer converts string tokens to tokenizer.Token
func (u *Unigram) tokensToTokenizer(tokens []string, sequence string) []tokenizer.Token {
	var result []tokenizer.Token
	var offset int

	for _, token := range tokens {
		length := len(token)
		id, ok := u.TokenToId(token)
		if !ok {
			// Handle unknown token
			if u.unkID != nil {
				id = *u.unkID
			} else {
				// Skip if no unkID is defined
				offset += length
				continue
			}
		}

		result = append(result, tokenizer.Token{
			Id:      id,
			Value:   token,
			Offsets: []int{offset, offset + length},
		})

		offset += length
	}

	return result
}

// tokenizeWithViterbi implements the Viterbi algorithm for tokenization
func (u *Unigram) tokenizeWithViterbi(sequence string) ([]string, error) {
	if len(sequence) == 0 {
		return []string{}, nil
	}

	// Create a lattice for the Viterbi algorithm
	// Each position in the lattice represents the best segmentation up to that point
	type latticeNode struct {
		bestScore float64
		bestIndex int
		id        int // Store the token ID for backtracking
	}

	// Initialize the lattice with one node per character in the sequence
	n := len(sequence)
	lattice := make([]latticeNode, n+1)

	// The first position has a score of 0 (base case)
	lattice[0] = latticeNode{bestScore: 0, bestIndex: 0, id: -1}

	// Initialize the rest of the lattice with negative infinity
	for i := 1; i <= n; i++ {
		lattice[i] = latticeNode{bestScore: math.Inf(-1), bestIndex: -1, id: -1}
	}

	// Constant for unknown token penalty, matching the Rust implementation
	const kUnkPenalty float64 = 10.0
	unkScore := u.getMinScore() - kUnkPenalty

	// For each position in the lattice
	for i := 0; i < n; i++ {
		if lattice[i].bestScore == math.Inf(-1) {
			// If we can't reach this position, skip it
			continue
		}

		// Try all possible tokens starting at this position
		// Limit the maximum token length to 20 characters for efficiency
		maxLen := 20
		if i+maxLen > n {
			maxLen = n - i
		}

		// Track if we've found a single character token at this position
		hasSingleNode := false

		// Get the length of the current character in UTF-8
		charLen := 1
		if i < n {
			r, _ := utf8.DecodeRuneInString(sequence[i:])
			charLen = utf8.RuneLen(r)
		}

		// Try to match tokens of all lengths
		// We don't need to iterate from longest to shortest anymore
		// because the Viterbi algorithm will find the optimal path
		for j := i + 1; j <= i+maxLen && j <= n; j++ {
			// Extract the substring from i to j
			token := sequence[i:j]

			// Find the token in the vocabulary
			id, ok := u.TokenToId(token)

			if !ok {
				// Skip if token is not in vocabulary
				continue
			}

			// If the token is in the vocabulary, use its score
			score := u.vocab[id].Score

			// Add a bonus for longer tokens to match the Rust implementation
			// This helps prioritize longer tokens like "ab" over "a" + "b"
			tokenLength := j - i
			lengthBonus := float64(tokenLength) * 0.1

			// Calculate the new score for this path
			newScore := lattice[i].bestScore + score + lengthBonus

			// If this is a better path to position j, update the lattice
			if newScore > lattice[j].bestScore {
				lattice[j].bestScore = newScore
				lattice[j].bestIndex = i
				lattice[j].id = id
			}

			// Check if we've found a single character token
			if j == i+charLen {
				hasSingleNode = true
			}
		}

		// If we haven't found a single character token and we have an unknown token ID,
		// add an unknown token for this character
		if !hasSingleNode && u.unkID != nil {
			j := i + charLen
			if j <= n {
				newScore := lattice[i].bestScore + unkScore
				if newScore > lattice[j].bestScore {
					lattice[j].bestScore = newScore
					lattice[j].bestIndex = i
					lattice[j].id = *u.unkID
				}
			}
		}
	}

	// If we couldn't reach the end of the sequence, handle unknown tokens
	if lattice[n].bestScore == math.Inf(-1) {
		// If we have an unknown token ID, return the whole sequence as unknown
		if u.unkID != nil {
			return []string{sequence}, nil
		}

		// If we're using byte fallback, handle it here
		if u.bytesFallback {
			return u.tokenizeWithByteFallback(sequence), nil
		}

		return nil, fmt.Errorf("could not tokenize sequence with Viterbi algorithm")
	}

	// Backtrack to find the best segmentation
	if u.fuseUnk && u.unkID != nil {
		// When fuseUnk is true, we need to fuse consecutive unknown tokens
		var tokens []string
		var unkTokens []string
		pos := n

		for pos > 0 {
			start := lattice[pos].bestIndex
			if start < 0 {
				// This shouldn't happen if the algorithm is correct
				return nil, fmt.Errorf("invalid backtracking index: %d", start)
			}

			token := sequence[start:pos]
			id := lattice[pos].id

			// Check if this is an unknown token
			if id == *u.unkID {
				// Add to the list of unknown tokens
				unkTokens = append([]string{token}, unkTokens...)
			} else {
				// If we have accumulated unknown tokens, add them as a single token
				if len(unkTokens) > 0 {
					fusedUnk := strings.Join(unkTokens, "")
					tokens = append([]string{fusedUnk}, tokens...)
					unkTokens = nil
				}
				// Add the current token
				tokens = append([]string{token}, tokens...)
			}

			pos = start
		}

		// If we have any remaining unknown tokens, add them
		if len(unkTokens) > 0 {
			fusedUnk := strings.Join(unkTokens, "")
			tokens = append([]string{fusedUnk}, tokens...)
		}

		return tokens, nil
	} else {
		// Standard backtracking without fusing unknown tokens
		var tokens []string
		pos := n

		for pos > 0 {
			start := lattice[pos].bestIndex
			if start < 0 {
				// This shouldn't happen if the algorithm is correct
				return nil, fmt.Errorf("invalid backtracking index: %d", start)
			}

			token := sequence[start:pos]
			tokens = append([]string{token}, tokens...)
			pos = start
		}

		return tokens, nil
	}
}

// tokenizeWithByteFallback tokenizes a string by representing each byte as a separate token
func (u *Unigram) tokenizeWithByteFallback(sequence string) []string {
	var tokens []string

	// Convert the string to bytes
	bytes := []byte(sequence)

	// For each byte, create a token in the format "<0xXX>"
	for _, b := range bytes {
		// Format the byte as a token
		token := fmt.Sprintf("<0x%02X>", b)

		// Check if the token is in the vocabulary
		_, ok := u.TokenToId(token)
		if ok {
			// If the token is in the vocabulary, use it
			tokens = append(tokens, token)
		} else if u.unkID != nil {
			// If the token is not in the vocabulary and we have an unknown token ID,
			// use the unknown token
			unkToken, _ := u.IdToToken(*u.unkID)
			tokens = append(tokens, unkToken)
		} else {
			// If we don't have an unknown token ID, use the byte token format
			tokens = append(tokens, token)
		}
	}

	return tokens
}
