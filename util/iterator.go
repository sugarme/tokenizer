package util

import (
	"reflect"
)

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

// Iter contains data and methods for an interator
type Iter struct {
	items []interface{}
	next  int
}

// NewIter creates a Iter from an input slice. Otherwise
// it will panic.
func NewIter(data interface{}) *Iter {
	if reflect.TypeOf(data).Name() != "slice" {
		panic("Invalid input 'data'. It must be a slice.")
	}

	next := -1
	if len(data.([]interface{})) > 0 {
		next = 0
	}

	return &Iter{
		items: data.([]interface{}),
		next:  next,
	}
}

// Next implements a iterator functionality for Iter
func (it *Iter) Next() (retVal interface{}, ok bool) {
	if it.next == -1 {
		return retVal, false
	}

	retVal = it.items[it.next]

	if it.next+1 > len(it.items) {
		it.next = -1
	} else {
		it.next += 1
	}

	return retVal, true
}
