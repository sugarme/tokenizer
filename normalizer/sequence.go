package normalizer

import ()

// Sequence wraps a slice of normalizers to normalize
// string in sequence.
type Sequence struct {
	Normalizers []Normalizer `json:"normalizers"`
}

func NewSequence(norms []Normalizer) *Sequence {
	return &Sequence{norms}
}

// Implement Normalizer for Sequence
func (s *Sequence) Normalize(normalized *NormalizedString) (*NormalizedString, error) {
	input := normalized
	var err error
	for _, n := range s.Normalizers {
		input, err = n.Normalize(input)
		if err != nil {
			return nil, err
		}
	}

	return input, nil
}
