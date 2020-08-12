package tokenizer

import (
	"fmt"
	"github.com/sugarme/tokenizer/normalizer"
	"log"
	"regexp"
	"unicode"
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

// type MatchingSet struct{}

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
	// TODO. continue
}
