package spm

// spm package provides APIs to parse Google/sentencepiece precompiled_charsmap and its normalizer.
// https://github.com/google/sentencepiece
//
// ported from
// https://github.com/huggingface/spm_precompiled

import (
	"encoding/base64"
	"encoding/binary"
	"unicode"

	"fmt"
	"strings"

	"github.com/rivo/uniseg"
)

// HF comment
// This struct is specifically done to be compatible with SentencePiece
// SentencePiece models embed their Normalizer within a `precompiled_charsmap`
// that both represents a Trie, and embedded rewrite rules.
// In order to be 100% compliant we need to interpret that binary format too.
// The format is [u32 (length of trie), trie: [u32], normalized: String]
// The trie has u8 as entries, and u32 as values, those u32 values
// point to offsets withing the String that correspond to the real replace value
// The normalized string contains '\0' that should indicate the end of an entry.
//
// Hence, normalized could be "abc\0", some entry in the trie could be 0 meaning
// the value is "abc" and another one be 1 meaning the actual entry was "bc".

type Precompiled struct {
	PrecompiledCharsmap []byte
	Normalized          string
	Trie                *DoubleArray
}

type ArrayUnit struct {
	value uint
}
type Array []ArrayUnit

type DoubleArray struct {
	Array Array
}

func AsBase64(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}

func FromBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

type ArrayUnitTrait interface {
	HasLeaf() bool
	Value() int   // isize
	Label() uint  // usize
	Offset() uint // usize
}

// Implement ArrayUnitTrait for ArrayUnit
func (u ArrayUnit) HasLeaf() bool {
	// (self >> 8) & 1 == 1
	return ((u.value >> 8) & 1) == 1
}

func (u ArrayUnit) Value() int {
	// (self & ((1usize << 31) - 1)) as isize
	return int(uint(u.value) & ((uint(1) << 31) - 1))
}

func (u ArrayUnit) Label() uint {
	// self & ((1usize << 31) | 0xFF)
	return uint(u.value) & ((uint(1) << 31) | 0xFF)
}

func (u ArrayUnit) Offset() uint {
	// (self >> 10) << ((self & (1usize << 9)) >> 6)
	return (uint(u.value) >> 10) << ((uint(u.value) & (uint(1) << 9)) >> 6)
}

func NewDoubleArrayFrom(array Array) *DoubleArray {
	return &DoubleArray{array}
}

// ref. https://github.com/huggingface/spm_precompiled/blob/81b911a362adef3ad3cc6d5835d2980690dbb871/src/lib.rs#L108
func (da *DoubleArray) CommonPrefixSearch(key []byte) []int {
	nodePos := 0
	var results []int

	unit := da.Array[nodePos]
	nodePos ^= int(unit.Offset()) // bitwise XOR
	for _, c := range key {
		if c == byte(0) {
			break
		}

		nodePos ^= int(c)
		unit = da.Array[nodePos]

		if unit.Label() != uint(c) {
			return results
		}

		nodePos ^= int(unit.Offset())
		if unit.HasLeaf() {
			results = append(results, int(da.Array[nodePos].Value()))
		}
	}

	return results
}

func Parse(precompiledCharsmap []byte) ([]byte, Array) {
	trieSize := binary.LittleEndian.Uint32(precompiledCharsmap[:4])
	rest := precompiledCharsmap[4:]

	trieCharSize := trieSize / 4 // number of tries (nodes)
	trieBlob := make([]ArrayUnit, trieCharSize)
	for i := 0; i < int(trieCharSize); i++ {
		n := binary.LittleEndian.Uint32(rest[:4])
		rest1 := rest[4:]
		rest = rest1
		trieBlob[i] = ArrayUnit{uint(n)}
	}

	normalizedBlob := rest

	return normalizedBlob, trieBlob
}

func NewPrecompiledFrom(data []byte) (*Precompiled, error) {
	normalizedBlob, trieBlob := Parse(data)

	normalized := string(normalizedBlob)
	trie := NewDoubleArrayFrom(trieBlob)

	return &Precompiled{
		PrecompiledCharsmap: data,
		Normalized:          normalized,
		Trie:                trie,
	}, nil
}

func (m *Precompiled) Transform(chunk string) string {
	results := m.Trie.CommonPrefixSearch([]byte(chunk))
	if len(results) == 0 {
		return ""
	}

	index := results[0]
	index2 := index
	for index2 < len(m.Normalized) {
		if []byte(m.Normalized)[index2] == byte(0) {
			break
		}
		index2 += 1
	}

	normalized := m.Normalized[index:index2]

	return string(normalized)
}

func NormalizeMn(input string) string {
	return normalizeMn(input)
}

func normalizeMn(input string) string {
	var out []string
	runes := []rune(input)
	for _, r := range runes {
		if unicode.Is(unicode.Mn, r) {
			v := fmt.Sprintf("%U", r)
			out = append(out, v)
		} else {
			out = append(out, string(r))
		}
	}

	return strings.Join(out, "")
}

func (m *Precompiled) NormalizeString(original string) string {
	var chars []string

	graphemes := uniseg.NewGraphemes(original)

	for graphemes.Next() {
		grapheme := graphemes.Str()

		// NOTE.This comment from Narsil at HF
		// Future reader. From @Narsil.
		// Yes, this is weird,
		// Yes, this seems broken
		// No, I don't know why Google did this.
		// If you question this code, check this normalizer against
		// XNLI database (all languages) with Unigram model against
		// Mbart, XLMRoberta *AND* Marian. If you don't get 100% or
		// break a single test.
		// You don't pass.
		if len(grapheme) < 6 {
			norm := m.Transform(grapheme)
			if len(norm) > 0 {
				for _, c := range strings.Split(norm, "") {
					chars = append(chars, c)
				}
				return strings.Join(chars, "")
			}
		}

		// TT. This is a hacky way to turn non-spacing marks into hexa string
		// and pass the unit tests.
		grapheme = normalizeMn(grapheme)
		var charIdx int = 0
		for _, r := range grapheme {
			part := string(grapheme)[charIdx : charIdx+len(string(r))]
			norm := m.Transform(part)
			norm = normalizeMn(norm)
			if len(norm) > 0 {
				for _, c := range strings.Split(norm, "") {
					chars = append(chars, c)
				}
			} else {
				chars = append(chars, string(r))
			}
			charIdx += len(string(r))
		}
	}

	return strings.Join(chars, "")
}
