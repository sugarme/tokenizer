package normalizer

type Strip struct {
	stripLeft  bool
	stripRight bool
}

func NewStrip(stripLeft, stripRight bool) *Strip {
	return &Strip{
		stripLeft:  stripLeft,
		stripRight: stripRight,
	}
}

// Implement Normalizer interface for Strip:
// =========================================

func (s *Strip) Normalize(normalized *NormalizedString) (*NormalizedString, error) {

	if s.stripLeft && s.stripRight {
		return normalized.Strip(), nil
	} else {
		if s.stripLeft {
			return normalized.LStrip(), nil
		}

		if s.stripRight {
			return normalized.RStrip(), nil
		}
	}

	return normalized, nil
}

type StripAccents struct{}

func NewStripAccents() *StripAccents {
	return new(StripAccents)
}

func (sa *StripAccents) Normalize(normalized *NormalizedString) (*NormalizedString, error) {
	return normalized.RemoveAccents(), nil
}
