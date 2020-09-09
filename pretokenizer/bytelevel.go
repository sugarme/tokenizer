package pretokenizer

import (
	// "fmt"
	"regexp"
	"strings"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
	slice "github.com/sugarme/tokenizer/util/slice"
)

// Regular epxression to split string to `word` token
// including prefix whitespace. Contractions and punctuation
// will be split as well.
// Ref.https://regex101.com/r/pf5XJv
// const splitRegStr = `'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+`
// const splitRegStr = `'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+`
// TODO: this RE does not cover the case with trailing whitespace!!!
const splitRegStr = `'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+`

var splitRE = regexp.MustCompile(splitRegStr)

var BytesChar map[uint8]string = GenerateBytesChar()

var CharBytes map[string]uint8 = func() map[string]uint8 {
	var bc = GenerateBytesChar()
	var cb map[string]uint8 = make(map[string]uint8)
	for b, c := range bc {
		cb[c] = b
	}
	return cb
}()

// BytesChar maps first 0-255 (byte) to first 0-255 `char` in unicode
// Ref. https://en.wikipedia.org/wiki/List_of_Unicode_characters
func GenerateBytesChar() map[uint8]string {
	var bs []uint8 // byte
	var bc map[uint8]string = make(map[uint8]string)

	// Basic latin
	for i := 33; i <= 126; i++ {
		bs = append(bs, uint8(i))
		bc[uint8(i)] = string(uint8(i))
	}

	// latin-1 supplement (excluding `173`)
	for i := 161; i <= 172 && i != 173; i++ {
		bs = append(bs, uint8(i))
		bc[uint8(i)] = string(i)
	}

	// Append `control` byte (0-32) and (127-160) and 173
	// Due to the control byte, first 256 runes will be shifted right 256
	var n = 0
	for i := 0; i <= 255; i++ {
		if !slice.Contain(uint8(i), bs) {
			// if !contain(uint8(i), bs) {
			bs = append(bs, uint8(i))
			bc[uint8(i)] = string(256 + n)
			n += 1
		}
	}

	return bc
}

// ByteLevel provides all the neccessary steps to handle the
// BPE tokenization at byte-level. It takes care of all the required
// processing steps to transform a utf-8 string as needed before and
// after the BPE model does it job.
type ByteLevel struct {
	// whether to add a leading space to the first word.
	// It allows to treat the leading word just as any other words.
	AddPrefixSpace bool

	// Whether the post processing step should trim offsets
	// to avoid including whitespaces.
	TrimOffsets bool
}

// NewByteLevel returns a default ByteLevel with both
// AddPrefixSpace and TrimOffsets set true
func NewByteLevel() *ByteLevel {
	return &ByteLevel{
		AddPrefixSpace: true,
		TrimOffsets:    true,
	}
}

// Alphabet returns set of first 256 unicode `char`
func (bl *ByteLevel) Alphabet() map[string]struct{} {
	var ab = make(map[string]struct{})
	for _, c := range BytesChar {
		ab[c] = struct{}{}
	}

	return ab
}

// SetAddPrefixSpace set `AddPrefixSpace` property
func (bl *ByteLevel) SetAddPrefixSpace(v bool) {
	bl.AddPrefixSpace = v
}

// SetTrimOffsets set `TrimOffsets` property
func (bl *ByteLevel) SetTrimOffsets(v bool) {
	bl.TrimOffsets = v
}

// Implement `PreTokenizer` methods for `ByteLevel`:
// =================================================

// PreTokenizer, as a `PreTokenizer`, `ByteLevel` is in charge of transforming all the unicode characters into
// their byte-level counterpart. It also splits the input according to the configured regex.
func (bl *ByteLevel) PreTokenize(pretokenized *tokenizer.PreTokenizedString) (*tokenizer.PreTokenizedString, error) {

	pretok := pretokenized.Split(func(noop int, normalized *normalizer.NormalizedString) []tokenizer.SplitIdx {
		var newNormalized *normalizer.NormalizedString = normalized
		if bl.AddPrefixSpace && !strings.HasPrefix(normalized.GetNormalized(), " ") {
			newNormalized = normalized.Prepend(" ")
		}

		splitPattern := normalizer.NewRegexpPattern(splitRegStr)
		splits := newNormalized.Split(splitPattern, normalizer.IsolatediBehavior)

		var splitIdx []tokenizer.SplitIdx

		for _, s := range splits {
			split := s // NOTE: to deep copy variable otherwise its updated to the last item as we will pass its pointer.
			splitIdx = append(splitIdx, tokenizer.SplitIdx{Normalized: &split, Tokens: nil})
		}

		return splitIdx
	})

	finalPretok := pretok.Normalize(func(normalized *normalizer.NormalizedString) *normalizer.NormalizedString {
		s := normalized.GetNormalized()
		var changeMap []normalizer.ChangeMap
		for _, r := range s {
			bytes := []byte(string(r))
			for i, b := range bytes {
				change := 0
				if i > 0 {
					change = 1
				}

				char := BytesChar[b]
				changeMap = append(changeMap, normalizer.ChangeMap{RuneVal: char, Changes: change})
			}
		}

		return normalized.Transform(changeMap, 0)
	})

	return finalPretok, nil
}

// Implement Decoder for `ByteLevel`:
// ==================================

// Decode converts any byte-level characters to their unicode couterpart
// before merging everything back into a single string
func (bl *ByteLevel) Decode(tokens []string) string {
	s := strings.Join(tokens, "")
	chars := strings.Split(s, "")

	var bytes []byte

	for _, c := range chars {
		b := CharBytes[c]

		bytes = append(bytes, b)
	}

	return string(bytes)
}

// Implement PostProcessor for ByteLevel
func (bl *ByteLevel) AddedToken(isPair bool) uint {
	return 0
}

func (bl *ByteLevel) Process(encoding, pairEncoding *tokenizer.Encoding, addSpecialTokens bool) *tokenizer.Encoding {

	if !bl.TrimOffsets {
		return tokenizer.DefaultProcess(encoding, pairEncoding, addSpecialTokens)
	}

	var newEncoding *tokenizer.Encoding
	newEncoding = processOffsets(encoding, bl.AddPrefixSpace)

	overflowEncodings := newEncoding.GetOverflowing()
	var newOverflowEncodings []tokenizer.Encoding
	for _, e := range overflowEncodings {
		newEn := processOffsets(&e, bl.AddPrefixSpace)
		newOverflowEncodings = append(newOverflowEncodings, *newEn)
	}
	newEncoding.Overflowing = newOverflowEncodings

	var (
		newPairEncoding         *tokenizer.Encoding
		newOverflowPairEncoding []tokenizer.Encoding
	)

	if pairEncoding != nil {
		newPairEncoding = processOffsets(pairEncoding, bl.AddPrefixSpace)
		for _, en := range newPairEncoding.Overflowing {
			newEn := processOffsets(&en, bl.AddPrefixSpace)
			newOverflowPairEncoding = append(newOverflowPairEncoding, *newEn)
		}
		newPairEncoding.Overflowing = newOverflowPairEncoding
	}

	return tokenizer.DefaultProcess(newEncoding, newPairEncoding, addSpecialTokens)
}

// func processOffsets(isTrimOffsets bool, encoding *tokenizer.Encoding) *tokenizer.Encoding {
func processOffsets(encoding *tokenizer.Encoding, addPrefixSpace bool) *tokenizer.Encoding {

	type Modif struct {
		LeadingSpaces int
		TrailingSpace int
	}

	var modifs []Modif
	var newOffsets [][]int

	toks := encoding.GetTokens()
	for _, tok := range toks {
		var leadingSpaces int = 0
		chars := strings.Split(tok, "")
		for _, c := range chars {
			// if c != "Ġ" {
			if c != BytesChar[' '] && c != " " {
				break
			}
			leadingSpaces += 1
		}

		var trailingSpaces int = 0
		for i := len(chars) - 1; i >= 0; i-- {
			// if chars[i] != "Ġ" {
			if chars[i] != BytesChar[' '] && chars[i] != " " {
				break
			}
			trailingSpaces += 1
		}

		if leadingSpaces > 0 || trailingSpaces > 0 {
			modifs = append(modifs, Modif{
				LeadingSpaces: leadingSpaces,
				TrailingSpace: trailingSpaces,
			})
		}
	}

	for i, m := range modifs {
		var offset0, offset1 int
		offsets := encoding.GetOffsets()[i]
		ld := m.LeadingSpaces
		offset0 = offsets[0]
		if m.LeadingSpaces > 0 {
			if i == 0 && addPrefixSpace && m.LeadingSpaces == 1 {
				// If we are processing the first pair of offsets, with `addPrefixSpace`,
				// then we shouldn't remove anything we added. If there are more than one
				// leading spaces though, it means we didn't add them, and they should be
				// removed.
				ld = 0
			}

			offset0 = offsets[0] + ld
			if offset0 > offsets[1] {
				offset0 = offsets[1]
			}
		}

		tl := m.TrailingSpace

		offset1 = offsets[1]
		if tl > 0 && offsets[1] >= tl {
			offset1 = offsets[1] - tl
			if offset1 < offset0 {
				offset1 = offset0
			}
		}

		newOffsets = append(newOffsets, []int{offset0, offset1})
	}

	encoding.Offsets = newOffsets

	return encoding
}
