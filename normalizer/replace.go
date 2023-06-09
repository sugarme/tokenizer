package normalizer

import (
	"fmt"
	"strings"
)

// Enum of different patterns that Replace can use.
type ReplacePattern int

const (
	String ReplacePattern = iota
	Regex
)

type Replace struct {
	PatternType ReplacePattern `json:"pattern_type"`
	Pattern     Pattern        `json:"pattern"`
	Content     string         `json:"content"`
}

var _ Normalizer = new(Replace)

func NewReplace(patternType ReplacePattern, pattern string, content string) *Replace {
	var pat Pattern
	switch patternType {
	case String:
		pat = NewStringPattern(pattern)
	case Regex:
		pat = NewRegexpPattern(pattern)
	default:
		msg := fmt.Sprintf("Not supported ReplacePattern %q", patternType)
		panic(msg)
	}

	return &Replace{
		PatternType: patternType,
		Pattern:     pat,
		Content:     content,
	}
}

// Implement Normalizer for Replace
func (r *Replace) Normalize(normalized *NormalizedString) (*NormalizedString, error) {
	return normalized.Replace(r.Pattern, r.Content), nil
}

// Implement Decoder for Replace
func (r *Replace) DecodeChain(tokens []string) []string {
	var out []string
	for _, token := range tokens {
		var newTokParts []string
		offsetMatches := r.Pattern.FindMatches(token)
		for _, offsetMatch := range offsetMatches {
			if offsetMatch.Match {
				newTokParts = append(newTokParts, r.Content)
			} else {
				start := offsetMatch.Offsets[0]
				end := offsetMatch.Offsets[1]
				newTokParts = append(newTokParts, token[start:end])
			}
		}
		newTok := strings.Join(newTokParts, "")
		out = append(out, newTok)
	}

	return out
}

func (r *Replace) Decode(tokens []string) string {
	return strings.Join(r.DecodeChain(tokens), "")
}
