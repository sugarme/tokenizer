package util

import (
	"reflect"
)

// Contain checks whether a element is in a slice (same type)
func Contain(item interface{}, slice interface{}) bool {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("Invalid data-type: item and slice are not with the same type.")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

// Reverse reverses any slice
// Source: https://stackoverflow.com/questions/28058278
func Reverse(s interface{}) (retVal interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}

	return s
}
