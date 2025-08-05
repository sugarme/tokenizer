package main

import (
	"fmt"
	"log"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/model/unigram"
	"github.com/sugarme/tokenizer/util"
)

func main() {
	// Create a simple Unigram model for testing
	vocab := []unigram.TokenScore{
		{Token: "hello", Score: -1.0},
		{Token: "world", Score: -2.0},
		{Token: "!", Score: -3.0},
		{Token: "how", Score: -4.0},
		{Token: "are", Score: -5.0},
		{Token: "you", Score: -6.0},
		{Token: "today", Score: -7.0},
		{Token: "?", Score: -8.0},
	}

	// Set unk_id to 0 (the first token)
	unkID := 0
	opts := util.NewParams(nil)
	opts.Set("unk_id", unkID)
	opts.Set("byte_fallback", false)

	// Create the Unigram model
	model, err := unigram.New(vocab, opts)
	if err != nil {
		log.Fatalf("Failed to create Unigram model: %v", err)
	}

	// Test tokenization
	testTokenization(model, "hello world!")
	testTokenization(model, "how are you today?")
	testTokenization(model, "unknown token")
}

func testTokenization(model tokenizer.Model, text string) {
	fmt.Printf("Tokenizing: %q\n", text)
	
	tokens, err := model.Tokenize(text)
	if err != nil {
		fmt.Printf("Error tokenizing: %v\n", err)
		return
	}
	
	fmt.Printf("Tokens (%d):\n", len(tokens))
	for i, token := range tokens {
		fmt.Printf("  %d. ID=%d, Value=%q, Offsets=%v\n", i, token.Id, token.Value, token.Offsets)
	}
	fmt.Println()
}
