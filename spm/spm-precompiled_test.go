package spm

import (
	"log"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	buff, err := FromBase64(compiledString)
	if err != nil {
		panic(err)
	}
	precompiled1, err := NewPrecompiledFrom(buff)
	if err != nil {
		panic(err)
	}

	log.Printf("number of tries: %v\n", len(precompiled1.Trie.Array))

	bytes := NmtNfkc()
	precompiled2, err := NewPrecompiledFrom(bytes)
	if err != nil {
		panic(err)
	}

	originalBytes := []byte{0xd8, 0xa7, 0xd9, 0x93}
	got := precompiled2.Trie.CommonPrefixSearch(originalBytes)
	want := []int{4050}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}

	norm := precompiled2.Normalized
	got1 := string(norm[4050:4053])
	want1 := "آ\x00"
	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("want %s, got %s\n", want1, got1)
	}

}

func TestCommonPrefixSearch(t *testing.T) {
	m, err := NewPrecompiledFrom(NmtNfkc())
	if err != nil {
		panic(err)
	}

	buf := []byte("\ufb01")
	got := m.Trie.CommonPrefixSearch(buf)
	want := []int{2130}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}

	// check the null termination
	got1 := string(m.Normalized[2130:2133])
	want1 := "fi\x00"
	if !reflect.DeepEqual(want1, got1) {
		t.Errorf("want %q, got %q\n", want1, got1)
	}

	got2 := m.Trie.CommonPrefixSearch([]byte(" "))
	var want2 []int = nil
	if !reflect.DeepEqual(want2, got2) {
		t.Errorf("want %v, got %v\n", want2, got2)
	}

	got3 := m.Trie.CommonPrefixSearch([]byte("𝔾"))
	want3 := []int{1786}
	if !reflect.DeepEqual(want3, got3) {
		t.Errorf("want %v, got %v\n", want3, got3)
	}

	// Transform
	got4 := m.Transform("𝔾")
	want4 := "G"
	if !reflect.DeepEqual(want4, got4) {
		t.Errorf("want %v, got %v\n", want4, got4)
	}
	got4 = m.Transform("𝕠")
	want4 = "o"
	if !reflect.DeepEqual(want4, got4) {
		t.Errorf("want %v, got %v\n", want4, got4)
	}
	got4 = m.Transform("\u200d")
	want4 = " "
	if !reflect.DeepEqual(want4, got4) {
		t.Errorf("want %v, got %v\n", want4, got4)
	}
}

func TestPrecompiled_NormalizeString(t *testing.T) {
	m, err := NewPrecompiledFrom(NmtNfkc())
	if err != nil {
		panic(err)
	}

	originalBytes := []byte{0xd8, 0xa7, 0xd9, 0x93}
	original := string(originalBytes) //  "آ"
	log.Printf("original: %s\n", original)

	normalized := "آ" // this grapheme is 2 runes
	got := m.NormalizeString(original)
	want := normalized
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}

	// Thai
	original = "เขาไม่ได้พูดสักคำ"
	// normalized := `เขาไม\u{e48}ได\u{e49}พ\u{e39}ดส\u{e31}กค\u{e4d}า`
	normalized = `เขาไมU+0E48ไดU+0E49พU+0E39ดสU+0E31กคU+0E4Dา`

	got = m.NormalizeString(original)
	want = normalized
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %s, got %s\n", want, got)
	}

	// Hindi
	original = `ड़ी दुख`
	// normalized = `ड\u{93c}ी द\u{941}ख`
	normalized = `डU+093Cी दU+0941ख`
	got = m.NormalizeString(original)
	want = normalized
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %s, got %s\n", want, got)
	}
}
