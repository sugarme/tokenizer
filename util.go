package tokenizer

import (
	"errors"
	"log"
)

type TruncationParams struct {
	MaxLength int
	Strategy  TruncationStrategy
	Stride    int
}

type PaddingParams struct {
	Strategy  PaddingStrategy
	Direction PaddingDirection
	PadId     int
	PadTypeId int
	PadToken  string
}

// PaddingStrategy is a enum of either
// - string `BatchLongest`
// - or a func type `Fixed(uint)` which return a uint
// Example:
//
//	func main() {
//	    var ps PaddingStrategy
//	    ps = NewPaddingStrategy(WithFixed(3))
//	    fmt.Println(ps.Value)
//	}
type PaddingStrategy struct {
	Value interface{}
	Name  string
}

type PaddingStrategyOption func(*PaddingStrategy)

func WithBatchLongest() PaddingStrategyOption {
	return func(ps *PaddingStrategy) {
		ps.Value = "BatchLongest"
		ps.Name = "BatchLongest"
	}
}

func WithFixed(size int) PaddingStrategyOption {
	return func(ps *PaddingStrategy) {
		ps.Value = size
		ps.Name = "Fixed"
	}
}

func NewPaddingStrategy(opts ...PaddingStrategyOption) *PaddingStrategy {
	const defaultVal = "BatchLongest"

	ps := &PaddingStrategy{
		Value: defaultVal,
		Name:  defaultVal,
	}

	for _, opt := range opts {
		opt(ps)
	}

	return ps

}

// TruncationStrategy is enum of int type represents truncation strategy
type TruncationStrategy int

const (
	LongestFirst TruncationStrategy = iota
	OnlyFirst
	OnlySecond
)

const (
	SecondSequenceNotProvided = "Truncation error: Second sequence not provided"
	SequenceTooShort          = "Truncation error: Sequence to truncate too short to respect the provided max_length"
)

func TruncateEncodings(encoding, pairEncoding *Encoding, params *TruncationParams) (tEncoding, tPairEncoding *Encoding) {
	var (
		totalLength int
		toRemove    int
		err         error
	)

	if params.MaxLength == 0 {
		return encoding, pairEncoding
	}

	totalLength = len(encoding.GetIds())
	if pairEncoding != nil {
		totalLength = len(encoding.GetIds()) + len(pairEncoding.GetIds())
	}

	if totalLength < params.MaxLength {
		return encoding, pairEncoding
	}

	toRemove = totalLength - params.MaxLength

	switch params.Strategy {
	case LongestFirst:
		nFirst := len(encoding.GetIds())
		nSecond := len(pairEncoding.GetIds())

		for i := 0; i < toRemove; i++ {
			if nFirst > nSecond {
				nFirst -= 1
			}
			nSecond -= 1
		}

		encoding.Truncate(nFirst, params.Stride)
		if pairEncoding != nil {
			pairEncoding.Truncate(nSecond, params.Stride)
		}

	case OnlyFirst, OnlySecond:
		var truncateFunc = func(target *Encoding) (*Encoding, error) {
			targetLength := len(target.GetIds())
			if targetLength > toRemove {
				target.Truncate(targetLength-toRemove, params.Stride)
				return target, nil
			} else {
				err := errors.New(SequenceTooShort)
				return nil, err
			}
		}

		if params.Strategy == OnlyFirst {
			encoding, err = truncateFunc(encoding)
		} else if pairEncoding != nil {
			pairEncoding, err = truncateFunc(pairEncoding)
		} else {
			err = errors.New(SecondSequenceNotProvided)
		}

	}

	if err != nil {
		log.Fatal(err)
	}

	return encoding, pairEncoding
}

func PadEncodings(encodings []Encoding, params PaddingParams) []Encoding {
	if len(encodings) == 0 {
		return encodings
	}

	var padLength int

	switch params.Strategy.Name {
	case "Fixed":
		padLength = params.Strategy.Value.(int)
	case "BatchLongest":
		var max int = 0
		for _, encoding := range encodings {
			if len(encoding.GetIds()) > max {
				max = len(encoding.GetIds())
			}
		}
		padLength = max
	}

	// TODO: implement concurrency with for loop
	var newEncodings []Encoding
	for _, e := range encodings {
		en := e
		paddedEn := en.Pad(padLength, params.PadId, params.PadTypeId, params.PadToken, params.Direction)
		newEncodings = append(newEncodings, *paddedEn)
	}

	return newEncodings
}

type Range []int

func NewRange(start, end int) Range {
	if start < 0 {
		panic("Invalid 'start' for NewRange()")
	}
	if end < 0 || end <= start {
		panic("Invalid 'end' for NewRange()")
	}

	var r []int
	for i := start; i < end; i++ {
		r = append(r, i)
	}

	return r
}

func (r Range) Len() int {
	return len(r)
}

func (r Range) Contains(item int) bool {
	for _, v := range r {
		if v == item {
			return true
		}
	}

	return false
}

func (r Range) IsEmpty() bool {
	return len(r) == 0
}

// TODO. more methods of Range
