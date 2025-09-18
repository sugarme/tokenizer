package unigram

import (
	"reflect"
	"testing"

	"github.com/sugarme/tokenizer/util"
)

// Test cases ported from Rust implementation:
// E:\code\hf\tokenizers\tokenizers\src\models\unigram\model.rs

// test_encode tests basic encoding functionality
func TestEncode(t *testing.T) {
	// In Rust:
	// let sentencepieces = vec![
	//     ("<unk>".to_string(), 0.0),
	//     ("a".to_string(), 0.0),
	//     ("b".to_string(), 0.0),
	//     ("c".to_string(), 0.0),
	//     ("d".to_string(), 0.0),
	//     ("cd".to_string(), 1.0),
	//     ("ab".to_string(), 2.0),
	//     ("abc".to_string(), 5.0),
	//     ("abcd".to_string(), 10.0),
	// ];
	pieces := []TokenScore{
		{Token: "<unk>", Score: 0.0},
		{Token: "a", Score: 0.0},
		{Token: "b", Score: 0.0},
		{Token: "c", Score: 0.0},
		{Token: "d", Score: 0.0},
		{Token: "cd", Score: 1.0},
		{Token: "ab", Score: 2.0},
		{Token: "abc", Score: 5.0},
		{Token: "abcd", Score: 10.0},
	}
	params := util.NewParams(map[string]interface{}{
		"unk_id":        0,
		"byte_fallback": false,
	})
	model, err := New(pieces, params)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	// Test tokenization result for "abcd"
	tokens, err := model.Tokenize("abcd")
	if err != nil {
		t.Fatalf("Failed to tokenize: %v", err)
	}
	if got, want := len(tokens), 1; got != want {
		t.Errorf("Wrong number of tokens: got %d, want %d", got, want)
	}
	if got, want := tokens[0].Value, "abcd"; got != want {
		t.Errorf("Wrong token value: got %q, want %q", got, want)
	}
}

// test_encode2 tests more complex encoding scenarios
func TestEncode2(t *testing.T) {
	// In Rust:
	// let sentencepieces = vec![
	//     ("<unk>".to_string(), 0.0),
	//     ("ab".to_string(), 0.0),
	//     ("cd".to_string(), -0.1),
	//     ("abc".to_string(), -0.2),
	//     ("a".to_string(), -0.3),
	//     ("b".to_string(), -0.4),
	//     ("c".to_string(), -0.5),
	//     ("ABC".to_string(), -0.5),
	//     ("abcdabcd".to_string(), 20.0), // User defined just max the scores.
	//     ("q".to_string(), 20.5),
	//     ("r".to_string(), 20.5),
	//     ("qr".to_string(), -0.5),
	// ];
	pieces := []TokenScore{
		{Token: "<unk>", Score: 0.0},
		{Token: "ab", Score: 0.0},
		{Token: "cd", Score: -0.1},
		{Token: "abc", Score: -0.2},
		{Token: "a", Score: -0.3},
		{Token: "b", Score: -0.4},
		{Token: "c", Score: -0.5},
		{Token: "ABC", Score: -0.5},
		{Token: "abcdabcd", Score: 20.0},
		{Token: "q", Score: 20.5},
		{Token: "r", Score: 20.5},
		{Token: "qr", Score: -0.5},
	}
	params := util.NewParams(map[string]interface{}{
		"unk_id":        0,
		"fuse_unk":      true,
		"byte_fallback": false,
	})
	model, err := New(pieces, params)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	// Test various tokenization scenarios
	testCases := []struct {
		input    string
		expected []string
	}{
		{"abc", []string{"abc"}},
		{"AB", []string{"AB"}},
		{"abcd", []string{"ab", "cd"}},
		{"abcc", []string{"abc", "c"}},
		{"xabcabaabcdd", []string{"x", "abc", "ab", "a", "ab", "cd", "d"}},
		{"xyz東京", []string{"xyz東京"}},
		{"ABC", []string{"ABC"}},
		{"abABCcd", []string{"ab", "ABC", "cd"}},
		{"ababcdabcdcd", []string{"ab", "abcdabcd", "cd"}},
		{"abqrcd", []string{"ab", "q", "r", "cd"}},
	}
	for _, tc := range testCases {
		tokens, err := model.Tokenize(tc.input)
		if err != nil {
			t.Errorf("Failed to tokenize %q: %v", tc.input, err)
			continue
		}

		// Extract just the token values for comparison
		tokenValues := make([]string, len(tokens))
		for i, token := range tokens {
			tokenValues[i] = token.Value
		}

		if !reflect.DeepEqual(tokenValues, tc.expected) {
			t.Errorf("Failed for input %q: got %v, want %v", tc.input, tokenValues, tc.expected)
		}
	}

	params = util.NewParams(map[string]interface{}{
		"unk_id":        0,
		"fuse_unk":      false,
		"byte_fallback": false,
	})
	model, err = New(pieces, params)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	testCases = []struct {
		input    string
		expected []string
	}{
		{"AB", []string{"A", "B"}},
		{"xyz東京", []string{"x", "y", "z", "東", "京"}}, // set_fuse_unk
	}

	for _, tc := range testCases {
		tokens, err := model.Tokenize(tc.input)
		if err != nil {
			t.Errorf("Failed to tokenize %q: %v", tc.input, err)
			continue
		}

		// Extract just the token values for comparison
		tokenValues := make([]string, len(tokens))
		for i, token := range tokens {
			tokenValues[i] = token.Value
		}

		if !reflect.DeepEqual(tokenValues, tc.expected) {
			t.Errorf("Failed for input %q: got %v, want %v", tc.input, tokenValues, tc.expected)
		}
	}
}

// test_unigram_bytefallback tests the byte fallback functionality
func TestUnigramByteFallback(t *testing.T) {
	// In Rust:
	// let sentencepieces = vec![
	//     ("<unk>".to_string(), 0.0),
	//     ("<0xC3>".to_string(), -0.01),
	//     ("<0xA9>".to_string(), -0.03),
	// ];
	pieces := []TokenScore{
		{Token: "<unk>", Score: 0.0},
		{Token: "<0xC3>", Score: -0.01},
		{Token: "<0xA9>", Score: -0.03},
	}
	params := util.NewParams(map[string]interface{}{
		"unk_id":        0,
		"byte_fallback": true,
	})
	model, err := New(pieces, params)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	// Test tokenization of "é" (UTF-8: C3 A9)
	tokens, err := model.Tokenize("é")
	if err != nil {
		t.Fatalf("Failed to tokenize: %v", err)
	}
	if got, want := len(tokens), 2; got != want {
		t.Errorf("Wrong number of tokens: got %d, want %d", got, want)
	}
	if got, want := tokens[0].Value, "<0xC3>"; got != want {
		t.Errorf("Wrong first token value: got %q, want %q", got, want)
	}
	if got, want := tokens[1].Value, "<0xA9>"; got != want {
		t.Errorf("Wrong second token value: got %q, want %q", got, want)
	}

	// Test tokenization of "?é"
	tokens, err = model.Tokenize("?é")
	if err != nil {
		t.Fatalf("Failed to tokenize: %v", err)
	}
	if got, want := len(tokens), 3; got != want {
		t.Errorf("Wrong number of tokens: got %d, want %d", got, want)
	}
	if got, want := tokens[0].Value, "<unk>"; got != want {
		t.Errorf("Wrong first token value: got %q, want %q", got, want)
	}
}
