package pretokenizer

import (
	"regexp"
	"strings"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/normalizer"
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

func encodeUTF8(r rune) []byte {
	const (
		// first byte of a 2-byte encoding starts 110 and carries 5 bits of data
		b2Lead = 0xC0 // 1100 0000
		b2Mask = 0x1F // 0001 1111

		// first byte of a 3-byte encoding starts 1110 and carries 4 bits of data
		b3Lead = 0xE0 // 1110 0000
		b3Mask = 0x0F // 0000 1111

		// first byte of a 4-byte encoding starts 11110 and carries 3 bits of data
		b4Lead = 0xF0 // 1111 0000
		b4Mask = 0x07 // 0000 0111

		// non-first bytes start 10 and carry 6 bits of data
		mbLead = 0x80 // 1000 0000
		mbMask = 0x3F // 0011 1111
	)

	switch i := uint32(r); {
	case i <= 1<<7-1: // max code point that encodes into a single byte
		return []byte{byte(r)}
	case i <= 1<<11-1: // into two bytes
		return []byte{
			b2Lead | byte(r>>6),
			mbLead | byte(r)&mbMask}
	case i <= 1<<16-1: // three
		return []byte{
			b3Lead | byte(r>>12),
			mbLead | byte(r>>6)&mbMask,
			mbLead | byte(r)&mbMask}
	default:
		return []byte{
			b4Lead | byte(r>>18),
			mbLead | byte(r>>12)&mbMask,
			mbLead | byte(r>>6)&mbMask,
			mbLead | byte(r)&mbMask}
	}
}

// BytesChar maps first 0-255 (byte) to first 0-255 `char` in unicode
// Ref. https://en.wikipedia.org/wiki/List_of_Unicode_characters
// Ref. https://rosettacode.org/wiki/UTF-8_encode_and_decode
// See example: https://play.golang.org/p/_1W0ni2uZWm
func GenerateBytesChar() map[uint8]string {
	var bc map[uint8]string = make(map[uint8]string)

	// NOTE: control codes are: 0-32, 127-160 and 173

	// 0 ('Ā') - 32 ('Ġ') - control codes
	n := 0
	for i := 256; i <= 288; i++ {
		r := rune(i)
		b := encodeUTF8(r)
		bc[uint8(n)] = string(b)
		n++
	}

	// 33 ('!') - 126 ('~')
	for i := 33; i <= 126; i++ {
		r := rune(i)
		b := encodeUTF8(r)
		bc[uint8(n)] = string(b)
		n++
	}

	// 127 ('ġ') - 160 ('ł') - control codes
	for i := 289; i <= 322; i++ {
		r := rune(i)
		b := encodeUTF8(r)
		bc[uint8(n)] = string(b)
		n++
	}

	// 161 ('¡') - 172 ('¬')
	for i := 161; i <= 172; i++ {
		r := rune(i)
		b := encodeUTF8(r)
		bc[uint8(n)] = string(b)
		n++
	}

	// 173 - ('Ń') - control code
	if n == 173 {
		r := rune(323)
		b := encodeUTF8(r)
		bc[uint8(n)] = string(b)
		n++
	}

	// 174 ('®') - 255 ('ÿ')
	for i := 174; i <= 255; i++ {
		r := rune(i)
		b := encodeUTF8(r)
		bc[uint8(n)] = string(b)
		n++
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
		splits := newNormalized.Split(splitPattern, normalizer.IsolatedBehavior)

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

var _ tokenizer.Decoder = new(ByteLevel)

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

func (bl *ByteLevel) DecodeChain(tokens []string) []string {
	out := make([]string, len(tokens))
	for i, s := range tokens {
		chars := strings.Split(s, "")
		var bytes []byte
		for _, c := range chars {
			b := CharBytes[c]

			bytes = append(bytes, b)
		}

		out[i] = string(bytes)
	}

	return out
}

// Implement PostProcessor for ByteLevel
// =====================================

func (bl *ByteLevel) AddedToken(isPair bool) int {
	return 0
}

func (bl *ByteLevel) Process(encoding, pairEncoding *tokenizer.Encoding, addSpecialTokens bool) *tokenizer.Encoding {
	encodings := tokenizer.PrepareEncodings(encoding, pairEncoding)
	var newEncodings []tokenizer.Encoding
	if bl.TrimOffsets {
		for _, enc := range encodings {
			processedEnc := processOffsets(&enc, bl.AddPrefixSpace)
			var overflowing []tokenizer.Encoding
			for _, of := range processedEnc.GetOverflowing() {
				processedOF := processOffsets(&of, bl.AddPrefixSpace)
				overflowing = append(overflowing, *processedOF)
			}
			processedEnc.Overflowing = overflowing

			newEncodings = append(newEncodings, *processedEnc)
		}
	} else {
		newEncodings = encodings
	}

	for i, enc := range newEncodings {
		enc.SetSequenceIds(i)
	}

	if pairEncoding != nil {
		return tokenizer.MergeEncodings(newEncodings, false)
	} else {
		return &newEncodings[0]
	}
}

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
		chars := []rune(tok)
		for _, c := range chars {
			if string(c) != "Ġ" {
				break
			}
			leadingSpaces += 1
		}

		var trailingSpaces int = 0
		for i := len(chars) - 1; i >= 0; i-- {
			if string(chars[i]) != "Ġ" {
				break
			}
			trailingSpaces += 1
		}

		modifs = append(modifs, Modif{
			LeadingSpaces: leadingSpaces,
			TrailingSpace: trailingSpaces,
		})
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

func ProcessOffsets(encoding *tokenizer.Encoding, addPrefixSpace bool) *tokenizer.Encoding {
	return processOffsets(encoding, addPrefixSpace)
}
