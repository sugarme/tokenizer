package tokenizer

// tokenizer represents a tokenization pipeline
// TODO: full description

import (
	// "bufio"
	// "context"
	// "fmt"
	"log"
	// "math"
	// "os"
	"reflect"
	"strings"
	// "regexp"
	"sync"

	// progressbar "github.com/schollz/progressbar/v2"
	// "golang.org/x/sync/errgroup"

	"github.com/sugarme/tokenizer/normalizer"
	// "github.com/sugarme/tokenizer/util"
)

const mb = 1024 * 1024
const gb = 1024 * mb

type Offsets struct {
	Start int
	End   int
}

type PreToken struct {
	Value   string
	Offsets Offsets
}

type Token struct {
	Id      int
	Value   string
	Offsets Offsets
}

// PreTokenizer processes strings before going to the model
// It splits the given string into multiple substrings and keeps track
// of offsets of split substrings from the `NormalizedString`. In some
// occasion, the `PreTokenizer` might need to modify the given `NormalizedString`
// to ensure it entirely keeps track of the offsets and the mapping with
// the original string.
type PreTokenizer interface {
	// PreTokenize(pretokenized PreTokenizedString) (retVal []PreToken)
	PreTokenize(pretokenized PreTokenizedString) (retVal PreTokenizedString, err error)
}

// Model represents a model used during tokenization (i.e., BPE, Word, or Unigram)
type Model interface {
	// Tokenize tokenizes the given sequence into multiple underlying `Token`
	// The `offsets` on the `Token` are expected to be relative to the given
	// sequence
	// Tokenize(tokens []PreToken) ([]Token, error)
	Tokenize(sequence string) ([]Token, error)
	// TokenToId finds the ID associated with a string token
	TokenToId(token string) (id int, ok bool)
	// IdToToken find the string token associated with an ID
	IdToToken(id int) (token string, ok bool)
	// GetVocab retrieves the entire vocabulary mapping (token -> Id)
	GetVocab() map[string]int
	// GetVocabSize retrieves the entire vocabulary mapping(map[token]id)
	GetVocabSize() int
	// Save saves the current `Model` in the given folder, using the
	// given `prefixOpt` for various files that need to be saved.
	Save(path string, prefixOpt ...string) error
}

// PostProcessor is in charge of post-processing an encoded output of
// the `Tokenizer`.
// It adds any special tokens that a language model would require.
type PostProcessor interface {
	// AddedTokens returns the number of tokens that will be added during the processing step
	AddedTokens(isPair bool) int
	// Process processes both encodings and returns a new merged one
	// NOTE: pairEncoding is optional
	Process(encoding, pairEncoding *Encoding, addSpecialTokens bool) *Encoding
}

// DefaultProcess is a helper function of PostProcessor's Process method
// It helps to fast track by just merging encoding and its pair.
func DefaultProcess(encoding, pairEncoding *Encoding, addSpecialTokens bool) *Encoding {
	if pairEncoding != nil {
		return encoding.MergeWith(pairEncoding, false)
	}

	return encoding
}

// Decoder takes care of (merges) the given slice of tokens to string
type Decoder interface {
	Decode(tokens []string) string
}

// Trainer is responsible for training a model. It takes lines/sentences
// and returns a tokenizer `Model` when done.
type Trainer interface {
	// Whether showing progress bar or not
	WithProgressBar() bool
	// Actual training method. It will return a trained model and
	// a list of `special tokens` to be added directly to the tokenizer
	// along with the model
	Train(words map[string]int) (Model, []AddedToken)
	// ProcessTokens processes a bunch of tokens and counts them as relevant
	ProcessTokens(words map[string]int, tokens []string)
}

// Implement methods for `Token`
// NewToken generate new token from input data
func NewToken(id int, value string, offsets Offsets) Token {
	return Token{
		Id:      id,
		Value:   value,
		Offsets: offsets,
	}
}

// InputSequence :
// ===============

type InputType int

const (
	RawInput = iota
	PretokenizedInput
)

type InputSequence struct {
	input     []string
	inputType InputType
}

// NewInputSequence creates a new InputSequence from input
// A valid input can be a string type (RawInput) or slice of string (PretokenizedInput)
func NewInputSequence(input interface{}) (retVal InputSequence) {

	switch reflect.TypeOf(input).Name() {
	case "string":
		return InputSequence{
			input:     []string{input.(string)},
			inputType: RawInput,
		}
	case "slice":
		if reflect.TypeOf(input).Elem().Name() != "string" {
			log.Fatalf("Invalid input type: %v. Expect type of 'string' or '[]string'\n", reflect.TypeOf(input).Name())
		}
		return InputSequence{
			input:     input.([]string),
			inputType: PretokenizedInput,
		}
	default:
		log.Fatalf("Invalid input type: %v. Expect type of 'string' or '[]string'\n", reflect.TypeOf(input).Name())
	}

	return
}

type Single struct {
	Sentence InputSequence
}
type Dual struct {
	Sentence InputSequence
	Pair     InputSequence
}

type EncodeInput interface{}

func NewSingleEncodeInput(sentence InputSequence) (retVal EncodeInput) {
	return Single{sentence}
}

func NewDualEncodeInput(sentence, pairSentence InputSequence) (retVal EncodeInput) {
	return Dual{sentence, pairSentence}
}

// Tokenizer represents a tokenization pipeline.
// It can implement any encoding or decoding of any text.
type Tokenizer struct {
	// Parts
	normalizer    *normalizer.Normalizer // optional
	preTokenizer  *PreTokenizer          // optional
	model         *Model
	postProcessor *PostProcessor // optional
	decoder       *Decoder       // optional

	// Added vocabulary capability
	addedVocabulary AddedVocabulary

	// General processing parameters
	trunc   *TruncationParams // optional
	padding *PaddingParams    // optional
}

// Implementing methods for Tokenizer
func NewTokenizer(model Model) Tokenizer {
	return Tokenizer{
		normalizer:      nil,
		preTokenizer:    nil,
		model:           &model,
		postProcessor:   nil,
		decoder:         nil,
		addedVocabulary: NewAddedVocabulary(),
		trunc:           nil,
		padding:         nil,
	}
}

func (t *Tokenizer) WithNormalizer(n normalizer.Normalizer) {
	t.normalizer = &n
}

func (t *Tokenizer) GetNormalizer() *normalizer.Normalizer {
	return t.normalizer
}

func (t *Tokenizer) WithPreTokenizer(preTokenizer PreTokenizer) {
	t.preTokenizer = &preTokenizer
}

func (t *Tokenizer) GetPreTokenizer() *PreTokenizer {
	return t.preTokenizer
}

func (t *Tokenizer) WithPostProcessor(postProcessor PostProcessor) {
	t.postProcessor = &postProcessor
}

func (t *Tokenizer) GetPostProcessor() *PostProcessor {
	return t.postProcessor
}

func (t *Tokenizer) WithDecoder(decoder *Decoder) {
	t.decoder = decoder
}

func (t *Tokenizer) GetDecoder() *Decoder {
	return t.decoder
}

func (t *Tokenizer) WithModel(model *Model) {
	t.model = model
}

func (t *Tokenizer) GetModel() *Model {
	return t.model
}

func (t *Tokenizer) WithTruncation(trunc *TruncationParams) {
	t.trunc = trunc
}

func (t *Tokenizer) GetTruncation() *TruncationParams {
	return t.trunc
}

func (t *Tokenizer) WithPadding(padding *PaddingParams) {
	t.padding = padding
}

func (t *Tokenizer) GetPadding() (retVal *PaddingParams) {
	return t.padding
}

// GetVocab get the vocabulary
func (t *Tokenizer) GetVocab(withAddedTokens bool) map[string]int {
	finalVocab := (*t.model).GetVocab()
	if withAddedTokens {
		addedVocab := t.addedVocabulary.GetVocab()
		if len(addedVocab) > 0 {
			for k, v := range addedVocab {
				finalVocab[k] = v
			}
		}
	}

	return finalVocab
}

// GetVocabSize get the size of vocabulary
func (t *Tokenizer) GetVocabSize(withAddedTokens bool) int {
	if !withAddedTokens {
		return (*t.model).GetVocabSize()
	}

	return (*t.model).GetVocabSize() + t.addedVocabulary.Len()
}

// TokenToId converts a token to a corresponding id
func (t *Tokenizer) TokenToId(token string) (id int, ok bool) {
	id, ok = t.addedVocabulary.TokenToId(token, *t.model)
	return id, ok
}

// IdToToken converts an Id to a corresponding token
func (t *Tokenizer) IdToToken(id int) (token string, ok bool) {
	token, ok = t.addedVocabulary.IdToToken(id, *t.model)
	return token, ok
}

// Normalize normalizes the given sentence and return the corresponding normalized string
func (t *Tokenizer) Normalize(sentence string) (retVal *normalizer.NormalizedString, err error) {

	var subs []*normalizer.NormalizedString
	isPairs := t.addedVocabulary.ExtractAndNormalize(sentence, t.normalizer)
	for _, isPair := range isPairs {
		if isPair.Id != -1 { // id is optional
			return isPair.Substring.Normalized, nil
		} else {
			// The PreTokenizers can still manipulate the normalized strings
			// so we do this anyway and will merge it back a NormalizedString
			preTok, err := t.doPreTokenize(sentence)
			if err != nil {
				return retVal, err
			}

			subs = append(subs, preTok.IntoMerged())
		}
	}

	retVal = subs[0]
	for i := 1; i < len(subs); i++ {
		retVal = retVal.MergeWith(subs[i])
	}

	return retVal, nil
}

// EncodeSingleSequence encodes a single sequence
func (t *Tokenizer) EncodeSingleSequence(sequence InputSequence, typeId int) (retVal *Encoding, err error) {

	var subseqEncodings []*Encoding

	for subseqIdx, subseq := range sequence.input {
		// isPairs is slice of pair (id, substring)
		isPairs := t.addedVocabulary.ExtractAndNormalize(subseq, t.normalizer)
		var encodings []*Encoding
		for _, isPair := range isPairs {
			if isPair.Id != -1 { // it's an added token, no need to tokenize. We have an ID
				encoding := NewEncodingFromTokens([]Token{
					NewToken(isPair.Id, isPair.Substring.Normalized.GetNormalized(), isPair.Substring.OriginalOffsets),
				}, typeId)

				encoding.SetWord(0, 0)
				// return encoding, nil
				encodings = append(encodings, encoding)

			} else { // let's tokenize
				preTok, err := t.doPreTokenize(isPair.Substring.Normalized.GetNormalized())
				if err != nil {
					return retVal, err
				}

				encoding, err := t.doTokenize(preTok, isPair.Substring.OriginalOffsets, typeId)
				if err != nil {
					return nil, err
				}

				encodings = append(encodings, encoding)
			}
		}

		// At this point, the `words` are good for each sub encoding independently,
		// but we need to make them grow sequentially.
		var subseqEncoding *Encoding
		subseqEncoding = encodings[0]
		for i := 1; i < len(encodings); i++ {
			e := encodings[i]
			wordIds := subseqEncoding.GetWords()
			lastWordId := wordIds[len(wordIds)-1]
			for _, id := range e.GetWords() {
				e.SetWord(id, id+lastWordId)
			}
			subseqEncoding = subseqEncoding.MergeWith(e, false)
		}

		// If we are handling already pre-tokenized input, each word should have the
		// relevant index from the given input, not determined by the pre-tokenization step
		if sequence.inputType == PretokenizedInput {
			wordIds := subseqEncoding.GetWords()

			for _, id := range wordIds {
				subseqEncoding.SetWord(id, subseqIdx)
			}
		}

		subseqEncodings = append(subseqEncodings, subseqEncoding)
	}

	retVal = DefaultEncoding()
	for _, e := range subseqEncodings {
		retVal = retVal.MergeWith(e, true) // TODO. double-check whether growing offsets???
	}

	return retVal, nil
}

// Encode the given input. This method accepts both single sequences, as well as pair
// sequences. Also, a sequence can be a string, or already pre-tokenized input directly:
func (t *Tokenizer) Encode(input EncodeInput, addSpecialTokens bool) (retVal *Encoding, err error) {
	var encoding, pairEncoding *Encoding

	// Encode and Postprocess
	switch reflect.TypeOf(input).Name() {
	case "Single":
		seq := input.(Single).Sentence
		encoding, err = t.EncodeSingleSequence(seq, 0)
		if err != nil {
			return retVal, err
		}

	case "Dual":
		seq := input.(Dual).Sentence
		encoding, err = t.EncodeSingleSequence(seq, 0)
		if err != nil {
			return retVal, err
		}
		pairSeq := input.(Dual).Pair
		pairEncoding, err = t.EncodeSingleSequence(pairSeq, 1)
		if err != nil {
			return retVal, err
		}

	default:
		log.Fatalf("Invalid input type - '%v'. \n", reflect.TypeOf(input).Name())
	}

	return t.PostProcess(encoding, pairEncoding, addSpecialTokens), nil
}

// Decode decodes the given ids, back to a String
func (t *Tokenizer) Decode(ids []int, skipSpecialTokens bool) (retVal string) {

	var tokens []string
	for _, id := range ids {
		if tok, ok := t.addedVocabulary.IdToToken(id, *t.model); ok {
			if !skipSpecialTokens || !t.addedVocabulary.IsSpecialToken(tok) {
				tokens = append(tokens, tok)
			}
		}
	}

	if t.decoder != nil {
		return (*t.decoder).Decode(tokens)
	}

	return strings.Join(tokens, " ")
}

// AddSpecialTokens registers the given tokens as special tokens. This is especially useful for removing
// these special tokens while decoding
func (t *Tokenizer) AddSpecialTokens(tokens []AddedToken) (retVal int) {
	return t.addedVocabulary.AddSpecialTokens(tokens, *t.model, t.normalizer)
}

// AddTokens adds the given tokens to the added vocabulary
func (t *Tokenizer) AddTokens(tokens []AddedToken) (retVal int) {
	return t.addedVocabulary.AddTokens(tokens, *t.model, t.normalizer)
}

// doNormalize does Normalization logic, go through all normalizers
func (t *Tokenizer) doNormalize(s string) (retVal *normalizer.NormalizedString, err error) {

	normalized := normalizer.NewNormalizedFrom(s)
	if t.normalizer != nil {
		normalized, err = (*t.normalizer).Normalize(normalized)
		if err != nil {
			return retVal, err
		}
	}

	return normalized, nil
}

// doPreTokenize does the pretokenization logic, handling the case where there is no PreTokenizer set
func (t *Tokenizer) doPreTokenize(sentence string) (retVal PreTokenizedString, err error) {

	pretokenized := NewPreTokenizedString(sentence)

	if t.preTokenizer != nil {
		pretokenized, err = (*t.preTokenizer).PreTokenize(pretokenized)
		if err != nil {
			return retVal, err
		}
	}

	return pretokenized, nil
}

// doTokenize does Tokenization logic, makes the bridge between the pre-tokenization phase and the real
// tokenization phase, and converting offsets back to the original referential.
func (t *Tokenizer) doTokenize(pretokenized PreTokenizedString, originalOffsets Offsets, typeId int) (retVal *Encoding, err error) {

	var substrings []SubString
	var encodings []*Encoding

	// Exclude all empty `normalized` substring
	for _, sub := range pretokenized.parts {
		if !sub.Normalized.IsEmpty() {
			substrings = append(substrings, sub)
		}
	}

	for wordIdx, substr := range substrings {
		tokens, err := (*t.model).Tokenize(substr.Normalized.GetNormalized())
		if err != nil {
			return nil, err
		}

		// We convert the normalized offsets back to the original
		for _, token := range tokens {
			oRange := substr.Normalized.ConvertOffset(normalizer.NewRange(token.Offsets.Start, token.Offsets.End, normalizer.NormalizedTarget))
			var convertedOffsets Offsets

			if oRange.Start() == -1 || oRange.End() == -1 {
				convertedOffsets = token.Offsets
			}

			convertedOffsets = Offsets{
				Start: originalOffsets.Start + substr.OriginalOffsets.Start + oRange.Start(),
				End:   originalOffsets.Start + substr.OriginalOffsets.Start + oRange.End(),
			}

			encoding := DefaultEncoding()
			encoding.Ids = []int{token.Id}
			encoding.TypeIds = []int{typeId}
			encoding.Tokens = []string{token.Value}
			encoding.Offsets = []Offsets{convertedOffsets}
			encoding.Words = []int{wordIdx}

			encodings = append(encodings, encoding)
		}
	}

	mergedEncoding := DefaultEncoding()
	retVal = mergedEncoding.Merge(encodings, false)

	return retVal, nil
}

// PostProcess does post-processing logic, handling the case where there is no PostProcessor set
func (t *Tokenizer) PostProcess(encoding, pairEncoding *Encoding, addSpecialTokens bool) (retVal *Encoding) {

	var tEncoding, tPairEncoding *Encoding

	// 1. Truncate if needed
	if t.trunc == nil {
		tEncoding, tPairEncoding = encoding, pairEncoding
	} else {
		trunc := t.trunc
		var nAddedTokens int = 0 // number of AddedToken
		if t.postProcessor != nil {
			processor := *t.postProcessor
			nAddedTokens = processor.AddedTokens(pairEncoding != nil)
		}

		if addSpecialTokens && nAddedTokens > 0 {
			params := trunc
			params.MaxLength = trunc.MaxLength - nAddedTokens

			tEncoding, tPairEncoding = TruncateEncodings(encoding, pairEncoding, params)
		} else {
			tEncoding, tPairEncoding = TruncateEncodings(encoding, pairEncoding, trunc)
		}
	}

	// 2. Post-process
	var finalEncoding *Encoding
	if t.postProcessor != nil {
		processor := *t.postProcessor
		finalEncoding = processor.Process(tEncoding, tPairEncoding, addSpecialTokens)
	} else {
		finalEncoding = DefaultProcess(tEncoding, tPairEncoding, addSpecialTokens)
	}

	// 3. Pad if needed
	if t.padding == nil {
		return finalEncoding
	}

	var padEncodings []*Encoding
	encodings := []*Encoding{finalEncoding}
	padEncodings = PadEncodings(encodings, *t.padding)
	if len(padEncodings) <= 1 {
		return finalEncoding
	} else {
		return padEncodings[0].Merge(padEncodings[1:], true)
	}
}

// EncodeBatch encodes all sentences in concurrency
func (t *Tokenizer) EncodeBatch(inputs []EncodeInput, addSpecialTokens bool) (retVal []*Encoding, err error) {
	var encodings []*Encoding
	var wg sync.WaitGroup

	wg.Add(len(inputs))

	// Encoding concurrently
	for i := 0; i < len(inputs); i++ {
		go func(i int) {
			defer wg.Done()

			e, err := t.Encode(inputs[i], addSpecialTokens)
			if err != nil {
				log.Fatal(err)
			}
			encodings = append(encodings, e)

		}(i)
	}

	wg.Wait()

	// Do padding if included
	if t.padding != nil {
		PadEncodings(encodings, *t.padding)
	}

	return encodings, nil
}

// DecodeBatch decodes all sentences in concurrency
func (t *Tokenizer) DecodeBatch(sentences [][]int, skipSpecialTokens bool) []string {
	var decodings []string
	var wg sync.WaitGroup

	wg.Add(len(sentences))

	// Decoding concurrently
	for i := 0; i < len(sentences); i++ {
		go func(i int) {
			defer wg.Done()

			s := t.Decode(sentences[i], skipSpecialTokens)
			decodings = append(decodings, s)

		}(i)
	}

	wg.Wait()

	return decodings
}

// wordCount returns a map of word and its count
func (t *Tokenizer) wordCount(trainer Model, files []string) (retVal map[string]int) {

	// TODO: implement
	return
}

// Train trains a model and return a new Tokenizer, using the given Trainer
func (t *Tokenizer) Train(trainer Model, files []string) (retVal *Tokenizer) {

	// TODO: implement

	return
}

// Train a model and replace our current Model, using the given Trainer
func (t *Tokenizer) TrainAndReplace(trainer Model, files []string) (err error) {

	// TODO: implement
	return
}

// NewTokenizerFromFile instantiates a new Tokenizer from the given file
func NewTokenizerFromFile(file string) (retVal *Tokenizer) {

	// TODO: implement
	return
}

// Serialize serializes current Tokenizer to string
func (t *Tokenizer) Serialize(pretty bool) (retVal string) {

	// TODO: implement
	return
}

// Save saves the current tokenizer at the given path
func (t *Tokenizer) Save(path string, pretty bool) (err error) {

	// TODO: implement
	return
}

/*
 * // Train trains a model and replaces the current model using a given trainer
 * // The tokenizer does the following steps
 * // 1. Concurrently, reads training data (text) from files, normalizes text using
 * // 		specified normalizer, and generates a slice of words and their frequency (count)
 * // 2. Train tokenizer model using specified tokenizer configuration on slice of word-count
 * //		generated from previous step to create `vocab` and `merges` data (files)
 * // 3. Update current tokenizer with newly generated model (`vocab` and `merges` data)
 * func (t *Tokenizer) Train(trainer Trainer, files []string) error {
 *   type Job struct {
 *     File     string
 *     Progress *progressbar.ProgressBar
 *   }
 *
 *   var jobs []Job
 *   wChan := make(chan map[string]uint32)
 *
 *   // channel to signal the main thread that all the words have been
 *   doneChan := make(chan (bool), 1)
 *   dict := make(map[string]uint32)
 *
 *   scanWG := new(sync.WaitGroup)
 *
 *   for _, f := range files {
 *     fsize, err := util.FileSize(f)
 *     if err != nil {
 *       log.Fatal(err)
 *     }
 *     bar := progressbar.New(int(fsize))
 *
 *     jobs = append(jobs, Job{f, bar})
 *   }
 *
 *   // Step 1. scan text files by chunks in goroutines. In each goroutine,
 *   // scan line by line, chop into tokens with (value, count) and
 *   // queue them up in a channel for next step.
 *   // We will set up a wait group to wait for all done.
 *   // For each file do:
 *   // 1. Create a goroutine to read file by chunks
 *   // 2. Read line by line
 *   // 3. Pre-tokenize line of text to tokens
 *   // 4. Process tokens into its value and count
 *   // 5. Send result to a channel for further processing.
 *   for i := 0; i < len(jobs); i++ {
 *     currentJob := i
 *
 *     file := jobs[currentJob].File
 *     // current is the counter for bytes of the file.
 *     var current int64 = 0
 *     var limit int64 = 100 * mb
 *
 *     fi, err := os.Stat(file)
 *     if err != nil {
 *       return err
 *     }
 *     fsize := float64(fi.Size())
 *
 *     chunkNum := int(math.Ceil(fsize / float64(limit)))
 *
 *     // Setup some workers to process
 *     for n := 1; n <= chunkNum; n++ {
 *       scanWG.Add(1)
 *
 *       go func(n int, file string) {
 *         // start reading file chunk by chunk
 *         current = t.processChunk(current, limit, file, wChan, trainer)
 *         fmt.Printf("File chunk %d has been completed\n", n)
 *         scanWG.Done()
 *       }(n, file)
 *     }
 *   }
 *
 *   // Read all incoming words from the channel and add to the dict
 *   go func() {
 *     fmt.Println("Start collecting words...")
 *     for words := range wChan {
 *       for w, c := range words {
 *         count, ok := dict[w]
 *         // word exists, sum up frequency
 *         if ok {
 *           dict[w] = count + c
 *         } else {
 *           // word not exist, let add it
 *           dict[w] = c
 *         }
 *       }
 *     }
 *
 *     // signal the main thread all done with this goroutine
 *     doneChan <- true
 *   }()
 *
 *   // wait for all goroutines to complete
 *   scanWG.Wait()
 *   close(wChan)
 *
 *   // Wait for dictionary to process all words then close
 *   <-doneChan
 *   close(doneChan)
 *
 *   fmt.Printf("Dictionary length: %v words\n", len(dict))
 *   // // Print out some samples
 *   // var count = 0
 *   // for k, _ := range dict {
 *   // if count <= 5 {
 *   // fmt.Println(k)
 *   // count++
 *   // }
 *   // }
 *
 *   // Training model
 *   fmt.Println("Start training...")
 *   model, specialTokens := trainer.Train(dict)
 *
 *   // Replace with trained model
 *   t.Model = &model
 *   t.AddSpecialTokens(specialTokens)
 *
 *   return nil
 * }
 *  */

/*
 * // processChunk reads file chunk and processes it to word-count and sends off to channel
 * // offset: start bound
 * // limit: end bound
 * // filename: file path includes file name
 * // channel: channel to send proccessed words to.
 * // current: cummulative point where the file processing stops.
 * // trainer: Trainer to process tokens
 * func (t *Tokenizer) processChunk(offset int64, limit int64, filename string, channel chan (map[string]uint32), trainer Trainer) (current int64) {
 *   file, err := os.Open(filename)
 *   if err != nil {
 *     panic(err)
 *   }
 *   defer file.Close()
 *
 *   // move the pointer of the file to the start of designated chunk
 *   file.Seek(offset, 0) // 0 means relative to the origin of file
 *
 *   scanner := bufio.NewScanner(file)
 *   buf := make([]byte, 0, 1*gb) // initial buffer
 *   scanner.Buffer(buf, 2*gb)    // max buffer size = 2GB
 *
 *   var cummulativeSize int64
 *
 *   for scanner.Scan() {
 *     // Stop if read size has exceed the chunk size
 *     cummulativeSize += int64(len(scanner.Bytes()))
 *     if cummulativeSize > limit {
 *       break
 *     }
 *
 *     // line words
 *     lwords := make(map[string]uint32)
 *     var line string
 *     line = scanner.Text()
 *     // NOTE: io.scanner returns line w/o `\n`. We add it back manually.
 *     // line = fmt.Sprintf("%v\n", line)
 *
 *     normalized := t.normalize(line)
 *     // NOTE: if there are no preTokenizer, the default `preTokenize`
 *     // will return the whole line without modification. Hence,
 *     // token will be a line string. In that case, we may need to strip
 *     // white spaces in the next step.
 *     preTokenized := t.preTokenize(normalized.GetNormalized())
 *     var tokens []string
 *     for _, tok := range preTokenized {
 *       tokens = append(tokens, tok.Value)
 *     }
 *     // process tokens
 *     trainer.ProcessTokens(lwords, tokens)
 *     // send to channel for further process
 *     channel <- lwords
 *
 *   }
 *
 *   return cummulativeSize
 *
 * }
 *
 *  */

/*
 *
 * func (t *Tokenizer) CTrain(trainer Trainer, files []string) error {
 *   type Job struct {
 *     File     string
 *     Progress *progressbar.ProgressBar
 *   }
 *
 *   var jobs []Job
 *
 *   for _, f := range files {
 *     fsize, err := util.FileSize(f)
 *     if err != nil {
 *       log.Fatal(err)
 *     }
 *     bar := progressbar.New(int(fsize))
 *
 *     jobs = append(jobs, Job{f, bar})
 *   }
 *
 *   // Doing jobs concurrently
 *
 *   g, ctx := errgroup.WithContext(context.Background())
 *   lnChan := make(chan map[string]uint32)
 *
 *   for i := 0; i < len(jobs); i++ {
 *     current := i
 *     g.Go(func() error {
 *       // Now, do the job
 *       file, err := os.Open(jobs[current].File)
 *       if err != nil {
 *         return err
 *       }
 *       defer file.Close()
 *
 *       var line string
 *       words := make(map[string]uint32)
 *
 *       scanner := bufio.NewScanner(file)
 *       for scanner.Scan() {
 *         line = scanner.Text()
 *         // io.scanner returns line w/o `\n`. We add it back manually.
 *         line = fmt.Sprintf("%v\n", line)
 *
 *         normalized := t.normalize(line)
 *         preTokenized := t.preTokenize(normalized.GetNormalized())
 *         var tokens []string
 *         for _, tok := range preTokenized {
 *           tokens = append(tokens, tok.Value)
 *         }
 *         trainer.ProcessTokens(words, tokens)
 *
 *         // Pass processed data to channel
 *         lnChan <- words
 *
 *         select {
 *         case lnChan <- words:
 *         // Keep going
 *         case <-ctx.Done():
 *           return ctx.Err()
 *         }
 *       }
 *
 *       if err := scanner.Err(); err != nil {
 *         return err
 *       }
 *
 *       return nil
 *
 *     })
 *   }
 *
 *   // Close out the channel when the first error occurs or
 *   // when processing is successful.
 *   go func() {
 *     g.Wait()
 *     close(lnChan)
 *   }()
 *
 *   err := g.Wait()
 *
 *   // as long as an error occurs, return it.
 *   if err != nil {
 *     return g.Wait()
 *   }
 *
 *   // Handle result coming from channel
 *   // words is a dictionary of words and their frequency
 *   words := make(map[string]uint32)
 *
 *   // calculate frequency and create a final map
 *   for result := range lnChan {
 *     fmt.Printf("Result: %v\n", result)
 *     for w, c := range result {
 *       count, ok := words[w]
 *       // word exists, sum up frequency
 *       if ok {
 *         words[w] = count + c
 *       }
 *       // word not exist, let add it
 *       words[w] = c
 *     }
 *   }
 *
 *   // Training model
 *   model, specialTokens := trainer.Train(words)
 *
 *   // Replace with trained model
 *   t.Model = &model
 *   t.AddSpecialTokens(specialTokens)
 *
 *   return nil
 * }
 *
 *  */
