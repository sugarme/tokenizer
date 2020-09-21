package bpe

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/emirpasic/gods/trees/binaryheap"
	"github.com/emirpasic/gods/utils"
	// 800 stars
	progressbar "github.com/schollz/progressbar/v2"
	// 2.2 stars
	// progressbar "github.com/cheggaaa/pb/v3"

	"github.com/sugarme/tokenizer"
)

// Map with no value
// Ref: https://stackoverflow.com/questions/57620170
type UintSet map[int]struct{}

type CharSet map[string]struct{}

type TMerge struct {
	Pair  Pair
	Count int
	Pos   UintSet
	Time  time.Time
}

// NOTE: there exists `Config`
type TConfig struct {
	MinFrequency            int
	VocabSize               int
	ShowProgress            bool
	SpecialTokens           []tokenizer.AddedToken
	LimitAlphabet           *int
	InitialAlphabet         CharSet
	ContinuingSubwordPrefix *string
	EndOfWordSuffix         *string
}

// BpeTrainerBuilder can be used to create a `BpeTrainer`
// with a custom configuration
type BpeTrainerBuilder struct {
	Config *TConfig
}

func NewBPETrainerBuilder() *BpeTrainerBuilder {
	config := TConfig{
		MinFrequency:            0,
		VocabSize:               30000,
		ShowProgress:            true,
		SpecialTokens:           nil,
		LimitAlphabet:           nil,
		InitialAlphabet:         nil,
		ContinuingSubwordPrefix: nil,
		EndOfWordSuffix:         nil,
	}
	return &BpeTrainerBuilder{
		Config: &config,
	}
}

// MinFequency set minimum frequency
func (btb *BpeTrainerBuilder) MinFrequency(freq int) {
	btb.Config.MinFrequency = freq
}

// VocabSize set the vocabulary size
func (btb *BpeTrainerBuilder) VocabSize(size int) {
	btb.Config.VocabSize = size
}

// ShowProgress set whether to show progress
func (btb *BpeTrainerBuilder) ShowProgress(show bool) {
	btb.Config.ShowProgress = show
}

// SpecialToken set special tokens
func (btb *BpeTrainerBuilder) SpecialTokens(tokens []tokenizer.AddedToken) {
	btb.Config.SpecialTokens = tokens
}

//LimitAlphabet set the alphabet limit
func (btb *BpeTrainerBuilder) LimitAlphabet(limit int) {
	btb.Config.LimitAlphabet = &limit
}

// InitialAlphabet set the initial alphabet
func (btb *BpeTrainerBuilder) InitialAlphabet(alphabet CharSet) {
	btb.Config.InitialAlphabet = alphabet
}

// ContinuingSubwordPrefix set the ContinuingSubwordPrefix
func (btb *BpeTrainerBuilder) ContinuingSubwordPrefix(prefix string) {
	btb.Config.ContinuingSubwordPrefix = &prefix
}

// EndOfWordSuffix set the EndOfWordSuffix
func (btb *BpeTrainerBuilder) EndOfWordSuffix(suffix string) {
	btb.Config.EndOfWordSuffix = &suffix
}

// Build constructs the final BpeTrainer
func (btb *BpeTrainerBuilder) Build() *BpeTrainer {
	return &BpeTrainer{
		MinFrequency:            btb.Config.MinFrequency,
		VocabSize:               btb.Config.VocabSize,
		ShowProgress:            btb.Config.ShowProgress,
		SpecialTokens:           btb.Config.SpecialTokens,
		LimitAlphabet:           btb.Config.LimitAlphabet,
		InitialAlphabet:         btb.Config.InitialAlphabet,
		ContinuingSubwordPrefix: btb.Config.ContinuingSubwordPrefix,
		EndOfWordSuffix:         btb.Config.EndOfWordSuffix,
	}
}

// BpeTrainer is in charge of training a `BPE` model from a
// mapping of words to word counts.
//
// Example:
// wordCounts := map[string]int = {
// 	{"Hello", 1},
// 	{"World", 1},
// }
// trainer := NewBPETrainer()
// model, specialTokens := trainer.Train(wordCounts)
type BpeTrainer struct {
	// The minimum frequency a pair must have to produce a merge operation
	MinFrequency int
	// The target vocabulary size
	VocabSize int
	// Whether to show progress while training
	ShowProgress bool
	// A list of special tokens that the model should know of
	SpecialTokens []tokenizer.AddedToken
	// Whether to limit the number of initial tokens that can be kept before
	// computing merges
	LimitAlphabet *int // TODO: replace with int and `None` value = -1
	// The initial alphabet we want absolutely to include. This allows to cover
	// some characters that are not necessarily in the training set
	InitialAlphabet CharSet
	// An optional prefix to use on any subword that exist only behind another one
	ContinuingSubwordPrefix *string
	// An optional suffix to characterize and end-of-word subword
	EndOfWordSuffix *string
}

func NewBpeTrainer(minFreq int, vocabSize int) *BpeTrainer {
	btb := NewBPETrainerBuilder()
	bpeTrainer := btb.Build()

	bpeTrainer.MinFrequency = minFreq
	bpeTrainer.VocabSize = vocabSize

	return bpeTrainer

}

func (bt *BpeTrainer) setupProgress() interface{} {
	if bt.ShowProgress {
		// TODO: setup progress bar
	}
	return nil
}

// set the progress bar in the finish state
func (bt *BpeTrainer) finalizeProgress(pb interface{}, finalLen int) interface{} {
	if pb != nil {
		// TODO:
		// set length
		// finish up
	}

	return nil
}

// updateProgress update the progress bar with the new provided length and msg
func (bt *BpeTrainer) updateProgress(p interface{}, len int, msg string) {
	// TODO: update progress bar
}

// addSpecialTokens adds the provided special tokens to the initial vocabulary
func (bt *BpeTrainer) addSpecialTokens(w2id map[string]int, id2w []string) {
	for _, tok := range bt.SpecialTokens {
		if _, ok := w2id[tok.Content]; !ok {
			id2w = append(id2w, tok.Content)
			w2id[tok.Content] = len(id2w) - 1
		}
	}
}

// computeAlphabet generate maps of `chars` from input words and limit it if relevant
func (bt *BpeTrainer) computeAlphabet(wc map[string]int) (wordToId map[string]int, IdToWord []string) {
	// compute the alphabet from seen words
	var (
		alphabet map[string]int = make(map[string]int)
		w2id     map[string]int = make(map[string]int)
		id2w     []string
	)

	for word, count := range wc {
		chars := strings.Split(word, "")
		for _, char := range chars {
			// if char not existing, newCount will be zero
			if newCount, ok := alphabet[char]; ok {
				newCount += count
				alphabet[char] = newCount
			} else {
				alphabet[char] = newCount
			}
		}
	}

	// Also, include anything from the provided intial alphabet
	// NOTE: InitialAlphabet is CharSet which is map[string]struct{}
	for initChar, _ := range bt.InitialAlphabet {
		// asign a uint max as frequency
		alphabet[initChar] = math.MaxUint32
	}

	type keptItem struct {
		Char string
		Freq int
	}
	var kept []keptItem
	// NOTE: alphabet map need to be sorted first
	keys := make([]string, 0)
	for k, _ := range alphabet {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, char := range keys {
		freq := alphabet[char]
		kept = append(kept, keptItem{char, freq})
	}

	// compute the number of chars to remove from the alphabet
	// if `limitAlphabet` < `len(initialAlphabet)` some of these
	// initial characters will be removed.
	// TODO: how we sort these characters before cutting off?
	var toRemove int = 0
	var limit int
	if bt.LimitAlphabet != nil {
		limit = int(*bt.LimitAlphabet)
		if len(alphabet) > int(*bt.LimitAlphabet) {
			toRemove = len(alphabet) - limit
		}
	}

	// remove the unwanted `chars`
	if toRemove > 0 {
		fmt.Println("We are going to remove some chars...")
		// 1. Sort `kept` by char alphabetically?
		// TODO: double-check this (sort by char or freq? asc or desc)
		sort.Slice(kept, func(i, j int) bool {
			return kept[i].Char < kept[j].Char
		})
		// 2. Remove the unwanted chars
		kept = kept[:toRemove]
	}

	// // Keep the initial alphabet (sorted by determinism)
	// sort.Slice(kept, func(i, j int) bool {
	// // sort by freq
	// return kept[i].Freq > kept[j].Freq
	// })

	sort.Slice(kept, func(i, j int) bool {
		return kept[i].Char < kept[j].Char
	})

	for _, k := range kept {
		if _, ok := w2id[k.Char]; !ok {
			id2w = append(id2w, k.Char)
			w2id[k.Char] = len(id2w) - 1
		}
	}

	// // Print out `char` and its `freq` to check
	// // TODO: delete this block
	// for _, v := range kept {
	// fmt.Printf("char: %v %v\n", v.Char, v.Freq)
	// }
	// fmt.Println(w2id)
	// fmt.Println(id2w)

	return w2id, id2w
}

// tokenizerWord tokenizes words and adds subwords (prefix, suffix) to the vocabulary when relevant
func (bt *BpeTrainer) tokenizeWords(wc map[string]int, w2id map[string]int, id2w []string, pb interface{}) ([]Word, []int, map[string]int, []string) {
	// NOTE: bp is progress bar.
	// TODO: update bp to specific progress bar type

	// words := make([]Word, len(wc))
	// counts := make([]int, len(wc))
	var words []Word
	var counts []int
	var sortedWords []string

	keys := sortedKeys(wc)

	for _, word := range keys {
		sortedWords = append(sortedWords, word)
		// for word, count := range wc {
		count := wc[word]
		var currentWord Word

		counts = append(counts, count)

		chars := strings.Split(word, "")

		for i, c := range chars {
			var s string
			if _, ok := w2id[c]; ok {
				// Add the `continuingSubwordPrefix` if relevant
				if i == 0 { // first `char`
					if prefix := bt.ContinuingSubwordPrefix; prefix != nil {
						s = fmt.Sprintf("%v%v", &prefix, c)
					} else {
						s = c
					}
				}
				if i > 0 && i < len(chars)-1 {
					s = c
				}
				// Add the `endOfWordSuffix` if relevant
				if i == len(chars)-1 { // last `char`
					if suffix := bt.EndOfWordSuffix; suffix != nil {
						s = fmt.Sprintf("%v%v", &suffix, c)
					} else {
						s = c
					}
				}

				// Insert the new formed string if neccessary
				if _, ok := w2id[s]; !ok {
					id2w = append(id2w, s)
					w2id[s] = len(id2w) - 1
					currentWord.Add(w2id[s], len(s))
				} else {
					currentWord.Add(w2id[s], len(s))
				}
			}

		} // end loop of `chars`

		words = append(words, currentWord)

		// TODO: update progress bar to 1

	} // end loop of `wc`

	// fmt.Printf("Sorted original words: %v\n", sortedWords)
	// fmt.Println(w2id)
	// fmt.Println(id2w)
	// fmt.Println(words)
	// fmt.Println(counts)

	return words, counts, w2id, id2w

}

// coutPairs counts frequency of pairs (char pair) from input words and put into maps
func (bt *BpeTrainer) countPairs(words []Word, counts []int, progress interface{}) (map[Pair]int, map[Pair]UintSet) {

	type pcResult struct {
		PC map[Pair]int
		WT map[Pair]UintSet
	}

	var pairCounts map[Pair]int = make(map[Pair]int, bt.VocabSize*2)
	var whereToUpdate map[Pair]UintSet = make(map[Pair]UintSet, bt.VocabSize*2)

	// Divide w into work units that take ~100μs-1ms to compute.
	n := len(words)
	// size := int(1000000 / n)
	// if size < 8 {
	// size = 8
	// }

	// batch size
	var workerNum int
	if n < 10000 {
		workerNum = 1
	} else {
		workerNum = 4
	}

	size := int(n / workerNum)
	if size < 1 {
		size = 1
	}

	resChan := make(chan pcResult)
	doneChan := make(chan (bool), 1)

	var pcWG = new(sync.WaitGroup)

	// pcWG.Add(n)
	for i, j := 0, size; i < n; i, j = j, j+size {
		if j > n {
			j = n
		}

		pcWG.Add(1)

		// fmt.Printf("i: %v - j: %v\n", i, j)

		go func(i, j int, words []Word, counts []int) {
			defer pcWG.Done()
			var pc map[Pair]int = make(map[Pair]int)
			var wt map[Pair]UintSet = make(map[Pair]UintSet)

			for k := i; k < j; k++ {
				// Do individual task here with index k
				// which is a single word at index k in slice `words`
				word := words[k]
				count := counts[k]

				// fmt.Printf("Word: %v\n", word)
				var window = 2
				for x := 0; i < len(word.Symbols)-1; x += window - 1 {
					y := x + window
					if y > len(word.Symbols) {
						// TODO: should we stop when last chunk < chunk size or we just return it
						break
						// y = len(word.Symbols)
					}

					w := word.Symbols[x:y]
					pair := Pair{
						C1: w[0].C,
						C2: w[1].C,
					}

					// Initialize pairCounts and whereToUpdate for this pair
					// if we just seen it
					var curCount int = count
					if c, ok := pc[pair]; !ok {
						pc[pair] = 0
					} else {
						curCount += c
					}

					// Then update counts
					pc[pair] = curCount

					// hashset map[int]struct{}
					var hs UintSet = make(map[int]struct{})
					if h, ok := wt[pair]; ok {
						h[k] = struct{}{} // found. Modify it
					} else {
						// create a new
						hs[k] = struct{}{}
						wt[pair] = hs
					}

				}

				// TODO: update progress bar

			}
			// Send off result to channel

			resChan <- pcResult{pc, wt}

		}(i, j, words, counts)
	}

	// Setup a goroutine to aggregate results send from pairCount workers
	// via resChan channel
	go func() {
		for res := range resChan {
			for pair, count := range res.PC {
				if _, ok := pairCounts[pair]; !ok {
					pairCounts[pair] = count
				} else {
					c := pairCounts[pair]
					pairCounts[pair] = c + count
				}
			}

			for pair, hashSet := range res.WT {
				whereToUpdate[pair] = hashSet
			}
		}

		doneChan <- true
	}()

	pcWG.Wait()
	close(resChan)

	<-doneChan

	// TODO: test whether having a data race??? as goroutines update pairCounts and whereToUpdate
	return pairCounts, whereToUpdate

}

// countPairsM counts frequency of pairs not using concurrency/paralellism
func (bt *BpeTrainer) countPairsM(words []Word, counts []int, progress interface{}) (map[Pair]int, map[Pair]UintSet) {

	var pairCounts map[Pair]int = make(map[Pair]int, bt.VocabSize*2)
	var whereToUpdate map[Pair]UintSet = make(map[Pair]UintSet, bt.VocabSize*2)

	for i := 0; i < len(words); i++ {
		word := words[i]
		var window = 2
		for x := 0; x < len(word.Symbols)-1; x += window - 1 {
			y := x + window
			if y > len(word.Symbols) {
				// TODO: should we stop when last chunk < chunk size or we just return it
				break
				// y = len(word.Symbols)
			}

			w := word.Symbols[x:y]
			pair := Pair{
				C1: w[0].C,
				C2: w[1].C,
			}

			// Initialize pairCounts and whereToUpdate for this pair if it is not existing
			// 1. `pairCounts` map
			if _, ok := pairCounts[pair]; !ok {
				pairCounts[pair] = 0
			}
			// Then update counts
			count := counts[i]
			pairCounts[pair] += count

			// 2. `whereToUpdate` is a map of with
			// - key	: a pair
			// - value: a hashset (map[int]struct{}). It keeps the index/indices of the word
			// (in input words) where pair comes from.
			// Hence: `whereToUpdate` maps which `word` from input words the pair comes from.
			var hs UintSet = make(map[int]struct{})
			if h, ok := whereToUpdate[pair]; !ok {
				// Not exisitng: create a new
				hs[i] = struct{}{}
				whereToUpdate[pair] = hs
			} else {
				// Found one: append one more index to it
				h[i] = struct{}{}
				whereToUpdate[pair] = h
			}

		}
	} // end of `for` loop

	// fmt.Printf("pairCounts: %v\n", pairCounts)
	// fmt.Printf("whereToUpdate: %v\n", whereToUpdate)

	return pairCounts, whereToUpdate

}

// Implement Trainer interface. It has the following methods:
// 1. WithProgressBar() bool
// 2. Train(words map[string]int) (Model, []string)
// 3. ProcessTokens(words map[string]int, tokens []string)

func (bt *BpeTrainer) WithProgressBar() bool {
	// TODO: implement a progress bar
	return false
}

// Train trains bpe model on input wordCounts and returns
// 1. BPE model; 2. merges
// func (bt *BpeTrainer) Train(wordCounts map[string]int) (BPE, []string) {
func (bt *BpeTrainer) Train(wordCounts map[string]int) (tokenizer.Model, []tokenizer.AddedToken) {

	// fmt.Printf("Word Counts: %v\n", wordCounts)

	bpe, merges := bt.train(wordCounts)

	return bpe, merges
}

// Process a bunch of tokens, counting them
func (bt *BpeTrainer) ProcessTokens(words map[string]int, tokens []string) {
	for _, token := range tokens {
		c, _ := words[token]
		c += 1
		words[token] = c
	}
}

// Train a BPE model
// func (bt *BpeTrainer) train(wordCounts map[string]int) (BPE, []string) {
func (bt *BpeTrainer) train(wordCounts map[string]int) (BPE, []tokenizer.AddedToken) {
	// return bt.Train(wordCounts)
	var (
		wordToId map[string]int = make(map[string]int)
		idToWord []string
	)

	var progress = bt.setupProgress()

	// NOTE: temporary add progress bar for counting trained words ONLY
	// TODO: setup progress bar for the whole training process.
	pb := progressbar.New(int(bt.VocabSize))

	// 1. Add all special tokens to the vocabular
	fmt.Printf("1. Adding special tokens...\n")
	bt.addSpecialTokens(wordToId, idToWord)

	// 2. Compute the initial alphabet (create maps of `chars`)
	// These maps will be updated if `prefix`, `suffix` are added
	// in the following steps
	fmt.Printf("2. Creating maps of 'chars'...\n")
	wordToId, idToWord = bt.computeAlphabet(wordCounts)
	// fmt.Printf("Before id2Word: length %v - values:  %v\n", len(idToWord), idToWord)
	// fmt.Printf("Before word2Id: length %v - %v\n", len(wordToId), wordToId)

	// 3. Tokenize words (add prefix, suffix to the map if relevant)
	// NOTE: `char` maps (wordToId, idToWord) will be updated if added prefix and/or suffix
	fmt.Printf("3. Tokenizing words...\n")
	bt.updateProgress(progress, len(wordCounts), "Tokenize word")

	words, counts, wordToId, idToWord := bt.tokenizeWords(wordCounts, wordToId, idToWord, progress)

	// fmt.Printf("Words: %v\n", idToWord)

	bt.finalizeProgress(progress, len(words))

	// 4. Count pairs in words
	// words will be split to `char`, paired and count their frequency.
	// The result will be a map of (pairs and their frequency) and
	// a map of (pairs and their int hashset - which is a map of key with no value)
	// represent a position to update pair.
	fmt.Printf("4. Pairing and counting co-occurence...\n")
	bt.updateProgress(progress, len(words), "Count pairs")

	var (
		pairCounts    map[Pair]int     = make(map[Pair]int)
		whereToUpdate map[Pair]UintSet = make(map[Pair]UintSet)
	)

	pairCounts, whereToUpdate = bt.countPairs(words, counts, progress)
	// pairCounts, whereToUpdate = bt.countPairsM(words, counts, progress)

	// 5. Do merges
	fmt.Printf("5. Merging pairs from top count down...\n")

	// countComparator sort heap descendingly by `Count` field of TMerge struct
	countComparator := func(a, b interface{}) int {
		c1 := a.(TMerge).Count
		c2 := b.(TMerge).Count

		if c1 == c2 {
			aTime := a.(TMerge).Time
			bTime := b.(TMerge).Time

			return utils.TimeComparator(bTime, aTime)
		}

		return utils.IntComparator(c2, c1)
	}

	var queue = binaryheap.NewWith(countComparator)

	// insert them to the queue
	for pair, pos := range whereToUpdate {
		if count, ok := pairCounts[pair]; ok {
			// char1 := idToWord[pair.C1]
			// char2 := idToWord[pair.C2]
			// fmt.Printf("pair chars: %v%v - pair: %v - count: %v - pos: %v\n", char1, char2, pair, count, pos)
			queue.Push(TMerge{
				Pair:  pair,
				Count: count,
				Pos:   pos,
			})
		} else {
			fmt.Println("Something wrong here...")
		}
	}

	bt.finalizeProgress(progress, len(words))
	bt.updateProgress(progress, bt.VocabSize, "Compute merges")

	type TMerges struct {
		Pair    Pair
		PairVal int
	}

	var merges []TMerges

	// fmt.Printf("Word2ID before merge loop: %v\n", wordToId)

	// // Print out token and freq after tokenization
	// for n := 0; n <= queue.Size(); n++ {
	// if v, ok := queue.Pop(); ok {
	// merge := v.(TMerge)
	// count := merge.Count
	// tok := merge.Pair
	// fmt.Printf("Token: %v - Freq: %v\n", tok, count)
	// }
	// }

	// for i := 0; i < 50; i++ {
	for {
		// fmt.Println(len(wordToId))
		// Stop as soon as we have a big enough vocabulary
		if len(wordToId) >= bt.VocabSize {
			pb.Finish()
			fmt.Printf("\nVocab has enough words (%d)! Done.\n", bt.VocabSize)
			break
		}

		if queue.Empty() {
			fmt.Println("queue is empty")
			break
		}

		t, _ := queue.Pop()
		var top TMerge = t.(TMerge)

		// fmt.Printf("Top: count = %v | pair: %v\n", top.Count, top.Pair)

		if top.Count != pairCounts[top.Pair] {
			pairCounts[top.Pair] = top.Count
			queue.Push(top)
			// fmt.Println("Not found. Push new one...")

			continue
		}

		if top.Count < 1 || top.Count < bt.MinFrequency {
			fmt.Printf("\nWe stop because top count is hit limit\n")
			break
		}

		// If `pair` is no longer mapped in its origin word, skip merging
		var skip bool = true
		for p, _ := range top.Pos {
			w := words[p]
			c1 := top.Pair.C1
			c2 := top.Pair.C2

			for i := 0; i < len(w.Symbols); i++ {
				// found it
				if w.Symbols[i].C == c1 && (i+1) < len(w.Symbols) && w.Symbols[i+1].C == c2 {
					skip = false
				}
			}
		}

		if skip {
			// fmt.Println("Skipped because `pair` is no longer mapped in its origin word!")
			continue
		}

		partA := idToWord[top.Pair.C1]
		partB := idToWord[top.Pair.C2]

		// Build new token
		if prefix := bt.ContinuingSubwordPrefix; prefix != nil {
			if strings.HasPrefix(partB, *prefix) {
				// strip prefix
				partB = strings.TrimPrefix(partB, *prefix)
			}
		}

		newToken := fmt.Sprintf("%v%v", partA, partB)
		// fmt.Printf("new token: %v\n", newToken)

		// Insert new token
		newTokenId := len(idToWord)
		idToWord = append(idToWord, newToken)
		wordToId[newToken] = newTokenId
		merges = append(merges, TMerges{top.Pair, newTokenId})

		type TChange struct {
			WChange WChange
			WIndex  int
		}
		var changes []TChange
		// Merge the new pair in word(s) that contains the current pair
		for i, _ := range top.Pos {
			// NOTE: words []Word
			// TODO: merge each of these words concurrently
			w := words[i]
			wChanges, err := w.Merge(top.Pair.C1, top.Pair.C2, newTokenId)
			if err != nil {
				fmt.Println(err)
			}
			// update back `words` list. Because after merging, word map may be changed.
			words[i] = w
			for _, wc := range wChanges {
				changes = append(changes, TChange{wc, int(i)})
			}
		}

		// Introduce new formed pairs
		// NOTE: reset `whereToUpdate` first
		whereToUpdate = make(map[Pair]UintSet)
		for _, tc := range changes {
			if tc.WChange.Change > 0 {
				count := tc.WChange.Change * counts[tc.WIndex]

				pair := Pair{tc.WChange.C1, tc.WChange.C2}

				if _, ok := pairCounts[pair]; !ok {
					pairCounts[pair] = count
					// fmt.Printf("pair: %v - count: %v\n", pair, count)
				} else {
					c, _ := pairCounts[pair]
					c += count
					pairCounts[pair] = c
					// fmt.Printf("pair: %v - count: %v\n", pair, c)
				}

				var hs UintSet = make(map[int]struct{})
				if h, ok := whereToUpdate[pair]; !ok {
					// if not existing, we create new one anyway
					hs[tc.WIndex] = struct{}{}
					whereToUpdate[pair] = hs
				} else {
					// Existing, append one
					h[tc.WIndex] = struct{}{}
					whereToUpdate[pair] = h
				}
			}
		}

		// fmt.Printf("length of whereToUpdate: %v\n", len(whereToUpdate))
		// fmt.Println(whereToUpdate)

		for pair, pos := range whereToUpdate {
			count := pairCounts[pair]
			// char1 := idToWord[pair.C1]
			// char2 := idToWord[pair.C2]
			// fmt.Printf("pair chars: '%v%v' - pair: %v - count: %v - pos: %v\n", char1, char2, pair, count, pos)
			if count > 0 {
				queue.Push(TMerge{
					pair, count, pos, time.Now(),
				})
			}
		}

		// TODO: update progress bar by 1

		pb.Add(1)

	} // end of `for` loop

	bt.finalizeProgress(progress, len(merges))

	var builder *BpeBuilder
	builder = NewBpeBuilder()

	var newMerges Merges = make(map[Pair]PairVal)

	for i, m := range merges {
		pairVal := PairVal{
			i,
			m.PairVal,
		}
		newMerges[m.Pair] = pairVal
	}

	builder.VocabAndMerges(wordToId, newMerges)

	if prefix := bt.ContinuingSubwordPrefix; prefix != nil {
		builder.ContinuingSubwordPrefix(*prefix)
	}

	if suffix := bt.EndOfWordSuffix; suffix != nil {
		builder.EndOfWordSuffix(*suffix)
	}

	bpe, err := builder.Build()

	if err != nil {
		fmt.Println(err)
	}

	return *bpe, bt.SpecialTokens
}

// Whether we should show progress
func (bt *BpeTrainer) shouldShowProgress() bool {
	return bt.ShowProgress
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
