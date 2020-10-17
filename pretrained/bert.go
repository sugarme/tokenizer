package pretrained

import (
	"log"

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
// It uses special tokens: `[UNK]`, `[SEP]`, `[CLS]` and `[MASK]`
// Its normalizer configued with: clean text, lower-case, handle Chinese characters and
// strip accents.
func BertBaseUncased() *tokenizer.Tokenizer {

	util.CdToThis()

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
