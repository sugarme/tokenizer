package pretrained

import (
	"log"
	"os"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/decoder"
	"github.com/sugarme/tokenizer/model/wordpiece"
	"github.com/sugarme/tokenizer/normalizer"
	"github.com/sugarme/tokenizer/pretokenizer"
	"github.com/sugarme/tokenizer/processor"
	"github.com/sugarme/tokenizer/util"
)

// BertBaseUncase loads pretrained BERT tokenizer.
//
// Special tokens:
// - unknown token: "[UNK]"
// - sep token: "[SEP]"
// - cls token: "[CLS]"
// - mask token: "[MASK]"
// Its normalizer configued with: clean text, lower-case, handle Chinese characters and
// strip accents.
//
// Source:
// "https://cdn.huggingface.co/bert-base-uncased-vocab.txt"
func BertBaseUncased() *tokenizer.Tokenizer {

	currDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	util.CdToThis()
	defer util.CdBack(currDir)

	vocabFile := "model/bert-base-uncased-vocab.txt"
	model, err := wordpiece.NewWordPieceFromFile(vocabFile, "[UNK]")
	if err != nil {
		log.Fatal(err)
	}

	tk := tokenizer.NewTokenizer(model)

	bertNormalizer := normalizer.NewBertNormalizer(true, true, true, true)
	tk.WithNormalizer(bertNormalizer)

	bertPreTokenizer := pretokenizer.NewBertPreTokenizer()
	tk.WithPreTokenizer(bertPreTokenizer)

	sepId, ok := tk.TokenToId("[SEP]")
	if !ok {
		log.Fatalf("Cannot find ID for [SEP] token.\n")
	}
	sep := processor.PostToken{Id: sepId, Value: "[SEP]"}

	clsId, ok := tk.TokenToId("[CLS]")
	if !ok {
		log.Fatalf("Cannot find ID for [CLS] token.\n")
	}
	cls := processor.PostToken{Id: clsId, Value: "[CLS]"}

	postProcess := processor.NewBertProcessing(sep, cls)
	tk.WithPostProcessor(postProcess)

	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("[MASK]", true)})
	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("[SEP]", true)})
	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("[CLS]", true)})

	wpDecoder := decoder.DefaultWordpieceDecoder()
	tk.WithDecoder(wpDecoder)

	return tk
}

// BertLargeCasedWholeWordMaskingSquad loads pretrained BERT large case whole-word masking tokenizer
// finetuned on SQuAD dataset.
//
// Source:
// https://cdn.huggingface.co/bert-large-cased-whole-word-masking-finetuned-squad-vocab.txt
func BertLargeCasedWholeWordMaskingSquad() *tokenizer.Tokenizer {

	util.CdToThis()
	vocabFile := "model/bert-large-cased-whole-word-masking-finetuned-squad-vocab.txt"
	model, err := wordpiece.NewWordPieceFromFile(vocabFile, "[UNK]")
	if err != nil {
		log.Fatal(err)
	}

	tk := tokenizer.NewTokenizer(model)

	// NOTE. set lowercase to false for cased
	bertNormalizer := normalizer.NewBertNormalizer(true, false, true, true)
	tk.WithNormalizer(bertNormalizer)

	bertPreTokenizer := pretokenizer.NewBertPreTokenizer()
	tk.WithPreTokenizer(bertPreTokenizer)

	sepId, ok := tk.TokenToId("[SEP]")
	if !ok {
		log.Fatalf("Cannot find ID for [SEP] token.\n")
	}
	sep := processor.PostToken{Id: sepId, Value: "[SEP]"}

	clsId, ok := tk.TokenToId("[CLS]")
	if !ok {
		log.Fatalf("Cannot find ID for [CLS] token.\n")
	}
	cls := processor.PostToken{Id: clsId, Value: "[CLS]"}

	postProcess := processor.NewBertProcessing(sep, cls)
	tk.WithPostProcessor(postProcess)

	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("[MASK]", true)})
	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("[SEP]", true)})
	tk.AddSpecialTokens([]tokenizer.AddedToken{tokenizer.NewAddedToken("[CLS]", true)})

	wpDecoder := decoder.NewWordPieceDecoder("##", true)
	tk.WithDecoder(wpDecoder)

	return tk
}
