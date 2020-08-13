package tokenizer

import (
	"fmt"
	"log"
	"regexp"
	"unicode"

	"github.com/sugarme/regexpset"
	"github.com/sugarme/tokenizer/normalizer"
)

// AddedToken represents a token added by the user on top of the
// existing model vocabulary.
//
// AddedToken can be configured to specify the behaviour they should
// have in various situations. I.e.,:
// - Whether they should only match single words
// - Whether to include any whitespace on its left or right
type AddedToken struct {
	// Content is the content of added token
	Content string
	// whether this token is single word or break words
	SingleWord bool
	// Whether this token should strip whitespace on its left
	LStrip bool
	// Whether this token should strip whitespace on its right
	RStrip bool
	// Whether this token should be normalized
	Normalized bool
}

// DefaultAddedToken initiates a default AddedToken
func DefaultAddedToken() (retVal AddedToken) {
	return AddedToken{
		Content:    "",
		SingleWord: false,
		LStrip:     false,
		RStrip:     false,
		Normalized: true,
	}
}

type ATOption func(at *AddedToken)

// WithSingleWord specifies whether this token should only match on whole
// single words, and never part of a word.
func WithSingleWord(singleWord bool) ATOption {
	return func(at *AddedToken) {
		at.SingleWord = singleWord
	}
}

// WithLStrip specify whether this token should include all the whitespaces
// on its left in order to strip them out.
func WithLStrip(lstrip bool) ATOption {
	return func(at *AddedToken) {
		at.LStrip = lstrip
	}
}

// WithRStrip specify whether this token should include all the whitespaces
// on its right in order to strip them out.
func WithRStrip(rstrip bool) ATOption {
	return func(at *AddedToken) {
		at.RStrip = rstrip
	}
}

// WithNormalized specifies whether this token should be normalized and match against its normalized
// version in the input text.
func WithNormalized(normalized bool) ATOption {
	return func(at *AddedToken) {
		at.Normalized = normalized
	}
}

// NewAddedToken builds an AddedToken from given content
// specifying whether it is intended to be a special token.
// NOTE. Special token ar not normalized by default.
func NewAddedToken(s string, special bool, opts ...ATOption) (retVal AddedToken) {
	addedTok := DefaultAddedToken()
	addedTok.Content = s
	addedTok.Normalized = !special

	for _, opt := range opts {
		opt(&addedTok)
	}

	return addedTok
}

// GetPattern retrieves the pattern built for this token, according to all the specified parameters.
func (at AddedToken) GetPattern(n normalizer.Normalizer) (retVal string) {
	var reStr string // regular expression pattern

	if at.SingleWord {
		var firstB, lastB string
		runes := []rune(at.Content)
		firstChar := runes[0]
		lastChar := runes[len(runes)-1]
		if isWordCharacter(firstChar) {
			firstB = `\b`
		} else {
			firstB = ``
		}
		if isWordCharacter(lastChar) {
			lastB = `\b`
		} else {
			lastB = ``
		}

		// normalize the content
		content := normalizer.NewNormalizedFrom(at.Content)
		normalizedString, err := n.Normalize(content)
		if err != nil {
			log.Fatal(err)
		}
		normalized := normalizedString.GetNormalized()

		reStr = fmt.Sprintf("%v%v%v", firstB, regexp.QuoteMeta(normalized), lastB)

	} else {
		reStr = regexp.QuoteMeta(at.Content)
	}

	if at.LStrip && at.RStrip {
		reStr = fmt.Sprintf("(\\s)?%v(\\s)?", reStr)
	} else if at.LStrip {
		reStr = fmt.Sprintf("(\\s)?%v", reStr)
	} else if at.RStrip {
		reStr = fmt.Sprintf("%v(\\s)?", reStr)
	}

	return reStr
}

func isWordCharacter(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsMark(r) || unicode.IsDigit(r) || unicode.IsControl(r) || unicode.IsPunct(r) {
		return true
	}
	return false
}

// matchingSet is a set of regular expression string
type matchingSet struct {
	regexSet regexpset.RegexpSet
	ids      []uint32
}

// AddedVocabulary is a vocabulary built on top of the Model
//
// This provides a way to add new vocabulary to a Tokenizer that has already been trained,
// in a previous process, maybe by someone else. This is especially interesting in the case
// of fine-tunings, where we want to finetune a model while adding some new functionalities
// using some new special tokens, or maybe add some tokens in the case of unknown tokens, etc.
//
// One of the reasons we need to handle these tokens outside of the model is simply that
// for many models, it is not possible to add new tokens after the training process. For example,
// using BPE, the training process generates merges pairs along the vocabulary, and any token
// in the vocabulary can be decomposed in other tokens, down to the original alphabet. If we
// were to add new tokens after this training process, we couldn't make sure the merges pairs
// exist as required.
type AddedVocabulary struct {
	// Contains the mapping from String (token content) to ID. This map contains both special
	// tokens and classic added tokens that were added to the this vocabulary.
	addedTokenMap map[string]uint32
	// Contains the mapping from ID to AddedToken for all the added tokens, both special
	// and classic.
	addedTokenMapR map[uint32]string
	// Contains only the classic AddedToken, in the specific order the user gave them.
	addedTokens []AddedToken
	// Contains only the special AddedToken, in the specific order the user gave them.
	specialTokens []AddedToken
	// A map, containing all the special token for easy access while decoding. This let's
	// us remove them easily with an O(1) complexity.
	specialTokensSet map[string]bool
	// A struct containing all the non-normalized patterns used to split on AddedTokens
	splitRe matchingSet
	// A struct containing all the normalized patterns used to split on AddedTokens
	splitNormalizedRe matchingSet
}

// Len returns size of the additional vocabulary
func (av AddedVocabulary) Len() int {
	return len(av.addedTokenMap)
}

// GetVocab gets the additional vocabulary
func (av AddedVocabulary) GetVocab() (retVal map[string]uint32) {
	return av.addedTokenMap
}

// Get the id matching one of our token if it exists
func (av AddedVocabulary) TokenToId(token string, model Model) (retVal uint32, ok bool) {

	retVal, ok = av.addedTokenMap[token]
	if !ok {
		return model.TokenToId(token)
	}

	return retVal, ok
}

// Get the token matching the given id if it exists
func (av AddedVocabulary) IdToToken(id uint32, model Model) (retVal string, ok bool) {
	retVal, ok = av.addedTokenMapR[id]
	if !ok {
		return model.IdToToken(id)
	}

	return retVal, ok
}

// Check if a token is a special token
func (av AddedVocabulary) IsSpecialToken(token string) (retVal bool) {
	_, retVal = av.specialTokensSet[token]
	return retVal
}

// Add some special tokens to the vocabulary
// It returns number of added tokens
func (av AddedVocabulary) AddSpecialTokens(tokens []AddedToken, model Model, normalizer normalizer.Normalizer) (retVal int) {

	for _, tok := range tokens {
		_, isExist := av.specialTokensSet[tok.Content]
		if tok.Content != "" && !isExist {
			av.specialTokens = append(av.specialTokens, tok)
			av.specialTokensSet[tok.Content] = true
		}
	}

	// Then we delegate to `add_tokens`, that will take care of refreshing added tokens too.
	return av.AddTokens(tokens, model, normalizer)
}

// Add some tokens to the vocabulary
// It returns number of added tokens
func (av AddedVocabulary) AddTokens(tokens []AddedToken, model Model, normalizer normalizer.Normalizer) (retVal int) {

	ignored := 0
	for _, token := range tokens {
		if token.Content == "" {
			ignored++
			continue
		}

		var id uint32
		if i, ok := av.TokenToId(token.Content, model); ok {
			ignored++
			id = i
		} else {
			id = uint32(model.GetVocabSize() + len(av.addedTokenMap))
			av.addedTokenMap[token.Content] = id
		}

		// Update the current revert operation
		av.addedTokenMapR[id] = token.Content
	}

	av.refreshAddedTokens(model, normalizer)

	// return the number of added tokens
	return len(tokens) - ignored
}

type tokenId struct {
	token AddedToken
	id    uint32
}

// refreshAddedTokens reconstructs our internal RegexSet when new tokens are added to the vocabulary.
//
// NOTE. We keep two different regular expression sets, one that will take care of matching against the
// non-normalized string, and one matching against the normalized one.
func (av AddedVocabulary) refreshAddedTokens(model Model, normalizer normalizer.Normalizer) {

	var normIds, nnormIds []uint32
	var normPatterns, nnormPatterns []string
	tokens := append(av.specialTokens, av.addedTokens...)
	for _, token := range tokens {

		id, ok := av.TokenToId(token.Content, model)
		if !ok {
			log.Fatalf("Missing additional token.\n")
		}

		pattern := token.GetPattern(normalizer)
		if token.Normalized {
			normIds = append(normIds, id)
			normPatterns = append(normPatterns, pattern)
		} else {
			nnormIds = append(nnormIds, id)
			nnormPatterns = append(nnormPatterns, pattern)
		}
	}

	normSet, err := regexpset.NewRegexpSet(normPatterns)
	if err != nil {
		log.Fatal(err)
	}
	nnormSet, err := regexpset.NewRegexpSet(nnormPatterns)
	if err != nil {
		log.Fatal(err)
	}

	av.splitNormalizedRe = matchingSet{*normSet, normIds}
	av.splitRe = matchingSet{*nnormSet, nnormIds}
}

type idOffsets struct {
	id      *uint32 // optional
	offsets Offsets
}

// findMatches finds any AddedToken in the given sentence, using the provided MatchingSet.
// This method returns a list of "splits", each of them being a pair of ByteOffsets(Offsets)
// and an optional ID if it is an AddedToken. The list of splits cover the entire input string.
func (av AddedVocabulary) findMatches(sentence string, splitRe matchingSet) (retVal []idOffsets) {

	if len(sentence) == 0 {
		return []idOffsets{{nil, Offsets{0, 0}}}
	}

	return
}
