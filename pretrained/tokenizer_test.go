package pretrained

import (
	"log"
	"reflect"
	"testing"

	"github.com/sugarme/tokenizer"
)

func TestFromFile(t *testing.T) {
	configFile, err := tokenizer.CachedPath("hf-internal-testing/llama-tokenizer", "tokenizer.json")
	if err != nil {
		panic(err)
	}

	tk, err := FromFile(configFile)
	if err != nil {
		panic(err)
	}

	log.Printf("tokenizer: %#v\n", tk)

	log.Printf("vocab size (excl. added tokens): %v\n", tk.GetVocabSize(false))
	log.Printf("vocab size (incl. added tokens): %v\n", tk.GetVocabSize(true))
	log.Printf("special tokens: %q\n", tk.GetSpecialTokens())
}

/*
{
  "version": "1.0",
  "truncation": null,
  "padding": null,
  "added_tokens": [
    {
      "id": 0,
      "content": "<unk>",
      "single_word": false,
      "lstrip": false,
      "rstrip": false,
      "normalized": true,
      "special": true
    },
    {
      "id": 1,
      "content": "<s>",
      "single_word": false,
      "lstrip": false,
      "rstrip": false,
      "normalized": true,
      "special": true
    },
    {
      "id": 2,
      "content": "</s>",
      "single_word": false,
      "lstrip": false,
      "rstrip": false,
      "normalized": true,
      "special": true
    }
  ],
  "normalizer": {
    "type": "Sequence",
    "normalizers": [
      {
        "type": "Prepend",
        "prepend": "▁"
      },
      {
        "type": "Replace",
        "pattern": {
          "String": " "
        },
        "content": "▁"
      }
    ]
  },
  "pre_tokenizer": null,
  "post_processor": {
    "type": "TemplateProcessing",
    "single": [
      {
        "SpecialToken": {
          "id": "<s>",
          "type_id": 0
        }
      },
      {
        "Sequence": {
          "id": "A",
          "type_id": 0
        }
      }
    ],
    "pair": [
      {
        "SpecialToken": {
          "id": "<s>",
          "type_id": 0
        }
      },
      {
        "Sequence": {
          "id": "A",
          "type_id": 0
        }
      },
      {
        "SpecialToken": {
          "id": "<s>",
          "type_id": 1
        }
      },
      {
        "Sequence": {
          "id": "B",
          "type_id": 1
        }
      }
    ],
    "special_tokens": {
      "<s>": {
        "id": "<s>",
        "ids": [
          1
        ],
        "tokens": [
          "<s>"
        ]
      }
    }
  },
  "decoder": {
    "type": "Sequence",
    "decoders": [
      {
        "type": "Replace",
        "pattern": {
          "String": "▁"
        },
        "content": " "
      },
      {
        "type": "ByteFallback"
      },
      {
        "type": "Fuse"
      },
      {
        "type": "Strip",
        "content": " ",
        "start": 1,
        "stop": 0
      }
    ]
  },
  "model": {
    "type": "BPE",
    "dropout": null,
    "unk_token": "<unk>",
    "continuing_subword_prefix": null,
    "end_of_word_suffix": null,
    "fuse_unk": true,
    "byte_fallback": true,
    "vocab": {
      "<unk>": 0,
      "<s>": 1,
      "</s>": 2,
      ...
    }
  }
}

*/

// Test Bert model `bert-base-uncased`
func TestBert(t *testing.T) {
	configFile, err := tokenizer.CachedPath("bert-base-uncased", "tokenizer.json")
	if err != nil {
		panic(err)
	}

	tk, err := FromFile(configFile)
	if err != nil {
		panic(err)
	}

	sentence := `The Gophers craft code using [MASK] language.`
	en, err := tk.EncodeSingle(sentence, true)
	if err != nil {
		log.Fatal(err)
	}

	gotTokens := en.Tokens
	wantTokens := []string{
		"[CLS]", "the", "go", "##pher", "##s", "craft", "code", "using", "[MASK]", "language", ".", "[SEP]",
	}
	if !reflect.DeepEqual(wantTokens, gotTokens) {
		t.Errorf("want %v,\ngot %v\n", wantTokens, gotTokens)
	}

	gotIds := en.Ids
	wantIds := []int{101, 1996, 2175, 27921, 2015, 7477, 3642, 2478, 103, 2653, 1012, 102}
	if !reflect.DeepEqual(wantIds, gotIds) {
		t.Errorf("want %v,\ngot %v\n", wantIds, gotIds)
	}

	gotOffsets := en.Offsets
	wantOffsets := [][]int{
		{0, 0},
		{0, 3},
		{4, 6},
		{6, 10},
		{10, 11},
		{12, 17},
		{18, 22},
		{23, 28},
		{29, 35},
		{36, 44},
		{44, 45},
		{0, 0},
	}
	if !reflect.DeepEqual(wantOffsets, gotOffsets) {
		t.Errorf("want %v,\ngot %v\n", wantOffsets, gotOffsets)
	}
}

// Test Roberta model `roberta-base`
func TestRoberta(t *testing.T) {
	configFile, err := tokenizer.CachedPath("roberta-base", "tokenizer.json")
	if err != nil {
		panic(err)
	}

	tk, err := FromFile(configFile)
	if err != nil {
		panic(err)
	}

	sentence := `The Gophers craft code using [MASK] language.`
	en, err := tk.EncodeSingle(sentence, true)
	if err != nil {
		log.Fatal(err)
	}
	/*
		sequence: The Gophers craft code using [MASK] language.
		tokens: ['<s>', 'The', 'ĠG', 'ophers', 'Ġcraft', 'Ġcode', 'Ġusing', 'Ġ[', 'MAS', 'K', ']', 'Ġlanguage', '.', '</s>']
		ids: [0, 133, 272, 30482, 6306, 3260, 634, 646, 32804, 530, 742, 2777, 4, 2]
		offsets: [(0, 0), (0, 3), (4, 5), (5, 11), (12, 17), (18, 22), (23, 28), (29, 30), (30, 33), (33, 34), (34, 35), (36, 44), (44, 45), (0, 0)]
	*/

	gotTokens := en.Tokens
	wantTokens := []string{
		"<s>", "The", "ĠG", "ophers", "Ġcraft", "Ġcode", "Ġusing", "Ġ[", "MAS", "K", "]", "Ġlanguage", ".", "</s>",
	}
	if !reflect.DeepEqual(wantTokens, gotTokens) {
		t.Errorf("\nwant %q,\ngot  %q\n", wantTokens, gotTokens)
	}

	gotIds := en.Ids
	wantIds := []int{0, 133, 272, 30482, 6306, 3260, 634, 646, 32804, 530, 742, 2777, 4, 2}
	if !reflect.DeepEqual(wantIds, gotIds) {
		t.Errorf("\nwant %v,\ngot  %v\n", wantIds, gotIds)
	}

	gotOffsets := en.Offsets
	wantOffsets := [][]int{
		{0, 0},
		{0, 3},
		{4, 5},
		{5, 11},
		{12, 17},
		{18, 22},
		{23, 28},
		{29, 30},
		{30, 33},
		{33, 34},
		{34, 35},
		{36, 44},
		{44, 45},
		{0, 0},
	}
	if !reflect.DeepEqual(wantOffsets, gotOffsets) {
		t.Errorf("\nwant %v,\ngot  %v\n", wantOffsets, gotOffsets)
	}
}

// Test Llama model `hf-internal-testing/llama-tokenizer`
func TestLlama(t *testing.T) {
	configFile, err := tokenizer.CachedPath("hf-internal-testing/llama-tokenizer", "tokenizer.json")
	if err != nil {
		panic(err)
	}

	tk, err := FromFile(configFile)
	if err != nil {
		panic(err)
	}

	sentence := `The Gophers craft code using [MASK] language.`
	en, err := tk.EncodeSingle(sentence, true)
	if err != nil {
		log.Fatal(err)
	}
	/*
		tokens: ['<s>', '▁The', '▁G', 'oph', 'ers', '▁craft', '▁code', '▁using', '▁[', 'MA', 'SK', ']', '▁language', '.']
		ids: [1, 450, 402, 3021, 414, 25554, 775, 773, 518, 1529, 16033, 29962, 4086, 29889]
		offsets: [(0, 0), (0, 3), (3, 5), (5, 8), (8, 11), (11, 17), (17, 22), (22, 28), (28, 30), (30, 32), (32, 34), (34, 35), (35, 44), (44, 45)]
	*/

	gotTokens := en.Tokens
	wantTokens := []string{
		"<s>", "▁The", "▁G", "oph", "ers", "▁craft", "▁code", "▁using", "▁[", "MA", "SK", "]", "▁language", ".",
	}
	if !reflect.DeepEqual(wantTokens, gotTokens) {
		t.Errorf("\nwant %q,\ngot  %q\n", wantTokens, gotTokens)
	}

	gotIds := en.Ids
	wantIds := []int{
		1, 450, 402, 3021, 414, 25554, 775, 773, 518, 1529, 16033, 29962, 4086, 29889,
	}
	if !reflect.DeepEqual(wantIds, gotIds) {
		t.Errorf("\nwant %v,\ngot  %v\n", wantIds, gotIds)
	}

	gotOffsets := en.Offsets
	wantOffsets := [][]int{
		{0, 0},
		{0, 3},
		{3, 5},
		{5, 8},
		{8, 11},
		{11, 17},
		{17, 22},
		{22, 28},
		{28, 30},
		{30, 32},
		{32, 34},
		{34, 35},
		{35, 44},
		{44, 45},
	}
	if !reflect.DeepEqual(wantOffsets, gotOffsets) {
		t.Errorf("\nwant %v,\ngot  %v\n", wantOffsets, gotOffsets)
	}
}

func TestLlamaDecode(t *testing.T) {
	configFile, err := tokenizer.CachedPath("HuggingFaceH4/tiny-random-LlamaForCausalLM", "tokenizer.json")
	if err != nil {
		panic(err)
	}

	tk, err := FromFile(configFile)
	if err != nil {
		panic(err)
	}

	generateIds := []int{
		1, 18637, 29892, 526, 366, 1136, 455, 2470, 29973, 1815,
		366, 5193, 304, 592, 29973, 11240, 5919, 12794, 18055, 19519,
		21315, 22401, 31302, 799, 127, 22164, 22410, 4839, 213, 23512,
	}

	got := tk.Decode(generateIds, true)
	want := "Hey, are you consciours? Can you talk to me? Monte:(illon entering MY рос championship提ear|)):compos header� тя"

	if !reflect.DeepEqual(want, got) {
		t.Errorf("\nwant %q\ngot  %q\n", want, got)
	}
}
