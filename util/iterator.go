package util

// RuneIter is rune iterator with Next() method.
type RuneIter struct {
	items []rune
	next  int
}

// NewRuneIter creates a RuneIter
func NewRuneIter(data []rune) *RuneIter {
	next := -1
	if len(data) > 0 {
		next = 0
	}

	return &RuneIter{
		items: data,
		next:  next,
	}
}

// Next implement interator for RuneIter
func (it *RuneIter) Next() (retVal rune, ok bool) {
	if it.next == -1 {
		return retVal, false
	}

	retVal = it.items[it.next]

	if it.next+1 >= len(it.items) {
		it.next = -1
	} else {
		it.next += 1
	}

	return retVal, true
}

// Len returns total items in RuneIter
func (it *RuneIter) Len() int {
	return len(it.items)
}

// CurrentIndex returns current index of RuneIter
func (it *RuneIter) CurrentIndex() int {
	if len(it.items) == 0 {
		return -1
	} else {
		return it.next - 1
	}
}

// Reset resets current index to first item if any.
func (it *RuneIter) Reset() {
	if len(it.items) == 0 {
		it.next = -1
	} else {
		it.next = 0
	}
}

// Iter contains data and methods for an iterator
type Iter[T any] struct {
	items []T
	next  int
}

// NewIter creates an Iter from an input slice.
func NewIter[T any](data []T) *Iter[T] {
	next := -1
	if len(data) > 0 {
		next = 0
	}

	return &Iter[T]{
		items: data,
		next:  next,
	}
}

// Next implements iterator functionality for Iter
func (it *Iter[T]) Next() (retVal T, ok bool) {
	if it.next == -1 || it.next >= len(it.items) {
		return retVal, false
	}

	retVal = it.items[it.next]
	it.next++

	return retVal, true
}
