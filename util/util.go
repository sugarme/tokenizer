package util

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/rivo/uniseg"
)

/*
 * // Makerange creates a sequence of number (range)
 * // Ref. https://stackoverflow.com/questions/39868029
 * func MakeRange(min, max int) []int {
 *   if min == max {
 *     return []int{}
 *   }
 *   a := make([]int, max-min+1)
 *   for i := range a {
 *     a[i] = min + i
 *   }
 *   return a
 * }
 *  */
func MakeRange(start, end int) []int {

	var r []int
	if end == start {
		return []int{}
	} else {
		for i := start; i < end; i++ {
			r = append(r, start)
		}
	}

	return r
}

// StringIndex returns index (start) for substring on a given string
// It returns -1 and error if not matching
func StringIndex(s, sub string) (index int, err error) {
	i := strings.Index(s, sub)
	if i <= -1 {
		err := errors.New("Index not found")
		return -1, err
	}
	return i, nil
}

// ToASCII converts string to ASCII form
// Ref. https://stackoverflow.com/questions/12668681
func ToASCII(s string) string {
	var as []string
	for _, r := range []rune(s) {
		quoted := strconv.QuoteRuneToASCII(r)
		as = append(as, quoted[1:len(quoted)-1])
	}

	return strings.Join(as, "")
}

// ToGrapheme converts string to grapheme
// TODO: should we include this func?
// Ref: https://github.com/golang/go/issues/14820
func ToGrapheme(s string) string {
	gr := uniseg.NewGraphemes(s)

	var str []string

	for gr.Next() {
		s := fmt.Sprintf("%x", gr.Runes())
		str = append(str, s)
	}

	return strings.Join(str, "")
}

// Ref. https://stackoverflow.com/questions/14000534
type RuneGen func() rune

// MapRune maps ...
func MapRune(g RuneGen, f func(rune) rune) RuneGen {
	return func() rune {
		return f(g())
	}
}

// MinMax returns min and max from input int array
func MinMax(array []int) (int, int) {
	var max int = array[0]
	var min int = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

// StringInSlice check whether given string is in a slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// FileSize returns length of given file using `os.Stat()`
// Ref. https://stackoverflow.com/questions/17133590
func FileSize(filepath string) (int64, error) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	// get the size
	return fi.Size(), nil
}

// ReadAllLn reads all line by line from a file using bufio.scanner
func ReadAllLn(filepath string, keepBreakLine bool) ([]string, error) {
	var lines []string
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		if keepBreakLine {
			lines = append(lines, fmt.Sprintf("%v\n", l))
		}
		lines = append(lines, l)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// Zip zips 2 slices in first and second arguments to third argument
// Ref: https://stackoverflow.com/questions/26957040
//
// Usage Example
// a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
// b := []int{0, 9, 8, 7, 6, 5, 4, 3, 2, 1}
// c := [][2]int{}
//
// e := zip(a, b, &c)
//
// if e != nil {
// 		fmt.Println(e)
//  	return
// }
//
// fmt.Println(c)
func Zip(a, b, c interface{}) error {

	ta, tb, tc := reflect.TypeOf(a), reflect.TypeOf(b), reflect.TypeOf(c)

	if ta.Kind() != reflect.Slice || tb.Kind() != reflect.Slice || ta != tb {
		return fmt.Errorf("zip: first two arguments must be slices of the same type")
	}

	if tc.Kind() != reflect.Ptr {
		return fmt.Errorf("zip: third argument must be pointer to slice")
	}

	for tc.Kind() == reflect.Ptr {
		tc = tc.Elem()
	}

	if tc.Kind() != reflect.Slice {
		return fmt.Errorf("zip: third argument must be pointer to slice")
	}

	eta, _, etc := ta.Elem(), tb.Elem(), tc.Elem()

	if etc.Kind() != reflect.Array || etc.Len() != 2 {
		return fmt.Errorf("zip: third argument's elements must be an array of length 2")
	}

	if etc.Elem() != eta {
		return fmt.Errorf("zip: third argument's elements must be an array of elements of the same type that the first two arguments are slices of")
	}

	va, vb, vc := reflect.ValueOf(a), reflect.ValueOf(b), reflect.ValueOf(c)

	for vc.Kind() == reflect.Ptr {
		vc = vc.Elem()
	}

	if va.Len() != vb.Len() {
		return fmt.Errorf("zip: first two arguments must have same length")
	}

	for i := 0; i < va.Len(); i++ {
		ea, eb := va.Index(i), vb.Index(i)
		tt := reflect.New(etc).Elem()
		tt.Index(0).Set(ea)
		tt.Index(1).Set(eb)
		vc.Set(reflect.Append(vc, tt))
	}

	return nil
}

// MinMax returns min and max from input int array
func MinMaxFloat64(array []float64) (float64, float64) {
	var max float64 = array[0]
	var min float64 = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}
