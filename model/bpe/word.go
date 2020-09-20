package bpe

import (
	"errors"
	"math/rand"
	"time"

	"github.com/emirpasic/gods/trees/binaryheap"
	"github.com/emirpasic/gods/utils"
)

const DefaultCacheCapacity int = 10000

type Merge struct {
	Pos   int
	Rank  int
	NewId int
	Time  time.Time
}

// Ordering is a enum of Less, Equal, and Greater
type Ordering int

const (
	Less    Ordering = -1
	Equal   Ordering = 0
	Greater Ordering = 1
)

// NOTE.Should  we implement comparing methods?
// - Eq
// - PartialCmp
// - Cmp
func (m *Merge) Eq(other *Merge) bool {
	return m.Rank == other.Rank && m.Pos == other.Pos
}

func (m *Merge) PartialCmp(other *Merge) (Ordering, error) {
	// First, compare rank
	if m.Rank != other.Rank {
		if other.Rank < m.Rank {
			return Less, nil
		} else if other.Rank > m.Rank {
			return Greater, nil
		}
	}
	// Then, compare pos
	if other.Pos < m.Pos {
		return Less, nil
	} else if other.Pos > m.Pos {
		return Greater, nil
	} else {
		return Equal, nil
	}
}

func (m *Merge) Cmp(other *Merge) Ordering {
	res, _ := m.PartialCmp(other)
	return res
}

type Symbol struct {
	C    int
	Prev int
	Next int
	Len  int
}

// Some slice methods to manipulate slice struct Symbol
type Symbols []Symbol

// Insert inserts a symbol to the slice at `i` index point
func (ss *Symbols) Insert(s Symbol, i int) error {
	var err error
	if i < 0 || i > len(*ss) {
		err = errors.New("`i` index is out of bound.")
		return err
	}
	*ss = append((*ss)[:i], append([]Symbol{s}, (*ss)[i:]...)...)
	return nil
}

// Remove removes a symbol from the slice at `i` index point
func (ss *Symbols) Remove(i int) error {
	var err error
	if i < 0 || i > len(*ss)-1 {
		err = errors.New("`i` index is out of bound.")
		return err
	}
	*ss = append((*ss)[:i], (*ss)[i+1:]...)
	return nil
}

func (s *Symbol) MergeWith(other *Symbol, newC int) {
	s.C = newC
	s.Len += other.Len
	s.Next = other.Next
}

type Word struct {
	Symbols Symbols
}

func NewWord() *Word {
	return &Word{
		// Symbols: Symbols{},
		Symbols: []Symbol{},
	}
}

func (w *Word) Add(c int, byteLen int) {

	var symbols []Symbol

	symLen := len(w.Symbols)

	if symLen == 0 {
		newSym := Symbol{
			C:    c,
			Prev: -1,
			Next: -1,
			Len:  byteLen,
		}
		symbols = append(symbols, newSym)
	} else {
		for i, s := range w.Symbols {
			// first
			if i == 0 {
				sym := &w.Symbols[i]
				sym.Next = 1
				sym.Prev = -1
				symbols = append(symbols, *sym)
			} else if i == symLen-1 { // last
				sym := &w.Symbols[i]
				sym.Next = symLen
				sym.Prev = symLen - 2
				symbols = append(symbols, *sym)
			} else {
				symbols = append(symbols, s)
			}
		}

		newSym := Symbol{
			C:    c,
			Prev: symLen - 1,
			Next: -1,
			Len:  byteLen,
		}
		symbols = append(symbols, newSym)
	}

	w.Symbols = symbols
}

type Pair struct {
	C1 int
	C2 int
}

// PairVal holds pair's rank and NewId
type PairVal struct {
	Rank  int
	NewId int
}

type WChange struct {
	C1     int
	C2     int
	Change int
}

// Merge finds any pairs of (c1, c2) and removes in place. It also maps changes depending
// on the position of the pair in word.
func (w *Word) Merge(c1, c2, replacement int) ([]WChange, error) {
	// fmt.Printf("before merge word symbols: %v\n", w.Symbols)
	// fmt.Printf("c1: %v - c2: %v- replacement: %v\n", c1, c2, replacement)
	var changes []WChange

	i := 0

	for {
		if i >= len(w.Symbols) {
			break
		}
		// found a pair
		if w.Symbols[i].C == c1 && (i+1) < len(w.Symbols) && w.Symbols[i+1].C == c2 {
			first := w.Symbols[i]
			second := w.Symbols[i+1]

			// If there's other characters before the pair
			if i > 0 {
				changes = append(changes, WChange{
					C1:     w.Symbols[i-1].C,
					C2:     first.C,
					Change: -1,
				})
				changes = append(changes, WChange{
					C1:     w.Symbols[i-1].C,
					C2:     replacement,
					Change: 1,
				})
			}

			// Remove in place
			newS := Symbol{
				C:    replacement,
				Prev: first.Prev,
				Next: second.Next,
				Len:  first.Len + second.Len,
			}

			// Insert replacement before first `char` of pair
			err := w.Symbols.Insert(newS, i)
			if err != nil {
				return nil, err
			}

			// Remove first `char` of pair
			err = w.Symbols.Remove(i + 1)
			if err != nil {
				return nil, err
			}
			// And then the second
			err = w.Symbols.Remove(i + 1)
			if err != nil {
				return nil, err
			}

			// If there are other `chars` after the pair
			if i > 0 && i < len(w.Symbols)-1 {
				// fmt.Println("Yes, there some char after the pair")
				changes = append(changes, WChange{
					C1:     second.C,
					C2:     w.Symbols[i+1].C,
					Change: -1,
				})
				changes = append(changes, WChange{
					C1:     replacement,
					C2:     w.Symbols[i+1].C,
					Change: 1,
				})
			}
		}

		i++

	} // End of `for` loop

	// fmt.Printf("After merge word symbols: %v\n", w.Symbols)

	// fmt.Printf("Num of changes: %v\n", len(changes))
	// fmt.Printf("They are: %v\n", changes)

	return changes, nil
}

func (w *Word) MergeAll(merges map[Pair]PairVal, dropoutOpt ...float32) {
	var dropout float32 = 0.0
	if dropoutOpt != nil {
		dropout = dropoutOpt[0]
	}

	// countComaparator return the `smaller` rank value
	// if both ranks are equal, then return one with smaller timestamp
	countComparator := func(a, b interface{}) int {
		c1 := a.(Merge).Rank
		c2 := b.(Merge).Rank

		if c1 == c2 {
			aTime := a.(Merge).Time
			bTime := b.(Merge).Time

			return utils.TimeComparator(aTime, bTime)
		}

		return utils.IntComparator(c1, c2)
	}

	var queue = binaryheap.NewWith(countComparator)

	// Load items to the heap
	var window = 2

	for i := 0; i < len(w.Symbols)-1; i += window - 1 {
		j := i + window
		if j >= len(w.Symbols) {
			j = len(w.Symbols)
		}
		slice := w.Symbols[i:j]
		pair := Pair{
			C1: slice[0].C,
			C2: slice[1].C,
		}

		// NOTE: if found, push to the queue. If not, continue
		m, ok := merges[pair] // m is PairVal type with pair's rank and newId values
		if ok {
			// log.Fatalf("Cannot find a 'merge' for the pair: %+v\n", pair)
			var merge Merge = Merge{
				Pos:   i,
				Rank:  m.Rank,
				NewId: m.NewId,
			}

			queue.Push(merge)
		}
	}

	var skip []Merge
	r := rand.New(rand.NewSource(99)) // use fixed seed to produce same output on every run.

	// Pop the queue until empty
	for {
		top, ok := queue.Pop()
		if !ok { // it's empty
			break
		}

		if dropout >= 0.0 && r.Float32() < dropout {
			// if dropout > 0.0 {
			skip = append(skip, top.(Merge))
		} else {
			// Re-insert the skipped elements
			for _, s := range skip {
				queue.Push(s)
			}

			if (w.Symbols[top.(Merge).Pos]).Len > 0 {
				if (w.Symbols[top.(Merge).Pos]).Next == -1 {
					// Do nothing if the last symbol
					continue // TODO: do we skip one from outer loop?
				}

				nextPos := w.Symbols[top.(Merge).Pos].Next
				right := w.Symbols[nextPos]

				// Make sure we are not processing an expired queue entry
				targetNewPair := Pair{
					C1: w.Symbols[top.(Merge).Pos].C,
					C2: right.C,
				}

				m, ok := merges[targetNewPair]
				if !ok || m.NewId != top.(Merge).NewId {
					continue
				}

				// Otherwise, let's merge
				w.Symbols[top.(Merge).Pos].MergeWith(&right, top.(Merge).NewId)
				// Tag the right part as removed
				w.Symbols[nextPos].Len = 0

				// Update `prev` on the new `next` to the current pos
				if right.Next > -1 && right.Next < len(w.Symbols) {
					// create a variable so that we can asign an address.
					pos := int(top.(Merge).Pos)
					w.Symbols[right.Next].Prev = pos
				}

				// Insert the new pair formed with the previous symbol
				current := w.Symbols[top.(Merge).Pos]
				if current.Prev >= 0 {
					prev := current.Prev
					prevSymbol := w.Symbols[prev]
					newPair := Pair{
						C1: prevSymbol.C,
						C2: current.C,
					}
					if m, ok := merges[newPair]; ok {
						queue.Push(Merge{
							Pos:   current.Prev,
							Rank:  m.Rank,
							NewId: m.NewId,
						})
					}
				}

				// Insert the new pair formed with the next symbol
				next := current.Next
				if int(next) < len(w.Symbols) && next > -1 {
					nextSymbol := w.Symbols[next]
					newPair := Pair{
						C1: current.C,
						C2: nextSymbol.C,
					}
					if m, ok := merges[newPair]; ok {
						queue.Push(Merge{
							Pos:   top.(Merge).Pos,
							Rank:  m.Rank,
							NewId: m.NewId,
						})
					}
				}
			}
		}
	} // End of `for` loop

	// Filter out the `marked to remove` symbols
	w.removeSymbols()
}

// removeSymbols removes all symbols with lenth == 0
func (w *Word) removeSymbols() {
	var filtered []Symbol
	for _, s := range w.Symbols {
		if s.Len != 0 {
			filtered = append(filtered, s)
		}
	}
	w.Symbols = filtered
}

func (w *Word) GetChars() []int {
	var res []int
	for _, s := range w.Symbols {
		res = append(res, s.C)
	}
	return res
}

func (w *Word) GetOffsets() [][]int {
	var offsets [][]int

	var pos int = 0
	for _, s := range w.Symbols {
		end := pos + s.Len
		offsets = append(offsets, []int{pos, end})
		pos += s.Len
	}

	return offsets
}
