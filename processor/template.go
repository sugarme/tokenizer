package processor

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/util"
)

type SequenceEnum int

const (
	A SequenceEnum = iota
	B
)

type Piece interface {
	// ExtractId(s string) Piece
	WithTypeId(typeId int)
}

type SequencePiece struct {
	Id     SequenceEnum `json:"id"`
	TypeId int          `json:"type_id"`
}

var _ Piece = new(SequencePiece)

type SpecialTokenPiece struct {
	Id     string `json:"id"`
	TypeId int    `json:"type_id"`
}

var _ Piece = new(SpecialTokenPiece)

func extractId(s string) (Piece, error) {
	var p Piece
	if strings.HasPrefix(s, "$") {
		rest := strings.TrimPrefix(s, "$")

		var isNum bool
		num, err := strconv.Atoi(rest)
		if err == nil {
			isNum = true
		}

		switch {
		case rest == "", rest == "A", rest == "a":
			p = &SequencePiece{
				Id:     A,
				TypeId: 0,
			}
		case rest == "B", rest == "b":
			p = &SequencePiece{
				Id:     B,
				TypeId: 0,
			}

		case isNum:
			p = &SequencePiece{
				Id:     A,
				TypeId: num,
			}

		default:
			err := fmt.Errorf("Cannot extract Id from input %q\n", s)
			return nil, err
		}
	} else {
		p = &SpecialTokenPiece{
			Id:     s,
			TypeId: 0,
		}
	}

	return p, nil
}

func NewSequencePiece(id string, typeId int) *SequencePiece {
	var seqEnum SequenceEnum
	if id == "A" {
		seqEnum = A
	} else {
		seqEnum = B
	}
	return &SequencePiece{
		Id:     seqEnum,
		TypeId: typeId,
	}
}

func NewSpecialTokenPiece(id string, typeId int) *SpecialTokenPiece {
	return &SpecialTokenPiece{
		Id:     id,
		TypeId: typeId,
	}
}

// Implement Piece for SequencePiece:
// ----------------------------------
func (p *SequencePiece) WithTypeId(v int) {
	p.TypeId = v
}

func (p *SpecialTokenPiece) WithTypeId(v int) {
	p.TypeId = v
}

func NewPiece(s string) (Piece, error) {
	parts := strings.Split(s, ":")

	var (
		p   Piece
		err error
	)
	switch len(parts) {
	case 2:
		typeId, err := strconv.Atoi(parts[1])
		if err != nil {
			err = fmt.Errorf("Cannot build Piece from string %q", s)
			return nil, err
		}

		p, err = extractId(parts[0])
		if err != nil {
			err = fmt.Errorf("%w. Cannot build Piece from string %q", err, s)
			return nil, err
		}

		p.WithTypeId(typeId)

	case 1:
		p, err = extractId(parts[0])
		if err != nil {
			err = fmt.Errorf("%w. Cannot build Piece from string %q", err, s)
			return nil, err
		}

	default:
		err = fmt.Errorf("Cannot build Piece from string %q", s)
		return nil, err
	}

	return p, nil
}

// Represents a bunch of tokens to be used in a template.
// Usually, special tokens have only one associated id/token but in
// some cases, it might be interesting to have multiple ids/tokens.
type SpecialToken struct {
	// A unique id used to identify this SpecialToken in the template
	Id string

	// The list of associated ids
	Ids []int

	// The list of associated tokens
	Tokens []string
}

func NewSpecialToken(id string, ids []int, tokens []string) *SpecialToken {
	return &SpecialToken{
		Id:     id,
		Ids:    ids,
		Tokens: tokens,
	}
}

func NewSpecialTokenFrom(s string, id int) *SpecialToken {
	return NewSpecialToken(s, []int{id}, []string{s})
}

type Template []Piece

func NewTemplateFromOne(s string) (Template, error) {
	parts := strings.Split(s, " ")

	return NewTemplateFromMulti(parts)
}

func NewTemplateFromMulti(parts []string) (Template, error) {
	var tpl []Piece
	for _, part := range parts {
		p, err := NewPiece(part)
		if err != nil {
			return nil, err
		}
		tpl = append(tpl, p)
	}

	return tpl, nil
}

func NewTemplate(v interface{}) (Template, error) {
	switch typ := v.(type) {
	case string:
		return NewTemplateFromOne(v.(string))
	case []string:
		return NewTemplateFromMulti(v.([]string))
	default:
		err := fmt.Errorf("Unsupported input type %v\n", typ)
		return nil, err
	}
}

// A bunch of [`SpecialToken`] represented by their ID.
type Tokens struct {
	TokenMap    map[string]SpecialToken // NOTE. HF is an ordered map
	orderedKeys []string                // order of the TokenMap
}

func DefaultTokens() *Tokens {
	return &Tokens{
		TokenMap:    make(map[string]SpecialToken),
		orderedKeys: nil,
	}
}

func NewTokensFrom(toks []SpecialToken) *Tokens {
	m := make(map[string]SpecialToken)
	var keys []string
	for _, tok := range toks {
		keys = append(keys, tok.Id)
		m[tok.Id] = tok
	}

	return &Tokens{
		TokenMap:    m,
		orderedKeys: keys,
	}
}

func NewTokensFromMap(m map[string]SpecialToken) *Tokens {
	// TODO. How to sort this map to get ordered map?
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	return &Tokens{
		TokenMap:    m,
		orderedKeys: keys,
	}
}

func NewTokens(toks []tokenizer.Token) *Tokens {
	m := make(map[string]SpecialToken)
	var keys []string
	for _, tok := range toks {
		spt := NewSpecialTokenFrom(tok.Value, tok.Id)
		keys = append(keys, tok.Value)
		m[tok.Value] = *spt
	}

	return &Tokens{
		TokenMap:    m,
		orderedKeys: keys,
	}
}

func (t *Tokens) GetItemByOrder(index int) (SpecialToken, bool) {
	k := t.orderedKeys[index]

	return t.GetItemByKey(k)
}

func (t *Tokens) GetItemByKey(id string) (SpecialToken, bool) {
	val, ok := t.TokenMap[id]
	return val, ok
}

// / This PostProcessor takes care of processing each input `Encoding` by applying
// / the corresponding template, before merging them in the final Encoding.
// /
// / A `Template` is actually a sequence of `Piece` that will be
// / concatenated together in the given order. Each `Piece` represents either
// / one of the input `Encoding` or a `SpecialToken`.
// /
// / ## Example
// / ```
// / # use tokenizers::processors::template::TemplateProcessing;
// / let template = TemplateProcessing::builder()
// /     .try_single("[CLS] $A [SEP]").unwrap()
// /     .try_pair("[CLS] $A [SEP] $B:1 [SEP]:1").unwrap()
// /     .special_tokens(vec![("[CLS]", 1), ("[SEP]", 0)])
// /     .build()
// /     .unwrap();
// / ```
// /
type TemplateProcessing struct {
	Single        Template
	Pair          Template
	AddedSingle   int
	AddedPair     int
	SpecialTokens *Tokens
}

type TemplateProcessingDeserializer struct {
	Single        Template
	Pair          Template
	SpecialTokens *Tokens
}

func NewTemplateProcessingFrom(t *TemplateProcessingDeserializer) *TemplateProcessing {
	addedSingle := countAdded(t.Single, t.SpecialTokens)
	addedPair := countAdded(t.Pair, t.SpecialTokens)

	return &TemplateProcessing{
		Single:        t.Single,
		Pair:          t.Pair,
		AddedSingle:   addedSingle,
		AddedPair:     addedPair,
		SpecialTokens: t.SpecialTokens,
	}
}

func NewTemplateProcessing(single, pair Template, specialTokens *Tokens) *TemplateProcessing {
	return NewTemplateProcessingFrom(&TemplateProcessingDeserializer{
		Single:        single,
		Pair:          pair,
		SpecialTokens: specialTokens,
	})
}

// Count the number of added tokens in the given template
func countAdded(container Template, specialTokens *Tokens) int {
	var count int
	for _, p := range container {
		typ := getType(p)
		switch typ {
		case "*SequencePiece":
			count += 0
		case "*SpecialTokenPiece":
			spt := p.(*SpecialTokenPiece)
			id := spt.Id
			specialToken, ok := specialTokens.GetItemByKey(id)
			if ok {
				count += len(specialToken.Ids)
			}
		default:
			msg := fmt.Sprintf("Unsupported typ %q for 'specialTokens' item\n", typ)
			panic(msg)
		}
	}

	return count
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

type TemplateProcessingBuilder struct {
	*TemplateProcessing
}

func (b *TemplateProcessingBuilder) updateAddedTokens() {
	b.AddedSingle = countAdded(b.Single, b.SpecialTokens)
	b.AddedPair = countAdded(b.Pair, b.SpecialTokens)
}

func (b *TemplateProcessingBuilder) NewSingle(v interface{}) {
	tpl, err := NewTemplate(v)
	if err != nil {
		panic("err")
	}

	b.Single = tpl
	b.updateAddedTokens()
}

func (b *TemplateProcessingBuilder) NewPair(v interface{}) {
	tpl, err := NewTemplate(v)
	if err != nil {
		panic("err")
	}

	b.Pair = tpl
	b.updateAddedTokens()
}

func (b *TemplateProcessingBuilder) NewSpecialTokens(tokens []tokenizer.Token) {
	b.SpecialTokens = NewTokens(tokens)
	b.updateAddedTokens()
}

func (b *TemplateProcessingBuilder) DefaultAdded(isSingle bool) int {
	var t Template
	if isSingle {
		t = b.Single
	} else {
		t = b.Pair
	}

	return countAdded(t, b.SpecialTokens)
}

func (b *TemplateProcessingBuilder) Validate() error {
	var pairHasBoth bool = true
	var (
		hasA bool
		hasB bool
	)

	for _, piece := range b.Pair {
		if piece.(*SequencePiece).Id == A {
			hasA = true
		}

		if piece.(*SequencePiece).Id == B {
			hasB = true
		}
	}

	pairHasBoth = hasA && hasB

	if !pairHasBoth {
		err := fmt.Errorf("Template for 'pair' must use both sequences.")
		return err
	}

	check := func(sp string) string {
		var exist bool
		tok, ok := b.SpecialTokens.GetItemByOrder(0)
		if !ok {
			exist = false
		} else {
			exist = util.Contains(tok.Tokens, sp)
		}

		if exist {
			return sp
		} else {
			return ""
		}
	}

	var missing []string
	var pieces []Piece = append(b.Single, b.Pair...)
	for _, p := range pieces {
		typ := getType(p)
		switch typ {
		case "*SequencePeice":
			// None
		case "*SpecialToken":
			id := p.(*SpecialTokenPiece).Id
			s := check(id)
			if s != "" {
				missing = append(missing, s)
			}
		}
	}

	if len(missing) > 0 {
		var msg string
		for _, s := range missing {
			v := fmt.Sprintf("Missing SpecialToken %q", s)
			msg = fmt.Sprintf("%s, %s", msg, v)
		}

		return fmt.Errorf(msg)
	}

	return nil
}

func DefaultTemplateProcessing() *TemplateProcessing {
	single, err := NewTemplateFromOne("$0")
	if err != nil {
		panic(err)
	}

	pair, err := NewTemplateFromOne("$1")
	if err != nil {
		panic(err)
	}

	specialTokens := DefaultTokens()

	return &TemplateProcessing{
		Single:        single,
		Pair:          pair,
		AddedSingle:   0,
		AddedPair:     0,
		SpecialTokens: specialTokens,
	}
}

func (tp *TemplateProcessing) Builder() *TemplateProcessingBuilder {
	return &TemplateProcessingBuilder{tp}
}

func (tp *TemplateProcessingBuilder) Build() *TemplateProcessing {
	return tp.TemplateProcessing
}

func (tp *TemplateProcessing) ApplyTemplate(template []Piece, encodings []tokenizer.Encoding, addSpecialTokens bool) []tokenizer.Encoding {
	var finalEncodings []tokenizer.Encoding

	for _, piece := range template {
		typ := getType(piece)

		switch typ {
		case "*SequencePiece":
			sp := piece.(*SequencePiece)
			id := sp.Id
			typeId := sp.TypeId
			i := 0
			if id != A {
				i = 1
			}
			encoding := encodings[id]
			typeIds := util.Repeat(typeId, encoding.Len())
			encoding.SetTypeIds(typeIds)
			encoding.SetSequenceIds(i)

			finalEncodings = append(finalEncodings, encoding)

		case "*SpecialTokenPiece":
			spt := piece.(*SpecialTokenPiece)
			id := spt.Id
			typeId := spt.TypeId
			if addSpecialTokens {
				tok, ok := tp.SpecialTokens.GetItemByKey(id)
				if !ok {
					msg := fmt.Sprintf("Token not found with key %q", id)
					panic(msg)
				}
				length := len(tok.Ids)

				ids := tok.Ids
				typeIds := util.Repeat(typeId, length)
				tokens := tok.Tokens
				offsets := util.Repeat([]int{0, 0}, length)
				specialTokenMask := util.Repeat(1, length)
				attentionMask := util.Repeat(1, length)
				var overflowing []tokenizer.Encoding = nil
				encoding := tokenizer.NewEncoding(ids, typeIds, tokens, offsets, specialTokenMask, attentionMask, overflowing)

				finalEncodings = append(finalEncodings, *encoding)
			}
		}
	}

	return finalEncodings
}

// Implement PostProcessor for TemplateProcessing:
// -----------------------------------------------

var _ tokenizer.PostProcessor = new(TemplateProcessing)

func (tp *TemplateProcessing) AddedTokens(isPair bool) int {
	if isPair {
		return tp.AddedPair
	}

	return tp.AddedSingle
}

func (tp *TemplateProcessing) Process(encoding, pairEncoding *tokenizer.Encoding, addSpecialTokens bool) *tokenizer.Encoding {
	encodings := tokenizer.PrepareEncodings(encoding, pairEncoding)
	var template Template
	switch len(encodings) {
	case 2:
		template = tp.Pair
	case 1:
		template = tp.Single
	default:
		panic("Shouldn't be here. 'encoding' must be != nil")
	}

	appliedEncodings := tp.ApplyTemplate(template, encodings, addSpecialTokens)

	return tokenizer.MergeEncodings(appliedEncodings, false)
}
