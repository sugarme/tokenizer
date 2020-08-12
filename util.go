package tokenizer

import (
	"errors"
)

type TruncationParams struct {
	MaxLength uint
	Strategy  TruncationStrategy
	Stride    uint
}

type PaddingParams struct {
	Strategy  PaddingStrategy
	Direction PaddingDirection
	PadId     uint32
	PadTypeId uint32
	PadToken  string
}

// PaddingStrategy is a enum of either
// - string `BatchLongest`
// - or a func type `Fixed(uint)` which return a uint
// Example:
// func main() {
//     var ps PaddingStrategy
//     ps = NewPaddingStrategy(WithFixed(3))
//     fmt.Println(ps.Value)
// }
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

func WithFixed(size uint) PaddingStrategyOption {
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

func TruncateEncodings(e Encoding, params TruncationParams, pairOpt ...Encoding) (*Encoding, *Encoding, error) {
	var (
		encoding     *Encoding = &e
		pairEncoding *Encoding = nil
		totalLength  uint
		toRemove     uint
		err          error
	)

	if len(pairOpt) > 0 {
		pairEncoding = &pairOpt[0]
	}

	if params.MaxLength == 0 {
		return encoding, pairEncoding, nil
	}

	totalLength = uint(len(encoding.GetIds()))
	if pairEncoding != nil {
		totalLength = uint(len(encoding.GetIds()) + len(pairEncoding.GetIds()))
	}

	if totalLength < params.MaxLength {
		return encoding, pairEncoding, nil
	}

	toRemove = totalLength - params.MaxLength

	switch params.Strategy {
	case LongestFirst:
		nFirst := uint(len(encoding.GetIds()))
		nSecond := uint(len(pairEncoding.GetIds()))

		for i := 0; i < int(toRemove); i++ {
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
			targetLength := uint(len(target.GetIds()))
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
		return nil, nil, err
	}

	return encoding, pairEncoding, nil
}

func PadEncodings(encodings []Encoding, params PaddingParams) []Encoding {
	if len(encodings) == 0 {
		return encodings
	}

	var padLength uint

	switch params.Strategy.Name {
	case "Fixed":
		padLength = params.Strategy.Value.(uint)
	case "BatchLongest":
		var max int = 0
		for _, encoding := range encodings {
			if len(encoding.GetIds()) > max {
				max = len(encoding.GetIds())
			}
		}
		padLength = uint(max)
	}

	// TODO: implement concurrency with for loop
	var newEncodings []Encoding
	for _, e := range encodings {
		e.Pad(padLength, params.PadId, params.PadTypeId, params.PadToken, params.Direction)
		newEncodings = append(newEncodings, e)
	}

	return newEncodings
}
