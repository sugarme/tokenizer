package unigram

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"

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
	// Cache for tokenization
	cache map[string][]string
}

// Unigram implements the Unigram language model for tokenization
type Unigram struct {
	vocab         []TokenScore
	tokenToIDs    map[string]int
	unkID         *int
	bytesFallback bool
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
	}

	// Initialize the lattice with one node per character in the sequence
	n := len(sequence)
	lattice := make([]latticeNode, n+1)
	
	// The first position has a score of 0 (base case)
	lattice[0] = latticeNode{bestScore: 0, bestIndex: 0}
	
	// Initialize the rest of the lattice with negative infinity
	for i := 1; i <= n; i++ {
		lattice[i] = latticeNode{bestScore: math.Inf(-1), bestIndex: -1}
	}
	
	// For each position in the lattice
	for i := 0; i < n; i++ {
		if lattice[i].bestScore == math.Inf(-1) {
			// If we can't reach this position, skip it
			continue
		}
		
		// Try all possible tokens starting at this position
		// Limit the maximum token length to 20 characters for efficiency
		maxLen := 20
		if i + maxLen > n {
			maxLen = n - i
		}
		
		for j := i + 1; j <= i + maxLen && j <= n; j++ {
			// Extract the substring from i to j
			token := sequence[i:j]
			
			// Find the token in the vocabulary
			id, ok := u.TokenToId(token)
			var score float64
			
			if ok {
				// If the token is in the vocabulary, use its score
				score = u.vocab[id].Score
			} else if u.unkID != nil && j == i+1 {
				// If we have an unknown token ID, use it with a penalty
				// Only for single characters to avoid unknown tokens for every substring
				score = -10.0 // Penalty for unknown tokens
			} else {
				// Skip this token
				continue
			}
			
			// Calculate the new score for this path
			newScore := lattice[i].bestScore + score
			
			// If this is a better path to position j, update the lattice
			if newScore > lattice[j].bestScore {
				lattice[j].bestScore = newScore
				lattice[j].bestIndex = i
			}
		}
	}
	
	// If we couldn't reach the end of the sequence, handle unknown tokens
	if lattice[n].bestScore == math.Inf(-1) {
		// If we have an unknown token ID, return the whole sequence as unknown
		if u.unkID != nil {
			return []string{sequence}, nil
		}
		
		// If we're using byte fallback, we would handle it here
		if u.bytesFallback {
			// Implement byte fallback here if needed
			return []string{sequence}, nil
		}
		
		return nil, fmt.Errorf("could not tokenize sequence with Viterbi algorithm")
	}
	
	// Backtrack to find the best segmentation
	var tokens []string
	pos := n
	
	for pos > 0 {
		start := lattice[pos].bestIndex
		if start < 0 {
			// This shouldn't happen if the algorithm is correct
			return nil, fmt.Errorf("invalid backtracking index: %d", start)
		}
		tokens = append([]string{sequence[start:pos]}, tokens...)
		pos = start
	}
	
	return tokens, nil
}
